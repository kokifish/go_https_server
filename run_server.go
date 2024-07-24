package main

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"

	"github.com/kokifish/go_https_server/utils"
)

const REPORT_PATH = "report_upload"
const REPORT_HTML_PATH = "report_html"
const MAN_POC_PATH = "manual_poc"
const RECV_ZIP_PATH = "recv_zip"
const RECV_UNZIP_PATH = "recv_unzip"

func uploadV8ZipHandler(w http.ResponseWriter, r *http.Request) {
	r.ParseMultipartForm(1024 * 1024 * 1024) // 1GB?
	file, handler, err := r.FormFile("V8zip")
	if !utils.CheckIfNoError(err, "ERROR! Retrieving the File") {
		return
	}
	defer file.Close()

	// === extract [tag][version] from req, version is optional
	tag := r.FormValue("tag")
	if tag == "" {
		http.Error(w, "ERROR! Missing required parameter: tag", http.StatusBadRequest)
		return
	}
	version := r.FormValue("version")
	if version != "" {
		tag = tag + version
	}
	fmt.Println("recv a V8 zip from", tag)

	// === save recv zip to a local file
	tmp_save_zip_fname := filepath.Join(RECV_ZIP_PATH, tag+".zip")
	dst, err := os.Create(tmp_save_zip_fname)
	if err != nil {
		http.Error(w, "ERROR! Unable to create the file for writing.", http.StatusInternalServerError)
		return
	}
	defer dst.Close()
	_, err = io.Copy(dst, file)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// === unzip to recv_unzip/[tag][version] if it's a valid zip
	unzipPath := filepath.Join(RECV_UNZIP_PATH, tag)

	// Check if unzipPath exists, if so, delete it
	if _, err := os.Stat(unzipPath); !os.IsNotExist(err) {
		err = os.RemoveAll(unzipPath)
		if err != nil {
			http.Error(w, "ERROR! Failed to delete "+unzipPath+err.Error(), http.StatusInternalServerError)
			return
		}
	}
	err = utils.Unzip2Dir(tmp_save_zip_fname, unzipPath)
	if err != nil {
		http.Error(w, "ERROR! Failed to unzip file: "+err.Error(), http.StatusInternalServerError)
		return
	}
	w.Write([]byte(fmt.Sprintf("succ upload %s with tag %s", handler.Filename, tag)))
}

func uploadReportHandler(w http.ResponseWriter, r *http.Request) {
	r.ParseMultipartForm(10 * 1024 * 1024) // 10MB?
	file, handler, err := r.FormFile("report_file")
	if err == nil {
		defer file.Close()

		dst, err := os.Create(filepath.Join(REPORT_PATH, handler.Filename))
		if err != nil {
			http.Error(w, "ERROR! Unable to create the file for writing.", http.StatusInternalServerError)
			return
		}
		defer dst.Close()

		_, err = io.Copy(dst, file)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		w.Write([]byte(fmt.Sprintf("succ uploaded report_file: %s;", handler.Filename)))
	}

	htmlFile, htmlHandler, err := r.FormFile("report_html")
	if err == nil {
		defer htmlFile.Close()

		dst, err := os.Create(filepath.Join(REPORT_HTML_PATH, htmlHandler.Filename))
		if err != nil {
			http.Error(w, "ERROR! Unable to create the file for writing.", http.StatusInternalServerError)
			return
		}
		defer dst.Close()

		_, err = io.Copy(dst, htmlFile)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		w.Write([]byte(fmt.Sprintf("succ uploaded report_html: %s\n", htmlHandler.Filename)))
	}
}

func main() {
	utils.CreateDirIfNotExits(REPORT_PATH)
	utils.CreateDirIfNotExits(REPORT_HTML_PATH)
	utils.CreateDirIfNotExits(MAN_POC_PATH)
	utils.CreateDirIfNotExits(RECV_ZIP_PATH)
	utils.CreateDirIfNotExits(RECV_UNZIP_PATH)
	http.HandleFunc("/uploadreport", uploadReportHandler)
	http.HandleFunc("/uploadV8Zip", uploadV8ZipHandler)

	err := http.ListenAndServeTLS(":443", "server.crt", "server.key", nil) // port cert key
	if err != nil {
		panic(err)
	}
}
