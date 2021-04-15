package main

import (
	"bufio"
	"flag"
	"fmt"
	"os"
	"strings"

	crud "github.com/thalaivar-subu/paylaterservice/crud"
	database "github.com/thalaivar-subu/paylaterservice/database"
)

func main() {
	flag.Set("alsologtostderr", "true")
	flag.Set("log_dir", "./log/")
	flag.Parse()
	db := database.ConnectMysql()
	defer db.Close()

	scanner := bufio.NewScanner(os.Stdin)
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
			text := scanner.Text()
			if len(text) <= 0 || text == "exit" {
				fmt.Println("Application Exiting ", text)
				// exit if user entered an empty string or exit
				break
			}

			inputArgs := strings.Split(text, " ")
			// flag, errorMsg := validation.Validate(inputArgs)
			// if flag == false {
			// 	fmt.Println("Input Validation Failed ", errorMsg)
			// 	break
			// } else {
			// 	fmt.Println("Validation Success")
			// }
			if inputArgs[0] == "new" {
				if inputArgs[1] == "user" {
					flag, result, errorMsg := crud.CreateUser(inputArgs[2], inputArgs[3], inputArgs[4], db)
					if !flag {
						fmt.Println(errorMsg)
					} else {
						fmt.Print(result)
					}
				} else if inputArgs[1] == "merchant" {
					flag, result, errorMsg := crud.CreateMerchant(inputArgs[2], inputArgs[3], inputArgs[4], db)
					if !flag {
						fmt.Println(errorMsg)
					} else {
						fmt.Print(result)
					}
				}
			}

		}
	}

	initApplication()
}
