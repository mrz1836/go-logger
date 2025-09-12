package logger

import (
	"encoding/json"
	"strings"
	"testing"
	"unicode/utf8"
)

func FuzzMakeParameter(f *testing.F) {
	f.Add("key", "value")
	f.Add("", "")
	f.Add("unicode-key-test", "unicode-value-test")
	f.Add("key-with-special-chars", "value\nwith\ttabs")
	f.Add(strings.Repeat("long-key", 100), strings.Repeat("long-value", 100))
	f.Add("null-byte-key\x00", "null-byte-value\x00")
	f.Add("json-chars", `{"nested": "json"}`)

	f.Fuzz(func(t *testing.T, key, value string) {
		defer func() {
			if r := recover(); r != nil {
				t.Errorf("MakeParameter panicked: key=%q, value=%q, panic: %v", key, value, r)
			}
		}()

		param := MakeParameter(key, value)

		if param == nil {
			t.Error("MakeParameter should return non-nil parameter")
			return
		}

		if param.K != key {
			t.Errorf("Parameter key mismatch: expected %q, got %q", key, param.K)
		}

		if param.V != value {
			t.Errorf("Parameter value mismatch: expected %q, got %q", value, param.V)
		}

		if param.Key() != key {
			t.Errorf("Parameter.Key() mismatch: expected %q, got %q", key, param.Key())
		}

		if param.Value() != value {
			t.Errorf("Parameter.Value() mismatch: expected %q, got %q", value, param.Value())
		}
	})
}

func FuzzMakeParameterWithVariousTypes(f *testing.F) {
	f.Add("int-key", int(42))
	f.Add("float-key", int(314))
	f.Add("bool-key", int(1))
	f.Add("nil-key", int(0))

	f.Fuzz(func(t *testing.T, key string, intVal int) {
		defer func() {
			if r := recover(); r != nil {
				t.Errorf("MakeParameter with various types panicked: key=%q, intVal=%d, panic: %v",
					key, intVal, r)
			}
		}()

		testCases := []interface{}{
			intVal,
			float64(intVal) / 100.0,
			intVal > 0,
			nil,
			[]string{"a", "b"},
			map[string]int{"test": intVal},
		}

		for i, value := range testCases {
			testParameterWithValue(t, key, value, i)
		}
	})
}

func testParameterWithValue(t *testing.T, key string, value interface{}, caseIndex int) {
	param := MakeParameter(key, value)

	if param == nil {
		t.Errorf("MakeParameter should return non-nil parameter for test case %d", caseIndex)
		return
	}

	if param.Key() != key {
		t.Errorf("Test case %d: Parameter.Key() mismatch: expected %q, got %q", caseIndex, key, param.Key())
	}

	validateParameterValue(t, param.Value(), value, caseIndex)
}

func validateParameterValue(t *testing.T, paramValue, expectedValue interface{}, caseIndex int) {
	switch v := expectedValue.(type) {
	case []string:
		validateSliceValue(t, paramValue, v, caseIndex)
	case map[string]int:
		validateMapValue(t, paramValue, v, caseIndex)
	default:
		if paramValue != expectedValue {
			t.Errorf("Test case %d: Parameter.Value() mismatch: expected %v, got %v", caseIndex, expectedValue, paramValue)
		}
	}
}

func validateSliceValue(t *testing.T, paramValue interface{}, expectedSlice []string, caseIndex int) {
	pv, ok := paramValue.([]string)
	if !ok || len(pv) != len(expectedSlice) {
		t.Errorf("Test case %d: Parameter.Value() slice mismatch: expected %v, got %v", caseIndex, expectedSlice, paramValue)
	}
}

func validateMapValue(t *testing.T, paramValue interface{}, expectedMap map[string]int, caseIndex int) {
	pv, ok := paramValue.(map[string]int)
	if !ok || len(pv) != len(expectedMap) {
		t.Errorf("Test case %d: Parameter.Value() map mismatch: expected %v, got %v", caseIndex, expectedMap, paramValue)
	}
}

func FuzzParameterString(f *testing.F) {
	f.Add("key", "value")
	f.Add("", "")
	f.Add("unicode-key", "unicode-value-test")
	f.Add("special", "value\nwith\ttabs\rand\"quotes")
	f.Add("long-key", strings.Repeat("x", 10000))
	f.Add("json-key", `{"already": "json"}`)
	f.Add("control-chars", "value\x00\x01\x02\x1F")

	f.Fuzz(func(t *testing.T, key, value string) {
		defer func() {
			if r := recover(); r != nil {
				t.Errorf("Parameter.String() panicked: key=%q, value=%q, panic: %v", key, value, r)
			}
		}()

		param := MakeParameter(key, value)
		result := param.String()

		if result == "" && (key != "" || value != "") {
			t.Errorf("Parameter.String() returned empty string for non-empty input: key=%q, value=%q",
				key, value)
		}

		if result == "" {
			return
		}

		var parsed Parameter
		err := json.Unmarshal([]byte(result), &parsed)
		if err != nil {
			t.Errorf("Parameter.String() produced invalid JSON: %q, error: %v", result, err)
			return
		}

		// Handle invalid UTF-8 sequences: Go's JSON marshaler replaces invalid UTF-8
		// with the Unicode replacement character (\ufffd)
		expectedKey := key
		expectedValue := value
		if !utf8.ValidString(key) {
			// For invalid UTF-8, expect replacement characters in the unmarshaled result
			expectedKey = string([]rune(key)) // This converts invalid bytes to replacement chars
		}
		if !utf8.ValidString(value) {
			expectedValue = string([]rune(value))
		}

		if parsed.K != expectedKey {
			t.Errorf("JSON key mismatch: expected %q, got %q", expectedKey, parsed.K)
		}
		if parsed.V != expectedValue {
			t.Errorf("JSON value mismatch: expected %q, got %q", expectedValue, parsed.V)
		}
	})
}

func validateComplexParameterJSON(t *testing.T, key string, value interface{}, result string, caseIndex int) {
	if result == "" {
		return
	}

	var jsonData map[string]interface{}
	err := json.Unmarshal([]byte(result), &jsonData)
	if err != nil {
		t.Errorf("Test case %d: Parameter.String() produced invalid JSON: %q, error: %v",
			caseIndex, result, err)
		return
	}

	// Handle invalid UTF-8 sequences: Go's JSON marshaler replaces invalid UTF-8
	// with the Unicode replacement character (\ufffd)
	expectedKey := key
	if !utf8.ValidString(key) {
		// For invalid UTF-8, expect replacement characters in the unmarshaled result
		expectedKey = string([]rune(key)) // This converts invalid bytes to replacement chars
	}

	if jsonData["key"] != expectedKey {
		t.Errorf("Test case %d: JSON key mismatch: expected %q, got %v", caseIndex, expectedKey, jsonData["key"])
	}

	if jsonData["value"] == nil && value != nil {
		t.Errorf("Test case %d: JSON value should not be nil when input is not nil", caseIndex)
	}
}

func FuzzParameterStringWithComplexValues(f *testing.F) {
	f.Add("array-key")
	f.Add("object-key")
	f.Add("number-key")
	f.Add("bool-key")

	f.Fuzz(func(t *testing.T, key string) {
		defer func() {
			if r := recover(); r != nil {
				t.Errorf("Parameter.String() with complex values panicked: key=%q, panic: %v", key, r)
			}
		}()

		testValues := []interface{}{
			[]string{"a", "b", "c"},
			map[string]interface{}{"nested": "value", "number": 42},
			123.456,
			true,
			false,
			nil,
			struct{ Field string }{"test"},
		}

		for i, value := range testValues {
			param := MakeParameter(key, value)
			result := param.String()
			validateComplexParameterJSON(t, key, value, result, i)
		}
	})
}

func FuzzParameterKeyValue(f *testing.F) {
	f.Add("test-key", "test-value")
	f.Add("", "")
	f.Add("unicode-test", "value-test")
	f.Add(strings.Repeat("k", 1000), strings.Repeat("v", 1000))

	f.Fuzz(func(t *testing.T, key, value string) {
		defer func() {
			if r := recover(); r != nil {
				t.Errorf("Parameter Key/Value methods panicked: key=%q, value=%q, panic: %v",
					key, value, r)
			}
		}()

		param := &Parameter{
			K: key,
			V: value,
		}

		if param.Key() != key {
			t.Errorf("Parameter.Key() mismatch: expected %q, got %q", key, param.Key())
		}

		if param.Value() != value {
			t.Errorf("Parameter.Value() mismatch: expected %q, got %q", value, param.Value())
		}

		var keyValue KeyValue = param
		if keyValue.Key() != key {
			t.Errorf("KeyValue.Key() mismatch: expected %q, got %q", key, keyValue.Key())
		}

		if keyValue.Value() != value {
			t.Errorf("KeyValue.Value() mismatch: expected %q, got %q", value, keyValue.Value())
		}
	})
}
