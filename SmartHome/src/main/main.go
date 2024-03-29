package main

import (
	"TimeSeriesData/core"
	"fmt"
	"os"
	"strconv"
	"time"
)

func usage() {
	fmt.Println("Usage: SmartHome config_file_name\n  test_json date\n  test date\n  migrate source_folder")
}

func main() {
	l := len(os.Args)
	if l < 3 || l > 4 {
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
			testJson(s, os.Args[3])
		}
	case "test":
		if l != 4 {
			usage()
		} else {
			testBinary(s, os.Args[3])
		}
	case "migrate":
		if l != 4 {
			usage()
		} else {
			migrate(s, os.Args[3])
		}
	default:
		usage()
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
	return db
}

func initDatabase(s settings, dbConfiguration dBConfiguration) *dB {
	fmt.Println("Initializing database...")
	start := time.Now()
	db, err := initDB(s, dbConfiguration)
	fmt.Printf("%v elapsed.\n", time.Since(start))
	if err != nil {
		panic(err)
	}
	return db
}

func migrate(s settings, sourceFolder string) {
	destFolder := s.DataFolderPath
	s.DataFolderPath = sourceFolder
	db := buildDB(s, jsonDBConfiguration{})
	err := db.saveTo(destFolder, binaryDBConfiguration{})
	if err != nil {
		panic(err)
	}
}

func testJson(s settings, dateString string) {
	date, err := strconv.Atoi(dateString)
	if err != nil {
		panic(err)
	}
	db := buildDB(s, jsonDBConfiguration{})
	db.printStats(date)
}

func testBinary(s settings, dateString string) {
	date, err := strconv.Atoi(dateString)
	if err != nil {
		panic(err)
	}
	db := initDatabase(s, binaryDBConfiguration{})
	db.printStats(date)
}
