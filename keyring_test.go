package goaxle

import (
	"testing"
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
		t.Errorf("Unable to save keyRing: %v", err)
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

/* ex: set noexpandtab: */
