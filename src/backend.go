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
	// Day           uint16 // Is this thing needed?
	ExcerciseName []string
	Rep           []Excercise
	CurrentWeight []float32
}

type Excercise struct {
	Sets int
	Min  int
	Max  int
}

// func fillDay (table TableContent, train []TrainingDay) {
// 	for i := 0; i < int(table.RowsNum); i++{
// 		for j := 0; j < int(table.ColsNum); j++{

// 		}
// 	}
// }

func fillTrainingDay(table TableContent) []TrainingDay {
	// Determine the number of training days based on columns
	numTrainingDays := int(table.ColsNum) / 4
	fmt.Print(numTrainingDays)
	train := make([]TrainingDay, numTrainingDays) // Initialize slice with appropriate size

	// Use helper functions to fill in each TrainingDay
	fillExcerciseName(table, train)
	fillRep(table, train)
	fillCurrentWeight(table, train)

	return train
}

func fillExcerciseName(table TableContent, train []TrainingDay) {
	for i := 0; i < len(train); i++ { // iterating through each training day
		for j := 0; j < int(table.RowsNum); j++ { // iterating through rows
			for k := 0; k < int(table.ColsNum); k += 4 { // iterating every fourth column for exercise names
				exerciseName := table.TableMatrix[j][k]
				train[i].ExcerciseName = append(train[i].ExcerciseName, exerciseName)
				// fmt.Printf("Filling ExcerciseName for TrainingDay %d with [%d][%d]: %v\n", i+1, j, k, exerciseName)
			}
		}
	}
}

func fillRep(table TableContent, train []TrainingDay) {
	for i := 0; i < len(train); i++ {
		for j := 0; j < int(table.RowsNum); j++ {
			for k := 2; k < int(table.ColsNum); k += 4 { // iterating over every fourth column, starting from the third
				sets, min, max, err := parseExerciseFormat(table.TableMatrix[j][k])
				if err != nil {
					fmt.Printf("Error parsing Rep format at [%d][%d]: %v\n", j, k, err)
					continue
				}
				// Create and append a new Excercise struct to the Rep slice
				exercise := Excercise{Sets: sets, Min: min, Max: max}
				train[i].Rep = append(train[i].Rep, exercise)
				// fmt.Printf("Filling Rep for TrainingDay %d with [%d][%d]: %+v\n", i+1, j, k, exercise)
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
			for k := 1; k < int(table.ColsNum); k += 4 { // Iterating over every fourth column, starting from the second
				str := strings.ReplaceAll(table.TableMatrix[j][k], ",", ".")
				weight, err := strconv.ParseFloat(str, 32)
				if err != nil {
					fmt.Printf("Error parsing weight at [%d][%d]: %v\n", j, k, err)
					continue
				}
				// Append the parsed weight to the CurrentWeight slice
				train[i].CurrentWeight = append(train[i].CurrentWeight, float32(weight))
				// fmt.Printf("Filling CurrentWeight for TrainingDay %d with [%d][%d]: %.2f\n", i+1, j, k, weight)
			}
		}
	}
}

func main() {
	// Initialize TableContent
	var table TableContent
	table.TableMatrix = make([][]string, 0) // Initialize the TableMatrix as a slice of slices

	// Import data from the Excel file
	training := excelImportData(table) // Pass the address of table and training to modify them

	// Print all details of each TrainingDay
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
		row := make([]string, 0, table.ColsNum) // Create a new slice for the current row
		for j := 0; j < int(table.ColsNum); j++ {
			cellValue, err := f.GetCellValue(sheetName, colIndex(j+1)+strconv.Itoa(i+1)) // Use i+1 for 1-based row index
			if err != nil {
				log.Printf("Error getting cell value at row %d, column %d: %v", i, j, err)
				continue
			}
			row = append(row, cellValue) // Append the cell value to the current row
		}
		table.TableMatrix = append(table.TableMatrix, row) // Append the current row to TableMatrix
	}

	return fillTrainingDay(table)
}

func colIndex(num int) string {
	// Make AA BB ETC

	return string(rune('A' + (num - 1)))
}
