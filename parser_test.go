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
			want: &url.URL{
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
			flag := &gflag{flagType: tt.flagType}
			got, err := parseFlagValue(flag, tt.value)
			if (err != nil) != tt.wantErr {
				t.Errorf("parseFlagValue() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if err == nil && reflect.TypeOf(got) != reflect.TypeOf(tt.want) {
				t.Errorf("parseFlagValue() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestParseFlagValueError(t *testing.T) {
	tests := []struct {
		name  string
		flag  *gflag
		value string
	}{
		{
			name:  "unsupported flag type",
			flag:  &gflag{flagType: flagType(100)},
			value: "value",
		},
		{
			name:  "invalid int value",
			flag:  &gflag{flagType: flagInt},
			value: "invalid",
		},
		{
			name:  "invalid int64 value",
			flag:  &gflag{flagType: flagInt64},
			value: "invalid",
		},
		{
			name:  "invalid float32 value",
			flag:  &gflag{flagType: flagFloat32},
			value: "invalid",
		},
		{
			name:  "invalid float64 value",
			flag:  &gflag{flagType: flagFloat64},
			value: "invalid",
		},
		{
			name:  "invalid bool value",
			flag:  &gflag{flagType: flagBool},
			value: "invalid",
		},
		{
			name:  "invalid rune value",
			flag:  &gflag{flagType: flagRune},
			value: "invalid",
		},
		{
			name:  "invalid duration value",
			flag:  &gflag{flagType: flagDuration},
			value: "invalid",
		},
		{
			name:  "invalid time value",
			flag:  &gflag{flagType: flagTime},
			value: "invalid",
		},
		{
			name:  "invalid IP value",
			flag:  &gflag{flagType: flagIP},
			value: "invalid",
		},
		{
			name:  "invalid file path",
			flag:  &gflag{flagType: flagFilePath},
			value: "invalid",
		},
		{
			name:  "invalid directory path",
			flag:  &gflag{flagType: flagDirPath},
			value: "invalid",
		},
		{
			name:  "invalid URL",
			flag:  &gflag{flagType: flagURL},
			value: "invalid",
		},
		{
			name:  "invalid email",
			flag:  &gflag{flagType: flagEmail},
			value: "invalid",
		},
		{
			name:  "invalid host:port pair",
			flag:  &gflag{flagType: flagHostPortPair},
			value: "invalid",
		},
		{
			name:  "invalid MAC address",
			flag:  &gflag{flagType: flagMAC},
			value: "invalid",
		},
		{
			name:  "invalid UUID",
			flag:  &gflag{flagType: flagUUID},
			value: "invalid",
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			_, err := parseFlagValue(test.flag, test.value)
			if err == nil {
				t.Errorf("expected error, but got nil")
			}
		})
	}
}
