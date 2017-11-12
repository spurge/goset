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
)

type Settings struct {
	inVal  []chan map[string]interface{}
	inErr  chan error
	values map[string]interface{}
	err    error
	wait   bool
	done   chan bool
}

// Create a new goset settings instance
func New() Settings {
	return Settings{
		inErr: make(chan error),
	}
}

// Load local JSON file and merge it's values into the settings instance
func (s *Settings) Load(filename string) *Settings {
	go s.load(filename, s.getChan())
	return s
}

// Set and merge values into the settings instance
func (s *Settings) Set(values map[string]interface{}) *Settings {
	go s.set(values, s.getChan())
	return s
}

// Get values from your settings
func (s *Settings) Get(path string) (interface{}, error) {
	for _, inVal := range s.inVal {
		select {
		case values := <-inVal:
			s.values = merge(&s.values, &values)
		case err := <-s.inErr:
			s.err = err
		}
	}

	s.inVal = make([]chan map[string]interface{}, 0)

	if s.err != nil {
		return nil, s.err
	}

	return extract(s.values, path)
}

func (s *Settings) getChan() chan map[string]interface{} {
	inVal := make([]chan map[string]interface{}, len(s.inVal)+1)
	copy(inVal, s.inVal)
	s.inVal = inVal

	index := len(s.inVal) - 1

	if index < 0 {
		index = 0
	}

	s.inVal[index] = make(chan map[string]interface{})

	return s.inVal[index]
}

func (s *Settings) load(filename string, inVal chan map[string]interface{}) {
	data, err := ioutil.ReadFile(filename)

	if err != nil {
		s.inErr <- err
	}

	values, err := parse(data)

	if err != nil {
		s.inErr <- err
		return
	}

	s.set(values, inVal)
}

func (s *Settings) set(values map[string]interface{}, inVal chan map[string]interface{}) {
	inVal <- values
}
