package main

import (
	"bytes"
	"image"
	"image/png"
	"log"
	"net/http"
	"os"

	"github.com/fogleman/gg"
)

func createImage(text string) ([]byte, error) {
	// Open the existing image file
	existingImgFile, err := os.Open("valley.jpg")
	if err != nil {
		return nil, err
	}
	defer existingImgFile.Close()

	// Decode the existing image
	existingImg, _, err := image.Decode(existingImgFile)
	if err != nil {
		return nil, err
	}

	// Create a new drawing context using the existing image as the base
	dc := gg.NewContextForImage(existingImg)

	// Load the font file
	if err := dc.LoadFontFace("/Library/Fonts/Arial Unicode.ttf", 32); err != nil {
		return nil, err
	}

	// Set the text color
	dc.SetRGB(1, 1, 1) // White color

	// Draw the text on the image
	dc.DrawStringAnchored(text, 100, 100, 0.5, 0.5)

	buf := new(bytes.Buffer)
	err = png.Encode(buf, dc.Image())
	if err != nil {
		return nil, err
	}

	return buf.Bytes(), nil
}

func main() {
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("Hello, server is running!"))
	})

	http.HandleFunc("/image", func(w http.ResponseWriter, r *http.Request) {
		query := r.URL.Query()

		// Get the value for a specific key, "key"
		value := query.Get("text")

		img, err := createImage(value)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		// Write the value to the HTTP response
		//fmt.Fprintf(w, "Value: %s", value)

		w.Header().Set("Content-Type", "image/png")
		w.Write(img)
	})

	log.Println("Starting server on :8088")
	err := http.ListenAndServe(":8088", nil)
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}
