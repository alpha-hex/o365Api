package o365Api

import (
	"io"
	"fmt"
	"strings"
	"os"
	"path/filepath"
	"archive/zip"
	
)

type ZipRequest interface {
	Unzip(Zip) ([]string, error)
}

type Zip struct {
	source		string
	destination	string
}

func Unzip(request Zip) ([]string, error) {
	var filenames []string

	reader, err := zip.OpenReader(request.source)
	if err != nil {
		return []string{}, err
	}
	defer reader.Close()

	for _, file := range reader.File {
		path := filepath.Join(request.destination, file.Name)
		if !strings.HasPrefix(path, filepath.Clean(request.destination)+string(os.PathSeparator)) {
			return filenames, fmt.Errorf("%s: illegal file path", path)
		}

		filenames = append(filenames, path)

		if file.FileInfo().IsDir() {
			os.MkdirAll(path, os.ModePerm)
			continue
		}

		if err := os.MkdirAll(filepath.Dir(path), os.ModePerm); err != nil {
			return filenames, err
		}

		outFile, err := os.OpenFile(path, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, file.Mode())
		if err != nil {
			return filenames, err
		}

		rc, err := file.Open()
		if err != nil {
			return filenames, err
		}

		_, err = io.Copy(outFile, rc)
		outFile.Close()
		rc.Close()
		if err != nil {
			return filenames, err
		}
	}

	return filenames, nil
} 