package main

import (
	"fmt"
	"os"
	"strconv"
	"time"
)

func usage() {
	fmt.Println("Usage: HomeAccountingDB2 db_files_location\n  test_json date")
}

func main() {
	l := len(os.Args)
	if l < 3 || l > 5 {
		usage()
		return
	}
	switch os.Args[2] {
	case "test_json":
		if l != 4 {
			usage()
		} else {
			test(os.Args[1], JsonDBConfiguration{}, os.Args[3])
		}
	}
}

func test(dataFolderPath string, configuration DBConfiguration, dateString string) {
	date, err := strconv.Atoi(dateString)
	if err != nil {
		panic(err)
	}
	fmt.Println("Reading DB files...")
	start := time.Now()
	db, err := LoadDB(2012, 6, dataFolderPath, configuration, 1000)
	fmt.Printf("%v elapsed.\n", time.Since(start))
	if err != nil {
		panic(err)
	}
	fmt.Println("Calculating finance totals...")
	start = time.Now()
	err = db.BuildTotals(0)
	if err != nil {
		fmt.Println("BuildTotals error: " + err.Error())
		return
	}
	fmt.Printf("%v elapsed.\n", time.Since(start))
	db.PrintChanges(date)
}
