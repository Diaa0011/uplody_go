package handlersv1

import (
	"fmt"
	"net/http"
)

func ChunkUploadHandler(w http.ResponseWriter, r *http.Request) {
	err := r.ParseMultipartForm(10 << 20)

	if err != nil {
		http.Error(w, "Error parsing from", http.StatusBadRequest)
		return
	}

	file, _, err := r.FormFile("chunk")

	if err != nil {
		http.Error(w, "Error retrieving chunk", http.StatusBadRequest)
		return
	}
	defer file.Close()

	filename := r.FormValue("filename")
	chunkIndex := r.FormValue("chunkIndex")

	if err := uploadChunkToS3(file, filename, chunkIndex); err != nil {
		http.Error(w, "Error uploading to S3", http.StatusInternalServerError)
		fmt.Println("--------------------------------")
		fmt.Printf("AWS UPLOAD FAILED: %v\n", err)
		fmt.Println("--------------------------------")
		return
	}

	fmt.Fprintf(w, "Chunk %s uploaded successfully!", chunkIndex)

}
