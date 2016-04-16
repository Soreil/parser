package main

import (
	"bytes"
	"errors"
	"fmt"
	"image"
	_ "image/gif"
	_ "image/jpeg"
	_ "image/png"
	"log"
	"os"
	"path/filepath"

	"github.com/Soreil/audio"
)

var exts map[string]int

func main() {
	if len(os.Args) <= 1 {
		log.Fatalf("Usage:%s MEDIAPATH\n", os.Args[0])
	}

	_, err := os.Lstat(os.Args[1])
	if err != nil {
		log.Fatalln(err)
	}

	if err := filepath.Walk(os.Args[1], func(url string, info os.FileInfo, err error) error {
		if info == nil {
			return errors.New("Can't open a nil file")
		}
		if info.IsDir() {
			return nil
		}
		out, err := parse(url)
		if err != nil {
			log.Println(err)
			return nil
		}
		fmt.Println(out)
		return nil
	}); err != nil {
		log.Fatalln(err)
	}
}

func parse(fileName string) (string, error) {
	var out string
	f, err := os.Open(fileName)
	if err != nil {
		return out, err
	}
	defer func() {
		if err := f.Close(); err != nil {
			log.Fatal(err)
		}
	}()

	d, err := audio.NewDecoder(f)
	if err != nil {
		return out, err
	}
	defer d.Destroy()

	out += fmt.Sprintln(fileName)
	out += fmt.Sprintln(d.AudioFormat())
	out += fmt.Sprintln(d.Bitrate() / 1024)
	out += fmt.Sprintln(d.Duration())

	if d.HasImage() {
		out += fmt.Sprintln("Has an image")
		//_, img, err := imager.Thumbnail(bytes.NewBuffer(d.Picture()), imager.Sharp)
		_, img, err := image.DecodeConfig(bytes.NewBuffer(d.Picture()))
		if err != nil {
			return out, err
		}
		out += fmt.Sprintln(img)
	}
	return out, nil
}
