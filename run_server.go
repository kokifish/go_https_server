package main

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
)

const REPORT_PATH = "report_upload"
const REPORT_HTML_PATH = "report_html"

func mkdir_if_not_exits(path_name string) {
	if _, err := os.Stat(path_name); os.IsNotExist(err) {
		err := os.Mkdir(path_name, 0755)
		if err != nil {
			fmt.Println("Error %w", err)
			return
		}
	}
}

func uploadReportHandler(w http.ResponseWriter, r *http.Request) {
	r.ParseMultipartForm(10 * 1024 * 1024) // 10MB?
	file, handler, err := r.FormFile("report_file")
	if err == nil {
		defer file.Close()

		dst, err := os.Create(filepath.Join(REPORT_PATH, handler.Filename))
		if err != nil {
			http.Error(w, "Unable to create the file for writing. Check your write access privilege", http.StatusInternalServerError)
			return
		}
		defer dst.Close()

		_, err = io.Copy(dst, file)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		w.Write([]byte(fmt.Sprintf("Successfully uploaded report_file: %s;", handler.Filename)))
	}

	htmlFile, htmlHandler, err := r.FormFile("report_html")
	if err == nil {
		defer htmlFile.Close()

		dst, err := os.Create(filepath.Join(REPORT_HTML_PATH, htmlHandler.Filename))
		if err != nil {
			http.Error(w, "Unable to create the file for writing. Check your write access privilege", http.StatusInternalServerError)
			return
		}
		defer dst.Close()

		_, err = io.Copy(dst, htmlFile)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		w.Write([]byte(fmt.Sprintf("Successfully uploaded report_html: %s\n", htmlHandler.Filename)))
	}
}

func main() {
	mkdir_if_not_exits(REPORT_PATH)
	mkdir_if_not_exits(REPORT_HTML_PATH)
	http.HandleFunc("/uploadreport", uploadReportHandler)

	// 启动HTTPS服务
	err := http.ListenAndServeTLS(":443", "server.crt", "server.key", nil) // port cert key
	if err != nil {
		panic(err)
	}
}
