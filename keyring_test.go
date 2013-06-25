package goaxle

import (
	"testing"
	//"fmt"
	"time"
)

func testNewKeyRing(t *testing.T) (k *KeyRing) {
	keyRing := NewKeyRing(TEST_API_AXLE_SERVER, TEST_KEYRING_NAME)
	err := keyRing.Save()
	if err != nil {
		t.Errorf("Unable to save keyRing: %v", err)
		t.Fatal()
	}
	return keyRing
}

func testGetKeyRing(t *testing.T) {
	kr, err := GetKeyRing(TEST_API_AXLE_SERVER, TEST_KEYRING_NAME)
	if err != nil {
		t.Errorf("Unable to get KeyRing: %v", err)
		t.Fatal()
	}
	if kr == nil {
		t.Errorf("Unable to get keyring")
		t.Fatal()
	}
}

func testUpdateKeyRing(t *testing.T, k *KeyRing) {
	err := k.Save()
	if err == nil {
		t.Errorf("Not meant to be able to update keyrings :/")
		t.Fatal()
	}
}

func testDeleteKeyRing(t *testing.T) {
	err := DeleteKeyRing(TEST_API_AXLE_SERVER, TEST_KEYRING_NAME)
	if err != nil {
		t.Errorf("Unable to delete keyRing: %v", err)
	}
}

func testKeyRingKeys(t *testing.T, kr *KeyRing) {
	keys, err := KeyRingKeys(TEST_API_AXLE_SERVER, TEST_KEYRING_NAME, 0, 10)
	if err != nil {
		t.Errorf("Error getting keys for keyring: %v", err)
	}
	if len(keys) != 1 {
		t.Errorf("Wrong number of keys returned for keyring: %d", len(keys))
	}
}

func testKeyRingLinkKey(t *testing.T) {
	key, err := KeyRingLinkKey(TEST_API_AXLE_SERVER, TEST_KEYRING_NAME, TEST_KEY_NAME)
	if err != nil {
		t.Errorf("Error linking key to keyring: %v", err)
	}
	if key.Identifier != TEST_KEY_NAME {
		t.Errorf("Incorrect key identifier returned from linking: %v", key.Identifier)
	}
}

func testKeyRingStats(t *testing.T, kr *KeyRing) {
	anHourAgo, _ := time.ParseDuration("-1hr")
	stats, err := kr.Stats(
		time.Now().Add(anHourAgo),
		time.Now(),
		GRANULARITY_DAYS,
	)
	if err != nil {
		t.Errorf("Error getting keyring stats: %v", err)
	}
	if stats == nil {
		t.Errorf("No stats returned for %v", kr.Identifier)
	}
}

/* ex: set noexpandtab: */
