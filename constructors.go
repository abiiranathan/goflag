// Code generated by genc. DO NOT EDIT.

package goflag

import (
	"github.com/google/uuid"
	"net"
	"net/url"
	"time"
)

func Bool(name, shortName string, value bool, usage string, required bool) *gflag {
	return &gflag{
		name:      name,
		shortName: shortName,
		value:     value,
		flagType:  flagBool,
		usage:     usage,
		required:  required,
	}
}

func DirPath(name, shortName string, value string, usage string, required bool) *gflag {
	return &gflag{
		name:      name,
		shortName: shortName,
		value:     value,
		flagType:  flagDirPath,
		usage:     usage,
		required:  required,
	}
}

func Duration(name, shortName string, value time.Duration, usage string, required bool) *gflag {
	return &gflag{
		name:      name,
		shortName: shortName,
		value:     value,
		flagType:  flagDuration,
		usage:     usage,
		required:  required,
	}
}

func Email(name, shortName string, value string, usage string, required bool) *gflag {
	return &gflag{
		name:      name,
		shortName: shortName,
		value:     value,
		flagType:  flagEmail,
		usage:     usage,
		required:  required,
	}
}

func FilePath(name, shortName string, value string, usage string, required bool) *gflag {
	return &gflag{
		name:      name,
		shortName: shortName,
		value:     value,
		flagType:  flagFilePath,
		usage:     usage,
		required:  required,
	}
}

func Float32(name, shortName string, value float32, usage string, required bool) *gflag {
	return &gflag{
		name:      name,
		shortName: shortName,
		value:     value,
		flagType:  flagFloat32,
		usage:     usage,
		required:  required,
	}
}

func Float64(name, shortName string, value float64, usage string, required bool) *gflag {
	return &gflag{
		name:      name,
		shortName: shortName,
		value:     value,
		flagType:  flagFloat64,
		usage:     usage,
		required:  required,
	}
}

func HostPortPair(name, shortName string, value string, usage string, required bool) *gflag {
	return &gflag{
		name:      name,
		shortName: shortName,
		value:     value,
		flagType:  flagHostPortPair,
		usage:     usage,
		required:  required,
	}
}

func IP(name, shortName string, value net.IP, usage string, required bool) *gflag {
	return &gflag{
		name:      name,
		shortName: shortName,
		value:     value,
		flagType:  flagIP,
		usage:     usage,
		required:  required,
	}
}

func Int(name, shortName string, value int, usage string, required bool) *gflag {
	return &gflag{
		name:      name,
		shortName: shortName,
		value:     value,
		flagType:  flagInt,
		usage:     usage,
		required:  required,
	}
}

func Int64(name, shortName string, value int64, usage string, required bool) *gflag {
	return &gflag{
		name:      name,
		shortName: shortName,
		value:     value,
		flagType:  flagInt64,
		usage:     usage,
		required:  required,
	}
}

func IntSlice(name, shortName string, value []int, usage string, required bool) *gflag {
	return &gflag{
		name:      name,
		shortName: shortName,
		value:     value,
		flagType:  flagIntSlice,
		usage:     usage,
		required:  required,
	}
}

func MAC(name, shortName string, value net.HardwareAddr, usage string, required bool) *gflag {
	return &gflag{
		name:      name,
		shortName: shortName,
		value:     value,
		flagType:  flagMAC,
		usage:     usage,
		required:  required,
	}
}

func Rune(name, shortName string, value rune, usage string, required bool) *gflag {
	return &gflag{
		name:      name,
		shortName: shortName,
		value:     value,
		flagType:  flagRune,
		usage:     usage,
		required:  required,
	}
}

func String(name, shortName string, value string, usage string, required bool) *gflag {
	return &gflag{
		name:      name,
		shortName: shortName,
		value:     value,
		flagType:  flagString,
		usage:     usage,
		required:  required,
	}
}

func StringSlice(name, shortName string, value []string, usage string, required bool) *gflag {
	return &gflag{
		name:      name,
		shortName: shortName,
		value:     value,
		flagType:  flagStringSlice,
		usage:     usage,
		required:  required,
	}
}

func Time(name, shortName string, value time.Time, usage string, required bool) *gflag {
	return &gflag{
		name:      name,
		shortName: shortName,
		value:     value,
		flagType:  flagTime,
		usage:     usage,
		required:  required,
	}
}

func URL(name, shortName string, value *url.URL, usage string, required bool) *gflag {
	return &gflag{
		name:      name,
		shortName: shortName,
		value:     value,
		flagType:  flagURL,
		usage:     usage,
		required:  required,
	}
}

func UUID(name, shortName string, value uuid.UUID, usage string, required bool) *gflag {
	return &gflag{
		name:      name,
		shortName: shortName,
		value:     value,
		flagType:  flagUUID,
		usage:     usage,
		required:  required,
	}
}
