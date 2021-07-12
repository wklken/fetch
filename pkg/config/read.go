package config

import (
	"errors"
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
	switch strings.ToLower(getFileExt(path)) {
	case "yaml", "yml":
		return "yaml", nil
	case "json":
		return "json", nil
	case "toml":
		return "toml", nil
	default:
		return "", errors.New("not support yet")
	}
}

//fileType, err := getFileType(path)
//if err != nil {
//	return
//}
//viper.SetConfigType(fileType)

func ReadFromFile(path string) (v *viper.Viper, err error) {
	var reader *os.File
	reader, err = os.Open(path)
	if err != nil {
		return
	}

	//fmt.Println("path:", path)
	//
	//b1 := make([]byte, 200)
	//n1, err := reader.Read(b1)
	//fmt.Println(n1, err, string(b1))
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
	//fmt.Println("v.AllKeys", v.AllKeys())

	return v, nil
}
