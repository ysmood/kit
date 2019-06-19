package os

import (
	"os"
	"os/exec"
	"path"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/bmatcuk/doublestar"
	"github.com/karrick/godirwalk"
	gitignore "github.com/monochromegane/go-gitignore"
	"github.com/ysmood/gokit/pkg/utils"
)

// WalkContext ...
type WalkContext struct {
	dir                  string
	sort                 bool
	followSymbolicLinks  bool
	postChildrenCallback WalkFunc
	matcher              *Matcher

	callback WalkFunc
	patterns []string
}

// WalkGitIgnore special pattern to ignore all gitignore rules,
// including the ".gitignore" and ".gitignore_global"
const WalkGitIgnore = "!g"

// WalkIgnoreHidden special pattern to ignore all hidden files
const WalkIgnoreHidden = "!**" + string(os.PathSeparator) + ".[^.]*"

// WalkFunc ...
type WalkFunc = godirwalk.WalkFunc

// WalkDirent ...
type WalkDirent = *godirwalk.Dirent

// Walk Set up the walk, need to call Do to actually run it.
// If the pattern begins with "!", it will become a negative filter pattern.
// Each path will be tested against all pattern, each pattern will override the previous
// pattern's match result.
func Walk(patterns ...string) *WalkContext {
	return &WalkContext{
		dir:      ".",
		patterns: patterns,
	}
}

// Dir set dir
func (ctx *WalkContext) Dir(d string) *WalkContext {
	ctx.dir = d
	return ctx
}

// Sort whether to sort the result or not
func (ctx *WalkContext) Sort() *WalkContext {
	ctx.sort = true
	return ctx
}

// FollowSymbolicLinks ...
func (ctx *WalkContext) FollowSymbolicLinks() *WalkContext {
	ctx.followSymbolicLinks = true
	return ctx
}

// PostChildrenCallback ...
func (ctx *WalkContext) PostChildrenCallback(cb WalkFunc) *WalkContext {
	ctx.postChildrenCallback = cb
	return ctx
}

// Matcher ...
func (ctx *WalkContext) Matcher(m *Matcher) *WalkContext {
	ctx.matcher = m
	return ctx
}

// Do execute walk
func (ctx *WalkContext) Do(cb WalkFunc) error {
	ctx.callback = cb

	m := ctx.matcher
	if m == nil {
		m = NewMatcher(ctx.dir, ctx.patterns)
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

// MustList ...
func (ctx *WalkContext) MustList() []string {
	return utils.E(ctx.List())[0].([]string)
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
	cb WalkFunc,
) WalkFunc {
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

// Matcher ...
type Matcher struct {
	dir           string
	gitMatchers   map[string]gitignore.IgnoreMatcher
	gitSubmodules []string
	patterns      []string
}

// NewMatcher ...
func NewMatcher(dir string, patterns []string) *Matcher {
	dir, err := filepath.Abs(dir)
	utils.E(err)

	homeDir := HomeDir()
	gs := map[string]gitignore.IgnoreMatcher{}
	var submodules []string
	if hasWalkGitIgnore(patterns) {
		cmd := exec.Command("git", "rev-parse", "--show-toplevel")
		cmd.Dir = dir
		out, err := cmd.CombinedOutput()
		if err == nil {
			gitRoot := strings.TrimSpace(string(out))

			submodules = getGitSubmodules(dir)
			gPath := path.Join(homeDir, ".gitignore_global")
			g, err := gitignore.NewGitIgnore(gPath, dir)
			if err == nil {
				gs[gPath] = g
			}

			// check all parents
			p := dir
			for {
				addIgnoreFile(p, gs)
				if p == gitRoot || p == "/" {
					break
				}
				p = filepath.Dir(p)
			}
		}
	}

	return &Matcher{
		dir:           dir,
		gitMatchers:   gs,
		gitSubmodules: submodules,
		patterns:      patterns,
	}
}

var submoduleReg = regexp.MustCompile(`\A [a-f0-9]+ (.+) \(.+\)\z`)

func getGitSubmodules(dir string) []string {
	cmd := exec.Command("git", "submodule")
	cmd.Dir = dir
	out, err := cmd.CombinedOutput()
	if err != nil {
		return nil
	}

	list := []string{}
	for _, l := range strings.Split(string(out), "\n") {
		m := submoduleReg.FindStringSubmatch(l)

		if len(m) > 1 {
			p := filepath.Join(dir, m[1])
			list = append(list, p)
		}
	}

	return list
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

		addIgnoreFile(p, m.gitMatchers)
	}

	for _, g := range m.gitMatchers {
		if g.Match(p, isDir) {
			return true
		}
	}
	return false
}

func addIgnoreFile(p string, gs map[string]gitignore.IgnoreMatcher) {
	gPath := path.Join(p, ".gitignore")
	if _, has := gs[gPath]; !has {
		g, err := gitignore.NewGitIgnore(gPath)
		if err == nil {
			gs[gPath] = g
		}
	}
}

// Match ...
func (m *Matcher) Match(p string, isDir bool) (matched, negative bool, err error) {
	for _, pattern := range m.patterns {
		if pattern == WalkGitIgnore {
			if m.gitMatch(p, isDir) {
				matched = false
				negative = true
			}
			continue
		}

		mm, neg, e := pathMatch(pattern, m.dir, p)

		if e != nil {
			err = e
			return
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

	return
}

func pathMatch(pattern, dir, path string) (bool, bool, error) {
	name := path[len(dir):]
	nameLen := len(name)

	if nameLen > 0 && name[0] == os.PathSeparator {
		name = name[1:]
	}

	negative := false

	if pattern[0] == '!' {
		pattern = pattern[1:]
		negative = true
	}

	if nameLen == 0 {
		return false, negative, nil
	}

	matched, err := doublestar.PathMatch(pattern, name)
	if err != nil {
		return false, false, err
	}

	return matched, negative, nil
}
