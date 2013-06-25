package goaxle

import (
	"encoding/json"
	"fmt"
	"time"
	"net/url"
)

type Key struct {
	// Identifier is the name given to this Key.  Modification not supported.
	Identifier string `json:"-"`

	// The time this key was created.
	// Use of this field is discouraged, use ParseCreatedAt.
	CreatedAt float64 `json:"createdAt,omitempty"`
	// The time this key was updated.
	// Use of this field is discouraged, use ParseUpdatedAt.
	UpdatedAt float64 `json:"updatedAt,omitempty"`

	// A shared secret which is used when signing a call to the key.
	SharedSecret string `json:"sharedSecret,omitempty"`

	// Number of queries that can be called per day. Set to `-1` for no limit.
	Qpd int `json:"qpd"`

	// Number of queries that can be called per second. Set to `-1` for no limit.
	Qps int `json:"qps"`

	// Names of the Apis that this key belongs to.
	ForApis []string `json:"forApis,omitempty"`

	// Disable this Key causing errors when it's hit.
	Disabled bool `json:"disabled"`

	// address where this key is located
	axleAddress string
	// do need to create a new key on save?
	createOnSave bool
}

// NewKey creates a new Key object with defaults.
func NewKey(axleAddress string, identifier string) (out *Key) {
	out = &Key{
		Identifier:   identifier,
		Qpd:          172800,
		Qps:          2,
		Disabled:     false,
		axleAddress:  axleAddress,
		createOnSave: true,
	}
	return out
}

// Create / Update this Key on the ApiAxle server.
// To modify an existing Key, be sure to retrieve it with GetKey, otherwise
// the library will attempt to create a new Key of the same name.
func (this *Key) Save() (err error) {
	reqAddress := fmt.Sprintf(
		"%s%skey/%s",
		this.axleAddress,
		VERSION_ENDPOINT,
		url.QueryEscape(this.Identifier),
	)

	// update the updatedAt timestamp
	this.UpdatedAt = float64(time.Now().UnixNano() / (1000 * 1000))
	marshalled, err := json.Marshal(this)
	if err != nil {
		return fmt.Errorf("Unable to marshal Key: %s", err.Error())
	}

	httpMethod := "POST"
	if !this.createOnSave {
		httpMethod = "PUT"
	}

	body, err := doHttpRequest(httpMethod, reqAddress, marshalled)
	if err != nil {
		return err
	}

	if !this.createOnSave {
		err = populateKeyFromResponse(&this, body, []string{"results", "new"})
	} else {
		err = populateKeyFromResponse(&this, body, []string{"results"})
	}

	if err != nil {
		return err
	}

	this.createOnSave = false

	return nil
}

// GetKey retrieves an existing api object from the server.
func GetKey(axleAddress string, identifier string) (out *Key, err error) {

	reqAddress := fmt.Sprintf("%s%skey/%s", axleAddress, VERSION_ENDPOINT, url.QueryEscape(identifier))
	body, err := doHttpRequest("GET", reqAddress, nil)
	if err != nil {
		return nil, err
	}

	// unmarshal into our new key object
	key := NewKey(axleAddress, identifier)
	err = populateKeyFromResponse(&key, body, []string{"results"})
	if err != nil {
		return nil, err
	}
	key.createOnSave = false

	return key, err
}

// populateKeyFromResponse updates the provided Key pointer with the fields
// provided in the response map.
func populateKeyFromResponse(key **Key, body []byte, detailsLocation []string) (err error) {

	response := make(map[string]interface{})
	err = json.Unmarshal(body, &response)
	if err != nil {
		return fmt.Errorf(
			"Unable to unmarshal response: %s",
			err.Error(),
		)
	}

	// navigate to the correct spot in the response to read from
	for _, key := range detailsLocation {
		resultsInterface, exists := response[key]
		if !exists {
			return fmt.Errorf(
				"Response map did not contain expected key: %s",
				key,
			)
		}
		var isValidCast bool
		response, isValidCast = resultsInterface.(map[string]interface{})
		if !isValidCast {
			return fmt.Errorf(
				"key %s did not contain map",
				key,
			)
		}
	}

	// making use of json to populate the object
	jsonvalue, err := json.Marshal(response)
	if err != nil {
		return fmt.Errorf("Unable to decode key in response: %s", err.Error())
	}
	err = json.Unmarshal(jsonvalue, key)
	if err != nil {
		return fmt.Errorf("Unable to decode key in response: %s", err.Error())
	}
	return nil
}

// String provides a JSON-like formated representation of this Key object
func (this *Key) String() string {
	out, err := json.MarshalIndent(this, "", "    ")
	if err != nil {
		return "<nil>"
	}
	reqAddress := fmt.Sprintf(
		"%s%skey/%s",
		this.axleAddress,
		VERSION_ENDPOINT,
		url.QueryEscape(this.Identifier),
	)
	return fmt.Sprintf("Key - %s: %s", reqAddress, string(out))
}

// DeleteKey removes the identified Key.  Any existing objects represting this
// Key will error on Save().
func DeleteKey(axleAddress string, identifier string) (err error) {
	reqAddress := fmt.Sprintf("%s%skey/%s", axleAddress, VERSION_ENDPOINT, url.QueryEscape(identifier))

	body, err := doHttpRequest("DELETE", reqAddress, nil)
	if err != nil {
		return err
	}

	responseMap := make(map[string]interface{})
	err = json.Unmarshal(body, &responseMap)
	if err != nil {
		return fmt.Errorf(
			"Unable to unmarshal response from %s: %s",
			reqAddress,
			err.Error(),
		)
	}

	// in this case, our result is what is contained in the "results" key
	resultsInterface, exists := responseMap["results"]
	if !exists {
		return fmt.Errorf("Missing response from %s", reqAddress)
	}
	succeeded, isValidCast := resultsInterface.(bool)
	if !isValidCast {
		return fmt.Errorf(
			"Unable to extract response object from %s",
			reqAddress,
		)
	}

	if !succeeded {
		return fmt.Errorf("Delete of Key at %s failed", reqAddress)
	}

	return nil
}

// ApiCharts lists the top 100 apis for this key and their hit rate for time period granularity.
func (this *Key) ApiCharts(granularity Granularity) (out map[string]int, err error) {
	return KeyApiCharts(this.axleAddress, this.Identifier, granularity)
}
// KeyApiCharts lists the top 100 apis for the specified key and their hit rate for time period granularity.
func KeyApiCharts(axleAddress string, keyIdentifier string, granularity Granularity) (out map[string]int, err error) {
	reqAddress := fmt.Sprintf(
		"%s%skey/%s/apicharts?granularity=%s",
		axleAddress,
		VERSION_ENDPOINT,
		url.QueryEscape(keyIdentifier),
		granularity,
	)

	return doChartsRequest(reqAddress)
}

// Get the most used keys and their hit counts.
func KeysCharts(axleAddress string, granularity Granularity) (out map[string]int, err error) {
	reqAddress := fmt.Sprintf(
		"%s%skeys/charts?granularity=%s",
		axleAddress,
		VERSION_ENDPOINT,
		granularity,
	)

	return doChartsRequest(reqAddress)
}

// List apis belonging to a key.
func (this *Key) Apis() (out []*Api, err error) {
	return KeyApis(this.axleAddress, this.Identifier)
}
func KeyApis(axleAddress string, keyIdentifier string) (out []*Api, err error) {
	reqAddress := fmt.Sprintf(
		"%s%skey/%s/apis?resolve=true",
		axleAddress,
		VERSION_ENDPOINT,
		url.QueryEscape(keyIdentifier),
	)
	return doApisRequest(reqAddress, axleAddress)
}

// Get the real time hits for a key.
func (this *Key) Stats(from time.Time, to time.Time, granularity Granularity) (stats map[HitType]map[time.Time]map[int]int, err error) {
	return KeyStats(this.axleAddress, this.Identifier, from, to, "", granularity)
}
// Get the real time hits for a key.
func (this *Key) StatsForApi(from time.Time, to time.Time, forapi string, granularity Granularity) (stats map[HitType]map[time.Time]map[int]int, err error) {
	return KeyStats(this.axleAddress, this.Identifier, from, to, forapi, granularity)
}

// Get the real time hits for a key.
func KeyStats(axleAddress string, keyIdentifier string, from time.Time, to time.Time, forapi string, granularity Granularity) (stats map[HitType]map[time.Time]map[int]int, err error) {

	reqAddress := fmt.Sprintf(
		"%s%skey/%s/stats?from=%d&to=%d&granularity=%s",
		axleAddress,
		VERSION_ENDPOINT,
		url.QueryEscape(keyIdentifier),
		from.Unix(),
		to.Unix(),
		granularity,
	)

	if forapi != "" {
		reqAddress += "&forapi=" + url.QueryEscape(forapi)
	}

	return doStatsRequest(reqAddress)
}

// List all of the available keys.
func Keys(axleAddress string, from int, to int) (keys []*Key, err error) {

	reqAddress := fmt.Sprintf(
		"%s%skeys?resolve=true&from=%d&to=%d",
		axleAddress,
		VERSION_ENDPOINT,
		from,
		to,
	)

	return doKeysRequest(reqAddress, axleAddress)
}

/* ex: set noexpandtab: */
