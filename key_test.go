package goaxle

import (
	"testing"
	"time"
	//"fmt"
)

func testNewKey(t *testing.T) (k *Key) {
	key := NewKey(TEST_API_AXLE_SERVER, TEST_KEY_NAME)
	key.Qpd = 200000
	err := key.Save()
	if err != nil {
		t.Errorf("Unable to save key: %v", err)
		t.Fatal()
	}
	return key
}

func testGetKey(t *testing.T) {
	key, err := GetKey(TEST_API_AXLE_SERVER, TEST_KEY_NAME)
	if err != nil {
		t.Errorf("Unable to save key: %v", err)
		t.Fatal()
	}
	if key.Qpd != 200000 {
		t.Errorf("Loaded key is different")
		t.Fatal()
	}
}

func testUpdateKey(t *testing.T, k *Key) {
	origQpd := k.Qpd
	origQps := k.Qpd
	origUpdatedAt := k.UpdatedAt

	k.Qpd += 100
	k.Qps += 10

	time.Sleep(2 * time.Second)

	err := k.Save()
	if err != nil {
		t.Errorf("Error saving key: %v", err)
		t.Fatal()
	}

	if k.UpdatedAt == origUpdatedAt {
		t.Errorf("UpdatedAt didn't change from %v", origUpdatedAt)
		t.Fatal()
	}

	if k.Qpd == origQpd || k.Qps == origQps {
		t.Errorf("Qpd or Qps changes weren't saved: %v", err)
		t.Fatal()
	}

}

func testDeleteKey(t *testing.T) {
	err := DeleteKey(TEST_API_AXLE_SERVER, TEST_KEY_NAME)
	if err != nil {
		t.Errorf("Unable to delete key: %v", err)
	}
}

func testKeyApiCharts(t *testing.T) {
	_, err := KeyApiCharts(TEST_API_AXLE_SERVER, TEST_KEY_NAME, GRANULARITY_MINUTES)
	if err != nil {
		t.Errorf("Error getting chart: %v", err)
		t.Fatal()
	}
}

func testKeyApis(t *testing.T, k *Key) {

}

/* ex: set noexpandtab: */
