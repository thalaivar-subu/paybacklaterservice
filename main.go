package main

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/thalaivar-subu/paylaterservice/cmd/merchant"
	"github.com/thalaivar-subu/paylaterservice/cmd/txn"
	"github.com/thalaivar-subu/paylaterservice/cmd/user"
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
		defer func() {
			if err := recover(); err != nil {
				fmt.Println(err)
			}
		}()
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
			inValid := false
			if inputArgs[0] == "new" {
				if inputArgs[1] == "user" {
					result, errorMsg := user.CreateUser(inputArgs[2], inputArgs[3], inputArgs[4], db)
					if errorMsg != nil {
						output = "rejected!" + " (" + errorMsg.Error() + ")"
					} else {
						output = result
					}
				} else if inputArgs[1] == "merchant" {
					result, errorMsg := merchant.CreateMerchant(inputArgs[2], inputArgs[3], inputArgs[4], db)
					if errorMsg != nil {
						output = "rejected!" + " (" + errorMsg.Error() + ")"
					} else {
						output = result
					}
				} else if inputArgs[1] == "txn" {
					result, errorMsg := txn.CreateTxn(inputArgs[2], inputArgs[3], inputArgs[4], db)
					if errorMsg != nil {
						output = "rejected!" + " (" + errorMsg.Error() + ")"
					} else {
						output = result
					}
				} else {
					inValid = true
				}
			} else if inputArgs[0] == "report" {
				if inputArgs[1] == "users-at-credit-limit" {
					result, errorMsg := user.GetUsersAtCredLimit(db)
					if errorMsg != nil {
						output = "rejected!" + " (" + errorMsg.Error() + ")"
					} else {
						output = result
					}
				} else if inputArgs[1] == "total-dues" {
					result, errorMsg := user.GetTotalDues(db)
					if errorMsg != nil {
						output = "rejected!" + " (" + errorMsg.Error() + ")"
					} else {
						output = result
					}
				} else if inputArgs[1] == "discount" {
					result, errorMsg := merchant.GetDiscount(inputArgs[2], db)
					if errorMsg != nil {
						output = "rejected!" + " (" + errorMsg.Error() + ")"
					} else {
						output = result
					}
				} else if inputArgs[1] == "dues" {
					result, errorMsg := user.GetUserDues(inputArgs[2], db)
					if errorMsg != nil {
						output = "rejected!" + " (" + errorMsg.Error() + ")"
					} else {
						output = result
					}
				} else {
					inValid = true
				}
			} else if inputArgs[0] == "payback" {
				result, errorMsg := user.PayBack(inputArgs[1], inputArgs[2], db)
				if errorMsg != nil {
					output = "rejected!" + " (" + errorMsg.Error() + ")"
				} else {
					output = result
				}
			} else if inputArgs[0] == "update" {
				if inputArgs[1] == "merchant" {
					result, errorMsg := merchant.UpdateMerchantDiscount(inputArgs[2], inputArgs[3], db)
					if errorMsg != nil {
						output = "rejected!" + " (" + errorMsg.Error() + ")"
					} else {
						output = result
					}
				} else {
					inValid = true
				}
			} else {
				inValid = true
			}
			if inValid {
				output = "Not a Valid Command"
			}
			fmt.Fprintf(out, "%s", output)
			fmt.Println()
		}
	}
	initApplication()

	// handle error
	if scanner.Err() != nil {
		fmt.Println("Error: ", scanner.Err())
	}
}
