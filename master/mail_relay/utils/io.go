package utils

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
)

func ReadAndUnmarshal[V any](path string, value *V) error {
	file, err := os.Open(path)
	if err != nil {
		return fmt.Errorf("Error opening file: %v", err.Error())
	}
	defer file.Close()

	bytes, err := io.ReadAll(file)
	if err != nil {
		return fmt.Errorf("Error reading file: %v", err.Error())
	}

	if err := json.Unmarshal(bytes, value); err != nil {
		return fmt.Errorf("Error decoding JSON: %v", err.Error())
	}

	return nil
}

