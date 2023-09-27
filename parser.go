package goflag

import (
	"fmt"
	"net"
	"net/mail"
	"net/url"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"github.com/google/uuid"
)

// Returns the converted value from string based on flag type.
func parseFlagValue(flag *gflag, value string) (any, error) {
	switch flag.flagType {
	case flagString:
		return value, nil
	case flagInt:
		return ParseInt(value)
	case flagInt64:
		return ParseInt64(value)
	case flagFloat32:
		return ParseFloat32(value)
	case flagFloat64:
		return ParseFloat64(value)
	case flagBool:
		return ParseBool(value)
	case flagStringSlice:
		return ParseStringSlice(value)
	case flagIntSlice:
		return ParseIntSlice(value)
	case flagRune:
		return ParseRune(value)
	case flagDuration:
		return ParseDuration(value)
	case flagTime:
		return ParseTime(value)
	case flagIP:
		return ParseIP(value)
	case flagFilePath:
		return ParseFilePath(value)
	case flagDirPath:
		return ParseDirPath(value)
	case flagEmail:
		return ParseEmail(value)
	case flagURL:
		return ParseUrl(value)
	case flagUUID:
		return ParseUUID(value)
	case flagHostPortPair:
		return ParseHostPort(value)
	case flagMAC:
		return ParseMAC(value)
	}
	return nil, fmt.Errorf("unsupported flag type %s", flag.flagType.String())
}

// Parse a string to an int.
func ParseInt(value string) (int, error) {
	result, err := strconv.Atoi(value)
	if err != nil {
		return 0, fmt.Errorf("invalid int value %s", value)
	}
	return result, nil
}

// Parse a string to an int64.
func ParseInt64(value string) (int64, error) {
	result, err := strconv.ParseInt(value, 10, 64)
	if err != nil {
		return 0, fmt.Errorf("invalid int64 value %s", value)
	}
	return result, nil
}

// Parse a string to a float32.
func ParseFloat32(value string) (float32, error) {
	result, err := strconv.ParseFloat(value, 32)
	if err != nil {
		return 0.0, fmt.Errorf("invalid float32 value for flag %s", value)
	}
	return float32(result), nil
}

// Parse a string to a float64.
func ParseFloat64(value string) (float64, error) {
	result, err := strconv.ParseFloat(value, 64)
	if err != nil {
		return 0.0, fmt.Errorf("invalid float64 value for flag %s", value)
	}
	return result, nil
}

// Parse a string to a bool.
func ParseBool(value string) (bool, error) {
	result, err := strconv.ParseBool(value)
	if err != nil {
		return false, fmt.Errorf("invalid bool value for flag %s", value)
	}
	return result, nil
}

// Parse a comma-seperated string into a slice of strings.
func ParseStringSlice(value string) ([]string, error) {
	parts := strings.Split(value, ",")
	result := make([]string, len(parts))
	for index := range parts {
		result[index] = strings.TrimSpace(parts[index])
	}
	return result, nil
}

// Parse a comma-seperated string into a slice of ints.
func ParseIntSlice(value string) ([]int, error) {
	parts := strings.Split(value, ",")
	result := make([]int, len(parts))
	for index := range parts {
		intvalue, err := ParseInt(strings.TrimSpace(parts[index]))
		if err != nil {
			return nil, err
		}
		result[index] = intvalue

	}
	return result, nil
}

// Parse a string to a rune.
func ParseRune(value string) (rune, error) {
	if len(value) != 1 {
		return ' ', fmt.Errorf("expected one character")
	}
	return rune(value[0]), nil
}

// Parse a string to a duration.
// Uses time.ParseDuration. Supported units are "ns", "us" (or "µs"), "ms", "s", "m", "h".
// e.g 1h30m, 1h, 1m30s, 1m, 1m30s, 1ms, 1us, 1ns
func ParseDuration(value string) (time.Duration, error) {
	duration, err := time.ParseDuration(value)
	if err != nil {
		return 0, fmt.Errorf("invalid duration value for flag %s", value)
	}
	return duration, nil
}

// Parse a string to a time.Time.
// Uses time.Parse. Supported formats are:
// "2006-01-02T15:04 MST"
func ParseTime(value string) (time.Time, error) {
	layout := "2006-01-02T15:04 MST"
	t, err := time.Parse(layout, value)
	if err != nil {
		return t, fmt.Errorf("invalid time value for flag %s. Suppoted format: 2006-01-02T15:04 MST", value)
	}
	return t, nil
}

func ParseIP(value string) (net.IP, error) {
	ip := net.ParseIP(value)
	if ip == nil {
		return nil, fmt.Errorf("%s is not a valid IP address", value)
	}
	return ip, nil
}

// Resolve absolute file path and check that it exists.
func ParseFilePath(value string) (string, error) {
	filePath, err := filepath.Abs(value)
	if err != nil {
		return "", fmt.Errorf("unable to find absolute path to " + value)
	}

	f, err := os.Stat(filePath)
	if err != nil {
		return "", fmt.Errorf("can not stat: %s", err)
	}

	if f.IsDir() {
		return "", fmt.Errorf("%s is not a regular file", value)
	}
	return filePath, nil
}

// Resolve dirname from value and check that it exists.
func ParseDirPath(value string) (string, error) {
	filePath, err := filepath.Abs(value)
	if err != nil {
		return "", fmt.Errorf("unable to find absolute path to " + value)
	}

	f, err := os.Stat(filePath)
	if err != nil {
		return "", fmt.Errorf("can not stat: %s", err)
	}

	if !f.IsDir() {
		return "", fmt.Errorf("%s is not a directory", value)
	}
	return filePath, nil
}

// parse url from string with url.Parse.
func ParseUrl(value string) (*url.URL, error) {
	return url.ParseRequestURI(value)
}

// parse email from string with mail.Parse
func ParseEmail(value string) (string, error) {
	email, err := mail.ParseAddress(value)
	if err != nil {
		return "", fmt.Errorf("unable to parse email: %s", err)
	}
	return email.Address, nil
}

// parse host:port pair from value
// An empty string is considered a valid host. :)
// e.g ":8000" is a valid host-port pair.
func ParseHostPort(value string) (string, error) {
	hostPortPair := strings.SplitN(value, ":", 2)
	if len(hostPortPair) != 2 {
		return "", fmt.Errorf("invalid host:port pair")
	}

	// host := hostPortPair[0]
	port := hostPortPair[1]

	// make sure port is valid.
	portInt, err := strconv.Atoi(port)
	if err != nil {
		return "", fmt.Errorf("%s is not a valid port", port)
	}

	// Make sure port is in range.
	const maxPort int = (1 << 16) - 1
	if portInt < 0 || portInt > maxPort {
		return "", fmt.Errorf("port %s is out of range", port)
	}
	return value, nil
}

func ParseMAC(value string) (net.HardwareAddr, error) {
	mac, err := net.ParseMAC(value)
	if err != nil {
		return nil, fmt.Errorf("inavlid MAC address: %s", err)
	}
	return mac, nil
}

// parse UUID using the github.com/google/uuid package.
func ParseUUID(value string) (uuid.UUID, error) {
	id, err := uuid.Parse(value)
	if err != nil {
		return uuid.Nil, fmt.Errorf("invalid UUID: %s", err)
	}
	return id, nil
}
