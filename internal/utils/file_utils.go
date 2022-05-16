package utils

import (
	"archive/tar"
	"compress/gzip"
	"errors"
	"io"
	"log"
	"os"
	"path/filepath"
)

func ExtractTarGz(gzipStream io.Reader) error {
	log.Print("[DEBUG] untaring file")

	uncompressedStream, err := gzip.NewReader(gzipStream)
	if err != nil {
		return errors.New("ExtractTarGz: NewReader failed")
	}

	tarReader := tar.NewReader(uncompressedStream)

	for true {
		header, err := tarReader.Next()

		if err == io.EOF {
			break
		}

		if err != nil {
			log.Fatalf("ExtractTarGz: Next() failed: %s", err.Error())
		}

		switch header.Typeflag {
		case tar.TypeDir:
			if err := os.Mkdir(header.Name, 0755); err != nil {
				return errors.New("ExtractTarGz: NewReader failed")

			}
		case tar.TypeReg:
			outFile, err := os.Create(header.Name)
			if err != nil {
				// We were unable to create the file because there was no
				// Tar Directory entry coresponding to this file's directory.
				// Let's try to create the directory and try again.
				if err = os.MkdirAll(filepath.Dir(header.Name), 0755); err != nil {
					return errors.New("ExtractTarGz: NewReader failed")
				}
				outFile, err = os.Create(header.Name)
				if err != nil {
					return errors.New("ExtractTarGz: NewReader failed")
				}

			}
			defer outFile.Close()
			if _, err := io.Copy(outFile, tarReader); err != nil {
				return errors.New("ExtractTarGz: NewReader failed")

			}
		default:
			return errors.New("ExtractTarGz: NewReader failed")
		}
	}
	return nil
}
