package corpus

import (
	"os"
	"path/filepath"
)

type Walker struct {
	Exclude, Include []string
	MinSize          int64
}

func matchAny(s string, ref []string) bool {
	for _, p := range ref {
		if ok, err := filepath.Match(p, s); err == nil && ok {
			return ok
		}
	}
	return false
}

func matchIncludeExclude(s string, incl, excl []string) bool {
	if incl != nil && !matchAny(s, incl) {
		return false
	}
	if excl != nil && matchAny(s, excl) {
		return false
	}
	return true
}

func (w *Walker) matches(path string, info os.FileInfo) bool {
	if !matchIncludeExclude(path, w.Include, w.Exclude) {
		return false
	}
	if info.Mode().IsRegular() && info.Size() < w.MinSize {
		return false
	}
	return true
}

func (w *Walker) Walk(root string, walkFn filepath.WalkFunc) error {
	return filepath.Walk(root, func(path string, info os.FileInfo, err error) error {
		if info.IsDir() {
			// Only prune directories if they are
			// excluded. But don't skip the root
			// directory, even if it matches an exclude
			// pattern (think ".").
			if path != root && matchAny(filepath.Base(path), w.Exclude) {
				return filepath.SkipDir
			}
			return nil
		}
		if !w.matches(filepath.Base(path), info) {
			return nil
		}

		if !info.Mode().IsRegular() {
			return nil
		}

		walkFn(path, info, err)
		return nil
	})
}
