package function

import (
	"bufio"
	"os"
	"path/filepath"
	"strings"
)

func ParseNginxConfig(filePath string) ([]string, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return nil, err
	}
	defer file.Close()
	var domains []string
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if strings.HasPrefix(line, "server_name ") {
			parts := strings.Fields(line)
			if len(parts) > 1 {
				for _, domain := range parts[1:] {
					domain = strings.TrimRight(domain, ";")
					if domain != "" && domain != "_" && domain != "127.0.0.1" && domain != "phpmyadmin" {
						domains = append(domains, domain)
					}
				}
			}
		}
		if strings.HasPrefix(line, "include") {
			parts := strings.Fields(line)
			if len(parts) > 1 {
				includePath := parts[1]
				includePath = strings.TrimRight(includePath, ";")
				if !filepath.IsAbs(includePath) {
					includePath = filepath.Join(filepath.Dir(filePath), includePath)
				}
				matches, err := filepath.Glob(includePath)
				if err != nil {
					return nil, err
				}
				for _, match := range matches {
					subDomains, err := ParseNginxConfig(match)
					if err != nil {
						return nil, err
					}
					domains = append(domains, subDomains...)
				}
			}
		}
	}
	if err := scanner.Err(); err != nil {
		return nil, err
	}
	return domains, nil
}
