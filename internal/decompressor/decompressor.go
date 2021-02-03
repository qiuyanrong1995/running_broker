package decompressor

import (
	"archive/tar"
	"archive/zip"
	"compress/gzip"
	"errors"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
)

func DeCompressor(fName string) error {
	fArray := strings.Split(fName, ".")
	suf := fArray[len(fArray) - 1]
	switch suf {
	case "gz":

	case "zip":

	default:
		return errors.New(fmt.Sprintf("not support compress type `%s`", suf))
	}
	return nil
}

func GzDeCompress(fName string, savePath string) error {
	srcFile, err := os.Open(fName)
	if err != nil {
		return err
	}
	defer func() { _ = srcFile.Close() }()
	gr, err := gzip.NewReader(srcFile)
	if err != nil {
		return err
	}
	defer func() { _ = gr.Close() }()
	tr := tar.NewReader(gr)
	for {
		hdr, err := tr.Next()
		if err != nil {
			if err == io.EOF {
				break
			} else {
				return err
			}
		}

		var fileName string
		if strings.HasSuffix(savePath, string(os.PathSeparator)) {
			fileName = savePath + hdr.Name
		} else {
			fileName = savePath + string(os.PathSeparator) + hdr.Name
		}
		f, err := os.Create(fileName)
		if err != nil {
			return err
		}
		_, err = io.Copy(f, tr)
		if err != nil {
			return err
		}
	}
	return nil
}

func ZipDeCompress(fName string, savePath string) error {
	zr, err := zip.OpenReader(fName)
	if err != nil {
		return err
	}
	defer func() { _ = zr.Close() }()

	for _, f := range zr.File {
		fPath := filepath.Join(savePath, f.Name)
		if f.FileInfo().IsDir() {
			err = os.MkdirAll(fPath, os.ModePerm)
			if err != nil {
				return err
			}
		} else {
			if err = os.MkdirAll(filepath.Dir(fPath), os.ModePerm); err != nil {
				return err
			}
			inFile, err := f.Open()
			if err != nil {
				return err
			}

			outFile, err := os.Create(fPath)
			if err != nil {
				return err
			}
			_, err = io.Copy(outFile, inFile)
			if err != nil {
				return  err
			}
		}
	}
	return nil
}