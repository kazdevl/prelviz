package prelviz

import (
	"reflect"
	"testing"
)

func TestPrelviz_nodeInfoMap(t *testing.T) {
	type fields struct {
		projectModuleName string
		packageInfoMap    map[string]*PackageInfo
		config            *Config
	}
	tests := []struct {
		name   string
		fields fields
		want   map[string]*NodeInfo
	}{
		{
			name: "normal: with no config",
			fields: fields{
				projectModuleName: "mod",
				packageInfoMap: map[string]*PackageInfo{
					"sample/src": {
						Name:          "src",
						DirectoryPath: "sample/src",
					},
					"sample/dst1": {
						Name:          "dst1",
						DirectoryPath: "sample/dst1",
					},
					"sample/dst2": {
						Name:          "dst2",
						DirectoryPath: "sample/dst2",
					},
				},
				config: &Config{
					NgRelationMap:          make(map[string]map[string]struct{}),
					GroupingDirectoryPaths: make([]string, 0),
					ExcludePackageMap:      make(map[string]struct{}),
				},
			},
			want: map[string]*NodeInfo{
				"mod/sample/src": {
					Name:               "src",
					DirectoryPath:      "sample/src",
					IsGrouping:         false,
					ContainsPackageNum: 1,
				},
				"mod/sample/dst1": {
					Name:               "dst1",
					DirectoryPath:      "sample/dst1",
					IsGrouping:         false,
					ContainsPackageNum: 1,
				},
				"mod/sample/dst2": {
					Name:               "dst2",
					DirectoryPath:      "sample/dst2",
					IsGrouping:         false,
					ContainsPackageNum: 1,
				},
			},
		},
		{
			name: "normal: exclude package",
			fields: fields{
				projectModuleName: "mod",
				packageInfoMap: map[string]*PackageInfo{
					"sample/src": {
						Name:          "src",
						DirectoryPath: "sample/src",
					},
					"sample/dst1": {
						Name:          "dst1",
						DirectoryPath: "sample/dst1",
					},
					"sample/dst2": {
						Name:          "dst2",
						DirectoryPath: "sample/dst2",
					},
				},
				config: &Config{
					NgRelationMap:          make(map[string]map[string]struct{}),
					GroupingDirectoryPaths: make([]string, 0),
					ExcludePackageMap: map[string]struct{}{
						"mod/sample/dst2": {},
					},
				},
			},
			want: map[string]*NodeInfo{
				"mod/sample/src": {
					Name:               "src",
					DirectoryPath:      "sample/src",
					IsGrouping:         false,
					ContainsPackageNum: 1,
				},
				"mod/sample/dst1": {
					Name:               "dst1",
					DirectoryPath:      "sample/dst1",
					IsGrouping:         false,
					ContainsPackageNum: 1,
				},
			},
		},
		{
			name: "normal: grouping packages",
			fields: fields{
				projectModuleName: "mod",
				packageInfoMap: map[string]*PackageInfo{
					"sample/src": {
						Name:          "src",
						DirectoryPath: "sample/src",
					},
					"sample/grouping/dst1": {
						Name:          "dst1",
						DirectoryPath: "sample/grouping/dst1",
					},
					"sample/grouping/dst2": {
						Name:          "dst2",
						DirectoryPath: "sample/grouping/dst2",
					},
				},
				config: &Config{
					NgRelationMap:          make(map[string]map[string]struct{}),
					GroupingDirectoryPaths: []string{"sample/grouping"},
					ExcludePackageMap:      make(map[string]struct{}),
				},
			},
			want: map[string]*NodeInfo{
				"mod/sample/src": {
					Name:               "src",
					DirectoryPath:      "sample/src",
					IsGrouping:         false,
					ContainsPackageNum: 1,
				},
				"mod/sample/grouping": {
					DirectoryPath:      "sample/grouping",
					IsGrouping:         true,
					ContainsPackageNum: 2,
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := &Prelviz{
				projectModuleName: tt.fields.projectModuleName,
				packageInfoMap:    tt.fields.packageInfoMap,
				config:            tt.fields.config,
			}
			if got := m.nodeInfoMap(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("nodeInfoMap() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestPrelviz_nodeRelationCountMap(t *testing.T) {
	type fields struct {
		projectModuleName string
		packageInfoMap    map[string]*PackageInfo
		config            *Config
	}
	tests := []struct {
		name   string
		fields fields
		want   map[string]map[string]int
	}{
		{
			name: "normal: with no config",
			fields: fields{
				projectModuleName: "mod",
				packageInfoMap: map[string]*PackageInfo{
					"sample/src": {
						Name:          "src",
						DirectoryPath: "sample/src",
						ImportUsageMap: map[string]map[string]struct{}{
							"mod/sample/dst1": {"Sample1": {}, "Sample2": {}},
							"mod/sample/dst2": {"Sample3": {}, "Sample4": {}, "Sample5": {}},
						},
					},
					"sample/dst1": {
						Name:          "dst1",
						DirectoryPath: "sample/dst1",
						ImportUsageMap: map[string]map[string]struct{}{
							"mod/sample/dst2": {"Sample6": {}},
						},
					},
					"sample/dst2": {
						Name:           "dst2",
						DirectoryPath:  "sample/dst2",
						ImportUsageMap: map[string]map[string]struct{}{},
					},
				},
				config: &Config{
					NgRelationMap:          make(map[string]map[string]struct{}),
					GroupingDirectoryPaths: make([]string, 0),
					ExcludePackageMap:      make(map[string]struct{}),
				},
			},
			want: map[string]map[string]int{
				"mod/sample/src": {
					"mod/sample/dst1": 2,
					"mod/sample/dst2": 3,
				},
				"mod/sample/dst1": {
					"mod/sample/dst2": 1,
				},
			},
		},
		{
			name: "normal: include not target package",
			fields: fields{
				projectModuleName: "mod",
				packageInfoMap: map[string]*PackageInfo{
					"sample/src": {
						Name:          "src",
						DirectoryPath: "sample/src",
						ImportUsageMap: map[string]map[string]struct{}{
							"mod/sample/dst1": {"Sample1": {}, "Sample2": {}},
							"mod/sample/dst2": {"Sample3": {}, "Sample4": {}, "Sample5": {}},
							"fmt":             {"Println": {}, "Printf": {}},
						},
					},
					"sample/dst1": {
						Name:          "dst1",
						DirectoryPath: "sample/dst1",
						ImportUsageMap: map[string]map[string]struct{}{
							"mod/sample/dst2": {"Sample6": {}},
						},
					},
					"sample/dst2": {
						Name:          "dst2",
						DirectoryPath: "sample/dst2",
						ImportUsageMap: map[string]map[string]struct{}{
							"fmt": {"Println": {}, "Printf": {}},
						},
					},
				},
				config: &Config{
					NgRelationMap:          make(map[string]map[string]struct{}),
					GroupingDirectoryPaths: make([]string, 0),
					ExcludePackageMap:      make(map[string]struct{}),
				},
			},
			want: map[string]map[string]int{
				"mod/sample/src": {
					"mod/sample/dst1": 2,
					"mod/sample/dst2": 3,
				},
				"mod/sample/dst1": {
					"mod/sample/dst2": 1,
				},
			},
		},
		{
			name: "normal: exclude package",
			fields: fields{
				projectModuleName: "mod",
				packageInfoMap: map[string]*PackageInfo{
					"sample/src": {
						Name:          "src",
						DirectoryPath: "sample/src",
						ImportUsageMap: map[string]map[string]struct{}{
							"mod/sample/dst1": {"Sample1": {}, "Sample2": {}},
							"mod/sample/dst2": {"Sample3": {}, "Sample4": {}, "Sample5": {}},
							"fmt":             {"Println": {}, "Printf": {}},
						},
					},
					"sample/dst1": {
						Name:          "dst1",
						DirectoryPath: "sample/dst1",
						ImportUsageMap: map[string]map[string]struct{}{
							"mod/sample/dst2": {"Sample6": {}},
						},
					},
					"sample/dst2": {
						Name:          "dst2",
						DirectoryPath: "sample/dst2",
						ImportUsageMap: map[string]map[string]struct{}{
							"fmt": {"Println": {}, "Printf": {}},
						},
					},
				},
				config: &Config{
					NgRelationMap:          make(map[string]map[string]struct{}),
					GroupingDirectoryPaths: make([]string, 0),
					ExcludePackageMap: map[string]struct{}{
						"mod/sample/dst2": {},
					},
				},
			},
			want: map[string]map[string]int{
				"mod/sample/src": {
					"mod/sample/dst1": 2,
				},
			},
		},
		{
			name: "normal: grouping packages",
			fields: fields{
				projectModuleName: "mod",
				packageInfoMap: map[string]*PackageInfo{
					"sample/src": {
						Name:          "src",
						DirectoryPath: "sample/src",
						ImportUsageMap: map[string]map[string]struct{}{
							"mod/sample/grouping/dst1": {"Sample1": {}, "Sample2": {}},
							"mod/sample/grouping/dst2": {"Sample3": {}, "Sample4": {}, "Sample5": {}},
							"fmt":                      {"Println": {}, "Printf": {}},
						},
					},
					"sample/grouping/dst1": {
						Name:          "dst1",
						DirectoryPath: "sample/grouping/dst1",
						ImportUsageMap: map[string]map[string]struct{}{
							"mod/sample/grouping/dst2": {"Sample6": {}},
						},
					},
					"sample/grouping/dst2": {
						Name:          "dst2",
						DirectoryPath: "sample/grouping/dst2",
						ImportUsageMap: map[string]map[string]struct{}{
							"fmt": {"Println": {}, "Printf": {}},
						},
					},
				},
				config: &Config{
					NgRelationMap:          make(map[string]map[string]struct{}),
					GroupingDirectoryPaths: []string{"sample/grouping"},
					ExcludePackageMap:      make(map[string]struct{}),
				},
			},
			want: map[string]map[string]int{
				"mod/sample/src": {
					"mod/sample/grouping": 5,
				},
			},
		},
		{
			name: "normal: grouping packages and exclude package",
			fields: fields{
				projectModuleName: "mod",
				packageInfoMap: map[string]*PackageInfo{
					"sample/src": {
						Name:          "src",
						DirectoryPath: "sample/src",
						ImportUsageMap: map[string]map[string]struct{}{
							"mod/sample/grouping/dst1": {"Sample1": {}, "Sample2": {}},
							"mod/sample/grouping/dst2": {"Sample3": {}, "Sample4": {}, "Sample5": {}},
							"fmt":                      {"Println": {}, "Printf": {}},
						},
					},
					"sample/grouping/dst1": {
						Name:          "dst1",
						DirectoryPath: "sample/grouping/dst1",
						ImportUsageMap: map[string]map[string]struct{}{
							"mod/sample/grouping/dst2": {"Sample6": {}},
						},
					},
					"sample/grouping/dst2": {
						Name:          "dst2",
						DirectoryPath: "sample/grouping/dst2",
						ImportUsageMap: map[string]map[string]struct{}{
							"fmt": {"Println": {}, "Printf": {}},
						},
					},
				},
				config: &Config{
					NgRelationMap:          make(map[string]map[string]struct{}),
					GroupingDirectoryPaths: []string{"sample/grouping"},
					ExcludePackageMap:      map[string]struct{}{"mod/sample/grouping/dst2": {}},
				},
			},
			want: map[string]map[string]int{
				"mod/sample/src": {
					"mod/sample/grouping": 2,
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := &Prelviz{
				projectModuleName: tt.fields.projectModuleName,
				packageInfoMap:    tt.fields.packageInfoMap,
				config:            tt.fields.config,
			}
			got := m.nodeRelationCountMap()
			for srcPkg, dstPkgInfos := range got {
				for dstPkg, num := range dstPkgInfos {
					if tt.want[srcPkg][dstPkg] != num {
						t.Errorf("got(map[%s][%s]%d), want num = %d", srcPkg, dstPkg, num, tt.want[srcPkg][dstPkg])
					}
				}
			}
		})
	}
}

func TestPrelviz_importPathNodeName(t *testing.T) {
	type fields struct {
		projectModuleName string
		config            *Config
	}
	type args struct {
		importPath string
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   string
	}{
		{
			name: "normal: input value is not grouped",
			fields: fields{
				projectModuleName: "mod",
				config: &Config{
					GroupingDirectoryPaths: []string{"hoge"},
				},
			},
			args: args{
				importPath: "mod/sample/hoge",
			},
			want: "mod/sample/hoge",
		},
		{
			name: "normal: input value is grouped",
			fields: fields{
				projectModuleName: "mod",
				config: &Config{
					GroupingDirectoryPaths: []string{"sample"},
				},
			},
			args: args{
				importPath: "mod/sample/hoge",
			},
			want: "mod/sample",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := &Prelviz{
				projectModuleName: tt.fields.projectModuleName,
				config:            tt.fields.config,
			}
			if got := m.importPathNodeName(tt.args.importPath); got != tt.want {
				t.Errorf("Prelviz.importPathNodeName() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestPrelviz_nodeName(t *testing.T) {
	type fields struct {
		projectModuleName string
		config            *Config
	}
	type args struct {
		pkgDirPath string
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   string
	}{
		{
			name: "normal: config.GroupingDirectoryPaths is empty slice",
			fields: fields{
				projectModuleName: "mod",
				config: &Config{
					GroupingDirectoryPaths: []string{},
				},
			},
			args: args{
				pkgDirPath: "sample/hoge",
			},
			want: "mod/sample/hoge",
		},
		{
			name: "normal: input value is grouped",
			fields: fields{
				projectModuleName: "mod",
				config: &Config{
					GroupingDirectoryPaths: []string{
						"sample",
					},
				},
			},
			args: args{
				pkgDirPath: "sample/hoge",
			},
			want: "mod/sample",
		},
		{
			name: "normal: input value is not grouped",
			fields: fields{
				projectModuleName: "mod",
				config: &Config{
					GroupingDirectoryPaths: []string{
						"fugafua",
					},
				},
			},
			args: args{
				pkgDirPath: "sample/hoge",
			},
			want: "mod/sample/hoge",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := &Prelviz{
				projectModuleName: tt.fields.projectModuleName,
				config:            tt.fields.config,
			}
			if got := m.nodeName(tt.args.pkgDirPath); got != tt.want {
				t.Errorf("Prelviz.nodeName() = %v, want %v", got, tt.want)
			}
		})
	}
}
