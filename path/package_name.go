package path

import (
	"fmt"
	"os"
	"path/filepath"
)

func GetFullPackageName(path string) (string, error) {
	gopath := os.Getenv("GOPATH")

	basePath := filepath.Dir(path)

	abs, err := filepath.Abs(basePath)
	if err != nil {
		return "", err
	}

	rel, err := filepath.Rel(gopath+"/src", abs)
	if err != nil {
		return "", err
	}

	if rel[0:2] == ".." {
		return "", fmt.Errorf("'%v' is outside GOPATH", path)
	}

	return rel, nil
}
