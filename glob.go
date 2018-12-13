package gokit

import (
	"os"
	"os/user"
	"path"
	"path/filepath"

	"github.com/bmatcuk/doublestar"
	"github.com/monochromegane/go-gitignore"
)

// GlobOptions ...
type GlobOptions struct {
	Dir string
}

// GlobGitIgnore special pattern to ignore all gitignore rules,
// including the ".gitignore" and ".gitignore_global"
const GlobGitIgnore = "!g"

// GlobHidden special pattern to match all hidden files
const GlobHidden = "!**/.[^.]*"

// Glob If the pattern begins with "!", it will become a negative filter pattern.
// Each path will be tested against all pattern, each pattern will override the previous
// pattern's match result.
func Glob(patterns []string, opts *GlobOptions) ([]string, error) {
	if opts == nil {
		opts = &GlobOptions{
			Dir: ".",
		}
	}

	dir := opts.Dir

	if dir == "" {
		dir = "."
	}

	list := []string{}
	var gitIgnorer *gitignorer

	patterns, list = removeAbsPaths(patterns)

	if hasGlobGitIgnore(patterns) {
		var err error
		gitIgnorer, err = newGitignorer(dir)
		if err != nil {
			Err(err)
		}
	}

	filepath.Walk(dir, func(p string, info os.FileInfo, err error) (resErr error) {
		if err != nil {
			return err
		}

		matched := false
		isDir := info.IsDir()
		for _, pattern := range patterns {
			if pattern == GlobGitIgnore {
				if gitIgnorer != nil && gitIgnorer.match(p, info.IsDir()) {
					if isDir {
						resErr = filepath.SkipDir
					}
					matched = false
				}
				continue
			}

			// TODO: for the best performance, this is just a partial hacky solution.
			// We need a syntax analyzer to detect if a path is a subpath of a pattern.
			if resErr == filepath.SkipDir {
				if isSubpath(p, pattern) {
					resErr = nil
				}
			}

			m, negative, err := pathMatch(pattern, p)

			if err != nil {
				return err
			}

			if m {
				if negative {
					matched = false
					if isDir {
						resErr = filepath.SkipDir
					}
				} else {
					matched = true
					resErr = nil
				}
			}
		}

		if matched {
			list = append(list, p)
		}
		return resErr
	})

	return list, nil
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

func removeAbsPaths(patterns []string) ([]string, []string) {
	leftPatterns := []string{}
	paths := []string{}
	for _, p := range patterns {
		if Exists(p) {
			paths = append(paths, p)
		} else {
			leftPatterns = append(leftPatterns, p)
		}
	}

	return leftPatterns, paths
}

func hasGlobGitIgnore(patterns []string) bool {
	for _, p := range patterns {
		if p == GlobGitIgnore {
			return true
		}
	}

	return false
}

// TODO: for the best performance, this is just a partial hacky solution.
// We need a syntax analyzer to detect if a path is a subpath of a pattern.
func isSubpath(p, pattern string) bool {
	return len(p) <= len(pattern) && p == pattern[:len(p)]
}
