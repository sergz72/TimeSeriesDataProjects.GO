package main

import (
	"TimeSeriesData/core"
	"fmt"
	"os"
	"strconv"
	"time"
)

func usage() {
	fmt.Println("Usage: HomeAccountingDB2 config_file_name\n  test_json date\n  migrate source_folder")
}

func main() {
	l := len(os.Args)
	if l < 3 || l > 5 {
		usage()
		return
	}
	s, err := core.LoadJson[settings](os.Args[1])
	if err != nil {
		panic(err)
	}
	switch os.Args[2] {
	case "test_json":
		if l != 4 {
			usage()
		} else {
			test(s, jsonDBConfiguration{}, os.Args[3])
		}
	case "migrate":
		if l != 4 {
			usage()
		} else {
			migrate(s, os.Args[3])
		}
	}
}

func buildDB(s settings, dbConfiguration dBConfiguration) *dB {
	fmt.Println("Reading DB files...")
	start := time.Now()
	db, err := loadDB(s, dbConfiguration)
	fmt.Printf("%v elapsed.\n", time.Since(start))
	if err != nil {
		panic(err)
	}
	fmt.Println("Calculating finance totals...")
	start = time.Now()
	err = db.buildTotals(0)
	if err != nil {
		panic("BuildTotals error: " + err.Error())
	}
	fmt.Printf("%v elapsed.\n", time.Since(start))
	return db
}

func migrate(s settings, sourceFolder string) {
	destFolder := s.DataFolderPath
	s.DataFolderPath = sourceFolder
	db := buildDB(s, jsonDBConfiguration{})
	db.saveTo(destFolder, binaryDBConfiguration{})
}

func test(s settings, dbConfiguration dBConfiguration, dateString string) {
	date, err := strconv.Atoi(dateString)
	if err != nil {
		panic(err)
	}
	db := buildDB(s, dbConfiguration)
	db.printChanges(date)
}
