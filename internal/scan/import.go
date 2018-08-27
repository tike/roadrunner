package scan

import (
	"fmt"
	"go/build"
	"log"
	"os"
	"path/filepath"
	"runtime"
	"strings"
)

func GetPkgs(importPath string) (PkgPack, string, error) {
	pack := make(PkgPack)
	pkg, err := findPkg(importPath, gopaths()...)
	if err != nil {
		return nil, "", err
	}
	pack[importPath] = pkg

	vendorFolder := detectVendorFolder(pkg.Dir)
	if vendorFolder != "" {
		log.Printf("vendor folder detected: %s", vendorFolder)
	} else {
		log.Printf("no vendor folder found.")
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
			return nil, "", err
		}
	}
	return pack, pkg.ImportPath, nil
}

func detectVendorFolder(path string) string {
	for prfx := path; filepath.Base(prfx) != "src"; prfx = filepath.Dir(prfx) {
		vendorFolder := filepath.Join(prfx, "vendor")
		fInfo, err := os.Stat(vendorFolder)
		if err != nil {
			continue
		}
		if fInfo.IsDir() {
			return vendorFolder
		}
		log.Println("vendorFolder", prfx)
	}
	return ""
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
	if importPath[:1] == "." {
		AbsImportPath, err := filepath.Abs(importPath)
		if err != nil {
			return nil, err
		}
		pkg, err := build.ImportDir(AbsImportPath, 0)
		if err != nil {
			return nil, err
		}
		return pkg, nil
	}
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
