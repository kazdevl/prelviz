package prelviz

import (
	"bufio"
	"errors"
	"os"
	"path/filepath"
	"strings"
)

func GetModuleName(rootPath string) (string, error) {
	f, err := os.Open(filepath.Join(rootPath, "go.mod"))
	if err != nil {
		return "", err
	}
	defer f.Close()

	scanner := bufio.NewScanner(f)
	scanner.Scan()
	if err = scanner.Err(); err != nil {
		return "", err
	}
	oneline := scanner.Text()
	moduleInfo := strings.Split(oneline, " ")
	if len(moduleInfo) != 2 {
		return "", errors.New("invalid go.mod")
	}
	return moduleInfo[1], nil
}
