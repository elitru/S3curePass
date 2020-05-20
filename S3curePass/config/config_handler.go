package config

import (
	"S3curePass/errors"
	"S3curePass/models"
	"encoding/json"
	"io/ioutil"
)

//Path to config file
const CONFIG_FILE = "./config/config.json"

//holder struct for config
var config models.Config = models.Config{}

func GetConfig() models.Config {
	//check if config has already been loaded
	if config == (models.Config{}) {
		//load json string from config file
		jsonString, err := ioutil.ReadFile(CONFIG_FILE)

		//checkf for errors while reading file
		errors.CheckFatal(err)
		//parse file content string to config struct
		json.Unmarshal([]byte(jsonString), &config)
	}

	return config
}
