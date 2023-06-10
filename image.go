package main

import (
	"bytes"
	"image"
	"image/png"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
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
	if err := dc.LoadFontFace("/Library/Fonts/Arial Unicode.ttf", 32); err != nil { // Specify the path to your desired font file
		return nil, err
	}

	// Set the text color
	dc.SetRGB(1, 1, 1) // White color

	// Draw the text on the image
	margin := 20.0
	width := float64(dc.Width()) - 2*margin
	x := margin
	y := margin
	dc.DrawStringWrapped(text, x, y, 0, 0, width, 1.5, gg.AlignLeft)

	buf := new(bytes.Buffer)
	err = png.Encode(buf, dc.Image())
	if err != nil {
		return nil, err
	}

	return buf.Bytes(), nil
}

func ListFiles(directory string) ([]string, error) {
	fileInfo, err := os.ReadDir(directory)
	if err != nil {
		return nil, err
	}

	var files []string
	for _, file := range fileInfo {
		files = append(files, file.Name())
	}

	return files, nil
}

func main() {
	// Create a file server which serves files out of the "images" directory.
	// Note: The file server is wrapped in the http.StripPrefix function to
	// remove the "/images" prefix when looking for files.
	fs := http.StripPrefix("/images/", http.FileServer(http.Dir("/images")))
	http.Handle("/images/", fs)

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("Hello, server is running!"))
	})

	http.HandleFunc("/list", func(w http.ResponseWriter, r *http.Request) {
		files, err := ListFiles("images")
		if err != nil {
			log.Fatal(err)
		}

		for _, file := range files {
			w.Write([]byte("<a href='" + file + "'>" + file + "</a><br>" + "\n"))
		}

		//w.Write([]byte("Hello, server is running!"))
	})
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "index.html")
	})

	http.HandleFunc("/image", func(w http.ResponseWriter, r *http.Request) {

		var text string
		if r.Method == http.MethodPost {
			body, err := ioutil.ReadAll(r.Body)
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
			values, err := url.ParseQuery(string(body))
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
			text = values.Get("text")
		} else {
			text = "Hello, World!"
		}

		img, err := createImage(text)

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
