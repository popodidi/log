package file

import "io"

// Single returns a single file writer.
func Single(filename string) (io.WriteCloser, error) {
	return newSingle(filename)
}

// Rotate returns a rotating file writer.
func Rotate(rotator Rotator) io.WriteCloser {
	return newRotate(rotator)
}
