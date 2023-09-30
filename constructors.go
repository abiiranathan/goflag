// Code generated by genc. DO NOT EDIT.

package goflag

import (
	"github.com/google/uuid"
	"net"
	"net/url"
	"time"
)

func Bool(name, shortName string, value bool, usage string, required ...bool) *gflag {
	flag := &gflag{
		name:      name,
		shortName: shortName,
		value:     value,
		flagType:  flagBool,
		usage:     usage,
	}
	if len(required) > 0 {
		flag.required = required[0]
	}
	return flag
}

func DirPath(name, shortName string, value string, usage string, required ...bool) *gflag {
	flag := &gflag{
		name:      name,
		shortName: shortName,
		value:     value,
		flagType:  flagDirPath,
		usage:     usage,
	}
	if len(required) > 0 {
		flag.required = required[0]
	}
	return flag
}

func Duration(name, shortName string, value time.Duration, usage string, required ...bool) *gflag {
	flag := &gflag{
		name:      name,
		shortName: shortName,
		value:     value,
		flagType:  flagDuration,
		usage:     usage,
	}
	if len(required) > 0 {
		flag.required = required[0]
	}
	return flag
}

func Email(name, shortName string, value string, usage string, required ...bool) *gflag {
	flag := &gflag{
		name:      name,
		shortName: shortName,
		value:     value,
		flagType:  flagEmail,
		usage:     usage,
	}
	if len(required) > 0 {
		flag.required = required[0]
	}
	return flag
}

func FilePath(name, shortName string, value string, usage string, required ...bool) *gflag {
	flag := &gflag{
		name:      name,
		shortName: shortName,
		value:     value,
		flagType:  flagFilePath,
		usage:     usage,
	}
	if len(required) > 0 {
		flag.required = required[0]
	}
	return flag
}

func Float32(name, shortName string, value float32, usage string, required ...bool) *gflag {
	flag := &gflag{
		name:      name,
		shortName: shortName,
		value:     value,
		flagType:  flagFloat32,
		usage:     usage,
	}
	if len(required) > 0 {
		flag.required = required[0]
	}
	return flag
}

func Float64(name, shortName string, value float64, usage string, required ...bool) *gflag {
	flag := &gflag{
		name:      name,
		shortName: shortName,
		value:     value,
		flagType:  flagFloat64,
		usage:     usage,
	}
	if len(required) > 0 {
		flag.required = required[0]
	}
	return flag
}

func HostPortPair(name, shortName string, value string, usage string, required ...bool) *gflag {
	flag := &gflag{
		name:      name,
		shortName: shortName,
		value:     value,
		flagType:  flagHostPortPair,
		usage:     usage,
	}
	if len(required) > 0 {
		flag.required = required[0]
	}
	return flag
}

func IP(name, shortName string, value net.IP, usage string, required ...bool) *gflag {
	flag := &gflag{
		name:      name,
		shortName: shortName,
		value:     value,
		flagType:  flagIP,
		usage:     usage,
	}
	if len(required) > 0 {
		flag.required = required[0]
	}
	return flag
}

func Int(name, shortName string, value int, usage string, required ...bool) *gflag {
	flag := &gflag{
		name:      name,
		shortName: shortName,
		value:     value,
		flagType:  flagInt,
		usage:     usage,
	}
	if len(required) > 0 {
		flag.required = required[0]
	}
	return flag
}

func Int64(name, shortName string, value int64, usage string, required ...bool) *gflag {
	flag := &gflag{
		name:      name,
		shortName: shortName,
		value:     value,
		flagType:  flagInt64,
		usage:     usage,
	}
	if len(required) > 0 {
		flag.required = required[0]
	}
	return flag
}

func IntSlice(name, shortName string, value []int, usage string, required ...bool) *gflag {
	flag := &gflag{
		name:      name,
		shortName: shortName,
		value:     value,
		flagType:  flagIntSlice,
		usage:     usage,
	}
	if len(required) > 0 {
		flag.required = required[0]
	}
	return flag
}

func MAC(name, shortName string, value net.HardwareAddr, usage string, required ...bool) *gflag {
	flag := &gflag{
		name:      name,
		shortName: shortName,
		value:     value,
		flagType:  flagMAC,
		usage:     usage,
	}
	if len(required) > 0 {
		flag.required = required[0]
	}
	return flag
}

func Rune(name, shortName string, value rune, usage string, required ...bool) *gflag {
	flag := &gflag{
		name:      name,
		shortName: shortName,
		value:     value,
		flagType:  flagRune,
		usage:     usage,
	}
	if len(required) > 0 {
		flag.required = required[0]
	}
	return flag
}

func String(name, shortName string, value string, usage string, required ...bool) *gflag {
	flag := &gflag{
		name:      name,
		shortName: shortName,
		value:     value,
		flagType:  flagString,
		usage:     usage,
	}
	if len(required) > 0 {
		flag.required = required[0]
	}
	return flag
}

func StringSlice(name, shortName string, value []string, usage string, required ...bool) *gflag {
	flag := &gflag{
		name:      name,
		shortName: shortName,
		value:     value,
		flagType:  flagStringSlice,
		usage:     usage,
	}
	if len(required) > 0 {
		flag.required = required[0]
	}
	return flag
}

func Time(name, shortName string, value time.Time, usage string, required ...bool) *gflag {
	flag := &gflag{
		name:      name,
		shortName: shortName,
		value:     value,
		flagType:  flagTime,
		usage:     usage,
	}
	if len(required) > 0 {
		flag.required = required[0]
	}
	return flag
}

func URL(name, shortName string, value *url.URL, usage string, required ...bool) *gflag {
	flag := &gflag{
		name:      name,
		shortName: shortName,
		value:     value,
		flagType:  flagURL,
		usage:     usage,
	}
	if len(required) > 0 {
		flag.required = required[0]
	}
	return flag
}

func UUID(name, shortName string, value uuid.UUID, usage string, required ...bool) *gflag {
	flag := &gflag{
		name:      name,
		shortName: shortName,
		value:     value,
		flagType:  flagUUID,
		usage:     usage,
	}
	if len(required) > 0 {
		flag.required = required[0]
	}
	return flag
}
