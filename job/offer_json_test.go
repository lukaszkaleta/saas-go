package job

import (
	"encoding/json"
	"strings"
	"testing"
)

func TestOfferModel_JSON(t *testing.T) {
	m := &OfferModel{
		Id: 123,
	}

	data, err := json.Marshal(m)
	if err != nil {
		t.Fatalf("failed to marshal: %v", err)
	}

	s := string(data)
	if strings.Contains(s, "Idable") {
		t.Errorf("JSON should not contain Idable, but got: %s", s)
	}
	if !strings.Contains(s, `"id":123`) {
		t.Errorf("JSON should contain id:123, but got: %s", s)
	}
}

func TestOfferRevisionModel_JSON(t *testing.T) {
	m := &OfferRevisionModel{
		Id: 456,
	}

	data, err := json.Marshal(m)
	if err != nil {
		t.Fatalf("failed to marshal: %v", err)
	}

	s := string(data)
	if strings.Contains(s, "Idable") {
		t.Errorf("JSON should not contain Idable, but got: %s", s)
	}
	if !strings.Contains(s, `"id":456`) {
		t.Errorf("JSON should contain id:456, but got: %s", s)
	}
}
