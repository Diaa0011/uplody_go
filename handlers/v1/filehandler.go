package handlersv1

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strings"
)

func FileUploadHandler(w http.ResponseWriter, r *http.Request) {

	r.Body = http.MaxBytesReader(w, r.Body, 10<<30)

	err := r.ParseMultipartForm(10 << 20)

	if err != nil {
		http.Error(w, "File too large", http.StatusRequestEntityTooLarge)
	}

	//Retrieve the file from data
	file, handler, err := r.FormFile("textfile")

	if err != nil {
		http.Error(w, "Error in file retrieval", http.StatusBadRequest)
		return
	}
	defer file.Close()

	fileBytes, err := ioutil.ReadAll(file)

	if err != nil {
		http.Error(w, "Invalid file", http.StatusBadRequest)
		return
	}

	if !isValidFileType(fileBytes) {
		http.Error(w, "Invalid file type", http.StatusBadRequest)
		return
	}

	if err := uploadToS3(fileBytes, handler.Filename); err != nil {
		http.Error(w, "Error uploading to S3", http.StatusInternalServerError)
		fmt.Println("--------------------------------")
		fmt.Printf("AWS UPLOAD FAILED: %v\n", err)
		fmt.Println("--------------------------------")
		return
	}

	dst, err := createFile(handler.Filename)

	if err != nil {
		http.Error(w, "Error saving the file", http.StatusInternalServerError)
		return
	}
	defer dst.Close()

	if _, err := dst.ReadFrom(file); err != nil {
		http.Error(w, "Error saving the local file", http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, "Uploaded File: %s\n", handler.Filename)
	fmt.Fprintf(w, "File Size: %d\n", handler.Size)
	fmt.Fprintf(w, "MIME Header: %v\n", handler.Header)
	fmt.Fprintf(w, "File successfully uploaded to S3 and saved locally!")
}

func isValidFileType(file []byte) bool {
	fileType := http.DetectContentType(file)
	log.Printf("Detected file type: %s\n", fileType)
	return strings.HasPrefix(fileType, "image/")
}

func createFile(filename string) (*os.File, error) {
	if _, err := os.Stat("uploads"); os.IsNotExist(err) {
		// this will create the directory if it does not exist
		// numbers is about permissions, 0755 means read/write/execute for owner and read/execute for group and others
		// check the table below for more details about permissions number
		os.Mkdir("uploads", 0755)
	}

	dst, err := os.Create(filepath.Join("uploads", filename))

	if err != nil {
		return nil, err
	}

	return dst, nil

}

// Source - https://stackoverflow.com/a/31151508
// Posted by Shannon Matthews, modified by community. See post 'Timeline' for change history
// Retrieved 2026-04-26, License - CC BY-SA 4.0

// +-----+---+--------------------------+
// | rwx | 7 | Read, write and execute  |
// | rw- | 6 | Read, write              |
// | r-x | 5 | Read, and execute        |
// | r-- | 4 | Read,                    |
// | -wx | 3 | Write and execute        |
// | -w- | 2 | Write                    |
// | --x | 1 | Execute                  |
// | --- | 0 | no permissions           |
// +------------------------------------+

// +------------+------+-------+
// | Permission | Octal| Field |
// +------------+------+-------+
// | rwx------  | 0700 | User  |
// | ---rwx---  | 0070 | Group |
// | ------rwx  | 0007 | Other |
// +------------+------+-------+
