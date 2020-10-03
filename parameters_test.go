package logger

import "testing"

// TestMakeParameter test making an error struct and MakeParameter() method
func TestMakeParameter(t *testing.T) {
	param := MakeParameter("myKey", "myValue")
	if param.Key() != "myKey" {
		t.Fatalf("expected value: %s, got: %s", "myKey", param.Key())
	}
	if param.Value() != "myValue" {
		t.Fatalf("expected value: %s, got: %s", "myValue", param.Value())
	}
	if param.String() != `{"key":"myKey","value":"myValue"}` {
		t.Fatalf("expected value: %s, got: %s", `{"key":"myKey","value":"myValue"}`, param.String())
	}
}

// BenchmarkMakeParameter benchmarks the MakeParameter() method
func BenchmarkMakeParameter(b *testing.B) {
	for i := 0; i < b.N; i++ {
		_ = MakeParameter("myKey", "myValue")
	}
}
