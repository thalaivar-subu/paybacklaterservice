package main

import (
	"bytes"
	"fmt"
	"testing"
)

func TestCreateUser(t *testing.T) {
	// Create Users
	testData := map[string]string{
		"new user user1 u1@users.com 300": "user1(300)",
		"new user user2 u2@users.com 400": "user2(400)",
		"new user user3 u3@users.com 500": "user3(500)",
	}
	for k, v := range testData {
		t.Run(k, func(t *testing.T) {
			var command, result bytes.Buffer
			t.Log(&command, k)
			fmt.Fprintf(&command, k)
			run(&command, &result)
			if result.String() != v {
				t.Error("Result '" + result.String() + "' not equal to expectedOutput '" + v + "'")
			} else {
				t.Log("Pass -> ", result.String())
			}
		})
	}

}

func TestMerchantUser(t *testing.T) {
	testData := map[string]string{
		"new merchant m1 m1@merchants.com 0.5%%":  "m1(0.5%)",
		"new merchant m2 m2@merchants.com 1.5%%":  "m2(1.5%)",
		"new merchant m3 m3@merchants.com 6.25%%": "m3(6.25%)",
	}
	for k, v := range testData {
		t.Run(k, func(t *testing.T) {
			var command, result bytes.Buffer
			t.Log(&command, k)
			fmt.Fprintf(&command, k)
			run(&command, &result)
			if result.String() != v {
				t.Error("Result '" + result.String() + "' not equal to expectedOutput '" + v + "'")
			} else {
				t.Log("Pass -> ", result.String())
			}
		})
	}
}
func TestTxn(t *testing.T) {
	testData := map[string]string{
		"new txn user2 m1 500":         "rejected! (reason: credit limit)",
		"new txn user1 m2 300":         "success!",
		"new txn user1 m3 10":          "rejected! (reason: credit limit)",
		"report users-at-credit-limit": "user1",
		"new txn user3 m3 200":         "success!",
		"new txn user3 m3 300":         "success!",
	}
	for k, v := range testData {
		t.Run(k, func(t *testing.T) {
			var command, result bytes.Buffer
			t.Log(&command, k)
			fmt.Fprintf(&command, k)
			run(&command, &result)
			if result.String() != v {
				t.Error("Result '" + result.String() + "' not equal to expectedOutput '" + v + "'")
			} else {
				t.Log("Pass -> ", result.String())
			}
		})
	}
}

// func TestReport(t *testing.T) {
// 	testData := map[string]string{
// 		"report users-at-credit-limit": "user1\nuser3",
// 		// "report users-at-credit-limit": "user1\nuser3",
// 		"report discount m3": "6.25",
// 		"report total-dues":  "user1: 300\nuser3: 100\ntotal: 400",
// 	}
// 	for k, v := range testData {
// 		t.Run(k, func(t *testing.T) {
// 			var command, result bytes.Buffer
// 			t.Log(&command, k)
// 			fmt.Fprintf(&command, k)
// 			run(&command, &result)
// 			if result.String() != v {
// 				t.Error("Result '" + result.String() + "' not equal to expectedOutput '" + v + "'")
// 			} else {
// 				t.Log("Pass -> ", result.String())
// 			}
// 		})
// 	}
// }
