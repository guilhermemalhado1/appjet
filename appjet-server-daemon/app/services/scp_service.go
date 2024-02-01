package services

import "fmt"

func PullCodeFromSCP(codeDir string) string {
	return fmt.Sprintf("COPY %s /app_builder", codeDir)
}
