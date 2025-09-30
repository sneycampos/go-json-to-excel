package internal

import (
	"encoding/json"
	"fmt"
	"os"
	"strconv"
	"time"

	"github.com/xuri/excelize/v2"
)

func GenerateExcelFromJson(filepath string) (string, error) {
	file, err := os.Open(filepath)
	if err != nil {
		return "", fmt.Errorf("failed to open file: %w", err)
	}
	defer file.Close()
	defer os.Remove(filepath)

	f := excelize.NewFile()
	defer func() {
		if err := f.Close(); err != nil {
			fmt.Println(err)
		}
	}()

	sw, err := f.NewStreamWriter("Sheet1")
	if err != nil {
		fmt.Println(err)
		return "", err
	}

	dec := json.NewDecoder(file)

	// Expect top-level object `{`
	t, err := dec.Token()
	if err != nil {
		return "", fmt.Errorf("failed to read token: %w", err)
	}
	if delim, ok := t.(json.Delim); !ok || delim != '{' {
		return "", fmt.Errorf("expected JSON object at top level")
	}

	// Iterate keys (row numbers) in streaming mode
	for dec.More() {
		// Row key
		t, err := dec.Token()
		if err != nil {
			return "", fmt.Errorf("failed to read row key: %w", err)
		}

		rowKey, ok := t.(string)
		if !ok {
			return "", fmt.Errorf("expected string row key")
		}

		rowNum, err := strconv.Atoi(rowKey)
		if err != nil {
			return "", fmt.Errorf("invalid row key: %s", rowKey)
		}

		// Decode the value: []map[string]interface{}
		var cells []map[string]interface{}
		if err := dec.Decode(&cells); err != nil {
			return "", fmt.Errorf("failed to decode row %s: %w", rowKey, err)
		}

		// Write each cell object directly to Excel
		for _, cellObj := range cells {
			// Build row slice dynamically (preserve natural order of JSON keys)
			row := make([]interface{}, 0, len(cellObj))
			for _, v := range cellObj {
				row = append(row, v)
			}

			cell, _ := excelize.CoordinatesToCellName(1, rowNum)
			if err := sw.SetRow(cell, row); err != nil {
				return "", fmt.Errorf("failed to write row %d: %w", rowNum, err)
			}
		}
	}

	// Final closing `}`
	if _, err := dec.Token(); err != nil {
		return "", fmt.Errorf("failed to close root object: %w", err)
	}

	if err := sw.Flush(); err != nil {
		return "", fmt.Errorf("failed to flush stream writer: %w", err)
	}

	fileName := fmt.Sprintf("temp/Generated_%d.xlsx", time.Now().Unix())

	if err := f.SaveAs(fileName); err != nil {
		return "", fmt.Errorf("failed to save Excel file: %w", err)
	}

	return fileName, nil
}
