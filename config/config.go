package config

import "github.com/thalaivar-subu/paylaterservice/constants"

var Config = getConfig()

func getConfig() map[string]string {
	config := map[string]map[string]string{}
	config["development"] = map[string]string{
		"mysql": "root:123456@tcp(127.0.0.1:3306)/paylaterservice?charset=utf8&parseTime=True"}
	return config[constants.ENV]
}
