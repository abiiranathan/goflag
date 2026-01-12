package goflag

import (
	"net"
	"net/url"
	"os"
	"path/filepath"
	"reflect"
	"testing"
	"time"

	"github.com/google/uuid"
)

func TestParseFlagValue(t *testing.T) {
	tests := []struct {
		name     string
		flagType flagType
		value    string
		want     any
		wantErr  bool
	}{
		{
			name:     "string",
			flagType: flagString,
			value:    "hello",
			want:     "hello",
			wantErr:  false,
		},
		{
			name:     "int",
			flagType: flagInt,
			value:    "42",
			want:     42,
			wantErr:  false,
		},
		{
			name:     "int64",
			flagType: flagInt64,
			value:    "9223372036854775807",
			want:     int64(9223372036854775807),
			wantErr:  false,
		},
		{
			name:     "float32",
			flagType: flagFloat32,
			value:    "3.14",
			want:     float32(3.14),
			wantErr:  false,
		},
		{
			name:     "float64",
			flagType: flagFloat64,
			value:    "3.14159265359",
			want:     3.14159265359,
			wantErr:  false,
		},
		{
			name:     "bool",
			flagType: flagBool,
			value:    "true",
			want:     true,
			wantErr:  false,
		},
		{
			name:     "string slice",
			flagType: flagStringSlice,
			value:    "hello,world",
			want:     []string{"hello", "world"},
			wantErr:  false,
		},
		{
			name:     "int slice",
			flagType: flagIntSlice,
			value:    "1,2,3",
			want:     []int{1, 2, 3},
			wantErr:  false,
		},
		{
			name:     "rune",
			flagType: flagRune,
			value:    "a",
			want:     'a',
			wantErr:  false,
		},
		{
			name:     "duration",
			flagType: flagDuration,
			value:    "1h30m",
			want:     time.Duration(1*time.Hour + 30*time.Minute),
			wantErr:  false,
		},
		{
			name:     "time",
			flagType: flagTime,
			value:    "2022-01-01T00:00 UTC",
			want:     time.Date(2022, 1, 1, 0, 0, 0, 0, time.UTC),
			wantErr:  false,
		},
		{
			name:     "IP",
			flagType: flagIP,
			value:    "127.0.0.1",
			want:     net.ParseIP("127.0.0.1"),
			wantErr:  false,
		},
		{
			name:     "file path",
			flagType: flagFilePath,
			value:    "parser_test.go",
			want:     filepath.Join(os.Getenv("PWD"), "parser_test.go"),
			wantErr:  false,
		},
		{
			name:     "dir path",
			flagType: flagDirPath,
			value:    filepath.Join(os.Getenv("PWD")),
			want:     filepath.Join(os.Getenv("PWD")),
			wantErr:  false,
		},
		{
			name:     "dir path",
			flagType: flagDirPath,
			value:    filepath.Join(os.Getenv("PWD"), "parser_test.go"),
			want:     "",
			wantErr:  true,
		},
		{
			name:     "email",
			flagType: flagEmail,
			value:    "test@example.com",
			want:     "test@example.com",
			wantErr:  false,
		},
		{
			name:     "URL",
			flagType: flagURL,
			value:    "https://example.com",
			want: url.URL{
				Scheme: "https",
				Host:   "example.com",
			},
			wantErr: false,
		},
		{
			name:     "UUID",
			flagType: flagUUID,
			value:    "123e4567-e89b-12d3-a456-426655440000",
			want:     uuid.MustParse("123e4567-e89b-12d3-a456-426655440000"),
			wantErr:  false,
		},
		{
			name:     "host port pair",
			flagType: flagHostPortPair,
			value:    "localhost:8080",
			want:     "localhost:8080",
			wantErr:  false,
		},
		{
			name:     "MAC",
			flagType: flagMAC,
			value:    "00:11:22:33:44:55",
			want:     net.HardwareAddr{0x00, 0x11, 0x22, 0x33, 0x44, 0x55},
			wantErr:  false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			flag := &Flag{flagType: tt.flagType, value: reflect.New(reflect.TypeOf(tt.want)).Interface(), required: true}
			err := parseFlagValue(flag, tt.value)
			if (err != nil) != tt.wantErr {
				t.Errorf("parseFlagValue() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			// dereference the pointer
			got := reflect.ValueOf(flag.value).Elem().Interface()
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("parseFlagValue() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestParseFlagWithInvalidValue(t *testing.T) {
	tests := []struct {
		name     string
		flagType flagType
		value    string
		want     any
		wantErr  bool
	}{
		{
			name:     "string",
			flagType: flagInt,
			value:    "hello",
			want:     "hello",
			wantErr:  true,
		},
		{
			name:     "int",
			flagType: flagInt,
			value:    "integer",
			want:     42,
			wantErr:  true,
		},
		{
			name:     "int64",
			flagType: flagInt64,
			value:    "abcd",
			want:     int64(9223372036854775807),
			wantErr:  true,
		},
		{
			name:     "float32",
			flagType: flagFloat32,
			value:    "abcd",
			want:     float32(3.14),
			wantErr:  true,
		},
		{
			name:     "float64",
			flagType: flagFloat64,
			value:    "abcd",
			want:     3.14159265359,
			wantErr:  true,
		},
		{
			name:     "bool",
			flagType: flagBool,
			value:    "true_false",
			want:     true,
			wantErr:  true,
		},

		{
			name:     "rune",
			flagType: flagRune,
			value:    "100",
			want:     'a',
			wantErr:  true,
		},
		{
			name:     "duration",
			flagType: flagDuration,
			value:    "10tmh_10",
			want:     time.Duration(1*time.Hour + 30*time.Minute),
			wantErr:  true,
		},
		{
			name:     "time",
			flagType: flagTime,
			value:    "2022-01-01T00:00 UTCCC AM",
			want:     time.Date(2022, 1, 1, 0, 0, 0, 0, time.UTC),
			wantErr:  true,
		},
		{
			name:     "IP",
			flagType: flagIP,
			value:    "127.0.0.1.6",
			want:     net.ParseIP("127.0.0.1"),
			wantErr:  true,
		},
		{
			name:     "file path",
			flagType: flagFilePath,
			value:    "parser_test_error.go",
			want:     filepath.Join(os.Getenv("PWD"), "parser_test.go"),
			wantErr:  true,
		},
		{
			name:     "dir path",
			flagType: flagDirPath,
			value:    "some path",
			want:     filepath.Join(os.Getenv("PWD")),
			wantErr:  true,
		},
		{
			name:     "email",
			flagType: flagEmail,
			value:    "string_email",
			want:     "string_email",
			wantErr:  true,
		},
		{
			name:     "URL",
			flagType: flagURL,
			value:    "random_string",
			want: url.URL{
				Scheme: "https",
				Host:   "example.com",
			},
			wantErr: true,
		},
		{
			name:     "UUID",
			flagType: flagUUID,
			value:    "123e4567-e89b-12d3-a456",
			want:     uuid.MustParse("123e4567-e89b-12d3-a456-426655440000"),
			wantErr:  true,
		},
		{
			name:     "host port pair",
			flagType: flagHostPortPair,
			value:    "localhost8080",
			want:     "localhost:8080",
			wantErr:  true,
		},
		{
			name:     "host port pair - Invalid port",
			flagType: flagHostPortPair,
			value:    "localhost:invalidport",
			want:     "0",
			wantErr:  true,
		},
		{
			name:     "host port pair - port out of range",
			flagType: flagHostPortPair,
			value:    "localhost:1234567890",
			want:     "0",
			wantErr:  true,
		},
		{
			name:     "MAC",
			flagType: flagMAC,
			value:    "00:11:22:33:44:55:100:566:49 12:234",
			want:     net.HardwareAddr{0x00, 0x11, 0x22, 0x33, 0x44, 0x55},
			wantErr:  true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			flag := &Flag{flagType: tt.flagType, value: reflect.New(reflect.TypeOf(tt.want)).Interface(), required: true}
			err := parseFlagValue(flag, tt.value)
			if err == nil {
				t.Errorf("parseFlagValue(%s) error expected to fail but got no error", flag.name)
				return
			}

			// dereference the pointer
			got := reflect.ValueOf(flag.value).Elem().Interface()
			if reflect.DeepEqual(got, tt.want) {
				t.Errorf("parseFlagValue() = %v should not equal %v", got, tt.want)
			}
		})
	}
}
