package main

import (
	"bytes"
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"net/url"
	"strconv"
	"strings"
)

type Response struct {
	IP string `json:"IP"`
}

// The IPGeolocationError struct stores errors
// returned by the IP2Location.io API.
type IPGeolocationError struct {
	Error struct {
		ErrorCode    int    `json:"error_code"`
		ErrorMessage string `json:"error_message"`
	} `json:"error"`
}

func MyPublicIP() string {
	myUrl := "https://ip2location.io/get-ip.json"
	res, err := http.Get(myUrl)

	if err != nil {
		return ""
	}

	var response Response
	json.NewDecoder(res.Body).Decode(&response)
	return response.IP
}

func PrettyString(str string) (string, error) {
	var prettyJSON bytes.Buffer
	if err := json.Indent(&prettyJSON, []byte(str), "", "    "); err != nil {
		return "", err
	}
	return prettyJSON.String(), nil
}

// LookUpJSON will return a JSON based on the queried IP address
func LookUpJSON(ip string, lang string) (string, error) {
	var res string
	var ex IPGeolocationError

	myUrl := "https://api.ip2location.io?ip=" + url.QueryEscape(ip) + "&source=sdk-cli-iplio&source_version=" + version

	if strings.TrimSpace(apiKey) != "" {
		myUrl = myUrl + "&key=" + url.QueryEscape(apiKey) + "&lang=" + url.QueryEscape(lang)
	}

	resp, err := http.Get(myUrl)

	if err != nil {
		return res, err
	}

	defer resp.Body.Close()

	if resp.StatusCode == http.StatusOK {
		bodyBytes, err := io.ReadAll(resp.Body)

		if err != nil {
			return res, err
		}

		res = string(bodyBytes[:])

		return res, nil
	} else if resp.StatusCode == http.StatusBadRequest || resp.StatusCode == http.StatusUnauthorized {
		bodyBytes, err := io.ReadAll(resp.Body)

		if err != nil {
			return res, err
		}

		bodyStr := string(bodyBytes[:])
		if strings.Contains(bodyStr, "error_message") {
			err = json.Unmarshal(bodyBytes, &ex)

			if err != nil {
				return res, err
			}
			return res, errors.New("Error: " + ex.Error.ErrorMessage)
		}
	}

	return res, errors.New("Error HTTP " + strconv.Itoa(int(resp.StatusCode)))
}

// LookUp will return all geolocation fields based on the queried IP address inside a map
func LookUpMap(ip string, lang string) (map[string]interface{}, error) {
	var res map[string]interface{}
	var ex IPGeolocationError

	myUrl := "https://api.ip2location.io?ip=" + url.QueryEscape(ip) + "&source=sdk-cli-iplio&source_version=" + version

	if strings.TrimSpace(apiKey) != "" {
		myUrl = myUrl + "&key=" + url.QueryEscape(apiKey) + "&lang=" + url.QueryEscape(lang)
	}

	resp, err := http.Get(myUrl)

	if err != nil {
		return res, err
	}

	defer resp.Body.Close()

	if resp.StatusCode == http.StatusOK {
		bodyBytes, err := io.ReadAll(resp.Body)

		if err != nil {
			return res, err
		}

		err = json.Unmarshal(bodyBytes, &res)

		if err != nil {
			return res, err
		}

		return res, nil
	} else if resp.StatusCode == http.StatusBadRequest || resp.StatusCode == http.StatusUnauthorized {
		bodyBytes, err := io.ReadAll(resp.Body)

		if err != nil {
			return res, err
		}

		bodyStr := string(bodyBytes[:])
		if strings.Contains(bodyStr, "error_message") {
			err = json.Unmarshal(bodyBytes, &ex)

			if err != nil {
				return res, err
			}
			return res, errors.New("Error: " + ex.Error.ErrorMessage)
		}
	}

	return res, errors.New("Error HTTP " + strconv.Itoa(int(resp.StatusCode)))
}
