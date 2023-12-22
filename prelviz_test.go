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
