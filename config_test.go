package prelviz

import (
	"reflect"
	"testing"
)

func TestNewConfig(t *testing.T) {
	type args struct {
		path       string
		moduleName string
	}
	tests := []struct {
		name    string
		args    args
		want    *Config
		wantErr bool
	}{
		{
			name: "normal: .prelviz.config.json do not exists",
			args: args{
				path: "not/exists",
			},
			want: &Config{
				NgRelationMap:          make(map[string]map[string]struct{}),
				GroupingDirectoryPaths: make([]string, 0),
				ExcludePackageMap:      make(map[string]struct{}),
			},
			wantErr: false,
		},
		{
			name: "normal: .prelviz.config.json exists",
			args: args{
				path:       "testdata/config_test/valid",
				moduleName: "mod",
			},
			want: &Config{
				NgRelationMap: map[string]map[string]struct{}{
					"sample1": {"sample2": {}, "sample3": {}},
					"sample2": {"sample3": {}},
				},
				GroupingDirectoryPaths: []string{"sample4", "sample5", "sample6"},
				ExcludePackageMap: map[string]struct{}{
					"mod/sample7": {},
					"mod/testdata/config_test/valid/exclude/nest1_1":         {},
					"mod/testdata/config_test/valid/exclude/nest1_1/nest2_1": {},
				},
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := NewConfig(tt.args.path, tt.args.moduleName)
			if (err != nil) != tt.wantErr {
				t.Errorf("NewConfig() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewConfig() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestConfigBinder_ToConfig(t *testing.T) {
	type fields struct {
		NgRelations            []NgRelation
		GroupingDirectoryPaths []string
		ExcludePackages        []string
		ExcludeDirectorys      []string
	}
	type args struct {
		moduleName string
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    *Config
		wantErr bool
	}{
		{
			name: "anomaly: parent-children relation path exists in grouping_directory_path",
			fields: fields{
				GroupingDirectoryPaths: []string{"sample1", "sample1/sample2"},
			},
			args:    args{moduleName: "mod"},
			want:    nil,
			wantErr: true,
		},
		{
			name: "normal: duplicate in grouping_directory_path",
			fields: fields{
				GroupingDirectoryPaths: []string{"sample1", "sample1", "sample2"},
			},
			args: args{moduleName: "mod"},
			want: &Config{
				NgRelationMap:          make(map[string]map[string]struct{}),
				GroupingDirectoryPaths: []string{"sample1", "sample2"},
				ExcludePackageMap:      make(map[string]struct{}),
			},
			wantErr: false,
		},
		{
			name: "normal: ng_relation is nil",
			fields: fields{
				NgRelations:            nil,
				GroupingDirectoryPaths: []string{"sample1", "sample2"},
				ExcludePackages:        []string{"mod/sample3", "mod/sample4"},
			},
			args: args{moduleName: "mod"},
			want: &Config{
				NgRelationMap:          make(map[string]map[string]struct{}),
				GroupingDirectoryPaths: []string{"sample1", "sample2"},
				ExcludePackageMap:      map[string]struct{}{"mod/sample3": {}, "mod/sample4": {}},
			},
			wantErr: false,
		},
		{
			name: "normal: ng_relation is not set",
			fields: fields{
				NgRelations:            []NgRelation{},
				GroupingDirectoryPaths: []string{"sample1", "sample2"},
				ExcludePackages:        []string{"mod/sample3", "mod/sample4"},
			},
			args: args{moduleName: "mod"},
			want: &Config{
				NgRelationMap:          make(map[string]map[string]struct{}),
				GroupingDirectoryPaths: []string{"sample1", "sample2"},
				ExcludePackageMap:      map[string]struct{}{"mod/sample3": {}, "mod/sample4": {}},
			},
			wantErr: false,
		},
		{
			name: "normal: grouping_directory_path is nil",
			fields: fields{
				NgRelations: []NgRelation{
					{
						From: "mod/sample1",
						To:   []string{"mod/sample2", "mod/sample3"},
					},
				},
				GroupingDirectoryPaths: nil,
				ExcludePackages:        []string{"mod/sample3", "mod/sample4"},
			},
			args: args{moduleName: "mod"},
			want: &Config{
				NgRelationMap: map[string]map[string]struct{}{
					"mod/sample1": {"mod/sample2": {}, "mod/sample3": {}},
				},
				GroupingDirectoryPaths: make([]string, 0),
				ExcludePackageMap:      map[string]struct{}{"mod/sample3": {}, "mod/sample4": {}},
			},
			wantErr: false,
		},
		{
			name: "normal: grouping_directory_path is not set",
			fields: fields{
				NgRelations: []NgRelation{
					{
						From: "mod/sample1",
						To:   []string{"mod/sample2", "mod/sample3"},
					},
				},
				GroupingDirectoryPaths: []string{},
				ExcludePackages:        []string{"mod/sample3", "mod/sample4"},
			},
			args: args{moduleName: "mod"},
			want: &Config{
				NgRelationMap: map[string]map[string]struct{}{
					"mod/sample1": {"mod/sample2": {}, "mod/sample3": {}},
				},
				GroupingDirectoryPaths: make([]string, 0),
				ExcludePackageMap:      map[string]struct{}{"mod/sample3": {}, "mod/sample4": {}},
			},
			wantErr: false,
		},
		{
			name: "normal: grouping_directory_path is empty str slice",
			fields: fields{
				NgRelations: []NgRelation{
					{
						From: "mod/sample1",
						To:   []string{"mod/sample2", "mod/sample3"},
					},
				},
				GroupingDirectoryPaths: []string{"", ""},
				ExcludePackages:        []string{"mod/sample3", "mod/sample4"},
			},
			args: args{moduleName: "mod"},
			want: &Config{
				NgRelationMap: map[string]map[string]struct{}{
					"mod/sample1": {"mod/sample2": {}, "mod/sample3": {}},
				},
				GroupingDirectoryPaths: []string{""},
				ExcludePackageMap:      map[string]struct{}{"mod/sample3": {}, "mod/sample4": {}},
			},
			wantErr: false,
		},
		{
			name: "normal: exclude_package is nil",
			fields: fields{
				NgRelations: []NgRelation{
					{
						From: "mod/sample1",
						To:   []string{"mod/sample2", "mod/sample3"},
					},
				},
				GroupingDirectoryPaths: []string{"sample4", "sample5"},
				ExcludePackages:        nil,
			},
			args: args{moduleName: "mod"},
			want: &Config{
				NgRelationMap: map[string]map[string]struct{}{
					"mod/sample1": {"mod/sample2": {}, "mod/sample3": {}},
				},
				GroupingDirectoryPaths: []string{"sample4", "sample5"},
				ExcludePackageMap:      make(map[string]struct{}),
			},
			wantErr: false,
		},
		{
			name: "normal: exclude_package is not set",
			fields: fields{
				NgRelations: []NgRelation{
					{
						From: "mod/sample1",
						To:   []string{"mod/sample2", "mod/sample3"},
					},
				},
				GroupingDirectoryPaths: []string{"sample4", "sample5"},
				ExcludePackages:        []string{},
			},
			args: args{moduleName: "mod"},
			want: &Config{
				NgRelationMap: map[string]map[string]struct{}{
					"mod/sample1": {"mod/sample2": {}, "mod/sample3": {}},
				},
				GroupingDirectoryPaths: []string{"sample4", "sample5"},
				ExcludePackageMap:      make(map[string]struct{}),
			},
			wantErr: false,
		},
		{
			name: "normal: exclude_package is empty str slice",
			fields: fields{
				NgRelations: []NgRelation{
					{
						From: "mod/sample1",
						To:   []string{"mod/sample2", "mod/sample3"},
					},
				},
				GroupingDirectoryPaths: []string{"sample4", "sample5"},
				ExcludePackages:        []string{"", ""},
			},
			args: args{moduleName: "mod"},
			want: &Config{
				NgRelationMap: map[string]map[string]struct{}{
					"mod/sample1": {"mod/sample2": {}, "mod/sample3": {}},
				},
				GroupingDirectoryPaths: []string{"sample4", "sample5"},
				ExcludePackageMap:      make(map[string]struct{}),
			},
			wantErr: false,
		},
		{
			name: "normal: exclude_directory is not set",
			fields: fields{
				NgRelations: []NgRelation{
					{
						From: "mod/sample1",
						To:   []string{"mod/sample2", "mod/sample3"},
					},
				},
				GroupingDirectoryPaths: []string{"sample4", "sample5"},
				ExcludePackages:        []string{"mod/sample6", "mod/sample7"},
				ExcludeDirectorys:      nil,
			},
			args: args{moduleName: "mod"},
			want: &Config{
				NgRelationMap: map[string]map[string]struct{}{
					"mod/sample1": {"mod/sample2": {}, "mod/sample3": {}},
				},
				GroupingDirectoryPaths: []string{"sample4", "sample5"},
				ExcludePackageMap:      map[string]struct{}{"mod/sample6": {}, "mod/sample7": {}},
			},
			wantErr: false,
		},
		{
			name: "normal: exclude_directory is empty str slice",
			fields: fields{
				NgRelations: []NgRelation{
					{
						From: "mod/sample1",
						To:   []string{"mod/sample2", "mod/sample3"},
					},
				},
				GroupingDirectoryPaths: []string{"sample4", "sample5"},
				ExcludePackages:        []string{"mod/sample6", "mod/sample7"},
				ExcludeDirectorys:      []string{},
			},
			args: args{moduleName: "mod"},
			want: &Config{
				NgRelationMap: map[string]map[string]struct{}{
					"mod/sample1": {"mod/sample2": {}, "mod/sample3": {}},
				},
				GroupingDirectoryPaths: []string{"sample4", "sample5"},
				ExcludePackageMap:      map[string]struct{}{"mod/sample6": {}, "mod/sample7": {}},
			},
			wantErr: false,
		},
		{
			name: "normal: exclude_directory is set",
			fields: fields{
				NgRelations: []NgRelation{
					{
						From: "mod/sample1",
						To:   []string{"mod/sample2", "mod/sample3"},
					},
				},
				GroupingDirectoryPaths: []string{"sample4", "sample5"},
				ExcludePackages:        []string{"mod/sample6", "mod/sample7"},
				ExcludeDirectorys:      []string{"testdata/config_test/valid/exclude/nest1_1"},
			},
			args: args{moduleName: "mod"},
			want: &Config{
				NgRelationMap: map[string]map[string]struct{}{
					"mod/sample1": {"mod/sample2": {}, "mod/sample3": {}},
				},
				GroupingDirectoryPaths: []string{"sample4", "sample5"},
				ExcludePackageMap: map[string]struct{}{
					"mod/sample6": {},
					"mod/sample7": {},
					"mod/testdata/config_test/valid/exclude/nest1_1":         {},
					"mod/testdata/config_test/valid/exclude/nest1_1/nest2_1": {},
				},
			},
			wantErr: false,
		},
		{
			name: "normal: exclude_directory is set and duplicate in exclude_package",
			fields: fields{
				NgRelations: []NgRelation{
					{
						From: "mod/sample1",
						To:   []string{"mod/sample2", "mod/sample3"},
					},
				},
				GroupingDirectoryPaths: []string{"sample4", "sample5"},
				ExcludePackages:        []string{"mod/testdata/config_test/valid/exclude/nest1_1/nest2_1"},
				ExcludeDirectorys: []string{
					"testdata/config_test/valid/exclude/nest1_1",
					"testdata/config_test/valid/exclude/nest1_2",
				},
			},
			args: args{moduleName: "mod"},
			want: &Config{
				NgRelationMap: map[string]map[string]struct{}{
					"mod/sample1": {"mod/sample2": {}, "mod/sample3": {}},
				},
				GroupingDirectoryPaths: []string{"sample4", "sample5"},
				ExcludePackageMap: map[string]struct{}{
					"mod/testdata/config_test/valid/exclude/nest1_1":         {},
					"mod/testdata/config_test/valid/exclude/nest1_1/nest2_1": {},
					"mod/testdata/config_test/valid/exclude/nest1_2":         {},
					"mod/testdata/config_test/valid/exclude/nest1_2/nest2_1": {},
					"mod/testdata/config_test/valid/exclude/nest1_2/nest2_2": {},
				},
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := ConfigBinder{
				NgRelations:            tt.fields.NgRelations,
				GroupingDirectoryPaths: tt.fields.GroupingDirectoryPaths,
				ExcludePackages:        tt.fields.ExcludePackages,
				ExcludeDirectorys:      tt.fields.ExcludeDirectorys,
			}
			got, err := c.ToConfig(tt.args.moduleName)
			if (err != nil) != tt.wantErr {
				t.Errorf("ToConfig() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("ToConfig() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestConfig_IsGroupingPackage(t *testing.T) {
	type fields struct {
		GroupingDirectoryPaths []string
	}
	type args struct {
		pkgDirPath string
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   bool
	}{
		{
			name: "normal: GroupingDirectoryPaths is empty slice",
			fields: fields{
				GroupingDirectoryPaths: []string{},
			},
			args: args{
				pkgDirPath: "app/sample",
			},
			want: false,
		},
		{
			name: "normal: empty string in GroupingDirectoryPaths",
			fields: fields{
				GroupingDirectoryPaths: []string{""},
			},
			args: args{
				pkgDirPath: "app/sample",
			},
			want: false,
		},
		{
			name: "normal: input value exists in GroupingDirectoryPaths",
			fields: fields{
				GroupingDirectoryPaths: []string{"app/sample", "app/hoge"},
			},
			args: args{
				pkgDirPath: "app/sample/fugafuga",
			},
			want: true,
		},
		{
			name: "normal: input value do not exists in GroupingDirectoryPaths",
			fields: fields{
				GroupingDirectoryPaths: []string{"app/sample", "app/hoge"},
			},
			args: args{
				pkgDirPath: "app/fuga/fugafuga",
			},
			want: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &Config{
				GroupingDirectoryPaths: tt.fields.GroupingDirectoryPaths,
			}
			if got := c.IsGroupingPackage(tt.args.pkgDirPath); got != tt.want {
				t.Errorf("Config.IsGroupingPackage() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestConfig_GroupingPackageDirectoryPath(t *testing.T) {
	type fields struct {
		GroupingDirectoryPaths []string
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
			name: "normal: GroupingDirectoryPaths is empty slice",
			fields: fields{
				GroupingDirectoryPaths: []string{},
			},
			args: args{
				pkgDirPath: "app/sample",
			},
			want: "app/sample",
		},
		{
			name: "normal: empty string in GroupingDirectoryPaths",
			fields: fields{
				GroupingDirectoryPaths: []string{""},
			},
			args: args{
				pkgDirPath: "app/sample",
			},
			want: "app/sample",
		},
		{
			name: "normal: input value exists in GroupingDirectoryPaths",
			fields: fields{
				GroupingDirectoryPaths: []string{"app/sample", "app/hoge"},
			},
			args: args{
				pkgDirPath: "app/sample/fugafuga",
			},
			want: "app/sample",
		},
		{
			name: "normal: input value do not exists in GroupingDirectoryPaths",
			fields: fields{
				GroupingDirectoryPaths: []string{"app/sample", "app/hoge"},
			},
			args: args{
				pkgDirPath: "app/fuga/fugafuga",
			},
			want: "app/fuga/fugafuga",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &Config{
				GroupingDirectoryPaths: tt.fields.GroupingDirectoryPaths,
			}
			if got := c.GroupingPackageDirectoryPath(tt.args.pkgDirPath); got != tt.want {
				t.Errorf("Config.GroupingPackageDirectoryPath() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestConfig_IsExcludePackage(t *testing.T) {
	type fields struct {
		ExcludePackageMap map[string]struct{}
	}
	type args struct {
		pkg string
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   bool
	}{
		{
			name: "normal: ExcludePackageMap is empty map",
			fields: fields{
				ExcludePackageMap: make(map[string]struct{}),
			},
			args: args{
				pkg: "github.com/sample_project/app/sample",
			},
			want: false,
		},
		{
			name: "normal: input value do not exists in ExcludePackageMap",
			fields: fields{
				ExcludePackageMap: map[string]struct{}{
					"github.com/sample_project/app/sample/hoge": {},
				},
			},
			args: args{
				pkg: "github.com/sample_project/app/sample",
			},
			want: false,
		},
		{
			name: "normal: input value exists in ExcludePackageMap",
			fields: fields{
				ExcludePackageMap: map[string]struct{}{
					"github.com/sample_project/app/sample": {},
				},
			},
			args: args{
				pkg: "github.com/sample_project/app/sample",
			},
			want: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &Config{
				ExcludePackageMap: tt.fields.ExcludePackageMap,
			}
			if got := c.IsExcludePackage(tt.args.pkg); got != tt.want {
				t.Errorf("IsExcludePackage() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestConfig_IsNgRelation(t *testing.T) {
	type fields struct {
		NgRelationMap map[string]map[string]struct{}
	}
	type args struct {
		from string
		to   string
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   bool
	}{
		{
			name: "normal: NgRelationMap is empty map",
			fields: fields{
				NgRelationMap: make(map[string]map[string]struct{}),
			},
			args: args{
				from: "sample1",
				to:   "sample2",
			},
			want: false,
		},
		{
			name: "normal: from do not exists",
			fields: fields{
				NgRelationMap: map[string]map[string]struct{}{
					"sample_hoge": {"sample_hoge": {}},
				},
			},
			args: args{
				from: "sample1",
				to:   "sample2",
			},
			want: false,
		},
		{
			name: "normal: to do not exists",
			fields: fields{
				NgRelationMap: map[string]map[string]struct{}{
					"sample1": {"sample_hoge": {}},
				},
			},
			args: args{
				from: "sample1",
				to:   "sample2",
			},
			want: false,
		},
		{
			name: "normal: from exists, to exists",
			fields: fields{
				NgRelationMap: map[string]map[string]struct{}{
					"sample1": {"sample2": {}},
				},
			},
			args: args{
				from: "sample1",
				to:   "sample2",
			},
			want: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &Config{
				NgRelationMap: tt.fields.NgRelationMap,
			}
			if got := c.IsNgRelation(tt.args.from, tt.args.to); got != tt.want {
				t.Errorf("IsNgRelation() = %v, want %v", got, tt.want)
			}
		})
	}
}
