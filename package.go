package prelviz

import (
	"go/ast"
	"go/parser"
	"go/token"
	"os"
	"path/filepath"
	"strings"

	"github.com/samber/lo"
)

type PackageInfo struct {
	Name           string
	DirectoryPath  string
	ImportUsageMap map[string]map[string]struct{}
}

type PackageInfoMap map[string]*PackageInfo

func NewPackageInfoMap(projectDirectoryPath string) (map[string]*PackageInfo, error) {
	filePaths, err := targetGoFilePaths(projectDirectoryPath)
	if err != nil {
		return nil, err
	}

	packageInfoMap := make(map[string]*PackageInfo)
	for _, filePath := range filePaths {
		var packageInfo *PackageInfo
		packageInfo, err = NewPackageInfo(filePath, projectDirectoryPath)
		if err != nil {
			return nil, err
		}

		if info, ok := packageInfoMap[packageInfo.DirectoryPath]; ok {
			info.ImportUsageMap = lo.Assign(info.ImportUsageMap, packageInfo.ImportUsageMap)
		} else {
			packageInfoMap[packageInfo.DirectoryPath] = packageInfo
		}
	}
	return packageInfoMap, nil
}

func NewPackageInfo(filePath, projectDirectoryPath string) (*PackageInfo, error) {
	fset := token.NewFileSet()
	f, err := parser.ParseFile(fset, filePath, nil, 0)
	if err != nil {
		return nil, err
	}

	importUsageMap := make(map[string]map[string]struct{})
	importUsageNameMap := make(map[string]string)
	ast.Inspect(f, func(n ast.Node) bool {
		switch x := n.(type) {
		case *ast.ImportSpec:
			importPath := strings.Trim(x.Path.Value, `"`)
			if x.Name != nil {
				importUsageNameMap[x.Name.Name] = importPath
			} else {
				importUsageNameMap[filepath.Base(importPath)] = importPath
			}
		case *ast.SelectorExpr:
			xIndent, ok := x.X.(*ast.Ident)
			if !ok {
				return true
			}
			var importPath string
			if importPath, ok = importUsageNameMap[xIndent.Name]; ok {
				if _, ok = importUsageMap[importPath]; ok {
					importUsageMap[importPath][x.Sel.Name] = struct{}{}
				} else {
					importUsageMap[importPath] = map[string]struct{}{x.Sel.Name: {}}
				}
			}

		}
		return true
	})

	relativeFilePath, err := filepath.Rel(projectDirectoryPath, filePath)
	if err != nil {
		return nil, err
	}
	return &PackageInfo{
		Name:           f.Name.Name,
		ImportUsageMap: importUsageMap,
		DirectoryPath:  filepath.Dir(relativeFilePath),
	}, nil
}

func targetGoFilePaths(dir string) ([]string, error) {
	filePaths := make([]string, 0)
	if err := filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {
		if strings.HasSuffix(path, "_test.go") {
			return nil
		}
		if strings.HasSuffix(path, ".go") {
			filePaths = append(filePaths, path)
		}
		return nil
	}); err != nil {
		return nil, err
	}
	return filePaths, nil
}
