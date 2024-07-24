package utils

import (
	"archive/zip"
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"
	"strings"
)

func CreateDirIfNotExits(path_name string) {
	if _, err := os.Stat(path_name); os.IsNotExist(err) {
		err := os.Mkdir(path_name, 0755)
		if err != nil {
			fmt.Println("Error %w", err)
			return
		}
	}
}

func CheckIfNoError(err error, message string) bool {
	if err != nil {
		log.Printf("%s: %v", message, err)
		return false
	}
	return true
}

func Unzip2Dir(zip_path string, dest_path string) error {
	r, err := zip.OpenReader(zip_path)
	if err != nil {
		return err
	}
	defer r.Close()

	os.MkdirAll(dest_path, 0755)

	// Iterate through the files in the archive,
	// creating directories as needed and writing files.
	for _, item := range r.File {
		rc, err := item.Open()
		if err != nil {
			return err
		}
		defer rc.Close()

		fpath := filepath.Join(dest_path, item.Name)
		if item.FileInfo().IsDir() {
			os.MkdirAll(fpath, item.Mode())
		} else {
			var fdir string
			if lastIndex := strings.LastIndex(fpath, string(os.PathSeparator)); lastIndex > -1 {
				fdir = fpath[:lastIndex]
			}

			err = os.MkdirAll(fdir, item.Mode())
			if err != nil {
				return err
			}
			item, err := os.OpenFile(fpath, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, item.Mode())
			if err != nil {
				return err
			}
			defer item.Close()

			_, err = io.Copy(item, rc)
			if err != nil {
				return err
			}
		}
	}
	return nil
}
