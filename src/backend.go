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

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/widget"

	"github.com/xuri/excelize/v2"
)

// Global map for exercise images
var ExerciseImages = map[string]string{
	"Bent Over Dumbbell Row": "https://encrypted-tbn0.gstatic.com/images?q=tbn:ANd9GcTp0MU-e7KDM7ZVh75PAHeaPSV2xRZb99iXYg&s",
	"Squat":                  "https://encrypted-tbn0.gstatic.com/images?q=tbn:ANd9GcTp0MU-e7KDM7ZVh75PAHeaPSV2xRZb99iXYg&s",
}

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

func createDayButtons(train []TrainingDay, w fyne.Window) *fyne.Container {
	var dayButtons []fyne.CanvasObject
	for i := range train {
		dayIndex := i
		dayButtons = append(dayButtons, widget.NewButton(fmt.Sprintf("Train %d", dayIndex+1), func() {
			println("Going deeper into Train", dayIndex+1)
			w.SetContent(createExerciseButtons(train[dayIndex], train, w))
		}))
	}
	return container.NewVBox(dayButtons...)
}

func createExerciseButtons(day TrainingDay, train []TrainingDay, w fyne.Window) *fyne.Container {
	var exerciseButtons []fyne.CanvasObject
	for _, exerciseName := range day.ExcerciseName {
		name := exerciseName
		exerciseButtons = append(exerciseButtons, widget.NewButton(name, func() {
			println("Exercise:", name)
			// Access the global ExerciseImages map to get the image URL
			imageURL := ExerciseImages[name]
			image := createImageFromURL(imageURL)
			if image != nil {
				w.SetContent(container.NewVBox(
					image, // Show the exercise image
					widget.NewButton("Back", func() {
						w.SetContent(createExerciseButtons(day, train, w))
					}),
				))
			}
		}))
	}

	// Back button at the bottom
	backButton := widget.NewButton("Back", func() {
		w.SetContent(createDayButtons(train, w))
	})

	// Start button just above Back button
	startButton := widget.NewButton("Start", func() {
		println("Start clicked")
	})

	return container.NewVBox(
		container.NewGridWithColumns(1, exerciseButtons...), // Exercises
		layout.NewSpacer(), // Spacer
		startButton,        // Start button
		backButton,         // Back button
	)
}

func createImageFromURL(imageURL string) fyne.CanvasObject {
	res, err := fyne.LoadResourceFromURLString(imageURL)
	if err != nil {
		return nil
	}
	img := canvas.NewImageFromResource(res)
	img.SetMinSize(fyne.NewSize(300, 300))

	return img
}

func front(train []TrainingDay) {
	app := app.New()
	w := app.NewWindow("Training Days")

	w.SetContent(createDayButtons(train, w))
	w.Resize(fyne.NewSize(400, 600))
	w.ShowAndRun()
}

func fillTrainingDay(table TableContent) []TrainingDay {
	numTrainingDays := int(table.ColsNum)/4 + 1
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

func fillRep(table TableContent, train []TrainingDay) {
	day := 0
	for i := 2; i < int(table.ColsNum); i += 4 {
		if day >= len(train) {
			break
		}

		for j := 0; j < int(table.RowsNum); j++ {
			sets, min, max, err := parseExerciseFormat(table.TableMatrix[j][i])
			if err != nil {
				continue // TODO: Error Handling
			}
			exercise := Excercise{
				Sets: sets,
				Min:  min,
				Max:  max,
			}
			train[day].Rep = append(train[day].Rep, exercise)
		}
		day++
	}
}

func fillCurrentWeight(table TableContent, train []TrainingDay) {
	day := 0
	for i := 1; i < int(table.ColsNum); i += 4 {
		if day >= len(train) {
			break
		}
		for j := 0; j < int(table.RowsNum); j++ {
			excercieWeight, err := strconv.ParseFloat(table.TableMatrix[j][i], 32)
			if err != nil {
				continue // TODO: Error Handling
			}
			train[day].CurrentWeight = append(train[day].CurrentWeight, float32(excercieWeight))
		}
		day++
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

func main() {
	var table TableContent
	table.TableMatrix = make([][]string, 0)

	training := excelImportData(table)

	front(training)
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
