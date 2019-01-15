package gokit

import (
	"os"
	"os/exec"
	"os/user"
	"path"
	"path/filepath"
	"strings"

	"github.com/bmatcuk/doublestar"
	"github.com/karrick/godirwalk"
	gitignore "github.com/monochromegane/go-gitignore"
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
const WalkHidden = "!**" + string(os.PathSeparator) + ".[^.]*"

// Walk If the pattern begins with "!", it will become a negative filter pattern.
// Each path will be tested against all pattern, each pattern will override the previous
// pattern's match result.
func Walk(patterns []string, cb godirwalk.WalkFunc, opts *WalkOptions) error {
	if opts == nil {
		opts = &WalkOptions{
			Dir: ".",
		}
	}

	m, err := newMatcher(opts.Dir, patterns)
	if err != nil {
		return err
	}

	return godirwalk.Walk(m.dir, &godirwalk.Options{
		Unsorted:             !opts.Sorted,
		FollowSymbolicLinks:  opts.FollowSymbolicLinks,
		Callback:             genMatchFn(m, cb),
		PostChildrenCallback: genMatchFn(m, opts.PostChildrenCallback),
	})
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
		if p == WalkGitIgnore {
			list = append(list, p)
			continue
		}
		if p[0] == '!' {
			list = append(list, path.Clean("!"+dir+string(os.PathSeparator)+p[1:]))
			continue
		}
		list = append(list, path.Clean(dir+string(os.PathSeparator)+p))
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

func genMatchFn(
	m *matcher,
	cb godirwalk.WalkFunc,
) godirwalk.WalkFunc {
	return func(p string, info *godirwalk.Dirent) (resErr error) {
		matched, negative, err := m.match(p, info.IsDir())
		if err != nil {
			return err
		}

		if matched && cb != nil {
			err := cb(p, info)
			if err != nil {
				return err
			}
		}

		if negative && info.IsDir() {
			return filepath.SkipDir
		}

		return nil
	}
}

type matcher struct {
	dir           string
	gitMatchers   []gitignore.IgnoreMatcher
	gitSubmodules []string
	patterns      []string
}

func newMatcher(dir string, patterns []string) (*matcher, error) {
	dir, err := filepath.Abs(dir)
	if err != nil {
		return nil, err
	}

	user, err := user.Current()
	if err != nil {
		return nil, err
	}

	gs := []gitignore.IgnoreMatcher{}
	var submodules []string
	if hasWalkGitIgnore(patterns) {
		submodules = getGitSubmodules()
		g, err := gitignore.NewGitIgnore(path.Join(user.HomeDir, ".gitignore_global"), dir)
		if err == nil {
			gs = append(gs, g)
		}
	}

	return &matcher{
		dir:           dir,
		gitMatchers:   gs,
		gitSubmodules: submodules,
		patterns:      normalizePatterns(dir, patterns),
	}, nil
}

func getGitSubmodules() []string {
	out, err := exec.Command("git", "rev-parse", "--show-toplevel").Output()
	if err != nil {
		return nil
	}

	root := strings.TrimSpace(string(out))

	p := filepath.Join(root, filepath.Join(".git", "modules", "*"))

	l, _ := filepath.Glob(p)

	for i, p := range l {
		l[i] = strings.Replace(p, filepath.Join(root, ".git", "modules"), root, 1)
	}

	return l
}

func (m *matcher) gitMatch(p string, isDir bool) bool {
	if isDir {
		if l := len(p); l > 4 && p[len(p)-4:] == ".git" {
			return true
		}

		if m.gitSubmodules != nil {
			for _, sub := range m.gitSubmodules {
				if sub == p {
					return true
				}
			}
		}

		g, err := gitignore.NewGitIgnore(path.Join(p, ".gitignore"))
		if err == nil {
			m.gitMatchers = append(m.gitMatchers, g)
		}
	}

	for _, g := range m.gitMatchers {
		if g.Match(p, isDir) {
			return true
		}
	}
	return false
}

func (m *matcher) match(p string, isDir bool) (matched, negative bool, err error) {
	for _, pattern := range m.patterns {
		if pattern == WalkGitIgnore {
			if m.gitMatch(p, isDir) {
				matched = false
				negative = true
			}
			continue
		}

		mm, neg, err := pathMatch(pattern, p)

		if err != nil {
			return matched, neg, err
		}

		if mm {
			if neg {
				negative = true
				matched = false
			} else {
				matched = true
			}
		}
	}

	return matched, negative, nil
}
