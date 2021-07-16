package config

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/spf13/viper"
	//"path/filepath"
)

func getFileExt(path string) string {
	ext := filepath.Ext(path)
	if len(ext) > 1 {
		return ext[1:]
	}
	return ""
}

func getConfigType(path string) (string, error) {
	ext := strings.ToLower(getFileExt(path))
	switch ext {
	case "yaml", "yml":
		return "yaml", nil
	case "json":
		return "json", nil
	case "toml":
		return "toml", nil
	//case "hcl":
	//	return "hcl", nil
	case "properties", "props", "prop":
		return "prop", nil
	case "ini":
		return "ini", nil
	default:
		return "", fmt.Errorf("filt type `%s` not support yet", ext)
	}
}

func ReadFromFile(path string) (v *viper.Viper, err error) {
	var reader *os.File
	reader, err = os.Open(path)
	if err != nil {
		return
	}

	configType, err := getConfigType(path)
	if err != nil {
		return
	}

	v = viper.New()
	v.SetConfigType(configType)
	err = v.ReadConfig(reader)
	if err != nil {
		return
	}

	return v, nil
}
