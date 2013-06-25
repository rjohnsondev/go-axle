// Package goaxle provides bindings to the ApiAxle management API.
package goaxle

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"net/http"
	"time"
	"encoding/json"
	"strconv"
)

// API protocol type.
type Protocol string

// API format type.
type ApiFormat string

// Time duration for granularity of reports.
type Granularity string

// Response class.
type HitType string

const (
	API_PROTOCOL_HTTP  Protocol = "http"
	API_PROTOCOL_HTTPS Protocol = "https"

	API_FORMAT_JSON ApiFormat = "json"
	API_FORMAT_XML  ApiFormat = "xml"

	GRANULARITY_SECONDS Granularity = "second"
	GRANULARITY_MINUTES Granularity = "minute"
	GRANULARITY_HOURS   Granularity = "hour"
	GRANULARITY_DAYS    Granularity = "day"

	HIT_TYPE_CACHED   HitType = "cached"
	HIT_TYPE_UNCACHED HitType = "uncached"
	HIT_TYPE_ERROR    HitType = "error"

	VERSION_ENDPOINT = "v1/"
)

func Info(axleAddress string) (info map[string]interface{}, err error) {
	reqAddress := fmt.Sprintf("%s%sinfo", axleAddress, VERSION_ENDPOINT)
	body, err := doHttpRequest("GET", reqAddress, nil)
	if err != nil {
		return nil, fmt.Errorf("Unable to get Axle info: %s", err)
	}
	out := make(map[string]interface{})
	err = json.Unmarshal(body, &out)
	if err != nil {
		return nil, fmt.Errorf("Unable to get unmarshal info: %s", err)
	}
	if resp, exists := out["results"]; exists {
		// cast it on
		info, isValidCast := resp.(map[string]interface{})
		if (!isValidCast) {
			return nil, fmt.Errorf("Unable to get axle info, results was not a map")
		}
		return info, nil
	}
	return nil, fmt.Errorf("Unable to get axle info, missing results in response")
}

func Ping(axleAddress string) (err error) {
	reqAddress := fmt.Sprintf("%s%sping", axleAddress, VERSION_ENDPOINT)
	resp, err := http.Get(reqAddress)
	if err != nil {
		return fmt.Errorf("Unable to ping server at %v: %v", axleAddress, err)
	}
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("Unable to ping server at %v: %v", axleAddress, err)
	}
	if string(body) != "pong" {
		return fmt.Errorf(
			"ApiAxle server at %v didn't respond with pong, but with \"%v\"",
			axleAddress,
			string(body),
		)
	}
	return nil
}

// doHttpRequest performs verb on reqAddress, optionally posting postData.
// It returns the full page contents as a slice, and / or an error object
// describing any issues encountered.
func doHttpRequest(verb string, reqAddress string, postData []byte) (body []byte, err error) {

	buf := bytes.NewBuffer(make([]byte, 0))
	var req *http.Request = nil
	if postData != nil {
		buf = bytes.NewBuffer(postData)
	}
	req, err = http.NewRequest(verb, reqAddress, buf)
	if err != nil {
		return nil, fmt.Errorf(
			"Unable to prepare %s - %s: %s",
			verb,
			reqAddress,
			err.Error(),
		)
	}

	req.Header = map[string][]string{
		"Content-type": {"application/json"},
	}
	client := new(http.Client)
	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf(
			"Unable to %s api at %s: %s",
			verb,
			reqAddress,
			err.Error(),
		)
	}
	defer resp.Body.Close()

	body, err = ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf(
			"Unable to read response from %s: %s",
			verb,
			err.Error(),
		)
	}

	if resp.StatusCode != 200 {
		return body, fmt.Errorf(
			"Unable to %s api at %s, server returned status \"%s\" (%s)",
			verb,
			reqAddress,
			resp.Status,
			string(body),
		)
	}

	return body, nil
}

// parseFloatToTime is a utility function to convert a Javascript number
// respresentation of a date to a Go time.
func parseFloatToTime(theTime float64) time.Time {
	// this.CreatedAt is a float representing number of milliseconds since epoch
	seconds := int64(theTime / 1000)
	milliSeconds := int64(theTime) % 1000
	nanoSeconds := milliSeconds * 1000 * 1000
	return time.Unix(seconds, nanoSeconds)
}


func doStatsRequest(reqAddress string) (stats map[HitType]map[time.Time]map[int]int, err error) {
	body, err := doHttpRequest("GET", reqAddress, nil)
	if err != nil {
		return nil, err
	}

	responseMap := make(map[string]interface{})
	err = json.Unmarshal(body, &responseMap)
	if err != nil {
		return nil, fmt.Errorf(
			"Unable to unmarshal response from %s: %s",
			reqAddress,
			err.Error(),
		)
	}

	resultsInterface, exists := responseMap["results"]
	if !exists {
		return nil, fmt.Errorf("Missing results from %s", reqAddress)
	}
	results, validCast := resultsInterface.(map[string]interface{})
	if !validCast {
		return nil, fmt.Errorf("Missing stat details from %s", reqAddress)
	}

	stats = make(map[HitType]map[time.Time]map[int]int)
	for hitType, value := range results {
		if _, exists := stats[HitType(hitType)]; !exists {
			stats[HitType(hitType)] = make(map[time.Time]map[int]int)
		}
		statsInterface, goodCast := value.(map[string]interface{})
		if !goodCast {
			return nil, fmt.Errorf("Bad stats object at %s", reqAddress)
		}
		for timeStampStr, value := range statsInterface {
			timeStamp, _ := strconv.Atoi(timeStampStr)
			timeGroup := time.Unix(int64(timeStamp), 0)
			if _, exists := stats[HitType(hitType)][timeGroup]; !exists {
				stats[HitType(hitType)][timeGroup] = make(map[int]int)
			}
			statsInterface, goodCast := value.(map[string]interface{})
			if !goodCast {
				return nil, fmt.Errorf("Bad stats object at %s", reqAddress)
			}
			for responseCodeStr, countInterface := range statsInterface {
				responseCode, _ := strconv.Atoi(responseCodeStr)
				countFloat, goodCast := countInterface.(float64)
				if !goodCast {
					return nil, fmt.Errorf("Bad stats object at %s", reqAddress)
				}
				stats[HitType(hitType)][timeGroup][responseCode] = int(countFloat)
			}
		}
	}

	return stats, nil
}


func doChartsRequest(reqAddress string) (out map[string]int, err error) {

	body, err := doHttpRequest("GET", reqAddress, nil)
	if err != nil {
		return nil, err
	}

	responseMap := make(map[string]interface{})
	err = json.Unmarshal(body, &responseMap)
	if err != nil {
		return nil, fmt.Errorf(
			"Unable to unmarshal response from %s: %s",
			reqAddress,
			err.Error(),
		)
	}

	resultsInterface, exists := responseMap["results"]
	if !exists {
		return nil, fmt.Errorf("Missing results from %s", reqAddress)
	}
	results, isValidCast := resultsInterface.(map[string]interface{})
	if !isValidCast {
		return nil, fmt.Errorf("Unable to cast to map from %s", reqAddress)
	}
	out = make(map[string]int, len(results))
	for key, count := range results {
		out[key] = int(count.(float64))
	}

	return out, nil
}

func doKeysRequest(reqAddress string, axleAddress string) (keys []*Key, err error) {

	body, err := doHttpRequest("GET", reqAddress, nil)
	if err != nil {
		return nil, err
	}

	responseMap := make(map[string]interface{})
	err = json.Unmarshal(body, &responseMap)
	if err != nil {
		return nil, fmt.Errorf(
			"Unable to unmarshal response from %s: %s",
			reqAddress,
			err.Error(),
		)
	}

	resultsInterface, exists := responseMap["results"]
	if !exists {
		return nil, fmt.Errorf("Missing results from %s", reqAddress)
	}
	results, isValidCast := resultsInterface.(map[string]interface{})
	if !isValidCast {
		return nil, fmt.Errorf("Unable to cast to list of keys from %s", reqAddress)
	}
	keys = make([]*Key, len(results))
	x := 0
	for identifier, keyInterface := range results {
		key := NewKey(axleAddress, identifier)
		jsonvalue, err := json.Marshal(keyInterface)
		if err != nil {
			return nil, fmt.Errorf("Unable to decode key in response: %s", err)
		}
		err = json.Unmarshal(jsonvalue, key)
		if err != nil {
			return nil, fmt.Errorf("Unable to decode key in response: %s", err)
		}
		key.createOnSave = false
		keys[x] = key
		x++
	}

	return keys, nil
}

func doApisRequest(reqAddress string, axleAddress string) (out []*Api, err error) {
	body, err := doHttpRequest("GET", reqAddress, nil)
	if err != nil {
		return nil, err
	}

	response := make(map[string]interface{})
	err = json.Unmarshal(body, &response)
	if err != nil {
		return nil, fmt.Errorf(
			"Unable to unmarshal response: %s",
			err.Error(),
		)
	}
	response, validCast := response["results"].(map[string]interface{})
	if !validCast {
		return nil, fmt.Errorf(
			"Unable to unmarshal response: %s",
			err.Error(),
		)
	}
	out = make([]*Api, len(response))
	x := 0
	for identifier, value := range response {
		api := NewApi(axleAddress, identifier, "")
		jsonvalue, err := json.Marshal(value)
		if err != nil {
			return nil, fmt.Errorf("Unable to decode api in response: %s", err.Error())
		}
		err = json.Unmarshal(jsonvalue, api)
		if err != nil {
			return nil, fmt.Errorf("Unable to decode api in response: %s", err.Error())
		}
		api.createOnSave = false
		out[x] = api
		x++
	}

	return out, nil
}

/* ex: set noexpandtab: */
