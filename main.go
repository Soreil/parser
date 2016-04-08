package main

import (
	"bytes"
	"fmt"
	"image"
	"log"
	"os"

	"github.com/Soreil/audio"
)

func readFiles(name string) (map[string][]os.FileInfo, error) {
	if name[len(name)-1] != os.PathSeparator {
		name += string(os.PathSeparator)
	}

	f, err := os.Open(name)
	if err != nil {
		return nil, err
	}
	fis, err := f.Readdir(0)
	if err != nil {
		return nil, err
	}
	dirfis := map[string][]os.FileInfo{
		name: fis,
	}
	for _, fi := range fis {
		if fi.IsDir() {
			newfis, err := readFiles(name + fi.Name())
			if err != nil {
				return nil, err
			}
			for k, v := range newfis {
				dirfis[k] = append(dirfis[k], v...)
			}
		}
	}
	return dirfis, nil
}

func main() {
	if len(os.Args) <= 1 {
		log.Fatalf("Usage:%s MEDIAPATH\n", os.Args[0])
	}

	fi, err := os.Lstat(os.Args[1])
	if err != nil {
		log.Fatalln(err)
	}
	fis := make(map[string][]os.FileInfo)
	if fi.IsDir() {
		fis, err = readFiles(os.Args[1])
		if err != nil {
			log.Fatalln(err)
		}
	} else {
		fis[os.Args[1]] = []os.FileInfo{fi}
	}

	var fileCount int
	for dir, fi := range fis {
		for _, fi := range fi {
			f, err := os.Open(dir + fi.Name())
			if err != nil {
				log.Fatalln(err)
			}

			if fi.IsDir() {
				continue
			}
			fmt.Println(dir + fi.Name())
			d, err := audio.NewDecoder(f)
			if err := f.Close(); err != nil {
				log.Fatal(err)
			}
			if err != nil {
				log.Println("audio error:", err)
				continue
			}
			fmt.Println(d.AudioFormat())
			fmt.Println(d.Bitrate() / 1024)
			fmt.Println(d.Duration())

			if d.HasImage() {
				fmt.Println("Has an image")
				//_, img, err := imager.Thumbnail(bytes.NewBuffer(d.Picture()), imager.Sharp)
				_, img, err := image.DecodeConfig(bytes.NewBuffer(d.Picture()))
				if err != nil {
					log.Println(err)
				} else {
					fmt.Println(img)
				}
			}
			fmt.Println()
			fileCount++
		}
	}
	fmt.Println("Filecount", fileCount)
}
