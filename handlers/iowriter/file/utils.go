package file

import "os"

// exist returns if the file exists
func exist(path string) (exist bool, err error) {
	_, err = os.Stat(path)
	if err == nil {
		exist = true
		return
	}
	if os.IsNotExist(err) {
		exist = false
		err = nil
		return
	}
	return
}
