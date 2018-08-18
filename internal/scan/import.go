package scan

import (
	"fmt"
	"go/build"
	"os"
	"path/filepath"
	"runtime"
	"strings"
)

func GetPkgs(importPath string) (PkgPack, error) {
	pack := make(PkgPack)
	pkg, err := findPkg(importPath, gopaths()...)
	if err != nil {
		return nil, err
	}
	pack[importPath] = pkg

	vendorFolder, err := detectVendorFolder(pkg.Dir)
	if err != nil {
		return nil, err
	}

	var gopath []string
	if vendorFolder != "" {
		gopath = append([]string{vendorFolder}, gopaths()...)
	} else {
		gopath = gopaths()
	}

	for _, dep := range pkg.Imports {
		pack, err = getPkgs(pack, dep, gopath...)
		if err != nil {
			return nil, err
		}
	}
	return pack, nil
}

func detectVendorFolder(path string) (string, error) {
	vendorFolder := filepath.Join(path, "vendor")
	fInfo, err := os.Stat(vendorFolder)
	if err != nil {
		fmt.Println("stat:", err)
		return "", err
	}
	if !fInfo.IsDir() {
		return "", nil
	}
	return vendorFolder, nil
}

func getPkgs(found PkgPack, importPath string, gopath ...string) (PkgPack, error) {
	if importPath == "C" {
		return found, nil
	}

	if _, ok := found[importPath]; ok {
		return found, nil
	}

	pkg, err := findPkg(importPath, gopath...)
	if err != nil {
		return found, err
	}
	//	if pkg.Goroot {
	//		return found, err
	//	}
	found[importPath] = pkg

	for _, imp := range pkg.Imports {
		if _, ok := found[importPath]; ok {
			continue
		}
		_, err := getPkgs(found, imp, gopath...)
		if err != nil {
			return found, err
		}
	}
	return found, nil
}

func findPkg(importPath string, gopath ...string) (*build.Package, error) {
	for _, gopath := range gopath {
		pkg, err := build.Import(importPath, gopath, 0)
		if err != nil {
			return nil, err
		}
		return pkg, nil
	}
	return nil, fmt.Errorf("scan.findpkg: %s", importPath)
}

func gopaths() []string {
	paths := os.Getenv("GOPATH")
	pathListSep := string(os.PathListSeparator)

	pathlist := []string{paths}
	if strings.Contains(paths, pathListSep) {
		pathlist = strings.Split(paths, pathListSep)
	}
	return append(pathlist, runtime.GOROOT())
}
