package gokit

import (
	"fmt"
	"os/user"
	"path"
	"path/filepath"

	"github.com/bmatcuk/doublestar"
	"github.com/karrick/godirwalk"
	"github.com/monochromegane/go-gitignore"
)

// WalkOptions ...
type WalkOptions struct {
	Dir                  string
	Sorted               bool
	FollowSymbolicLinks  bool
	PostChildrenCallback godirwalk.WalkFunc
}

// WalkGitIgnore special pattern to ignore all gitignore rules,
// including the ".gitignore" and ".gitignore_global"
const WalkGitIgnore = "!g"

// WalkHidden special pattern to match all hidden files
const WalkHidden = "!**/.[^.]*"

// Walk If the pattern begins with "!", it will become a negative filter pattern.
// Each path will be tested against all pattern, each pattern will override the previous
// pattern's match result.
func Walk(patterns []string, cb godirwalk.WalkFunc, opts *WalkOptions) error {
	if opts == nil {
		opts = &WalkOptions{
			Dir: ".",
		}
	}

	dir := opts.Dir

	if dir == "" {
		dir = "."
	}

	var gitIgnorer *gitignorer

	patterns = normalizePatterns(dir, patterns)

	if hasWalkGitIgnore(patterns) {
		var err error
		gitIgnorer, err = newGitignorer(dir)
		if err != nil {
			Err(err)
		}
	}

	return godirwalk.Walk(dir, &godirwalk.Options{
		Unsorted:             !opts.Sorted,
		FollowSymbolicLinks:  opts.FollowSymbolicLinks,
		Callback:             genMatchFn(patterns, gitIgnorer, cb),
		PostChildrenCallback: genMatchFn(patterns, gitIgnorer, opts.PostChildrenCallback),
	})
}

type gitignorer struct {
	gi  gitignore.IgnoreMatcher
	ggi gitignore.IgnoreMatcher
}

func newGitignorer(dir string) (*gitignorer, error) {
	user, err := user.Current()
	if err != nil {
		return nil, err
	}

	gi, err := gitignore.NewGitIgnore(path.Join(dir, ".gitignore"))
	if err != nil {
		return nil, err
	}

	ggi, err := gitignore.NewGitIgnore(path.Join(user.HomeDir, ".gitignore_global"), dir)
	if err != nil {
		return nil, err
	}

	return &gitignorer{gi, ggi}, nil
}

func (g *gitignorer) match(p string, isDir bool) bool {
	return g.gi.Match(p, isDir) || g.ggi.Match(p, isDir)
}

func pathMatch(pattern, name string) (bool, bool, error) {
	negative := false

	if pattern[0] == '!' {
		pattern = pattern[1:]
		negative = true
	}

	matched, err := doublestar.PathMatch(pattern, name)
	if err != nil {
		return false, false, err
	}

	return matched, negative, nil
}

func normalizePatterns(dir string, patterns []string) []string {
	list := []string{}
	for _, p := range patterns {
		list = append(list, path.Clean(fmt.Sprint(dir, "/", p)))
	}
	return list
}

func hasWalkGitIgnore(patterns []string) bool {
	for _, p := range patterns {
		if p == WalkGitIgnore {
			return true
		}
	}

	return false
}

// TODO: for the best performance, this is just a partial hacky solution.
// We need a syntax analyzer to detect if a path is a subpath of a pattern.
func isSubpath(p, pattern string) bool {
	if path.Clean(p) == "." {
		return true
	}
	return len(p) <= len(pattern) && p == pattern[:len(p)]
}

func genMatchFn(
	patterns []string,
	gitIgnorer *gitignorer,
	cb godirwalk.WalkFunc,
) godirwalk.WalkFunc {
	return func(p string, info *godirwalk.Dirent) (resErr error) {
		matched := false
		isDir := info.IsDir()
		for _, pattern := range patterns {
			if pattern == WalkGitIgnore {
				if gitIgnorer != nil && gitIgnorer.match(p, info.IsDir()) {
					if isDir {
						resErr = filepath.SkipDir
					}
					matched = false
				}
				continue
			}

			m, negative, err := pathMatch(pattern, p)

			if err != nil {
				return err
			}

			if m {
				if negative {
					matched = false
				} else {
					matched = true
					resErr = nil
				}
			}

			if isDir && !matched {
				resErr = filepath.SkipDir
			}

			// TODO: for the best performance, this is just a partial hacky solution.
			// We need a syntax analyzer to detect if a path is a subpath of a pattern.
			if resErr == filepath.SkipDir {
				if isSubpath(p, pattern) {
					resErr = nil
				}
			}
		}

		if matched && cb != nil {
			err := cb(p, info)
			if err != nil {
				return err
			}
		}
		return resErr
	}
}
