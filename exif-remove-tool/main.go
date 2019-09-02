package main

import (
	"bytes"
	"errors"
	"fmt"
	"image"

	"io/ioutil"
	"os"
	"path/filepath"
	"strings"

	"../../go-exif-remove"
)

func main() {

	if len(os.Args) == 1 {
		var files []string
		root := "img"
		err := filepath.Walk(root, func(path string, info os.FileInfo, err error) error {
			if path != "img" {
				files = append(files, path)
			}
			return nil
		})
		if err != nil {
			panic(err)
		}
		pass := 0
		fail := 0
		for _, file := range files {
			fmt.Println(file)
			if b, err := handleFile(file); err != nil {
				fail += 1
			} else {
				pass += 1
				f := filepath.Base(file)
				ioutil.WriteFile("img_output/"+f, b, 0644)
			}
			fmt.Println()
		}

		percentage := 100 * pass / (pass + fail)
		fmt.Printf("Results (%v%%): %v pass, %v fail \n", int(percentage), pass, fail)
	} else {
		path := os.Args[1]
		if b, err := handleFile(path); err != nil {
			fmt.Printf(err.Error())
		} else {
			file := filepath.Base(path)
			ioutil.WriteFile("img_output/"+file, b, 0644)
		}
	}

}

func handleFile(filepath string) ([]byte, error) {
	if data, err := ioutil.ReadFile(filepath); err != nil {
		fmt.Printf(err.Error())
		return nil, err
	} else {
		//_, err = jpeg.Decode(bytes.NewReader(data))
		//_, err = png.Decode(bytes.NewReader(data))
		//if err != nil  {
		//	fmt.Printf("ERROR: original image is corrupt" + err.Error() + "\n")
		//	return nil, err
		//}
		_, _, err := image.Decode(bytes.NewReader(data))
		if err != nil {
				fmt.Printf("ERROR: original image is corrupt" + err.Error() + "\n")
				return nil, err
		}
		filtered, err := exifremove.Remove(data)
		if err != nil {
			if !strings.EqualFold(err.Error(), "no exif data") && !strings.EqualFold(err.Error(), "file does not have EXIF") {
				fmt.Printf("* " + err.Error() + "\n")
				return nil, errors.New(err.Error())
			}
		}
		return filtered, nil
	}
}