package prelviz

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/samber/lo"
)

type ConfigBinder struct {
	NgRelations            []NgRelation `json:"ng_relation"`
	GroupingDirectoryPaths []string     `json:"grouping_directory_path"`
	ExcludePackages        []string     `json:"exclude_package"`
	ExcludeDirectorys      []string     `json:"exclude_directory"`
}

type Config struct {
	NgRelationMap          map[string]map[string]struct{}
	GroupingDirectoryPaths []string
	ExcludePackageMap      map[string]struct{}
}

type NgRelation struct {
	From string   `json:"from"`
	To   []string `json:"to"`
}

const configJsonName = ".prelviz.config.json"

func NewConfig(path, moduleName string) (*Config, error) {
	filePath := filepath.Join(path, configJsonName)
	if !fileExists(filePath) {
		return &Config{
			NgRelationMap:          make(map[string]map[string]struct{}),
			GroupingDirectoryPaths: make([]string, 0),
			ExcludePackageMap:      make(map[string]struct{}),
		}, nil
	}

	var cb ConfigBinder
	raw, err := os.ReadFile(filePath)
	if err != nil {
		return nil, err
	}
	if err = json.Unmarshal(raw, &cb); err != nil {
		return nil, err
	}

	c, err := cb.ToConfig(moduleName)
	if err != nil {
		return nil, err
	}
	return c, nil
}

func (c ConfigBinder) ToConfig(moduleName string) (*Config, error) {
	conf := &Config{
		NgRelationMap:          make(map[string]map[string]struct{}),
		GroupingDirectoryPaths: make([]string, 0),
		ExcludePackageMap:      make(map[string]struct{}),
	}

	if c.NgRelations != nil {
		m := make(map[string]map[string]struct{})
		for _, ngRelation := range c.NgRelations {
			if _, ok := m[ngRelation.From]; !ok {
				m[ngRelation.From] = make(map[string]struct{})
			}
			for _, to := range ngRelation.To {
				m[ngRelation.From][to] = struct{}{}
			}
		}
		conf.NgRelationMap = m
	}

	if c.ExcludePackages != nil {
		m := make(map[string]struct{})
		for _, excludePackage := range c.ExcludePackages {
			if excludePackage == "" {
				continue
			}
			m[excludePackage] = struct{}{}
		}
		conf.ExcludePackageMap = m
	}

	if c.ExcludeDirectorys != nil {
		for _, dir := range c.ExcludeDirectorys {
			if err := filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {
				if info.IsDir() {
					conf.ExcludePackageMap[filepath.Join(moduleName, path)] = struct{}{}
				}
				return nil
			}); err != nil {
				return nil, err
			}
		}
	}

	if c.GroupingDirectoryPaths != nil {
		for _, groupDirPath := range c.GroupingDirectoryPaths {
			if groupDirPath == "" {
				continue
			}
			if _, ok := lo.Find(c.GroupingDirectoryPaths, func(s string) bool {
				if s == groupDirPath {
					return false
				}
				return strings.HasPrefix(s, groupDirPath)
			}); ok {
				return nil, fmt.Errorf("don't set parent-child relationships. %s", groupDirPath)
			}
		}
		conf.GroupingDirectoryPaths = lo.Uniq(c.GroupingDirectoryPaths)
	}

	return conf, nil
}

func (c *Config) IsGroupingPackage(pkgDirPath string) bool {
	for _, groupDirPath := range c.GroupingDirectoryPaths {
		if groupDirPath == "" {
			continue
		}
		if strings.HasPrefix(pkgDirPath, groupDirPath) {
			return true
		}
	}
	return false
}

func (c *Config) GroupingPackageDirectoryPath(pkgDirPath string) string {
	for _, groupDirPath := range c.GroupingDirectoryPaths {
		if groupDirPath == "" {
			continue
		}
		if strings.HasPrefix(pkgDirPath, groupDirPath) {
			return groupDirPath
		}
	}
	return pkgDirPath
}

func (c *Config) IsExcludePackage(pkg string) bool {
	if _, ok := c.ExcludePackageMap[pkg]; ok {
		return true
	}
	return false
}

func (c *Config) IsNgRelation(from string, to string) bool {
	if _, ok := c.NgRelationMap[from]; !ok {
		return false
	}
	if _, ok := c.NgRelationMap[from][to]; !ok {
		return false
	}
	return true
}

func fileExists(filePath string) bool {
	_, err := os.Stat(filePath)
	return err == nil
}
