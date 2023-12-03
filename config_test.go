package prelviz

import (
	"reflect"
	"testing"
)

func TestNewConfig(t *testing.T) {
	type args struct {
		path string
	}
	tests := []struct {
		name    string
		args    args
		want    *Config
		wantErr bool
	}{
		{
			name: "anomaly: duplicate in grouping_directory_path",
			args: args{
				path: "testdata/config_test/duplicate",
			},
			want:    nil,
			wantErr: true,
		},
		{
			name: "anomaly: parent-children relation path exists in grouping_directory_path",
			args: args{
				path: "testdata/config_test/paretchildrelation",
			},
			want:    nil,
			wantErr: true,
		},
		{
			name: "normal: .prelviz.config.json do not exists",
			args: args{
				path: "not/exists",
			},
			want:    nil,
			wantErr: false,
		},
		{
			name: "normal: .prelviz.config.json exists",
			args: args{
				path: "testdata/config_test/valid",
			},
			want: &Config{
				NgRelations: []NgRelation{
					{
						From: "sample1",
						To:   []string{"sample2", "sample3"},
					},
					{
						From: "sample2",
						To:   []string{"sample3"},
					},
				},
				GroupingDirectoryPaths: []string{"sample4", "sample5", "sample6"},
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := NewConfig(tt.args.path)
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
			name: "normal: GroupingDirectoryPaths is nil",
			fields: fields{
				GroupingDirectoryPaths: nil,
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
			name: "normal: multiple empty string in GroupingDirectoryPaths",
			fields: fields{
				GroupingDirectoryPaths: []string{"", "", ""},
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
			name: "normal: GroupingDirectoryPaths is nil",
			fields: fields{
				GroupingDirectoryPaths: nil,
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
			name: "normal: multiple empty string in GroupingDirectoryPaths",
			fields: fields{
				GroupingDirectoryPaths: []string{"", "", ""},
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

func TestConfig_ToToNgPackageRelationMap(t *testing.T) {
	type fields struct {
		NgRelations []NgRelation
	}
	tests := []struct {
		name   string
		fields fields
		want   NgPackageRelationMap
	}{
		{
			name: "normal: NgRelations is nil",
			fields: fields{
				NgRelations: nil,
			},
			want: nil,
		},
		{
			name: "normal: NgRelations is empty",
			fields: fields{
				NgRelations: []NgRelation{},
			},
			want: NgPackageRelationMap{},
		},
		{
			name: "normal",
			fields: fields{
				NgRelations: []NgRelation{
					{
						From: "sample1",
						To:   []string{"sample2", "sample3"},
					},
					{
						From: "sample2",
						To:   []string{"sample3"},
					},
				},
			},
			want: NgPackageRelationMap{
				"sample1": {"sample2": {}, "sample3": {}},
				"sample2": {"sample3": {}},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &Config{
				NgRelations: tt.fields.NgRelations,
			}
			if got := c.ToNgPackageRelationMap(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Config.ToNgPackageRelationMap() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestNgPackageRelationMap_IsNgRelation(t *testing.T) {
	type field struct {
		ngPackageRelationMap NgPackageRelationMap
	}
	type args struct {
		from string
		to   string
	}
	tests := []struct {
		name  string
		field field
		args  args
		want  bool
	}{
		{
			name: "normal: from's value do not exists",
			field: field{
				ngPackageRelationMap: NgPackageRelationMap{
					"sample1": {"sample2": {}, "sample3": {}},
				},
			},
			args: args{
				from: "sample4",
				to:   "sample2",
			},
			want: false,
		},
		{
			name: "normal: to's value do not exists",
			field: field{
				ngPackageRelationMap: NgPackageRelationMap{
					"sample1": {"sample2": {}, "sample3": {}},
				},
			},
			args: args{
				from: "sample1",
				to:   "sample4",
			},
			want: false,
		},
		{
			name: "normal: from's value and to's value exists",
			field: field{
				ngPackageRelationMap: NgPackageRelationMap{
					"sample1": {"sample2": {}, "sample3": {}},
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
			n := NgPackageRelationMap(tt.field.ngPackageRelationMap)
			if got := n.IsNgRelation(tt.args.from, tt.args.to); got != tt.want {
				t.Errorf("NgPackageRelationMap.IsNgRelation() = %v, want %v", got, tt.want)
			}
		})
	}
}
