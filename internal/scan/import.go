package scan

import (
	"fmt"
	"go/build"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"regexp"
	"strings"
)

func GetPkgs(importPath string) (PkgPack, string, error) {
	pkgBase := gomodRoot(importPath)
	if pkgBase == "" {
		return gopathBuild(importPath)
	}
	return gomodBuild(pkgBase, importPath)
}

func gomodBuild(pkgBase, importPath string) (PkgPack, string, error) {
	if build.IsLocalImport(importPath) {
		importPath = filepath.Join(pkgBase, importPath)
	}

	pkg, err := build.Import(importPath, pkgBase, 0)
	if err != nil {
		return nil, "", err
	}

	pack := make(PkgPack)
	pack, err = gomodImports(pack, pkgBase, pkg)
	if err != nil {
		return nil, "", err
	}

	return pack, pkg.ImportPath, nil
}

func gomodImports(pack PkgPack, pkgBase string, pkg *build.Package) (PkgPack, error) {
	pack[pkg.ImportPath] = pkg
	for _, dep := range pkg.Imports {
		if !strings.HasPrefix(dep, pkgBase) {
			continue
		}
		if _, ok := pack[dep]; ok {
			continue
		}

		depPack, err := build.Import(dep, pkgBase, 0)
		if err != nil {
			fmt.Println(err)
			continue
		}
		gomodImports(pack, pkgBase, depPack)
	}
	return pack, nil
}

func gomodRoot(path string) string {
	for prfx := path; filepath.Base(prfx) != "src"; prfx = filepath.Dir(prfx) {
		gomod := filepath.Join(prfx, "go.mod")
		_, err := os.Stat(gomod)
		if err != nil {
			continue
		}
		log.Println("go.mod found:", gomod)

		dat, err := ioutil.ReadFile(gomod)
		if err != nil {
			panic(err)
		}

		matches := regexp.MustCompile("module ([^\n]+)").FindSubmatch(dat)
		if matches == nil || len(matches) < 2 {
			fmt.Println("matches", matches)
			os.Exit(1)
		}

		pkgbase := string(matches[1])
		log.Println("pkg basepath", pkgbase)
		return pkgbase

	}
	return ""
}

func gopathBuild(importPath string) (PkgPack, string, error) {
	pack := make(PkgPack)
	pkg, err := findPkg(importPath, gopaths()...)
	if err != nil {
		return nil, "", err
	}
	pack[importPath] = pkg

	gopath := gopaths()
	vendorFolder := detectVendorFolder(pkg.Dir)

	if vendorFolder != "" {
		log.Printf("vendor folder detected: %s", vendorFolder)
		gopath = append([]string{vendorFolder}, gopaths()...)
	}

	for _, dep := range pkg.Imports {
		pack, err = getPkgs(pack, dep, "", gopath...)
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

func getPkgs(found PkgPack, importPath string, goMod string, gopath ...string) (PkgPack, error) {
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

	found[importPath] = pkg

	for _, imp := range pkg.Imports {
		if _, ok := found[imp]; ok {
			continue
		}
		if goMod != "" && !strings.HasPrefix(imp, goMod) {
			continue
		}

		_, err := getPkgs(found, imp, goMod, gopath...)
		if err != nil {
			return found, err
		}
	}
	return found, nil
}

func findPkg(importPath string, gopath ...string) (*build.Package, error) {
	if build.IsLocalImport(importPath) {
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
	return pathlist
}
