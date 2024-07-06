package execute

import (
	"bytes"
	"context"
	"io"
	"os"
	"os/exec"
	"time"
)

func ExecuteCode(code string) (string, string) {
	// Создаем временный файл
	tmpFile, err := os.CreateTemp("", "*.go")
	if err != nil {
		return "", err.Error()
	}

	defer os.Remove(tmpFile.Name())

	// Записываем код в файл
	if _, err := io.WriteString(tmpFile, code); err != nil {
		return "", err.Error()
	}

	if err := tmpFile.Close(); err != nil {
		return "", err.Error()
	}

	// Выполняем файл
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	cmd := exec.CommandContext(ctx, "go", "run", tmpFile.Name())

	var out bytes.Buffer
	var stderr bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &stderr

	err = cmd.Run()
	if err != nil {
		return out.String(), stderr.String()
	}

	return out.String(), stderr.String()
}
