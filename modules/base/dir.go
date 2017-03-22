package base

import (
	"fmt"
	"os"
	"path"
	"sort"

	"github.com/Unknwon/com"
)

func GetFiles(dirPath string) ([]os.FileInfo, error) {
	if !com.IsDir(dirPath) {
		return nil, fmt.Errorf("path is not dir: %s", dirPath)
	}

	dir, err := os.Open(dirPath)
	if err != nil {
		return nil, err
	}

	files, err := dir.Readdir(0)
	if err != nil {
		return nil, err
	}
	return files, nil
}

func GetFileListSortByMTime(dirPath string) ([]string, error) {
	files, err := GetFiles(dirPath)
	if err != nil {
		return nil, err
	}

	sort.Slice(files, func(i, j int) bool {
		return files[i].ModTime().Before(files[j].ModTime())
	})

	fl := make([]string, 0, len(files))
	for _, f := range files {
		fl = append(fl, path.Join(dirPath, f.Name()))
	}
	return fl, nil
}
