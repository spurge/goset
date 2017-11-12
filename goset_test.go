package goset

import (
	"fmt"
	"testing"
)

const TEST_JSON_FILE = "./test/test.json"
const TEST_NONJSON_FILE = "./test/non-json"

func testSettingsValue(t *testing.T, s *Settings, prop string, match interface{}) {
	val, err := s.Get(prop)

	if err != nil {
		t.Error(err)
	}

	if val != match {
		t.Error(fmt.Sprintf("%s != %s", val, match))
	}
}

func TestJsonLoad(t *testing.T) {
	settings := New()
	settings.Load(TEST_JSON_FILE)

	testSettingsValue(t, &settings, "test", "value")
	testSettingsValue(t, &settings, "some other.with", "other value")
}

func TestFailedIOLoad(t *testing.T) {
	settings := New()
	settings.Load("nonexisting")

	_, err := settings.Get("test")

	if err == nil {
		t.Error("Should've got an io error")
	}
}

func TestFailedJson(t *testing.T) {
	settings := New()
	settings.Load(TEST_NONJSON_FILE)
	_, err := settings.Get("test")

	if err == nil {
		t.Error("Should've got an invalid json error")
	}
}

func TestSet(t *testing.T) {
	values := map[string]interface{}{
		"some-setting": "with-a-value",
		"an integer":   4,
		"bool":         false,
		"nested": map[string]interface{}{
			"property": 2,
		},
	}

	settings := New()
	settings.Set(values)

	for _, prop := range []string{"some-setting", "an integer", "bool"} {
		testSettingsValue(t, &settings, prop, values[prop])
	}

	testSettingsValue(t, &settings, "nested.property", 2)
}

func TestMergedValues(t *testing.T) {
	values := map[string]interface{}{
		"test": "a value",
		"nested": map[string]interface{}{
			"property": 2,
		},
	}

	settings := New()
	settings.Load(TEST_JSON_FILE).Set(values)

	testSettingsValue(t, &settings, "some other.with", "other value")
	testSettingsValue(t, &settings, "nested.property", 2)
}
