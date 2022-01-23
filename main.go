package main

import (
	"fmt"
	"os"
)

func main() {
	// apiUrl := "https://api.esv.org/v3/passage/text"

	apiKey := os.Getenv("ESV_API_KEY")

	fmt.Println("Hello, ESV! ", apiKey)
}
