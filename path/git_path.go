package path

import (
	"log"
	"os"
	"path/filepath"
)

func GitPath(path string) (string, error) {
	var err error
	var absPath string

	absPath, err = filepath.Abs(path)
	if err != nil {
		return "", err
	}

	recursiveDirPath := filepath.Dir(path)
	gitDirPath := recursiveDirPath
	for {
		if recursiveDirPath == "/" {
			break
		}
		if IsGitDirExists(recursiveDirPath) {
			gitDirPath = recursiveDirPath
			break
		}
		recursiveDirPath = filepath.Clean(recursiveDirPath + "/../")
	}

	gitDirPath, err = filepath.Abs(gitDirPath)
	if err != nil {
		return "", err
	}

	return filepath.Rel(gitDirPath, absPath)
}

func IsGitDirExists(dir string) bool {
	gitDir := filepath.Clean(dir + "/.git")

	info, err := os.Stat(gitDir)

	if os.IsNotExist(err) {
		return false
	}

	if os.IsPermission(err) {
		log.Printf("Permission denied when access %s", dir)
		return false
	}

	if err != nil {
		log.Printf("Got error when access %s: %v", dir, err)
		return false
	}

	if !info.IsDir() {
		return false
	}

	return true
}
