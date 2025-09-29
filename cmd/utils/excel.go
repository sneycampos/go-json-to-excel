package utils

import (
	"encoding/json"
	"fmt"
	"os"
	"sort"
	"strconv"
	"time"

	"github.com/xuri/excelize/v2"
)

func GenerateExcelFromJson(filepath string) (excelFilePath string, err error) {
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

	file, err := os.Open(filepath)
	if err != nil {
		fmt.Println(err)
		return "", err
	}
	defer file.Close()
	defer os.Remove(filepath) // clean up the temp file after processing

	dec := json.NewDecoder(file)
	var rowMap map[string][]map[string]interface{}
	if err := dec.Decode(&rowMap); err != nil {
		fmt.Println("Error decoding JSON:", err)
		return "", err
	}

	// Collect all column letters to determine max columns
	colSet := make(map[string]struct{})
	for _, cells := range rowMap {
		for _, cellObj := range cells {
			for col := range cellObj {
				colSet[col] = struct{}{}
			}
		}
	}

	// Sort columns alphabetically (A, B, C, ...)
	var columns []string
	for col := range colSet {
		columns = append(columns, col)
	}

	// Sort columns
	sort.Slice(columns, func(i, j int) bool { return columns[i] < columns[j] })

	// Sort row keys numerically - excelize requires rows in order to write using StreamWriter
	var rowKeys []int
	for k := range rowMap {
		if n, err := strconv.Atoi(k); err == nil {
			rowKeys = append(rowKeys, n)
		}
	}
	sort.Ints(rowKeys)

	for _, rowNum := range rowKeys {
		rowStr := strconv.Itoa(rowNum)
		cells := rowMap[rowStr]

		for _, cellObj := range cells {
			row := make([]interface{}, len(columns))
			for i, col := range columns {
				if val, ok := cellObj[col]; ok {
					row[i] = val
				} else {
					row[i] = nil
				}
			}
			cell, _ := excelize.CoordinatesToCellName(1, rowNum)

			if err := sw.SetRow(cell, row); err != nil {
				fmt.Println("Error writing row", rowNum, err)
				return "", err
			}
		}
	}

	if err := sw.Flush(); err != nil {
		fmt.Println(err)
		return "", err
	}

	fileName := fmt.Sprintf("temp/Generated_%d.xlsx", time.Now().Unix())
	if err := f.SaveAs(fileName); err != nil {
		fmt.Println(err)
		return "", err
	}

	return fileName, nil
}
