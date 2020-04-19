package path

import "testing"

func TestGetFullPackageName(t *testing.T) {
	tests := []struct {
		name      string
		inputPath string
		want      string
	}{
		{"Simple", "package_name.go", "github.com/egon12/ghcon/path"},
		{"Parent", "../", "github.com/egon12/ghcon"},
		{"Neighbor", "../review/interface.go", "github.com/egon12/ghcon/review"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			got, err := GetFullPackageName(tt.inputPath)

			if err != nil {
				t.Errorf("Unexpected Error: %v", err)
			}

			if got != tt.want {
				t.Errorf("Want '%v' got '%v'", tt.want, got)
			}
		})
	}
}

func TestGetFullPackageName_Error(t *testing.T) {
	tests := []struct {
		name      string
		inputPath string
		want      string
	}{
		{"Outside GoPath", "/", "'/' is outside GOPATH"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := GetFullPackageName(tt.inputPath)
			if err == nil {
				t.Errorf("Expect Error got %v", got)
				return
			}
			got = err.Error()
			if got != tt.want {
				t.Errorf("Want '%v' got '%v'", tt.want, got)
			}
		})
	}
}
