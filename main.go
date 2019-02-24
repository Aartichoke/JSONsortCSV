package main

import (
	"bufio"
	"encoding/csv"
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"strconv"
)

// Author: Daniel Foland
// Context: Programming Exercise for Radar
// Usage: Takes input csv/json, performs a specified sort, writes file to opposite type

// define data structure for input json/csv
type etl struct {
	id          int64
	name        string
	discovered  string
	description string
	status      string
}

var etlSlice []etl

func main() {
	// parse command line and check for valid inputs
	parseCommandLine()
	// read data from json/data and store in struct
	readInput()
	// sort previously read data based on fields in command line
	runSort()
	// write data to output file
	writeData()
}

// TODO add function descriptions
func parseCommandLine() {
	// define command line arguments
	inputPathPtr := flag.String("input", "", "Path to input .csv or .json")
	sortFieldPtr := flag.String("sortfield", "discovered", "Field to sort input data by")
	sortDirectionPtr := flag.String("sortdirection", "ascending", "Path to input CSV or JSON.")
	columnsPtr := flag.String("columns", "", "Columns/fields to use in output")
	flag.Parse()
	// check if path is valid
	if *inputPathPtr == "" {
		log.Fatal("Error: -input is not defined. Please use a valid .json or .csv file path")
	}
	if _, err := os.Stat(*inputPathPtr); os.IsNotExist(err) {
		log.Fatal("Error: -input  '", *inputPathPtr, "' is not a valid path.")
	}
	// check if input ends in .csv or .json
	if !checkFileExtension(*inputPathPtr) {
		log.Fatal("Error: -input extension '", filepath.Ext(*inputPathPtr), "' is not .json or .csv")
	}
	// check sort direction input
	if !checkSortDirection(*sortDirectionPtr) {
		log.Fatal("Error: -sortdirection input '", *sortDirectionPtr, "' is invalid, it must be either 'ascending' or 'descending'.")
	}
	// check sort field input
	if !checkSortField(*sortFieldPtr) {
		log.Fatal("Error: -sortfield input '", *sortFieldPtr, "' is invalid, it must be either 'status' or 'discovered'.")
	}
	// check sort field input
	if !checkColumnsField(*columnsPtr) {
		log.Fatal("Error:.", *columnsPtr)
	}
}

func readInput() {
	file := flag.Lookup("input").Value.String()
	if filepath.Ext(file) == ".json" {
		fmt.Println("Opening JSON file: ", file)
		jsonFile, err := os.Open(file)
		if err != nil {
			fmt.Println(err)
		}
		bytes, _ := ioutil.ReadAll(jsonFile)
		err = json.Unmarshal([]byte(bytes), &etlSlice)
		if err != nil {
			log.Fatal(err)
		}
	} else if filepath.Ext(file) == ".csv" {
		fmt.Println("Opening CSV file: ", file)
		f, _ := os.Open(file)
		r := csv.NewReader(bufio.NewReader(f))
		for {
			line, err := r.Read()
			if err == io.EOF {
				break
			}
			if err != nil {
				log.Fatal(err)
			}
			i64, err := strconv.ParseInt(line[0], 10, 64)
			if err != nil {
				log.Fatal(err)
			}
			etlSlice = append(etlSlice, etl{
				id:          i64,
				name:        line[1],
				discovered:  line[2],
				description: line[3],
				status:      line[4],
			})
		}
	}
	// TODO verify etlSlice is populated here
	fmt.Printf("Number of records processed: %d \n", len(etlSlice))
}

/*  runSort() performs a sort on the previously read data
 *  Inputs -
 *  Outputs - none
 */
func runSort() {
	// TODO determine ascending or descending

	// determine date or status

	/*	// parse int and date fields, fail on bad data format
		t, err := time.Parse("2006-01-02", lines["discovered"])
		if err != nil {
			log.Fatal(err)
		}
	*/
	// assign status field weights 0-2 ?

	//perform sort
}

func writeData() {
	// TODO detect if .json or .csv
	// write out accordingly
	// verify that file exists?

	fmt.Println("Output file: ")
}

/***** Helper functions below *****/
func checkStatus(status string) bool {
	switch status {
	case
		"New", "In progress", "Done":
		return true
	default:
		return false
	}
}

func checkSortDirection(direction string) bool {
	switch direction {
	case
		"ascending", "descending":
		return true
	default:
		return false
	}
}

func checkSortField(field string) bool {
	switch field {
	case
		"status", "discovered":
		return true
	default:
		return false
	}
}

func checkFileExtension(path string) bool {
	switch filepath.Ext(path) {
	case
		".json", ".csv":
		return true
	default:
		return false
	}
}

func checkColumnsField(cols string) bool {
	// TODO need to check cols here
	return true
}
