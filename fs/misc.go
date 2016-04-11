package fs;

import "os"

func FileExists(path string) (bool, error) {
	_, err := os.Stat(path)
	if (err == nil) {
		return true, nil
	} else if (os.IsNotExist(err)) {
		return false, nil
	}
	return true, err
}
