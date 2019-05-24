package os

import (
	"os"
	"os/exec"
	"path"
	"path/filepath"
	"strings"

	"github.com/bmatcuk/doublestar"
	"github.com/karrick/godirwalk"
	gitignore "github.com/monochromegane/go-gitignore"
)

type WalkContext struct {
	dir                  string
	sort                 bool
	followSymbolicLinks  bool
	postChildrenCallback godirwalk.WalkFunc
	matcher              *Matcher

	callback godirwalk.WalkFunc
	patterns []string
}

// WalkGitIgnore special pattern to ignore all gitignore rules,
// including the ".gitignore" and ".gitignore_global"
const WalkGitIgnore = "!g"

// WalkIgnoreHidden special pattern to ignore all hidden files
const WalkIgnoreHidden = "!**" + string(os.PathSeparator) + ".[^.]*"

// Walk If the pattern begins with "!", it will become a negative filter pattern.
// Each path will be tested against all pattern, each pattern will override the previous
// pattern's match result.
func Walk(patterns ...string) *WalkContext {
	return &WalkContext{
		dir:      ".",
		patterns: patterns,
	}
}

func (ctx *WalkContext) Dir(d string) *WalkContext {
	ctx.dir = d
	return ctx
}

func (ctx *WalkContext) Sort() *WalkContext {
	ctx.sort = true
	return ctx
}

func (ctx *WalkContext) FollowSymbolicLinks() *WalkContext {
	ctx.followSymbolicLinks = true
	return ctx
}

func (ctx *WalkContext) PostChildrenCallback(cb godirwalk.WalkFunc) *WalkContext {
	ctx.postChildrenCallback = cb
	return ctx
}

func (ctx *WalkContext) Matcher(m *Matcher) *WalkContext {
	ctx.matcher = m
	return ctx
}

func (ctx *WalkContext) Do(cb godirwalk.WalkFunc) error {
	ctx.callback = cb

	m := ctx.matcher
	if m == nil {
		var err error
		m, err = NewMatcher(ctx.dir, ctx.patterns)
		if err != nil {
			return err
		}
	}

	return godirwalk.Walk(m.dir, &godirwalk.Options{
		Unsorted:             !ctx.sort,
		FollowSymbolicLinks:  ctx.followSymbolicLinks,
		Callback:             genMatchFn(m, ctx.callback),
		PostChildrenCallback: genMatchFn(m, ctx.postChildrenCallback),
	})
}

// List walk and get list of the paths
func (ctx *WalkContext) List() ([]string, error) {
	list := []string{}
	return list, ctx.Do(func(p string, info *godirwalk.Dirent) error {
		list = append(list, p)
		return nil
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
	m *Matcher,
	cb godirwalk.WalkFunc,
) godirwalk.WalkFunc {
	return func(p string, info *godirwalk.Dirent) (resErr error) {
		matched, negative, err := m.Match(p, info.IsDir())
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

type Matcher struct {
	dir           string
	gitMatchers   map[string]gitignore.IgnoreMatcher
	gitSubmodules []string
	patterns      []string
}

func NewMatcher(dir string, patterns []string) (*Matcher, error) {
	dir, err := filepath.Abs(dir)
	if err != nil {
		return nil, err
	}

	homeDir := HomeDir()
	gs := map[string]gitignore.IgnoreMatcher{}
	var submodules []string
	if hasWalkGitIgnore(patterns) {
		submodules = getGitSubmodules()
		gPath := path.Join(homeDir, ".gitignore_global")
		g, err := gitignore.NewGitIgnore(gPath, dir)
		if err == nil {

			gs[gPath] = g
		}
	}

	return &Matcher{
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

func (m *Matcher) gitMatch(p string, isDir bool) bool {
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

		gPath := path.Join(p, ".gitignore")
		if _, has := m.gitMatchers[gPath]; !has {
			g, err := gitignore.NewGitIgnore(gPath)
			if err == nil {
				m.gitMatchers[gPath] = g
			}
		}
	}

	for _, g := range m.gitMatchers {
		if g.Match(p, isDir) {
			return true
		}
	}
	return false
}

func (m *Matcher) Match(p string, isDir bool) (matched, negative bool, err error) {
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
