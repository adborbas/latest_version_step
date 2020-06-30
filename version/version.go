package version

import (
	"strconv"
	"strings"
)

// Version ...
type Version struct {
	major, minor, patch int
}

// New ...
func New(raw string) *Version {
	values := strings.Split(raw, ".")

	if len(values) != 3 {
		return nil
	}

	major, err := strconv.Atoi(values[0])
	minor, err := strconv.Atoi(values[1])
	patch, err := strconv.Atoi(values[2])
	if err != nil {
		return nil
	}

	return &Version{
		major: major,
		minor: minor,
		patch: patch}
}
