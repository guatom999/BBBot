package utils

import (
	"bytes"
	"fmt"
	"io"
	"os"
)

func SaveBufferTofile(buff *bytes.Buffer, folderPath string, fileName string) (string, error) {

	if err := os.MkdirAll(folderPath, 0755); err != nil {
		return "", fmt.Errorf("failed to create folder :%s", err.Error())
	}

	filePath := fmt.Sprintf("%s/%s.png", folderPath, fileName)

	file, err := os.Create(filePath)
	if err != nil {
		return "", fmt.Errorf("failed to create file :%s", err.Error())
	}

	defer file.Close()

	if _, err := io.Copy(file, buff); err != nil {
		return "", fmt.Errorf("failed to write buffer to file :%s", err.Error())
	}

	return filePath, nil
}
