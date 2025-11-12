package main

import (
	"app/internal/avif"
	"io"
	"log"
	"os"
)

func main() {
	file, err := os.Open("/home/sangle/Downloads/vovkapanda.jpg")
	if err != nil {
		return
	}
	defer file.Close()
	buffer, err := io.ReadAll(file)
	if err != nil {
		panic(err)
	}
	imgBuffer, err := avif.EncodeImageToAVIF(buffer)
	if err != nil {
		log.Fatal(err.Error())
	}
	os.WriteFile("/home/sangle/Downloads/vovkapanda.avif", imgBuffer, 0644)
}
