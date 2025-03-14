package function

import (
	"encoding/json"
	"io"
	"net/http"
)

func GetExternalIP() (string, error) {
	resp, err := http.Get("https://ipinfo.io/json")
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	// 读取响应体
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}
	var ipInfo struct {
		Ip       string `json:"ip"`
		City     string `json:"city"`
		Region   string `json:"region"`
		Country  string `json:"country"`
		Loc      string `json:"loc"`
		Org      string `json:"org"`
		Postal   string `json:"postal"`
		Timezone string `json:"timezone"`
		Readme   string `json:"readme"`
	}
	err = json.Unmarshal(body, &ipInfo)
	if err != nil {
		return "", err
	}
	return ipInfo.Ip, nil
}
