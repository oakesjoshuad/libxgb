package xau

import (
	"os"
	"testing"
)

func TestXauth(t *testing.T) {
	localhost, err := os.Hostname()
	if err != nil {
		t.Error(err)
	}

	authNames := []string{"MIT-MAGIC-COOKIE-1"}
	if xai, err := GetAuthByAddr(FamilyLocal, localhost, "0", authNames[0]); err != nil {
		t.Errorf("Unable to find authinfo; %w", err)
	} else {
		t.Log("GetAuthByAddr: ", xai)
	}
	if xai, err := GetBestAuthByAddr(FamilyLocal, localhost, "0", authNames); err != nil {
		t.Errorf("Unable to find authinfo; %w", err)
	} else {
		t.Log("GetBestAuthByAddr: ", xai)
	}
}
