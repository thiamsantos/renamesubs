package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path"
	"path/filepath"
	"regexp"
	"strings"
)

const baseFolder = "."

func shouldRename(file os.FileInfo) bool {
	re := regexp.MustCompile(`S\d{2}E\d{2}`)
	name := file.Name()
	ext := filepath.Ext(name)
	supportedExtensions := []string{".mp4", ".srt"}

	for _, a := range supportedExtensions {
		if a == ext && re.MatchString(name) && (strings.Contains(name, "1080p") || strings.Contains(name, "720p")) {
			return true
		}
	}
	return false
}

func main() {
	files, err := ioutil.ReadDir(baseFolder)
	if err != nil {
		log.Fatal(err)
	}

	for _, file := range files {
		ext := filepath.Ext(file.Name())

		if !shouldRename(file) {
			continue
		}

		parts := strings.Split(file.Name(), ".")

		name := []string{}
		foundEpi := false
		foundQuality := false

		for _, part := range parts {
			re := regexp.MustCompile(`^S\d{2}E\d{2}$`)

			if !foundEpi && re.MatchString(part) {
				name = append(name, part)
				foundEpi = true
				continue
			}

			if !foundQuality && (part == "1080p" || part == "720p") {
				foundQuality = true
				continue
			}

			if foundEpi && !foundQuality {
				name = append(name, part)
				continue
			}
		}

		finalName := strings.Join(name, " ") + ext

		fmt.Println(file.Name() + " ->> " + finalName)

		from := path.Join(baseFolder, file.Name())
		to := path.Join(baseFolder, finalName)

		err = os.Rename(from, to)
		if err != nil {
			log.Fatal(err)
		}
	}
}
