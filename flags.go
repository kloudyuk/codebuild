package main

import (
	"fmt"
	"strings"
)

// CUSTOM FLAGS

// env is a slice of kv pairs each representing an environment variable
// it implements the flag.Value interface
type env map[string]string

// String returns a string representation of the env type
// This is required to satisfy the flag.Value interface
func (e env) String() string {
	var values []string
	for k, v := range e {
		values = append(values, fmt.Sprintf("%s=%s", k, v))
	}
	return strings.Join(values, ",")
}

// Set extracts the key/value pair from the flags raw string and stores them in the underlying map
// This is required to satisfy the flag.Value interface
func (e env) Set(raw string) error {
	vars := strings.Split(raw, ",")
	for _, v := range vars {
		parts := strings.Split(v, "=")
		e[parts[0]] = parts[1]
	}
	return nil
}
