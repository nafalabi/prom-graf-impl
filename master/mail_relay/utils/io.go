package utils

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
)

func Read(path string) (data []byte, err error) {
	file, err := os.Open(path)
	if err != nil {
		return data, fmt.Errorf("Error opening file: %v", err.Error())
	}
	defer file.Close()

	bytes, err := io.ReadAll(file)
	if err != nil {
		return data, fmt.Errorf("Error reading file: %v", err.Error())
	}

	data = bytes

	return data, nil
}

func ReadAndUnmarshal[V any](path string, value *V) error {
	bytes, err := Read(path)
	if err != nil {
		return err
	}

	if err := json.Unmarshal(bytes, value); err != nil {
		return fmt.Errorf("Error decoding JSON: %v", err.Error())
	}

	return nil
}
