package writer

import (
	"bufio"
	"fmt"
	"os"
	"passkeeper/internal/crypto"
	"strings"
)

func getFilePath(service string) string {
	return fmt.Sprintf("./passwords/%s.txt", service)
}

func getDecryptedFilePath(service string) string {
	return fmt.Sprintf("./passwords/%s-decrypt.txt", service)
}

func InitServiceStorage(service string) error {
	passwordsDir := "./passwords"
	if err := os.MkdirAll(passwordsDir, 0755); err != nil {
		return fmt.Errorf("error creating folder %s: %v", passwordsDir, err)
	}

	filePath := getFilePath(service)
	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		if _, err := os.Create(filePath); err != nil {
			return fmt.Errorf("error creating file %s: %v", filePath, err)
		}
	}
	return nil
}

func processLines(filePath string, generatedKey []byte, encrypt bool) ([]string, error) {
	content, err := os.ReadFile(filePath)
	if err != nil {
		return nil, err
	}

	var result []string
	for _, line := range strings.Split(string(content), "\n") {
		if strings.TrimSpace(line) == "" {
			continue
		}
		parts := strings.SplitN(line, ":", 2)
		if len(parts) != 2 {
			continue
		}

		var processed string
		if encrypt {
			processed, err = crypto.Encrypt(parts[1], generatedKey)
		} else {
			processed, err = crypto.Decrypt(parts[1], generatedKey)
		}
		if err != nil {
			return nil, err
		}
		result = append(result, fmt.Sprintf("%s:%s", parts[0], processed))
	}
	return result, nil
}

func EncryptFile(service string, _ string, key []byte) error {
	decryptedPath := getDecryptedFilePath(service)
	encryptedPath := getFilePath(service)

	lines, err := processLines(decryptedPath, key, true)
	if err != nil {
		return fmt.Errorf("error encrypting: %v", err)
	}

	return os.WriteFile(encryptedPath, []byte(strings.Join(lines, "\n")+"\n"), 0644)
}

func DecryptFile(service string, _ string, key []byte) error {
	decryptedPath := getDecryptedFilePath(service)
	encryptedPath := getFilePath(service)

	lines, err := processLines(encryptedPath, key, false)
	if err != nil {
		return fmt.Errorf("error decrypting: %v", err)
	}

	return os.WriteFile(decryptedPath, []byte(strings.Join(lines, "\n")+"\n"), 0644)
}

func ShowList(service string) error {
	filePath := getFilePath(service)

	content, err := os.ReadFile(filePath)
	if err != nil {
		return fmt.Errorf("error reading file: %v", err)
	}

	for _, line := range strings.Split(string(content), "\n") {
		if strings.TrimSpace(line) == "" {
			continue
		}
		parts := strings.SplitN(line, ":", 2)
		fmt.Println(parts[0])
	}

	return nil
}

func ShowValue(key, service, _ string, generatedKey []byte) (string, error) {
	filePath := getFilePath(service)

	file, err := os.Open(filePath)
	if err != nil {
		return "", fmt.Errorf("error opening file: %v", err)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		parts := strings.SplitN(scanner.Text(), ":", 2)
		if len(parts) == 2 && parts[0] == key {
			return crypto.Decrypt(parts[1], generatedKey)
		}
	}

	if err := scanner.Err(); err != nil {
		return "", fmt.Errorf("error reading file: %v", err)
	}
	return "", fmt.Errorf("key %s not found", key)
}

func WriteValue(key, value, service string, generatedKey []byte) error {
	filePath := getFilePath(service)

	encrypted, err := crypto.Encrypt(value, generatedKey)
	if err != nil {
		return fmt.Errorf("encryption failed: %v", err)
	}

	var lines []string
	if content, err := os.ReadFile(filePath); err == nil {
		for _, line := range strings.Split(string(content), "\n") {
			if !strings.HasPrefix(line, key+":") && strings.TrimSpace(line) != "" {
				lines = append(lines, line)
			}
		}
	}

	lines = append(lines, fmt.Sprintf("%s:%s", key, encrypted))
	return os.WriteFile(filePath, []byte(strings.Join(lines, "\n")+"\n"), 0644)
}

func RemoveValue(key, service string) error {
	filePath := getFilePath(service)

	content, err := os.ReadFile(filePath)
	if err != nil {
		return fmt.Errorf("error reading file: %v", err)
	}

	var lines []string
	for _, line := range strings.Split(string(content), "\n") {
		if !strings.HasPrefix(line, key+":") && strings.TrimSpace(line) != "" {
			lines = append(lines, line)
		}
	}

	return os.WriteFile(filePath, []byte(strings.Join(lines, "\n")+"\n"), 0644)
}
