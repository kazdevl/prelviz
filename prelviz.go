package prelviz

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"

	"github.com/awalterschulze/gographviz"
	"github.com/samber/lo"
)

type Prelviz struct {
	projectModuleName string
	packageInfoMap    map[string]*PackageInfo
	config            *Config
	output            io.Writer
	dotLayout         string
}

type NodeInfo struct {
	Name                  string
	DirectoryPath         string
	groupingDirectoryPath string
	ContainsPackageNum    int
}

func NewPrelviz(projectDirectoryPath, outputFilePath, dotLayout string) (*Prelviz, error) {
	config, err := NewConfig(projectDirectoryPath)
	if err != nil {
		return nil, err
	}

	name, err := GetModuleName(projectDirectoryPath)
	if err != nil {
		return nil, err
	}

	packageInfoMap, err := NewPackageInfoMap(projectDirectoryPath)
	if err != nil {
		return nil, err
	}

	var output io.Writer
	if outputFilePath == "" {
		output = os.Stdout
	} else {
		output, err = os.Create(outputFilePath)
		if err != nil {
			return nil, err
		}
	}

	return &Prelviz{
		projectModuleName: name,
		packageInfoMap:    packageInfoMap,
		config:            config,
		output:            output,
		dotLayout:         dotLayout,
	}, nil
}

func (m *Prelviz) Run() error {
	// add graph
	var (
		graphDefaultAttrs = map[string]string{
			"charset":   `"UTF-8"`,
			"label":     `"package relation"`,
			"labelloc":  `"t"`,
			"labeljust": `"c"`,
			"bgcolor":   `"#343434"`,
			"fontsize":  "18",
			"fontcolor": `"white"`,
			"style":     `"filled"`,
			"rankdir":   `"TB"`,
			"margin":    "0.5",
			"layout":    fmt.Sprintf(`"%s"`, m.dotLayout),
		}
		nodeDefaultAttrs = map[string]string{
			"shape":       `"record"`,
			"style":       `"solid,filled"`,
			"fontcolor":   "6",
			"fontsize":    "14",
			"colorscheme": `"spectral11"`,
		}
	)
	graphAst, _ := gographviz.ParseString(`digraph d {}`)
	graph := gographviz.NewGraph()
	if err := gographviz.Analyse(graphAst, graph); err != nil {
		return err
	}
	graphAttrs, err := gographviz.NewAttrs(graphDefaultAttrs)
	if err != nil {
		return err
	}
	graph.Attrs.Extend(graphAttrs)

	// add node
	for nodeName, info := range m.nodeInfoMap() {
		if info.isGroupingNode() {
			if graph.IsNode(nodeName) {
				continue
			}
			if err = graph.AddNode("G", m.toDotLangFormat(nodeName), lo.Assign(
				nodeDefaultAttrs,
				map[string]string{
					"fillcolor": "9",
					"label":     fmt.Sprintf(`"{path: %s|pkgNum: %d}"`, info.DirectoryPath, info.ContainsPackageNum),
				},
			)); err != nil {
				return err
			}
		} else {
			if err = graph.AddNode("G", m.toDotLangFormat(nodeName), lo.Assign(
				nodeDefaultAttrs,
				map[string]string{
					"fillcolor": "10",
					"label":     fmt.Sprintf(`"{pkg: %s|path: %s}"`, info.Name, info.DirectoryPath),
				},
			)); err != nil {
				return err
			}
		}
	}

	// add edge
	for srcNodeName, relationMap := range m.nodeRelationCountMap() {
		for dstNodeName, relationNum := range relationMap {
			if m.isNgRelation(srcNodeName, dstNodeName) {
				if err = graph.AddEdge(m.toDotLangFormat(srcNodeName), m.toDotLangFormat(dstNodeName), true, map[string]string{
					"color":     `"red"`,
					"weight":    fmt.Sprintf(`"%d"`, relationNum),
					"label":     fmt.Sprintf(`"dep:%d"`, relationNum),
					"fontcolor": `"white"`,
					"decorate":  `"true"`,
				}); err != nil {
					return err
				}
			} else {
				if err = graph.AddEdge(m.toDotLangFormat(srcNodeName), m.toDotLangFormat(dstNodeName), true, map[string]string{
					"color":     `"white"`,
					"weight":    fmt.Sprintf(`"%d"`, relationNum),
					"label":     fmt.Sprintf(`"dep:%d"`, relationNum),
					"fontcolor": `"white"`,
					"decorate":  `"true"`,
				}); err != nil {
					return err
				}
			}
		}
	}

	if _, err = fmt.Fprint(m.output, graph.String()); err != nil {
		return err
	}
	return nil
}

func (m *Prelviz) nodeRelationCountMap() map[string]map[string]int {
	nodeRelationCountMap := make(map[string]map[string]int)
	for pkgDirPath, info := range m.packageInfoMap {
		if m.isExcludePackage(pkgDirPath) {
			continue
		}

		nodeName := m.nodeName(pkgDirPath)
		for importPath, usageMap := range info.ImportUsageMap {
			if !m.isTargetPackage(importPath) {
				continue
			}
			importPathNodeName := m.importPathNodeName(importPath)
			if importPathNodeName == nodeName {
				continue
			}

			if _, ok := nodeRelationCountMap[nodeName]; !ok {
				nodeRelationCountMap[nodeName] = map[string]int{
					importPathNodeName: len(usageMap),
				}
			} else {
				nodeRelationCountMap[nodeName][importPathNodeName] += len(usageMap)
			}
		}
	}
	return nodeRelationCountMap
}

func (m *Prelviz) nodeInfoMap() map[string]*NodeInfo {
	nodeInfoMap := make(map[string]*NodeInfo)
	for pkgDirPath, info := range m.packageInfoMap {
		if m.isExcludePackage(pkgDirPath) {
			continue
		}

		nodeName := m.nodeName(pkgDirPath)
		if _, ok := nodeInfoMap[nodeName]; ok {
			nodeInfoMap[nodeName].ContainsPackageNum++
		} else {
			nodeInfoMap[nodeName] = &NodeInfo{
				Name:          info.Name,
				DirectoryPath: m.groupingPackageDirectoryPath(pkgDirPath),
			}
		}
	}
	return nodeInfoMap
}

func (m *Prelviz) importPathNodeName(importPath string) string {
	dirPath := strings.TrimPrefix(importPath, m.projectModuleName)
	dirPath = strings.TrimPrefix(dirPath, "/")
	return m.nodeName(dirPath)
}

func (m *Prelviz) isGroupingNode(pkgDirPath string) bool {
	if m.config == nil {
		return false
	}
	return m.config.IsGroupingPackage(pkgDirPath)
}

func (m *Prelviz) groupingPackageDirectoryPath(pkgDirPath string) string {
	return m.config.GroupingPackageDirectoryPath(pkgDirPath)
}

func (m *Prelviz) nodeName(pkgDirPath string) string {
	nodeName := filepath.Join(m.projectModuleName, pkgDirPath)
	if m.config.IsGroupingPackage(pkgDirPath) {
		return filepath.Join(m.projectModuleName, m.config.GroupingPackageDirectoryPath(pkgDirPath))
	}

	return nodeName
}

func (m *Prelviz) isTargetPackage(importPath string) bool {
	return strings.Contains(importPath, m.projectModuleName)
}

func (m *Prelviz) toDotLangFormat(in string) string {
	return fmt.Sprintf(`"%s"`, in)
}

func (m *Prelviz) isNgRelation(from, to string) bool {
	return m.config.IsNgRelation(from, to)
}

func (m *Prelviz) isExcludePackage(pkgDirPath string) bool {
	pkg := filepath.Join(m.projectModuleName, pkgDirPath)
	return m.config.IsExcludePackage(pkg)
}

func (m *NodeInfo) isGroupingNode() bool {
	return m.ContainsPackageNum > 0
}
