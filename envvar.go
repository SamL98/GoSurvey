package main

import (
	"os"
	"strings"
)

func getEnv() map[string]string {
	envvars := make(map[string]string)
	for _, item := range os.Environ() {
		splits := strings.Split(item, "=")
		key := splits[0]
		val := splits[1]
		envvars[key] = val
	}
	return envvars
}
