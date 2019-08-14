package o365Api

import (
	"archive/zip"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
)

type ZipRequest interface {
	Unzip(Zip) ([]string, error)
}

type Zip struct {
	Source      string
	Destination string
}

func (zipRequest Zip) Unzip() ([]string, error) {
	var filenames []string

	reader, err := zip.OpenReader(zipRequest.Source)
	if err != nil {
		return []string{}, err
	}
	defer reader.Close()

	for _, file := range reader.File {
		path := filepath.Join(zipRequest.Destination, file.Name)
		if !strings.HasPrefix(path, filepath.Clean(zipRequest.Destination)+string(os.PathSeparator)) {
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
