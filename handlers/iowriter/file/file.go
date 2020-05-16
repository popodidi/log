package file

import (
	"io"
	"strconv"
	"time"
)

// Single returns a single file writer.
func Single(filename string) io.WriteCloser {
	return newSingle(filename)
}

// Rotate returns a rotating file writer.
func Rotate(rotator Rotator) io.WriteCloser {
	return newRotate(rotator)
}

// PrefixSuffix returns a rotator with prefix and suffix
func PrefixSuffix(prefix, suffix string, rotator Rotator) Rotator {
	return &decorateNameRotator{
		rot: rotator,
		decorate: func(name string) string {
			return prefix + name + suffix
		},
	}
}

// Prefix returns a rotator with prefix.
func Prefix(prefix string, rotator Rotator) Rotator {
	return &decorateNameRotator{
		rot: rotator,
		decorate: func(name string) string {
			return prefix + name
		},
	}
}

// Suffix returns a rotator with suffix.
func Suffix(suffix string, rotator Rotator) Rotator {
	return &decorateNameRotator{
		rot: rotator,
		decorate: func(name string) string {
			return name + suffix
		},
	}
}

// SecondRotator returns a rotator that rotates every second with file name as
// {UNIX_SECOND}.{IDX} and file size limited to `size` bytes.
func SecondRotator(size int) Rotator {
	return &baseRotator{
		next: func() string {
			return strconv.Itoa(int(time.Now().Unix()))
		},
		size: size,
	}
}

const dateFmt = "20060102"

// DateRotator returns a rotator that rotates every day with file name as
// {YYYYMMDD}.{IDX} and file size limited to `size` bytes.
func DateRotator(size int) Rotator {
	return &baseRotator{
		next: func() string {
			return time.Now().Format(dateFmt)
		},
		size: size,
	}
}
