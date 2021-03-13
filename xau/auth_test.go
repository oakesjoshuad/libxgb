package xau

import (
	"testing"
)

func TestXauth(t *testing.T) {
	authNames := []string{"MIT-MAGIC-COOKIE-1"}
	if xai, err := GetAuthByAddr(FamilyLocal, "void", "0", authNames[0]); err != nil {
		t.Errorf("Unable to find authinfo; %w", err)
	} else {
		t.Log("GetAuthByAddr: ", xai)
	}
	if xai, err := GetBestAuthByAddr(FamilyLocal, "void", "0", authNames); err != nil {
		t.Errorf("Unable to find authinfo; %w", err)
	} else {
		t.Log("GetBestAuthByAddr: ", xai)
	}
}
