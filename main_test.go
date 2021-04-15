package main

import (
	"bytes"
	"testing"
)

func TestCli(t *testing.T) {
	var command, result bytes.Buffer
	testData := map[string]string{
		"new user user1 u1@users.com 300":         "user1(300)",
		"new user user2 u2@users.com 400":         "user2(400)",
		"new user user3 u3@users.com 500":         "user3(500)",
		"new merchant m1 m1@merchants.com 0.5%%":  "m1(0.5%)",
		"new merchant m2 m2@merchants.com 1.5%%":  "m2(1.5%)",
		"new merchant m3 m3@merchants.com 1.25%%": "m3(1.25%)",
		"new txn user2 m1 500":                    "rejected! (reason: credit limit)",
		"new txn user1 m2 300":                    "success!",
		"new txn user1 m3 10":                     "rejected! (reason: credit limit)",
		// "report users-at-credit-limit":           "user1",
		"new txn user3 m3 200": "success!",
		"new txn user3 m3 300": "success!",
		// "report users-at-credit-limit":           "user1\nuser3",
		// "report discount m3": "6.25",
		// "payback user3 400":  "user3(dues: 100)",
		// "report total-dues":  "user1: 300\nuser3: 100\ntotal: 400",
	}
	for k, _ := range testData {
		t.Log(&command, k)
		run(&command, &result)
		t.Log(result.String())
	}

}
