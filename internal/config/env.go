package config

import (
	"fmt"
	"net/url"
	"os"
	"strconv"
	"time"
)

func envString(key, fallback string) string {
	value := os.Getenv(key)
	if value == "" {
		return fallback
	}
	return value
}

func envDuration(key string, fallback time.Duration) time.Duration {
	value := os.Getenv(key)
	if value == "" {
		return fallback
	}
	durationVal, err := time.ParseDuration(value)
	if err != nil {
		panic(fmt.Sprintf("unable to convert %v to a duration value", value))
	}
	return durationVal
}

func envURL(key, fallback string) *url.URL {
	value := os.Getenv(key)
	if value == "" {
		value = fallback
	}

	u, err := url.Parse(value)
	if err != nil {
		panic(fmt.Sprintf("unable to convert %v to a url value", value))
	}

	return u
}

func envBool(key string, fallback bool) bool {
	value := os.Getenv(key)
	if value == "" {
		return fallback
	}
	boolVal, err := strconv.ParseBool(value)
	if err != nil {
		panic(fmt.Sprintf("unable to convert %v to a boolean value", value))
	}
	return boolVal
}

func envInt(key string, fallback int) int {
	value := os.Getenv(key)
	if value == "" {
		return fallback
	}
	intVal, err := strconv.Atoi(value)
	if err != nil {
		panic(fmt.Sprintf("unable to convert %v to an integer value", value))
	}
	return intVal
}
