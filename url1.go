package main

import (
	"crypto/sha1"
	"encoding/hex"
	"fmt"
	"net/http"
)

func longURL(w http.ResponseWriter, r *http.Request) {
	var input = "https://docs.google.com/document/d/1lydfEEGJ8UbFokoHJAPvN1wgnsEV4UDN6KLVeIA5Gvg/edit#heading=h.2sx5w06aktfj"
	fmt.Fprintln(w, "----URL Shortener program----")

	fmt.Fprintln(w, "Original URL : ", input)
	short_url := shortURL1(input)
	fmt.Fprintln(w, "Short URL is", short_url)

}

func shortURL1(input string) string {
	hasher := sha1.New()

	hasher.Write([]byte(input))
	hasherValue := hasher.Sum(nil)
	hasherString := hex.EncodeToString(hasherValue)
	return hasherString[:12]

}

func main() {

	http.HandleFunc("/", longURL)

	// Start the server on port 8080
	port := 8080
	fmt.Printf("Server listening on :%d\n", port)
	err := http.ListenAndServe(fmt.Sprintf(":%d", port), nil)
	if err != nil {
		fmt.Println("Error:", err)
	}
}
