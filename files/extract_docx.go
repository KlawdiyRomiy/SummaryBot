package files

import (
	"fmt"
	"github.com/chchench/textract"
)

func ExtractDocx(tempPath string) (string, error) {
	text, err := textract.RetrieveTextFromFile(tempPath)
	if err != nil {
		return "", fmt.Errorf("❌ Ошибка в ExtractDocx: %v", err)
	}
	return text, nil
}
