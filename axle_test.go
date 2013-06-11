package goaxle

import (
	"testing"
	//"fmt"
)

const (
	TEST_API_AXLE_SERVER = "http://localhost:28902/"
	TEST_API_NAME        = "goaxletestapi"
	TEST_KEY_NAME        = "goaxletestkey"
	TEST_KEYRING_NAME    = "goaxletestkeyring"
	TEST_API_ENDPOINT    = "localhost:80"
)

// as we rely on the state of the axle server for each test,
// it's just much more convient to do them all at once :/
func TestAll(t *testing.T) {

	// remove anything that might have been leftover from failed tests
	DeleteApi(TEST_API_AXLE_SERVER, TEST_API_NAME)
	DeleteKey(TEST_API_AXLE_SERVER, TEST_KEY_NAME)
	DeleteKeyRing(TEST_API_AXLE_SERVER, TEST_KEYRING_NAME)

	testGetNonExistentApi(t)
	testCreateApi(t)
	api := testGetApi(t)
	testUpdateApi(t, api)
	testApiStats(t, api)
	k := testNewKey(t)
	testGetKey(t)
	testUpdateKey(t, k)
	testLinkKey(t, api)
	testKeyCharts(t, api)
	testListKeys(t)
	testApis(t)
	testApiCharts(t)
	testUnlinkKey(t)
	testDeleteKey(t)
	kr := testNewKeyRing(t)
	testGetKeyRing(t)
	testUpdateKeyRing(t, kr)
	testDeleteKeyRing(t)
	testDeleteApi(t, api)

}

/* ex: set noexpandtab: */
