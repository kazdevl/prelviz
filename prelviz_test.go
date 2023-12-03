package prelviz

import (
	"testing"
)

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
			name: "normal: config is nil",
			fields: fields{
				projectModuleName: "mod",
				config:            nil,
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

func TestPrelviz_isNgRelation(t *testing.T) {
	type fields struct {
		config *Config
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
			name: "normal: config is nil",
			fields: fields{
				config: nil,
			},
			args: args{
				from: "from",
				to:   "to",
			},
			want: false,
		},
		{
			name: "normal: ngPackageRelationMap is nil",
			fields: fields{
				config: &Config{
					NgRelations: nil,
				},
			},
			args: args{
				from: "from",
				to:   "to",
			},
			want: false,
		},
		{
			name: "normal: ngPackageRelationMap is zero value",
			fields: fields{
				config: &Config{
					NgRelations: []NgRelation{},
				},
			},
			args: args{
				from: "from",
				to:   "to",
			},
			want: false,
		},
		{
			name: "normal: not ng relation",
			fields: fields{
				config: &Config{
					NgRelations: []NgRelation{
						{
							From: "from_",
							To:   []string{"to_"},
						},
					},
				},
			},
			args: args{
				from: "from",
				to:   "to",
			},
			want: false,
		},
		{
			name: "normal: ng relation",
			fields: fields{
				config: &Config{
					NgRelations: []NgRelation{
						{
							From: "from",
							To:   []string{"to"},
						},
					},
				},
			},
			args: args{
				from: "from",
				to:   "to",
			},
			want: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := &Prelviz{
				config: tt.fields.config,
			}
			if got := m.isNgRelation(tt.args.from, tt.args.to); got != tt.want {
				t.Errorf("Prelviz.isNgRelation() = %v, want %v", got, tt.want)
			}
		})
	}
}
