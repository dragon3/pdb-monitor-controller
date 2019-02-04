package log

import "testing"

func TestNew(t *testing.T) {
	if _, err := New("INVALID_LOG_LEVEL"); err == nil {
		t.Fatal("should be failed due to invalid log level")
	}
	if _, err := New("DEBUG"); err != nil {
		t.Fatal(err)
	}
	if _, err := New("INFO"); err != nil {
		t.Fatal(err)
	}
	if _, err := New("WARN"); err != nil {
		t.Fatal(err)
	}
	if _, err := New("ERROR"); err != nil {
		t.Fatal(err)
	}
}
