package scan

import (
	"go/build"
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
