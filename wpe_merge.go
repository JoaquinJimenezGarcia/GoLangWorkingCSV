package main

import (
	"encoding/csv"
	"encoding/json"
	"io"
	"log"
	"net/http"
	"os"
)

type CurrentData struct {
	AccountID   string
	AccountName string
	FirstName   string
	CreatedOn   string
}

type UpdatedData struct {
	AccountID int64  `json:"account_id"`
	Status    string `json:"status"`
	CreatedOn string `json:"created_on"`
}

func main() {
	inputFileName := os.Args[1]
	outputFileName := os.Args[2]

	openCSV(inputFileName, outputFileName)
}

func openCSV(inputFileName string, outputFileName string) {
	// Open the file
	csvfile, err := os.Open(inputFileName)
	if err != nil {
		log.Fatalln("Couldn't open the csv file", err)
	}

	// Parse the file
	r := csv.NewReader(csvfile)

	// Create the new CSV file
	file, err := os.Create(outputFileName)
	if err != nil {
		log.Fatalln("Couldn't create the csv file", err)
	}

	defer file.Close()

	writer := csv.NewWriter(file)
	defer writer.Flush()

	// Iterate through the records
	for {
		// Read each record from csv
		record, err := r.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			log.Fatal(err)
		}

		// Set an object with the current data
		currentData := CurrentData{
			AccountID:   record[0],
			AccountName: record[1],
			FirstName:   record[2],
			CreatedOn:   record[3],
		}

		// Check for the header line to add the new two fields
		if currentData.AccountID == "Account ID" {
			dataHeader := []string{currentData.AccountID, currentData.AccountName, currentData.FirstName, currentData.CreatedOn, "Status", "Status Set On"}

			// Write the first line as header
			err := writer.Write(dataHeader)
			if err != nil {
				log.Fatalln("Couldn't add the csv header", err)
			}

			continue
		}

		url := "http://interview.wpengine.io/v1/accounts/" + currentData.AccountID
		data := new(UpdatedData)

		// Make a request to the url to have the updated data
		err = makeRequest(url, data)
		if err != nil {
			log.Fatalln("Couldn't make the request", err)
		}

		// Set the new line for the file
		updatedData := []string{currentData.AccountID, currentData.AccountName, currentData.FirstName, currentData.CreatedOn, data.Status, data.CreatedOn}

		// Write the new line with old and updated data
		err = writer.Write(updatedData)
		if err != nil {
			log.Fatalln("Couldn't add the csv line with updated data", err)
		}

	}
}

// Function to make a request to the url. It must receive the url and the object to formate
func makeRequest(url string, data interface{}) error {
	resp, err := http.Get(url)
	if err != nil {
		return err
	}

	defer resp.Body.Close()

	// Return the data formated
	return (json.NewDecoder(resp.Body).Decode(data))
}
