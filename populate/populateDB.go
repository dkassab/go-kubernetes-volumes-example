package main

import (
	"encoding/csv"
	"flag"
	"log"
	"os"
	"path/filepath"

	"github.com/garyburd/redigo/redis"
)

// Employee information
// model that is output from the training.
type Employee struct {
	Name string `json:"name"`
}

func main() {

	// Declare the input directory flags.
	inDirPtr := flag.String("inDir", "", "The directory containing the data.")

	// Parse the command line flags.
	flag.Parse()

	// Open the Employee dataset file.
	f, err := os.Open(filepath.Join(*inDirPtr, "data.csv"))
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()

	// Create a new CSV reader reading from the opened file.
	reader := csv.NewReader(f)

	// Read in all of the CSV records
	reader.FieldsPerRecord = 1
	employeeData, err := reader.ReadAll()
	if err != nil {
		log.Fatal(err)
	}
	c, err := redis.Dial("tcp", "db:6379")
	if err != nil {
		panic(err)
	}
	defer c.Close()

	// Loop of records in the CSV, printing the employee data.
	for i, record := range employeeData {

		// Skip the header.
		if i == 0 {
			continue
		}
		//set
		c.Do("LPUSH", "employeelist", record)
		if err != nil {
			log.Fatal(err)
		}

	}
}
