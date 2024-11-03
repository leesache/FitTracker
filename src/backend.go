/*
MIT License

Copyright (c) 2024 Rushan Valimkhanov

Permission is hereby granted, free of charge, to any person obtaining a copy
of this software and associated documentation files (the "Software"), to deal
in the Software without restriction, including without limitation the rights
to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
copies of the Software, and to permit persons to whom the Software is
furnished to do so, subject to the following conditions:

The above copyright notice and this permission notice shall be included in all
copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
SOFTWARE.
*/

package main

import (
	"fmt"
	"log"
	"strconv"
	"strings"

	"github.com/xuri/excelize/v2"
	// "fyne.io/fyne/v2"
)

type TableContent struct {
	TableMatrix [][]string
	RowsNum     uint16
	ColsNum     uint16
}

type TrainingDay struct {
	ExcerciseName []string
	Rep           []Excercise
	CurrentWeight []float32
}

type Excercise struct {
	Sets int
	Min  int
	Max  int
}

func fillTrainingDay(table TableContent) []TrainingDay {
	numTrainingDays := int(table.ColsNum) / 4
	fmt.Print(numTrainingDays)
	train := make([]TrainingDay, numTrainingDays)

	fillExcerciseName(table, train)
	fillRep(table, train)
	fillCurrentWeight(table, train)

	return train
}

func fillExcerciseName(table TableContent, train []TrainingDay) {
	day := 0
	for i := 0; i < int(table.ColsNum); i += 4 {
		if day >= len(train) {
			break
		}

		for j := 0; j < int(table.RowsNum); j++ {
			exerciseName := table.TableMatrix[j][i]
			train[day].ExcerciseName = append(train[day].ExcerciseName, exerciseName)
		}
		day++
	}
}

// 0 0, 1 0, 2 0, 3 0, 4 0, 5 0, 6 0
// 0 4, 1 4, 2 4, 3 4, 4 4, 5 4, 6 4

func fillRep(table TableContent, train []TrainingDay) {
	for i := 0; i < len(train); i++ {
		for j := 0; j < int(table.RowsNum); j++ {
			for k := 2; k < int(table.ColsNum); k += 4 { // iterating over every fourth column, starting from the third
				sets, min, max, err := parseExerciseFormat(table.TableMatrix[j][k])
				if err != nil {
					fmt.Printf("Error parsing Rep format at [%d][%d]: %v\n", j, k, err)
					continue
				}
				exercise := Excercise{Sets: sets, Min: min, Max: max}
				train[i].Rep = append(train[i].Rep, exercise)
			}
		}
	}
}

func parseExerciseFormat(exercise string) (int, int, int, error) {
	parts := strings.FieldsFunc(exercise, func(r rune) bool {
		return r == 'x' || r == 'х'
	})

	if len(parts) != 2 {
		return 0, 0, 0, fmt.Errorf("invalid format: %s", exercise)
	}

	firstInt, err := strconv.Atoi(parts[0])
	if err != nil {
		return 0, 0, 0, fmt.Errorf("invalid number: %s", parts[0])
	}

	rangeParts := strings.Split(parts[1], "-")
	if len(rangeParts) != 2 {
		return 0, 0, 0, fmt.Errorf("invalid range format: %s", parts[1])
	}

	secondInt, err := strconv.Atoi(rangeParts[0])
	if err != nil {
		return 0, 0, 0, fmt.Errorf("invalid number: %s", rangeParts[0])
	}

	thirdInt, err := strconv.Atoi(rangeParts[1])
	if err != nil {
		return 0, 0, 0, fmt.Errorf("invalid number: %s", rangeParts[1])
	}

	return firstInt, secondInt, thirdInt, nil
}

func fillCurrentWeight(table TableContent, train []TrainingDay) {
	for i := 0; i < len(train); i++ {
		for j := 0; j < int(table.RowsNum); j++ {
			for k := 1; k < int(table.ColsNum); k += 4 {
				str := strings.ReplaceAll(table.TableMatrix[j][k], ",", ".")
				weight, err := strconv.ParseFloat(str, 32)
				if err != nil {
					fmt.Printf("Error parsing weight at [%d][%d]: %v\n", j, k, err)
					continue
				}
				train[i].CurrentWeight = append(train[i].CurrentWeight, float32(weight))
			}
		}
	}
}

func main() {
	var table TableContent
	table.TableMatrix = make([][]string, 0)

	training := excelImportData(table)

	for i, td := range training {
		fmt.Printf("Training Day %d:\n", i+1)
		for j, exerciseName := range td.ExcerciseName {
			fmt.Printf("  Exercise: %s, Sets: %d, Min: %d, Max: %d, Current Weight: %.2f\n",
				exerciseName, td.Rep[j].Sets, td.Rep[j].Min, td.Rep[j].Max, td.CurrentWeight[j])
		}
	}
}

func excelImportData(table TableContent) []TrainingDay {
	f, err := excelize.OpenFile("FitTrackerTable.xlsx") // make it work on any .xlsx file in folder.
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()

	sheetName := "Sheet1"

	rows, err := f.GetRows(sheetName)
	if err != nil {
		log.Fatal(err)
	}
	cols, err := f.GetCols(sheetName)
	if err != nil {
		log.Fatal(err)
	}
	table.RowsNum, table.ColsNum = uint16(len(rows)), uint16(len(cols))
	for i := 0; i < int(table.RowsNum); i++ {
		row := make([]string, 0, table.ColsNum)
		for j := 0; j < int(table.ColsNum); j++ {
			cellValue, err := f.GetCellValue(sheetName, colIndex(j+1)+strconv.Itoa(i+1))
			if err != nil {
				log.Printf("Error getting cell value at row %d, column %d: %v", i, j, err)
				continue
			}
			row = append(row, cellValue)
		}
		table.TableMatrix = append(table.TableMatrix, row)
	}

	return fillTrainingDay(table)
}

func colIndex(num int) string {
	return string(rune('A' + (num - 1)))
}
