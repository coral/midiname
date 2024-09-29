package main

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"

	"github.com/coral/midiname/ai"
	"github.com/coral/midiname/db"
	"github.com/coral/midiname/reader"
)

func main() {

	db, err := db.New("midiname.db")
	if err != nil {
		panic(err)
	}

	defer db.Close()

	ai, err := ai.New("prompt.txt")
	if err != nil {
		panic(err)
	}

	mPath := "/Users/coral/Downloads/MIDI"

	files, err := os.ReadDir(mPath)
	if err != nil {
		panic(err)
	}

	for _, file := range files {
		if !file.IsDir() {
			filePath := filepath.Join(mPath, file.Name())

			records, err := reader.ReadMidiFile(filePath)
			if err != nil {
				fmt.Println("invalid midi", err)
				continue
			}
			filtered := filter(records)

			resp, err := ai.TryFile(filePath, filtered)
			if err != nil {
				fmt.Println("invalid json", err)
				continue
			}
			// err = db.Add(resp, file.Name())
			// if err != nil {
			// 	fmt.Println("error adding to db", err)
			// 	continue
			// }
			fmt.Printf("%+v", resp)

		}
	}

}

func filter(records [][]string) string {
	var result struct {
		Text         []string `json:"Text,omitempty"`
		Title        []string `json:"Title,omitempty"`
		Lyric        []string `json:"Lyric,omitempty"`
		KeySignature []string `json:"KeySignature,omitempty"`
	}

	for _, record := range records {
		if len(record) < 3 {
			continue
		}

		switch record[2] {
		case "Text_t":
			if len(record) > 3 {
				result.Text = append(result.Text, record[3])
			}
		case "Title_t":
			if len(record) > 3 {
				result.Title = append(result.Title, record[3])
			}
		case "Lyric_t":
			if len(record) > 3 {
				result.Lyric = append(result.Lyric, record[3])
			}
		case "Key_signature":
			if len(record) > 3 {
				result.KeySignature = append(result.KeySignature, record[3])
			}
		}
	}

	jsonResult, err := json.Marshal(result)
	if err != nil {
		return ""
	}

	return string(jsonResult)
}
