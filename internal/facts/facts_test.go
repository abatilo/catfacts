package facts

import "testing"

func Test(t *testing.T) {
	s, completion := GenerateFact(0)
	if s == "" {
		t.Error("Expected a fact, got nothing")
	}

	if !completion {
		t.Error("Expected a fact to be complete, got incomplete")
	}

	t.Log(s)
}
