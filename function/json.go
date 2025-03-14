package function

import (
	"bytes"
	"encoding/json"
)

// HTML JSON 编码
func JsonHtmlEnCode(data any) ([]byte, error) {
	buffer := &bytes.Buffer{}
	encoder := json.NewEncoder(buffer)
	encoder.SetEscapeHTML(false)
	encoder.SetIndent("", "    ")
	err := encoder.Encode(data)
	return buffer.Bytes(), err
}
