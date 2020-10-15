package main

import "os"

// CheckDirExists check whether dir exist
func checkDirExists(path string) (bool, error) {
	_, err := os.Stat(path)
	if err != nil {
		if !os.IsNotExist(err) {
			return false, nil
		}
		return false, err
	}
	return true, nil
}
