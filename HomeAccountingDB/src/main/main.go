package main

import (
	"TimeSeriesData/core"
	"TimeSeriesData/crypto"
	"TimeSeriesData/network"
	"fmt"
	"os"
	"strconv"
	"time"
)

func usage() {
	fmt.Println("Usage: HomeAccountingDB2 config_file_name\n  test_json date\n  test date aes_key_file\n  migrate source_folder aes_key_file\n server")
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
			testJson(s, os.Args[3])
		}
	case "test":
		if l != 5 {
			usage()
		} else {
			testBinary(s, os.Args[3], os.Args[4])
		}
	case "migrate":
		if l != 5 {
			usage()
		} else {
			migrate(s, os.Args[3], os.Args[4])
		}
	case "server":
		if l != 3 {
			usage()
		} else {
			startServer(s)
		}
	default:
		usage()
	}
}

func startServer(s settings) {
	userData := tcpServerData{s: s}
	server, err := network.NewTcpServer[tcpServerData](s.ServerPort, s.Key, "HomeAccountingDB", &userData,
		func(request []byte, userData *tcpServerData) ([]byte, error) {
			return userData.handle(request)
		})
	if err != nil {
		panic(err)
	}
	err = server.Start()
	if err != nil {
		panic(err)
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

	fmt.Println("Building hints...")
	start = time.Now()
	err = db.buildHints()
	if err != nil {
		panic("BuildHints error: " + err.Error())
	}
	fmt.Printf("%v elapsed.\n", time.Since(start))

	return db
}

func initDatabase(s settings, dbConfiguration dBConfiguration) (*dB, error) {
	fmt.Println("Initializing database...")
	start := time.Now()
	db, err := initDB(s, dbConfiguration)
	fmt.Printf("%v elapsed.\n", time.Since(start))
	return db, err
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

func migrate(s settings, sourceFolder, aesKeyFile string) {
	destFolder := s.DataFolderPath
	s.DataFolderPath = sourceFolder
	db := buildDB(s, jsonDBConfiguration{})
	err := db.saveTo(destFolder, buildBinaryDbConfiguration(aesKeyFile))
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

func testBinary(s settings, dateString string, aesKeyFile string) {
	date, err := strconv.Atoi(dateString)
	if err != nil {
		panic(err)
	}
	db, err := initDatabase(s, buildBinaryDbConfiguration(aesKeyFile))
	if err != nil {
		panic(err)
	}
	db.printChanges(date)
}
