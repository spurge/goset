/*
Package goset is a settings library that provides async loading of settings
with a handy getter.

You can load and set with the same instance and it'll merge all your settings.
Event though all Load and Set actions are running async, they'll merge your
settings in the same order as you typed them.

Example of use:

	settings := Settings.New()
	// Set some default values
	settings.Set(map[string]interfaces{}{
		"some-key": map[string]interfaces{}{
			"hello": "world",
		},
	})
	// Then load some json data
	// For now, goset does only sopport JSON.
	settings.Load("some-file.json")
	// Fetch some values
	val, err := settings.Get("some-key.hello")

	// You can even chain them ...
	settings.Set(...).Load("file")
*/
package goset

import (
	"io/ioutil"
	"sync"
)

type Settings struct {
	sync.RWMutex
	values map[string]interface{}
	err    error
}

// Create a new goset settings instance
func New() Settings {
	return Settings{}
}

// Load local JSON file and merge it's values into the settings instance
func (s *Settings) Load(filename string) *Settings {
	s.Lock()

	go func() {
		defer s.Unlock()

		data, err := ioutil.ReadFile(filename)

		if err != nil {
			s.err = err
			return
		}

		values, err := parse(data)

		if err != nil {
			s.err = err
			return
		}

		s.mergeValues(values)
	}()

	return s
}

// Set and merge values into the settings instance
func (s *Settings) Set(values map[string]interface{}) *Settings {
	s.Lock()

	go func() {
		defer s.Unlock()

		s.mergeValues(values)
	}()

	return s
}

// Get values from your settings
func (s *Settings) Get(path string) (interface{}, error) {
	s.RLock()
	defer s.RUnlock()

	if s.err != nil {
		return nil, s.err
	}

	return extract(s.values, path)
}

func (s *Settings) mergeValues(values map[string]interface{}) {
	s.values = merge(&s.values, &values)
}
