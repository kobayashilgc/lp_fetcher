package cli

import (
	"encoding/json"
	"fmt"
	"os"
)

func writeJSONResult(output string, v any) error {
	data, err := json.Marshal(v)
	if err != nil {
		return err
	}
	return os.WriteFile(output, data, 0644)
}

func finishFetchCommand(output, status, msg string, v any) error {
	if err := writeJSONResult(output, v); err != nil {
		return err
	}
	if status == "failed" {
		return fmt.Errorf("%s", msg)
	}
	return nil
}
