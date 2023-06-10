package main

import (
	"log"
	"os"
	"runtime"
)

func getFontPath() string {
	if runtime.GOOS == "windows" {
		fontPath := os.Getenv("SystemRoot") + "\\Fonts\\Arial.ttf"
		if _, err := os.Stat(fontPath); err == nil {
			return fontPath
		}
	} else {
		fontPath := "/Library/Fonts/Arial Unicode.ttf"
		if _, err := os.Stat(fontPath); err == nil {
			return fontPath
		}
	}

	log.Fatal("Font file not found")
	return ""
}
