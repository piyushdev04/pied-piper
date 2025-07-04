package main

import (
	"fmt"
	"io"
	"net/http"
	"os"

	"github.com/klauspost/compress/zstd"
)

func main() {
	// Serve static files (frontend)
	fs := http.FileServer(http.Dir("../frontend"))
	http.Handle("/static/", http.StripPrefix("/static/", fs))
	http.HandleFunc("/", serveIndex)

	// API endpoints
	http.HandleFunc("/compress", handleCompress)
	http.HandleFunc("/decompress", handleDecompress)

	fmt.Println("Pied Piper server running on http://localhost:8080 ...")
	http.ListenAndServe(":8080", nil)
}

func serveIndex(w http.ResponseWriter, r *http.Request) {
	indexPath := "../frontend/index.html"
	if _, err := os.Stat(indexPath); os.IsNotExist(err) {
		http.Error(w, "index.html not found", 404)
		return
	}
	http.ServeFile(w, r, indexPath)
}

func handleCompress(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}
	file, _, err := r.FormFile("file")
	if err != nil {
		http.Error(w, "Failed to read file", 400)
		return
	}
	defer file.Close()

	w.Header().Set("Content-Type", "application/zstd")
	w.Header().Set("Content-Disposition", "attachment; filename=compressed.zst")

	encoder, err := zstd.NewWriter(w)
	if err != nil {
		http.Error(w, "Compression error", 500)
		return
	}
	defer encoder.Close()

	if _, err := io.Copy(encoder, file); err != nil {
		http.Error(w, "Compression failed", 500)
		return
	}
}

func handleDecompress(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}
	file, _, err := r.FormFile("file")
	if err != nil {
		http.Error(w, "Failed to read file", 400)
		return
	}
	defer file.Close()

	w.Header().Set("Content-Type", "application/octet-stream")
	w.Header().Set("Content-Disposition", "attachment; filename=decompressed")

	decoder, err := zstd.NewReader(file)
	if err != nil {
		http.Error(w, "Decompression error", 500)
		return
	}
	defer decoder.Close()

	if _, err := io.Copy(w, decoder); err != nil {
		http.Error(w, "Decompression failed", 500)
		return
	}
}
