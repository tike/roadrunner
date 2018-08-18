package scan

import (
	"go/build"
	"path/filepath"
	"sort"
)

type pkgPack map[string]*build.Package

func (p pkgPack) SortedNames() []string {
	var names = make(sort.StringSlice, 0, len(p))
	for name := range p {
		names = append(names, name)
	}

	names.Sort()
	return names
}

func (p pkgPack) FileList() []string {
	var names = make(sort.StringSlice, 0, 5*len(p))
	for _, pkg := range p {
		for _, file := range pkg.GoFiles {
			names = append(names, filepath.Join(pkg.ImportPath, file))
		}
	}

	names.Sort()
	return names
}
