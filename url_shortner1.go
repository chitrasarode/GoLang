package main

import (
	"crypto/sha1"
	"encoding/hex"
	"fmt"
	"os"
)

var shortURL = make(map[string]string)

func main() {

	var input string
	var ch int
	var short_url string

	for {
		fmt.Println("\n----URL Shortener Program----")
		fmt.Print("\n1.Get the short URL\n2.Expand the URL\n3.Exit\n")
		fmt.Println("Enter your choice : ")
		fmt.Scan(&ch)

		switch ch {
		case 1:
			fmt.Print("Enter the URL : ")
			fmt.Scan(&input)

			short_url = shorten(input)

			shortURL[short_url] = input
			fmt.Println("Shorten URL is : ", short_url)

		case 3:
			long_url, found := expand(short_url)
			if found {
				fmt.Printf("Expanded URL is : %s", long_url)
			} else {
				fmt.Println("Short URL not found.")
			}
		case 4:
			fmt.Println("Enter correct choice")
		default:
			os.Exit(1)

		}

	}

}

func shorten(input string) string {
	hasher := sha1.New()
	hasher.Write([]byte(input))
	hasherValue := hasher.Sum(nil)
	//fmt.Println(hasherValue)
	hasherString := hex.EncodeToString(hasherValue)
	return hasherString[:12]
}

func expand(short_url string) (string, bool) {
	longURL, found := shortURL[short_url]
	return longURL, found
}
