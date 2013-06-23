// Package goaxle provides bindings to the ApiAxle management API.
package goaxle

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"net/http"
	"time"
	"encoding/json"
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

	GRANULARITY_SECONDS Granularity = "seconds"
	GRANULARITY_MINUTES Granularity = "minutes"
	GRANULARITY_HOURS   Granularity = "hours"
	GRANULARITY_DAYS    Granularity = "days"

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

/* ex: set noexpandtab: */
