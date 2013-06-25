package goaxle

import (
	//"fmt"
	"testing"
	"time"
)

func testGetNonExistentApi(t *testing.T) {
	result, err := GetApi(TEST_API_AXLE_SERVER, TEST_API_NAME)
	if result != nil || err == nil {
		t.Errorf("Recieved result for non-existent api: %s", result)
		t.Fatal()
	}
}

func testCreateApi(t *testing.T) {
	// fail with missing endpoint
	api := NewApi(TEST_API_AXLE_SERVER, TEST_API_NAME, "")
	err := api.Save()
	if err == nil {
		t.Errorf("Save succeeded with missing endpoint")
		t.Fatal()
	}
	// now try for success
	api = NewApi(TEST_API_AXLE_SERVER, TEST_API_NAME, TEST_API_ENDPOINT)
	err = api.Save()
	if err != nil {
		t.Errorf("Error creating new endpoint: %s", err)
		t.Fatal()
	}
	// now try to create a duplicate endpoint where there already is one
	api = NewApi(TEST_API_AXLE_SERVER, TEST_API_NAME, TEST_API_ENDPOINT)
	err = api.Save()
	if err == nil {
		t.Errorf("Was able to save duplicate endpoint")
		t.Fatal()
	}
}

func testGetApi(t *testing.T) *Api {
	// fail with non-existent api
	result, err := GetApi(TEST_API_AXLE_SERVER, TEST_API_NAME+".non-existent")
	if result != nil || err == nil {
		t.Errorf("Was able to retrieve non-existent API!")
		t.Fatal()
	}
	// get the created one
	result, err = GetApi(TEST_API_AXLE_SERVER, TEST_API_NAME)
	if err != nil {
		t.Errorf("Error retrieving api: %s", err)
		t.Fatal()
	}
	return result
}

func testUpdateApi(t *testing.T, api *Api) {
	origUpdatedAt := api.UpdatedAt
	origTimeout := api.EndPointTimeout
	// sleep a bit to ensure the updatedAt changes
	time.Sleep(2 * time.Second)
	// change timeout
	api.EndPointTimeout += 10
	err := api.Save()
	if err != nil {
		t.Errorf("Unable to save api: %s", err)
		t.Fatal()
	}
	// ensure our object has been updated
	if api.UpdatedAt == origUpdatedAt {
		t.Errorf("UpdatedAt didn't change from %v", origUpdatedAt)
		t.Fatal()
	}
	if api.EndPointTimeout == origTimeout {
		t.Errorf("EndPointTimeout didn't change from %v", origTimeout)
		t.Fatal()
	}
}

func testDeleteApi(t *testing.T, api *Api) {
	err := DeleteApi(TEST_API_AXLE_SERVER, TEST_API_NAME)
	if err != nil {
		t.Errorf("Unable to delete: %v", err)
		t.Fatal()
	}
	// attempt double delete
	err = DeleteApi(TEST_API_AXLE_SERVER, TEST_API_NAME)
	if err == nil {
		t.Errorf("Managed to delete something that wasn't there!")
		t.Fatal()
	}
	// attempt a save
	err = api.Save()
	if err == nil {
		t.Errorf("Managed to save an API that no longer exists on the server!")
		t.Fatal()
	}
}

func testLinkKey(t *testing.T, api *Api) {
	key, err := api.LinkKey(TEST_KEY_NAME)
	if err != nil {
		t.Errorf("Error linking key: %v", err)
		t.Fatal()
	}
	if key == nil {
		t.Errorf("Key not returned")
		t.Fatal()
	}
}

func testUnlinkKey(t *testing.T) {
	_, err := ApiUnlinkKey(TEST_API_AXLE_SERVER, TEST_API_NAME, TEST_KEY_NAME)
	if err != nil {
		t.Errorf("Error unlinking key: %v", err)
		t.Fatal()
	}
}

func testApiKeys(t *testing.T) {
	keys, err := ApiKeys(TEST_API_AXLE_SERVER, TEST_API_NAME, 0, 10)
	if err != nil {
		t.Errorf("Error listing keys: %v", err)
		t.Fatal()
	}
	if keys[0].Identifier != TEST_KEY_NAME {
		t.Errorf("Key was not returned")
		t.Fatal()
	}
}

// TODO: Add checks for actual values, but this requires a
//		 mechanism to get usage stats into the server.
func testKeyCharts(t *testing.T, api *Api) {
	_, err := api.KeyCharts(GRANULARITY_MINUTES)
	if err != nil {
		t.Errorf("Error getting chart: %v", err)
		t.Fatal()
	}
}

func testApiStats(t *testing.T, api *Api) {
	anHourAgo, _ := time.ParseDuration("-1hr")
	stats, err := api.Stats(
		time.Now().Add(anHourAgo),
		time.Now(),
		GRANULARITY_DAYS,
	)
	if err != nil {
		t.Errorf("Error gettings stats: %v", err)
		t.Fatal()
	}
	if stats == nil {
		t.Errorf("Empty stats returned")
		t.Fatal()
	}
}

func testApis(t *testing.T) {
	apis, err := Apis(TEST_API_AXLE_SERVER, 0, 10)
	if err != nil {
		t.Errorf("Error getting apis: %v", err)
		t.Fatal()
	}
	if len(apis) <= 0 {
		t.Errorf("Wrong number of apis returned %d", len(apis))
		t.Fatal()
	}
}

func testApiCharts(t *testing.T) {

	_, err := ApiCharts(TEST_API_AXLE_SERVER, GRANULARITY_MINUTES)
	if err != nil {
		t.Errorf("Error getting api keys charts: %v", err)
		t.Fatal()
	}
}

/* ex: set noexpandtab: */
