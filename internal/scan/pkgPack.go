package scan

import (
	"go/build"
	"path/filepath"
	"sort"
)

type PkgPack map[string]*build.Package

func (p PkgPack) SortedNames() []string {
	var names = make(sort.StringSlice, 0, len(p))
	for name := range p {
		names = append(names, name)
	}

	names.Sort()
	return names
}

func (p PkgPack) FileList() []string {
	var names = make(sort.StringSlice, 0, 5*len(p))
	for _, pkg := range p {
		names = append(names, pkg.Dir)
		for _, file := range pkg.GoFiles {
			name := filepath.Join(pkg.Dir, file)
			names = append(names, name)
		}
	}

	names.Sort()
	return names
}
