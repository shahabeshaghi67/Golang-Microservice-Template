package testutils

import (
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"os"
	"path"
	"runtime"
)

// LoadJSON loads a json file from the res directory and returns the json as string. If target is not nil the
// json will be unmarshaled into the target parameter. LoadJSON panics on errors.
func LoadJSON(target interface{}, pathElem ...string) string {
	file := resFile(pathElem...)

	bytes, err := os.ReadFile(file)
	if err != nil {
		panic(fmt.Sprintf("unable to read test data file %v: %s", file, err))
	}

	if target != nil {
		UnmarshalJSON(bytes, target)
	}

	return string(bytes)
}

// SaveJSON saves a data structure as json to a file in the res directory. SaveJSON panics on errors.
func SaveJSON(data interface{}, pathElem ...string) {
	encoded, err := json.MarshalIndent(data, "", "  ")
	if err != nil {
		panic(fmt.Sprintf("unable to marshal test data: %s", err))
	}

	encoded = append(encoded, byte('\n'))
	file := resFile(pathElem...)

	err = os.WriteFile(file, encoded, os.FileMode(0600))
	if err != nil {
		panic(fmt.Sprintf("unable to write test data file %v: %s", file, err))
	}
}

// MarshalJSON simple helper that marshals data into a JSON string.
func MarshalJSON(data interface{}) string {
	encoded, err := json.MarshalIndent(data, "", "  ")
	if err != nil {
		panic(fmt.Sprintf("unable to marshal test data: %s", err))
	}

	return string(encoded)
}

// UnmarshalJSON is a simple helper that unmarshals data from JSON. The input can be a string, []byte, or
// an io.Reader.
func UnmarshalJSON(input interface{}, model interface{}) {
	var (
		data []byte
		err  error
	)

	switch input := input.(type) {
	case []byte:
		data = input
	case string:
		data = []byte(input)
	case io.Reader:
		data, err = io.ReadAll(input)
		if err != nil {
			panic(fmt.Sprintf("unable read json input: %s", err))
		}
	}

	err = json.Unmarshal(data, model)
	if err != nil {
		panic(fmt.Sprintf("unable to unmarshal test data: %s", err))
	}
}

// resFile returns a file from the resource directory.
func resFile(pathElem ...string) string {
	_, thisFile, _, _ := runtime.Caller(0)
	parentDir, _ := path.Split(path.Dir(thisFile))
	return path.Join(parentDir, "..", "res", path.Join(pathElem...))
}

// ShouldKeepTestDBContainer returns true if the `keep-test-db-container` flag
// is set. Otherwise, it returns false. The `keep-test-db-container` flag
// indicates the database docker container started by the tests should keep
// running afterward. Otherwise, the tests stop and remove the container when
// the tests are finished.
func ShouldKeepTestDBContainer() bool {
	return keepTestDBContainerFlag != nil && *keepTestDBContainerFlag
}

var keepTestDBContainerFlag = flag.Bool("keep-test-db-container", false, "keep test database docker container")
