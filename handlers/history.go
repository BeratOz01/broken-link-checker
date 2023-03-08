package handlers

import (
	"encoding/json"
	"fmt"
	"log"
	"os"

	"github.com/jedib0t/go-pretty/v6/table"
)

var (
	DbName     = "history.json" // The name of the JSON file that will be used to store the history
	Limit      = 5              // The number of results to be displayed on a page
	CurrentDir string           // The current directory of the project
	EmptyArray = []byte("[]")   // The initial value of the database
)

type History struct {
	Website     string `json:"website"`
	Date        string `json:"date"`
	TotalLinks  int    `json:"totalLinks"`
	IsCompleted bool   `json:"isCompleted"`
}

func PrintHistory(page int) {
	// Getting the history from the database
	history, err := getHistory()
	if err != nil {
		log.Fatal(err)
	}

	// Printing the history to the console with the help of go-pretty package
	t := table.NewWriter()
	t.SetOutputMirror(os.Stdout)

	// Setting the table headers
	t.AppendHeader(table.Row{"Website", "Total Links", "Is Completed", "Date"})

	if len(*history) == 0 {
		// I do not like this solution, but I couldn't find a better one
		t.AppendRow(table.Row{"NO HISTORY", "NO HISTORY", "NO HISTORY", "NO HISTORY"}, table.RowConfig{AutoMerge: true})
	} else {
		// Calculating the start and end index of the history
		offset := (page - 1) * 10
		end := offset + Limit
		if end > len(*history) {
			end = len(*history)
		}

		fmt.Println("Showing page", page, "(", offset+1, "-", end, "of", len(*history), "results)")

		// Appending the rows
		for i := offset; i < end; i++ {
			h := (*history)[i]
			t.AppendRow(table.Row{i + 1, h.Website, h.TotalLinks, h.IsCompleted, h.Date})
		}
	}

	// Setting the table style
	t.SetStyle(table.StyleLight)

	// Rendering the table
	t.Render()
}

/*
- This helper function reads the history from the database and returns it
- If the database doesn't exist, it will throw an error
*/
func getHistory() (*[]History, error) {
	var history []History
	db, err := getOrFail()

	if err != nil {
		return nil, err
	}

	err = json.NewDecoder(db).Decode(&history)
	if err != nil {
		return nil, err
	}

	// If everything is fine, we need to return the history with no error
	return &history, nil
}

/*
- This function is used to get the history from the database
- If the database doesn't exist, it will create it and write the initial value to it
*/
func getOrFail() (*os.File, error) {
	CurrentDir, err := os.Getwd()
	if err != nil {
		log.Fatal(err)
	}

	// Declaring the local variable for db and history
	var db *os.File

	if _, err := os.Stat(CurrentDir + "/" + DbName); os.IsNotExist(err) {
		db, err = os.Create(CurrentDir + "/" + DbName)
		if err != nil {
			return nil, err
		}
		// After creating the file, we need to write the initial value to it (empty array)
		db.Write(EmptyArray)
	}

	db, err = os.Open(CurrentDir + "/" + DbName)
	if err != nil {
		return nil, err
	}

	// If everything is fine, we need to return the db with no error
	return db, nil
}
