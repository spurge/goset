package goset

import (
	"fmt"
	"testing"
)

func testValue(t *testing.T, value interface{}, expected interface{}) {
	if value != expected {
		t.Error(fmt.Sprintf("%s != %s", value, expected))
	}
}

func testData() map[string]interface{} {
	return map[string]interface{}{
		"value": "some-value",
		"another value": map[string]interface{}{
			"data": "something",
		},
	}
}

func TestParse(t *testing.T) {
	data, err := parse([]byte(`{"str":"rts","num":8.4}`))

	if err != nil {
		t.Error(err)
	}

	if data["str"] != "rts" || data["num"] != 8.4 {
		t.Error("Parsed data don't match")
	}
}

func TestParseInvalidJson(t *testing.T) {
	data, err := parse([]byte(`{"sd:333"},`))

	if data != nil {
		t.Error("Data should be nil")
	}

	if err == nil {
		t.Error("No error on invalid json")
	}
}

func TestMerge(t *testing.T) {
	current := testData()
	next := map[string]interface{}{
		"value": "value",
		"some other": map[string]interface{}{
			"ok": "yes",
		},
	}

	merged := merge(&current, &next)

	testValue(t, merged["value"], next["value"])
	testValue(t, merged["another value"].(map[string]interface{})["data"], "something")
	testValue(t, merged["some other"].(map[string]interface{})["ok"], "yes")
}

func TestValidSingleExtract(t *testing.T) {
	data := testData()

	val, err := extract(data, "value")

	if err != nil {
		t.Error(err)
	}

	testValue(t, val, "some-value")
}

func TestValidNestedExtract(t *testing.T) {
	data := testData()

	val, err := extract(data, "another value.data")

	if err != nil {
		t.Error(err)
	}

	testValue(t, val, "something")
}

func TestInvalidSingleExtract(t *testing.T) {
	data := map[string]interface{}{}

	val, err := extract(data, "nonexisting")

	if val != nil {
		t.Error("Value shoud be nil")
	}

	testValue(t, fmt.Sprint(err), "nonexisting not found")
}

func TestInvalidNestedExtract(t *testing.T) {
	data := map[string]interface{}{
		"another value": map[string]interface{}{},
	}

	val, err := extract(data, "another value.nonexisting")

	if val != nil {
		t.Error("Value shoud be nil")
	}

	testValue(t, fmt.Sprint(err), "nonexisting not found")
}

func TestInvalidTypeInNestedExtract(t *testing.T) {
	data := map[string]interface{}{
		"another value": 3,
	}

	val, err := extract(data, "another value.key")

	if val != nil {
		t.Error("Value shoud be nil")
	}

	testValue(t, fmt.Sprint(err), "another value is not a map")
}
