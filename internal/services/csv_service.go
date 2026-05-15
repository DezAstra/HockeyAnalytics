package services

import (
	"encoding/csv"
	"os"
)

func ReadCSV(path string) ([][]string, error) {

	file, err := os.Open(path)

	if err != nil {
		return nil, err
	}

	defer file.Close()

	reader := csv.NewReader(file)

	rows, err := reader.ReadAll()

	if err != nil {
		return nil, err
	}

	return rows, nil
}
