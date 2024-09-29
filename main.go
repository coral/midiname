package main

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/coral/midiname/db"
	"github.com/coral/midiname/reader"
)

func main() {

	db, err := db.New("midiname.db")
	if err != nil {
		panic(err)
	}

	defer db.Close()

	// ai, err := ai.New("prompt.txt")
	// if err != nil {
	// 	panic(err)
	// }

	mPath := "/Users/coral/Downloads/MIDI"

	files, err := os.ReadDir(mPath)
	if err != nil {
		panic(err)
	}

	for _, file := range files {
		if !file.IsDir() {
			filePath := filepath.Join(mPath, file.Name())
			// resp, err := ai.TryFile(filePath)
			// if err != nil {
			// 	fmt.Println("invalid json", err)
			// 	continue
			// }
			// err = db.Add(resp, file.Name())
			// if err != nil {
			// 	fmt.Println("error adding to db", err)
			// 	continue
			// }
			// fmt.Printf("%s - %s\n", resp.Title, resp.Artist)

			records, err := reader.ReadMidiFile(filePath)
			fmt.Println(records, err)
		}
	}

}
