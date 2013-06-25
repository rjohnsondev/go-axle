package goaxle

import (
	"encoding/json"
	"fmt"
	"time"
	"net/url"
)

type Api struct {
	// Identifier is the name given to this API.  Modification not supported.
	Identifier string `json:"-"`

	// The time this api was created.
	// Use of this field is discouraged, use ParseCreatedAt.
	CreatedAt float64 `json:"createdAt,omitempty"`
	// The time this api was updated.
	// Use of this field is discouraged, use ParseUpdatedAt.
	UpdatedAt float64 `json:"updatedAt,omitempty"`

	// The time in seconds that every call under this API should be cached
	GlobalCache int `json:"globalCache"`

	// The protocol for the API, whether or not to use SSL
	Protocol Protocol `json:"protocol"`

	// The resulting data type of the endpoint.
	// This is redundant at the moment but will eventually support both XML too.
	ApiFormat ApiFormat `json:"apiFormat"`

	// The endpoint for the API. For example; "graph.facebook.com"
	EndPoint string `json:"endPoint,omitempty"`

	// Seconds to wait before timing out the connection
	EndPointTimeout int `json:"endPointTimeout"`

	// Max redirects that are allowed when endpoint called.
	EndPointMaxRedirects int `json:"endPointMaxRedirects"`

	// Regular expression used to extract API key from url.
	// Axle will use the **first** matched grouping and then apply that as the key.
	// Using the `api_key` or `apiaxle_key` will take precedence.
	ExtractKeyRegex string `json:"extractKeyRegex,omitempty"`

	// An optional path part that will always be called when the API is hit.
	DefaultPath string `json:"defaultPath,omitempty"`

	// Disable this API causing errors when it's hit.
	Disabled bool `json:"disabled"`

	// Set to true to require that SSL certificates be valid
	StrictSSL bool `json:"strictSSL"`

	// address where this api is located
	axleAddress string
	// do need to create a new api on save?
	createOnSave bool
}

func Apis(axleAddress string, from int, to int) (out []*Api, err error) {
	reqAddress := fmt.Sprintf(
		"%s%sapis?resolve=true&from=%d&to=%d",
		axleAddress,
		VERSION_ENDPOINT,
		from,
		to,
	)
	return doApisRequest(reqAddress, axleAddress)
}

// NewApi creates a new API object with defaults.
func NewApi(axleAddress string, identifier string, endPoint string) (out *Api) {
	out = &Api{
		Identifier:           identifier,
		Protocol:             API_PROTOCOL_HTTP,
		ApiFormat:            API_FORMAT_JSON,
		EndPoint:             endPoint,
		EndPointTimeout:      2,
		EndPointMaxRedirects: 2,
		StrictSSL:            true,
		createOnSave:         true,
		axleAddress:          axleAddress,
	}
	return out
}

// GetAPI retrieves an existing api object from the server.
func GetApi(axleAddress string, identifier string) (out *Api, err error) {

	reqAddress := fmt.Sprintf("%s%sapi/%s", axleAddress, VERSION_ENDPOINT, url.QueryEscape(identifier))
	body, err := doHttpRequest("GET", reqAddress, nil)
	if err != nil {
		return nil, err
	}
	// unmarshal into our new api object
	api := NewApi(axleAddress, identifier, "")
	err = populateApiFromResponse(&api, body, []string{"results"})
	if err != nil {
		return nil, err
	}
	api.createOnSave = false

	return api, err
}

// Create / Update this API on the ApiAxle server.
// To modify an existing API, be sure to retrieve it with GetApi, otherwise
// the library will attempt to create a new API of the same name.
func (this *Api) Save() (err error) {
	reqAddress := fmt.Sprintf("%s%sapi/%s", this.axleAddress, VERSION_ENDPOINT, url.QueryEscape(this.Identifier))

	// update the updatedAt timestamp
	this.UpdatedAt = float64(time.Now().UnixNano() / (1000 * 1000))
	marshalled, err := json.Marshal(this)
	if err != nil {
		return fmt.Errorf("Unable to marshal API: %s", err.Error())
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
		err = populateApiFromResponse(&this, body, []string{"results", "new"})
	} else {
		err = populateApiFromResponse(&this, body, []string{"results"})
	}

	if err != nil {
		return err
	}

	this.createOnSave = false

	return nil
}

// ParseCreatedAt returns the API created time as a Go time.Time.
func (this *Api) ParseCreatedAt() time.Time {
	return parseFloatToTime(this.CreatedAt)
}

// ParseUpdatedAt returns the updated time as a Go time.Time.
func (this *Api) ParseUpdatedAt() time.Time {
	return parseFloatToTime(this.UpdatedAt)
}

// String provides a JSON-like formated representation of this API object
func (this *Api) String() string {
	out, err := json.MarshalIndent(this, "", "    ")
	if err != nil {
		return "<nil>"
	}
	reqAddress := fmt.Sprintf(
		"%s%sapi/%s",
		this.axleAddress,
		VERSION_ENDPOINT,
		url.QueryEscape(this.Identifier),
	)
	return fmt.Sprintf("Api - %s: %s", reqAddress, string(out))
}

// LinkKey links the provided key with this API.
func (this *Api) LinkKey(keyIdentifier string) (key *Key, err error) {
	return LinkKey(this.axleAddress, this.Identifier, keyIdentifier)
}

// LinkKey links the provided key with this API.
func LinkKey(axleAddress string, apiIdentifier string, keyIdentifier string) (key *Key, err error) {
	reqAddress := fmt.Sprintf(
		"%s%sapi/%s/linkkey/%s",
		axleAddress,
		VERSION_ENDPOINT,
		url.QueryEscape(apiIdentifier),
		url.QueryEscape(keyIdentifier),
	)

	body, err := doHttpRequest("PUT", reqAddress, []byte("{}"))
	if err != nil {
		return nil, err
	}

	key = NewKey(axleAddress, keyIdentifier)
	err = populateKeyFromResponse(&key, body, []string{"results"})
	if err != nil {
		return nil, err
	}
	key.createOnSave = false

	return key, nil
}

// UnlinkKey disassociates the provided key with this API.
func (this *Api) UnlinkKey(keyIdentifier string) (key *Key, err error) {
	return UnlinkKey(this.axleAddress, this.Identifier, keyIdentifier)
}

// UnlinkKey disassociates the provided key with this API.
func UnlinkKey(axleAddress string, apiIdentifier string, keyIdentifier string) (key *Key, err error) {
	reqAddress := fmt.Sprintf(
		"%s%sapi/%s/unlinkkey/%s",
		axleAddress,
		VERSION_ENDPOINT,
		url.QueryEscape(apiIdentifier),
		url.QueryEscape(keyIdentifier),
	)

	body, err := doHttpRequest("PUT", reqAddress, []byte("{}"))
	if err != nil {
		return nil, err
	}

	key = NewKey(axleAddress, keyIdentifier)
	err = populateKeyFromResponse(&key, body, []string{"results"})
	if err != nil {
		return nil, err
	}
	key.createOnSave = false

	return key, nil
}

// Keys returns a listing of all the keys linked with this API
func (this *Api) Keys(from int, to int) (keys []*Key, err error) {
	return ApiKeys(this.axleAddress, this.Identifier, from, to)
}

// ApiKeys returns a listing of all the keys linked with this API
func ApiKeys(axleAddress string, apiIdentifier string, from int, to int) (keys []*Key, err error) {
	reqAddress := fmt.Sprintf(
		"%s%sapi/%s/keys?resolve=true&from=%d&to=%d",
		axleAddress,
		VERSION_ENDPOINT,
		url.QueryEscape(apiIdentifier),
		from,
		to,
	)

	return doKeysRequest(reqAddress, axleAddress)
}

// ApiCharts lists the top 100 keys and their hit rate for time period granularity.
func ApiCharts(axleAddress string, granularity Granularity) (out map[string]int, err error) {
	reqAddress := fmt.Sprintf(
		"%s%sapis/charts?granularity=%s",
		axleAddress,
		VERSION_ENDPOINT,
		granularity,
	)

	return doChartsRequest(reqAddress)
}

// KeyCharts lists the top 100 keys and their hit rate for time period granularity.
func (this *Api) KeyCharts(granularity Granularity) (results map[string]int, err error) {
	return ApiKeyCharts(this.axleAddress, this.Identifier, granularity)
}

// ApiKeyCharts lists the top 100 keys and their hit rate for time period granularity.
func ApiKeyCharts(axleAddress string, apiIdentifier string, granularity Granularity) (out map[string]int, err error) {
	reqAddress := fmt.Sprintf(
		"%s%sapi/%s/keycharts?granularity=%s",
		axleAddress,
		VERSION_ENDPOINT,
		url.QueryEscape(apiIdentifier),
		granularity,
	)

	return doChartsRequest(reqAddress)
}

func (this *Api) Stats(from time.Time, to time.Time, granularity Granularity) (stats map[HitType]map[time.Time]map[int]int, err error) {
	return ApiStats(this.axleAddress, this.Identifier, from, to, "", granularity)
}
func (this *Api) StatsForKey(from time.Time, to time.Time, forkey string, granularity Granularity) (stats map[HitType]map[time.Time]map[int]int, err error) {
	return ApiStats(this.axleAddress, this.Identifier, from, to, forkey, granularity)
}

func ApiStats(axleAddress string, apiIdentifier string, from time.Time, to time.Time, forkey string, granularity Granularity) (stats map[HitType]map[time.Time]map[int]int, err error) {

	reqAddress := fmt.Sprintf(
		"%s%sapi/%s/stats?from=%d&to=%d&granularity=%s",
		axleAddress,
		VERSION_ENDPOINT,
		url.QueryEscape(apiIdentifier),
		from.Unix(),
		to.Unix(),
		granularity,
	)

	if forkey != "" {
		reqAddress += "&forkey=" + forkey
	}

	return doStatsRequest(reqAddress)
}

// populateApiFromResponse updates the provided Api pointer with the fields
// provided in the response map.
func populateApiFromResponse(api **Api, body []byte, detailsLocation []string) (err error) {
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
				"key %s did not contain map: %s",
				key,
			)
		}
	}

	if _, exists := response["endPoint"]; !exists {
		return fmt.Errorf(
			"Unable to parse response into Api: Missing required field \"endPoint\"",
		)
	}
	// making use of json to populate the object
	jsonvalue, err := json.Marshal(response)
	if err != nil {
		return fmt.Errorf("Unable to decode api in response: %s", err.Error())
	}
	err = json.Unmarshal(jsonvalue, api)
	if err != nil {
		return fmt.Errorf("Unable to decode api in response: %s", err.Error())
	}
	return nil
}

// DeleteApi removes the identified API.  Any existing objects represting this
// API will error on Save().
func DeleteApi(axleAddress string, identifier string) (err error) {
	reqAddress := fmt.Sprintf("%s%sapi/%s", axleAddress, VERSION_ENDPOINT, url.QueryEscape(identifier))

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
		return fmt.Errorf("Delete of API at %s failed", reqAddress)
	}

	return nil
}

/* ex: set noexpandtab: */
