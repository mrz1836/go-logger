package logger

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

// TestMakeParameter test making an error struct and MakeParameter() method
func TestMakeParameter(t *testing.T) {
	param := MakeParameter("myKey", "myValue")
	assert.Equal(t, "myKey", param.Key())
	assert.Equal(t, "myValue", param.Value())
	assert.Equal(t, `{"key":"myKey","value":"myValue"}`, param.String())
}

// BenchmarkMakeParameter benchmarks the MakeParameter() method
func BenchmarkMakeParameter(b *testing.B) {
	for i := 0; i < b.N; i++ {
		_ = MakeParameter("myKey", "myValue")
	}
}
