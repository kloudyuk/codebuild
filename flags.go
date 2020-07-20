package main

import (
	"fmt"
	"os"
	"strings"

	flag "github.com/spf13/pflag"
)

// Flag Vars
var computeType string
var envVars env = make(map[string]string)
var roleARN string
var serviceRole string
var sourceLocation string
var sourceType string
var sourceVersion string
var follow bool
var wait bool

// Custom Flags

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

func (e env) Type() string {
	return "NAME=VALUE"
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

// ParseFlags parses the command line flags and populates flags from env if possible
func parseFlags() {
	flag.Parse()
	flagsFromEnv()
	// Ensure wait is true if follow is true
	if follow {
		wait = true
	}
}

func flagsFromEnv() {
	flag.VisitAll(func(f *flag.Flag) {
		// skip the flag for environment variables
		if f.Name == "env" {
			return
		}
		// skip flags set by user
		if f.Changed {
			return
		}
		// convert the flag name to the environment variable equivalent e.g. my-flag -> MY_FLAG
		envName := strings.ToUpper(strings.ReplaceAll(f.Name, "-", "_"))
		// if the environment variable exists and has a non-empty value, set the flag
		v := os.Getenv(envName)
		if v != "" {
			f.Value.Set(v)
		}
	})
	// parse the flags again to pick up any updated values
	flag.Parse()
}
