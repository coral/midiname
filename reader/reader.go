package reader

import (
	"bufio"
	"encoding/csv"
	"fmt"
	"os/exec"
	"strings"
)

func ReadMidiFile(filename string) ([][]string, error) {
	cmd := exec.Command("midicsv", filename)
	stdout, err := cmd.StdoutPipe()
	if err != nil {
		return nil, fmt.Errorf("error creating StdoutPipe: %v", err)
	}

	if err := cmd.Start(); err != nil {
		return nil, fmt.Errorf("error starting command: %v", err)
	}

	reader := csv.NewReader(bufio.NewReader(stdout))
	reader.LazyQuotes = true
	reader.FieldsPerRecord = -1 // Allow variable number of fields

	var records [][]string
	for {
		record, err := reader.Read()
		if err != nil {
			if err.Error() == "EOF" {
				break
			}
			fmt.Println(record, filename)
			return nil, fmt.Errorf("error reading CSV: %v", err)
		}
		// Trim whitespace from each field
		for i, field := range record {
			record[i] = strings.TrimSpace(field)
		}
		records = append(records, record)
	}

	if err := cmd.Wait(); err != nil {
		return nil, fmt.Errorf("command finished with error: %v", err)
	}

	return records, nil
}
