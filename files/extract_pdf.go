package files

import (
	"fmt"
	"github.com/chchench/textract"
)

func ExtractPDF(tempPath string) (string, error) {
	text, err := textract.RetrieveTextFromFile(tempPath)
	if err != nil {
		return "", fmt.Errorf("❌ Ошибка в ExtractPDF: %v", err)
	}
	return text, nil
}
