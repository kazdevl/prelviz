package prelviz

import (
	"reflect"
	"sort"
	"testing"
)

func Test_NewPackageInfoMap(t *testing.T) {
	type args struct {
		projectDirectoryPath string
	}
	tests := []struct {
		name    string
		args    args
		want    map[string]*PackageInfo
		wantErr bool
	}{
		{
			name: "normal",
			args: args{
				projectDirectoryPath: "testdata/package_test/valid",
			},
			want: map[string]*PackageInfo{
				"": {
					Name:          "valid",
					DirectoryPath: "",
					ImportUsageMap: map[string]map[string]struct{}{
						"fmt":  {"Println": {}, "Printf": {}},
						"time": {"Now": {}},
						"github.com/kazdevl/prelviz/testdata/package_test/valid/nest/sample": {
							"NewSample":    {},
							"SampleString": {},
							"SampleError":  {},
						},
					},
				},
				"nest/sample": {
					Name:          "sample",
					DirectoryPath: "nest/sample",
					ImportUsageMap: map[string]map[string]struct{}{
						"fmt":  {"Sprintf": {}, "Errorf": {}},
						"time": {"Time": {}, "DateOnly": {}},
					},
				},
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := NewPackageInfoMap(tt.args.projectDirectoryPath)
			if (err != nil) != tt.wantErr {
				t.Errorf("NewPackageInfoMap() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(len(got), len(tt.want)) {
				t.Errorf("NewPackageInfoMap() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_NewPackageInfo(t *testing.T) {
	type args struct {
		filePath             string
		projectDirectoryPath string
	}
	tests := []struct {
		name    string
		args    args
		want    *PackageInfo
		wantErr bool
	}{
		{
			name: "normal",
			args: args{
				filePath:             "testdata/package_test/valid/nest/sample/sample.go",
				projectDirectoryPath: "testdata/package_test/valid",
			},
			want: &PackageInfo{
				Name:          "sample",
				DirectoryPath: "nest/sample",
				ImportUsageMap: map[string]map[string]struct{}{
					"time": {"Time": {}, "DateOnly": {}},
					"fmt":  {"Sprintf": {}},
				},
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := NewPackageInfo(tt.args.filePath, tt.args.projectDirectoryPath)
			if (err != nil) != tt.wantErr {
				t.Errorf("NewPackageInfo() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewPackageInfo() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_targetGoFilePaths(t *testing.T) {
	type args struct {
		dir string
	}
	tests := []struct {
		name    string
		args    args
		want    []string
		wantErr bool
	}{
		{
			name: "normal: input dir do not exists",
			args: args{
				dir: "testdata/package_test/notexists",
			},
			want:    []string{},
			wantErr: false,
		},
		{
			name: "normal: go file do not exists",
			args: args{
				dir: "testdata/package_test/nogofile",
			},
			want:    []string{},
			wantErr: false,
		},
		{
			name: "normal",
			args: args{
				dir: "testdata/package_test/valid",
			},
			want: []string{
				"testdata/package_test/valid/sample.go",
				"testdata/package_test/valid/sample1.go",
				"testdata/package_test/valid/nest/sample/sample.go",
				"testdata/package_test/valid/nest/sample/sample1.go",
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := targetGoFilePaths(tt.args.dir)
			if (err != nil) != tt.wantErr {
				t.Errorf("targetGoFilePaths() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			sort.StringSlice(got).Sort()
			sort.StringSlice(tt.want).Sort()
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("targetGoFilePaths() = %v, want %v", got, tt.want)
			}
		})
	}
}
