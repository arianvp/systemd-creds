package server

import "testing"

func TestParsePeerNameExampleWorks(t *testing.T) {
	test := "\x00adf9d86b6eda275e/unit/foobar.service/credx"
	unitName, credID, ok := parsePeerName(test)
	if !ok {
		t.Errorf("Expected true, got false")
	}
	if unitName != "foobar.service" {
		t.Errorf("Expected foobar.service, got %s", unitName)
	}
	if credID != "credx" {
		t.Errorf("Expected credx, got %s", credID)
	}
}

func TestParsePeerNameExampleFails(t *testing.T) {
	test := "foobar"
	unitName, credID, ok := parsePeerName(test)
	if ok {
		t.Errorf("Expected false, got true")
	}
	if unitName != "" {
		t.Errorf("Expected empty string, got %s", unitName)
	}
	if credID != "" {
		t.Errorf("Expected empty string, got %s", credID)
	}
}
