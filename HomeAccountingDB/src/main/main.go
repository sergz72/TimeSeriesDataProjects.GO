package main

import (
	"TimeSeriesData/core"
	"TimeSeriesData/crypto"
	"fmt"
	"os"
	"strconv"
	"time"
)

func usage() {
	fmt.Println("Usage: HomeAccountingDB2 config_file_name\n  test_json date\n  test date\n  migrate source_folder")
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

func buildBinaryDbConfiguration(aesKeyFileName string) dBConfiguration {
	key, err := crypto.LoadAesKey(aesKeyFileName)
	if err != nil {
		panic(err)
	}
	aes, err := crypto.NewAesGcm(key)
	if err != nil {
		panic(err)
	}
	return newBinaryDBConfiguration(aes)
}

func migrate(s settings, sourceFolder string) {
	destFolder := s.DataFolderPath
	s.DataFolderPath = sourceFolder
	db := buildDB(s, jsonDBConfiguration{})
	err := db.saveTo(destFolder, buildBinaryDbConfiguration(s.Key))
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
	db.printChanges(date)
}

func testBinary(s settings, dateString string) {
	date, err := strconv.Atoi(dateString)
	if err != nil {
		panic(err)
	}
	db := initDatabase(s, buildBinaryDbConfiguration(s.Key))
	db.printChanges(date)
}
