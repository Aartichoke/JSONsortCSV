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
	"sort"
	"strconv"
	"strings"
	"time"
)

// Author: Daniel Foland
// Context: Programming Exercise for Radar
// Usage: Takes input csv/json, performs a specified sort, writes file to opposite type

// define data structure for input json/csv
type etl struct {
	// Remember: json requires capital field names (eg Id instead of id)
	Id          int
	Name        string
	Discovered  string
	Description string
	Status      string
}

// slice to hold input data
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

// parse command line inputs
func parseCommandLine() {
	// define command line arguments
	inputPathPtr := flag.String("input", "", "Path to input .csv or .json")
	sortFieldPtr := flag.String("sortfield", "discovered", "Field to sort input data by")
	sortDirectionPtr := flag.String("sortdirection", "ascending", "Path to input CSV or JSON.")
	columnsPtr := flag.String("columns", "Id,Name,Discovered,Description,Status", "Columns/fields to use in output. Capitalize each.")
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
	// check column field input
	columns := strings.Split(*columnsPtr, ",")
	for _, col := range columns {
		if !checkColumnsField(col) {
			log.Fatal("Error: Invalid column: ", col)
		}
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
			i, err := strconv.Atoi(line[0])
			if err != nil {
				log.Fatal(err)
			}
			etlSlice = append(etlSlice, etl{
				Id:          i,
				Name:        line[1],
				Discovered:  line[2],
				Description: line[3],
				Status:      line[4],
			})
		}
	}
	if len(etlSlice) == 0 {
		log.Fatal("Error: Data could not be stored from input file.")
	}
	fmt.Printf("Number of records processed: %d \n", len(etlSlice))
}

/*  runSort() performs a sort on the previously read data
 *  Inputs -
 *  Outputs - none
 */
func runSort() {
	if flag.Lookup("sortfield").Value.String() == "status" {
		// sort by first letter of status
		sort.Slice(etlSlice, func(i, j int) bool {
			if flag.Lookup("sortdirection").Value.String() == "ascending" {
				return etlSlice[i].Status[0] > etlSlice[j].Status[0]
			}
			return etlSlice[i].Status[0] < etlSlice[j].Status[0]
		})
	} else {
		// sort by date
		sort.Slice(etlSlice, func(i, j int) bool {
			t1, err := time.Parse("2006-01-02", etlSlice[i].Discovered)
			if err != nil {
				log.Fatal(err)
			}
			t2, err := time.Parse("2006-01-02", etlSlice[j].Discovered)
			if err != nil {
				log.Fatal(err)
			}
			if flag.Lookup("sortdirection").Value.String() == "ascending" {
				return t1.Sub(t2).Seconds() < 0
			}
			return t1.Sub(t2).Seconds() > 0
		})
	}
}

/*  writeData() write etl struct to opposite file type
 *  Inputs - global populated etlslice
 *  Outputs - writes file 'output'
 */
func writeData() {
	if filepath.Ext(flag.Lookup("input").Value.String()) == ".csv" {
		file, _ := json.MarshalIndent(etlSlice, "", " ")
		_ = ioutil.WriteFile("output.json", file, 0644)
		fmt.Println("Output file: output.json")
	} else {
		file, err := os.Create("output.csv")
		if err != nil {
			log.Fatal(err)
		}
		defer file.Close()
		writer := csv.NewWriter(file)
		defer writer.Flush()
		for _, thisetl := range etlSlice {
			columns := strings.Split(flag.Lookup("columns").Value.String(), ",")
			var output []string
			for _, col := range columns {
				// there's no dynamic struct access, so we have to do this manually because Id is not a string
				// Either assign field to an int, or handle each field individually
				switch col {
				case
					"Id":
					output = append(output, strconv.Itoa(thisetl.Id))
				case
					"Name":
					output = append(output, thisetl.Name)
				case
					"Discovered":
					output = append(output, thisetl.Discovered)
				case
					"Description":
					output = append(output, thisetl.Description)
				case
					"Status":
					output = append(output, thisetl.Status)
				}
			}
			if err := writer.Write(output); err != nil {
				log.Fatal(err)
			}
		}
		fmt.Println("Output file: output.csv")
	}
}

/***** Helper functions below *****/
/*  checkSortDirection() checks to if sort direction is valid
 *  Inputs - string to compare
 *  Outputs - true when string is valid, false otherwise
 */
func checkSortDirection(direction string) bool {
	switch direction {
	case
		"ascending", "descending":
		return true
	default:
		return false
	}
}

/*  checkSortField() checks to if sort field is valid
 *  Inputs - string to compare
 *  Outputs - true when string is valid, false otherwise
 */
func checkSortField(field string) bool {
	switch field {
	case
		"status", "discovered":
		return true
	default:
		return false
	}
}

/*  checkFileExtension() checks to if input file extension is valid
 *  Inputs - string to check for json or csv
 *  Outputs - true when csv/json matches
 */
func checkFileExtension(path string) bool {
	switch filepath.Ext(path) {
	case
		".json", ".csv":
		return true
	default:
		return false
	}
}

/*  checkColumnsField() checks to see if input string is a valid column
 *  Inputs - string to compare to valid columns
 *  Outputs - true when string matches a valid column, false there is no match
 */
func checkColumnsField(cols string) bool {
	switch cols {
	case
		"Id", "Name", "Discovered", "Description", "Status":
		return true
	default:
		return false
	}
}
