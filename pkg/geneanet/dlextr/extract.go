package dlextr

import (
	"archive/zip"
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"
)

const (
	errCreateOutputDir = "creating output directory: %w"
	errNewReader       = "zip NewReader: %w"
	errExtractZip      = "extracting zip archive: %w"
	errCreateDir       = "creating directory: %w"
	errOpenZipFile     = "opening zipped file: %w"
	errOpenDestFile    = "opening dest file: %w"
	errCopyZipData     = "copying data: %w"
)

func extract(file *zip.File, fileDst string) error {
	if file.FileInfo().IsDir() {
		path := filepath.Join(fileDst, file.Name) //nolint:gosec
		if err := os.MkdirAll(path, file.FileInfo().Mode()); err != nil {
			return fmt.Errorf(errCreateDir, err)
		}
	}

	zipFile, err := file.Open()
	if err != nil {
		return fmt.Errorf(errOpenZipFile, err)
	}

	defer zipFile.Close()

	f, err := os.Create(fileDst)
	if err != nil {
		return fmt.Errorf(errOpenDestFile, err)
	}

	defer f.Close()

	if _, err := io.Copy(f, zipFile); err != nil { //nolint:gosec
		return fmt.Errorf(errCopyZipData, err)
	}

	return nil
}

func (d *Download) Unzip() error {
	if err := os.MkdirAll(d.outputDir, os.ModePerm|os.ModeDir); err != nil {
		return fmt.Errorf(errCreateOutputDir, err)
	}

	zr, err := zip.NewReader(d.reader, d.size)
	if err != nil {
		return fmt.Errorf(errNewReader, err)
	}

	for _, file := range zr.File {
		log.Printf("Processing file: %s", file.Name)

		if err := extract(file, filepath.Join(d.outputDir, file.Name)); err != nil { //nolint:gosec
			return fmt.Errorf(errExtractZip, err)
		}
	}

	return nil
}
