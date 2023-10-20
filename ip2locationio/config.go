package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
)

type Config struct {
	APIKey string `json:"api_key"`
}

var config Config

func LoadConfig() {
	path, err := ConfigPath()

	if err != nil {
		// fmt.Println(err)
	} else {
		file, err := os.Open(path)

		if err != nil {
			// fmt.Println(err)
		} else {
			byteValue, _ := ioutil.ReadAll(file)
			json.Unmarshal(byteValue, &config)
			return
		}
	}
	// default config
	config.APIKey = ""
}

func UpdateAPIKey(apiKey string) {
	config.APIKey = apiKey
	SaveConfig()
}

func SaveConfig() {
	path, err := ConfigPath()

	if err != nil {
		// fmt.Println(err)
	} else {
		file, err := os.OpenFile(path, os.O_RDWR|os.O_CREATE, 0700)
		if err != nil {
			// fmt.Println(err)
		} else {
			byteValue, err := json.Marshal(&config)

			_, err = file.Write(byteValue)

			if err != nil {
				fmt.Println(err)
			}
		}

		if err := file.Close(); err != nil {
			// fmt.Println(err)
		}
	}
}

func ConfigPath() (string, error) {
	cd, err := os.UserConfigDir()
	if err != nil {
		return "", err
	}

	cd2 := filepath.Join(cd, "ip2locationio")
	if err := os.MkdirAll(cd2, 0700); err != nil {
		return "", err
	}

	return filepath.Join(cd2, "ip2locationio-config.json"), nil
}
