package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"strings"

	"github.com/rs/cors"
)

func baseHandler(w http.ResponseWriter, req *http.Request) {
	allFilesInJson := getAllFilesInCurrentDir()
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	w.Header().Set("Access-Control-Allow-Origin", "*")

	encoder := json.NewEncoder(w)

	encoder.Encode(allFilesInJson)
}

func main() {
	port := ":10001"
	fs := http.FileServer(http.Dir("."))

	mux := http.NewServeMux()
	mux.Handle("/images/", http.StripPrefix("/images", fs))

	mux.HandleFunc("/", baseHandler)

	handler := cors.Default().Handler(mux)
	fmt.Printf("I'm running at http://localhost%s\n", port)
	http.ListenAndServe(port, handler)
}

type PlatLocalResponse struct {
	Name            string
	ImageUrl        string
	WikiPage        string
	Classifications []string
}

func getAllFilesInCurrentDir() []PlatLocalResponse {
	files, err := os.ReadDir(".")
	if err != nil {
		panic("Shida bro na file")
	}

	var crops []PlatLocalResponse

	for _, file := range files {
		if strings.Contains(file.Name(), ".jpeg") || strings.Contains(file.Name(), ".webp") || strings.Contains(file.Name(), ".png") {
			parts := strings.Split(file.Name(), "-")
			crops = append(crops, PlatLocalResponse{
				Name:            file.Name(),
				ImageUrl:        "http://localhost:10001/images/" + file.Name(),
				Classifications: strings.Split(file.Name(), "-"),
				WikiPage:        "https://en.wikipedia.org/wiki/" + parts[0],
			})
		}
	}

	return crops
}
