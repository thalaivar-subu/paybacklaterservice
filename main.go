package main

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/thalaivar-subu/paylaterservice/crud"
	"github.com/thalaivar-subu/paylaterservice/database"
)

func main() {
	var f *os.File
	f = os.Stdin
	defer f.Close()
	run(os.Stdin, f)
}

func run(in io.Reader, out io.Writer) {
	database.ConnectMysql()
	db := database.Db
	defer db.Close()

	scanner := bufio.NewScanner(in)
	initApplication := func() {
		// defer func() {
		// 	if err := recover(); err != nil {
		// 		fmt.Println(err)
		// 	}
		// }()
		for {
			fmt.Print(">")
			// reads user input until \n by default
			scanner.Scan()
			// Holds the string that was scanned
			text := strings.TrimSpace(scanner.Text())
			if len(text) <= 0 || text == "exit" {
				fmt.Println("Application Exiting ", text)
				// exit if user entered an empty string or exit
				break
			}

			inputArgs := strings.Split(text, " ")
			output := ""
			if inputArgs[0] == "new" {
				if inputArgs[1] == "user" {
					flag, result, errorMsg := crud.CreateUser(inputArgs[2], inputArgs[3], inputArgs[4], db)
					if !flag {
						output = errorMsg.Error()
					} else {
						output = result
					}
				} else if inputArgs[1] == "merchant" {
					flag, result, errorMsg := crud.CreateMerchant(inputArgs[2], inputArgs[3], inputArgs[4], db)
					if !flag {
						output = errorMsg.Error()
					} else {
						output = result
					}
				} else if inputArgs[1] == "txn" {
					flag, result, errorMsg := crud.CreateTxn(inputArgs[2], inputArgs[3], inputArgs[4], db)
					if !flag {
						output = errorMsg.Error()
					} else {
						output = result
					}
				} else {
					output = "Not Valid new Command"
					os.Exit(1)
				}
				fmt.Print(output)
			}

		}
	}
	initApplication()
	// handle error
	if scanner.Err() != nil {
		fmt.Println("Error: ", scanner.Err())
	}
}
