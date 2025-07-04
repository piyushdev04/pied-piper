package main

import (
	"fmt"
	"io"
	"net/http"
	"os"

	"github.com/klauspost/compress/zstd"
)

// CORS middleware for API endpoints
func withCORS(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "https://pied-piper-terminal.vercel.app")
		w.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
		if r.Method == "OPTIONS" {
			w.WriteHeader(http.StatusOK)
			return
		}
		next(w, r)
	}
}

func main() {
	// Serve static files (frontend)
	fs := http.FileServer(http.Dir("../frontend"))
	http.Handle("/static/", http.StripPrefix("/static/", fs))
	http.HandleFunc("/", serveIndex)

	// API endpoints with CORS
	http.HandleFunc("/compress", withCORS(handleCompress))
	http.HandleFunc("/decompress", withCORS(handleDecompress))

	fmt.Println("Pied Piper server running on http://localhost:8080 ...")
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	http.ListenAndServe(":"+port, nil)
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
