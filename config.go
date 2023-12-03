package prelviz

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/samber/lo"
)

type Config struct {
	NgRelations            []NgRelation `json:"ng_relation"`
	GroupingDirectoryPaths []string     `json:"grouping_directory_path"`
}

type NgRelation struct {
	From string   `json:"from"`
	To   []string `json:"to"`
}

type NgPackageRelationMap map[string]map[string]struct{}

func NewConfig(path string) (*Config, error) {
	filePath := filepath.Join(path, ".prelviz.config.json")
	if !fileExists(filePath) {
		return nil, nil
	}

	var c Config
	raw, err := os.ReadFile(filePath)
	if err != nil {
		return nil, err
	}
	if err = json.Unmarshal(raw, &c); err != nil {
		return nil, err
	}
	if c.GroupingDirectoryPaths == nil {
		return &c, nil
	}

	duplicate := lo.FindDuplicates(c.GroupingDirectoryPaths)
	if len(duplicate) > 0 {
		return nil, errors.New("duplicate grouping directory path")
	}

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
	return &c, nil
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
	if c.GroupingDirectoryPaths == nil {
		return pkgDirPath
	}

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

func (c *Config) ToNgPackageRelationMap() NgPackageRelationMap {
	if c.NgRelations == nil {
		return nil
	}

	ngRelationMap := make(NgPackageRelationMap)
	for _, ngRelation := range c.NgRelations {
		if _, ok := ngRelationMap[ngRelation.From]; !ok {
			ngRelationMap[ngRelation.From] = make(map[string]struct{})
		}
		for _, to := range ngRelation.To {
			ngRelationMap[ngRelation.From][to] = struct{}{}
		}
	}
	return ngRelationMap
}

func (n NgPackageRelationMap) IsNgRelation(from string, to string) bool {
	if _, ok := n[from]; !ok {
		return false
	}
	if _, ok := n[from][to]; !ok {
		return false
	}
	return true
}

func fileExists(filePath string) bool {
	_, err := os.Stat(filePath)
	return err == nil
}
