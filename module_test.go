package prelviz

import "testing"

func TestGetModuleName(t *testing.T) {
	type args struct {
		rootPath string
	}
	tests := []struct {
		name    string
		args    args
		want    string
		wantErr bool
	}{
		{
			name: "anomaly: go.mod do not exists",
			args: args{
				rootPath: "testdata",
			},
			want:    "",
			wantErr: true,
		},
		{
			name: "anomaly: go.mod is invalid format",
			args: args{
				rootPath: "testdata/module_test/invalid",
			},
			want:    "",
			wantErr: true,
		},
		{
			name: "normal",
			args: args{
				rootPath: "testdata/module_test/valid",
			},
			want:    "sample",
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := GetModuleName(tt.args.rootPath)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetModuleName() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("GetModuleName() = %v, want %v", got, tt.want)
			}
		})
	}
}
