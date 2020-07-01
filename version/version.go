package version

import (
	"fmt"
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

// IsNewer ...
func (version Version) IsNewer(other Version) bool {
	if version.major > other.major {
		return true
	}

	if version.minor > other.minor {
		return true
	}

	if version.patch > other.patch {
		return true
	}

	return false
}

func (version Version) String() string {
	return fmt.Sprintf("%d.%d.%d", version.major, version.minor, version.patch)
}
