package utils

import (
	"archive/tar"
	"compress/gzip"
	"errors"
	"io"
	"log"
	"os"
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
				//log.Fatalf("ExtractTarGz: Mkdir() failed: %s", err.Error())
				return errors.New("ExtractTarGz: NewReader failed")

			}
		case tar.TypeReg:
			outFile, err := os.Create(header.Name)
			if err != nil {
				//log.Fatalf("ExtractTarGz: Create() failed: %s", err.Error())
				return errors.New("ExtractTarGz: NewReader failed")

			}
			defer outFile.Close()
			if _, err := io.Copy(outFile, tarReader); err != nil {
				//log.Fatalf("ExtractTarGz: Copy() failed: %s", err.Error())
				return errors.New("ExtractTarGz: NewReader failed")

			}
		default:
			return errors.New("ExtractTarGz: NewReader failed")
		}
	}
	return nil
}
