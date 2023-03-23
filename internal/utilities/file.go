package utilities

import "os"

func ReadFile(path string) ([]byte, error) {
	f, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	return f, nil
}
