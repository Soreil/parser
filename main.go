package main

import (
	"bytes"
	"fmt"
	"image"
	_ "image/gif"
	_ "image/jpeg"
	_ "image/png"
	"log"
	"net/http"
	_ "net/http/pprof"
	"os"
	"strings"

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
	go func() {
		log.Println(http.ListenAndServe("localhost:6060", nil))
	}()

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
		fis[os.Args[1][:strings.Index(os.Args[1], fi.Name())]] = []os.FileInfo{fi}
	}

	var fileCount int
	for dir, fi := range fis {
		for _, fi := range fi {
			if fi.IsDir() {
				continue
			}
			s, err := parse(fi, dir)
			if err != nil {
				log.Println(err)
			} else {
				fmt.Println(s)
			}
			fileCount++
		}
	}
	fmt.Println("Filecount", fileCount)
}

func parse(fi os.FileInfo, dir string) (string, error) {
	var out string

	out += fmt.Sprintln(dir + fi.Name())
	f, err := os.Open(dir + fi.Name())
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
