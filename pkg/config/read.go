package config

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/spf13/viper"
	"gopkg.in/yaml.v3"
	//"path/filepath"
)

func getFileExt(path string) string {
	ext := filepath.Ext(path)
	if len(ext) > 1 {
		return ext[1:]
	}
	return ""
}

const ConfigType = "yaml"

func getConfigType(path string) (string, error) {
	ext := strings.ToLower(getFileExt(path))
	switch ext {
	case "yaml", "yml":
		return ConfigType, nil
	default:
		return "", fmt.Errorf("file type `%s` not support yet", ext)
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

func ReadCasesFromFile(path string) (cases []*Case, err error) {
	// read lines, for display the failed asset line number
	fileLines, err := ReadLines(path)
	if err != nil {
		return
	}

	// only support yaml now
	_, err = getConfigType(path)
	if err != nil {
		return
	}

	// if configType != configType {
	// 	var reader *os.File
	// 	reader, err = os.Open(path)
	// 	if err != nil {
	// 		return
	// 	}

	// 	v := viper.New()
	// 	v.SetConfigType(configType)
	// 	err = v.ReadConfig(reader)
	// 	if err != nil {
	// 		return
	// 	}

	// 	var c Case
	// 	err = v.Unmarshal(&c)
	// 	if err != nil {
	// 		return
	// 	}
	// 	// set the content
	// 	c.FileLines = fileLines
	// 	c.AllKeys = v.AllKeys()
	// 	c.Path = path
	// 	c.Index = 1

	// 	cases = append(cases, &c)

	// 	return cases, nil
	// }

	b, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("error reading %s: %v", path, err)
	}

	decoder := yaml.NewDecoder(strings.NewReader(string(b)))
	count := 0
	for {
		count += 1
		var node yaml.Node
		err1 := decoder.Decode(&node)
		if err1 == io.EOF {
			break
		}
		if err1 != nil {
			return nil, fmt.Errorf("error decoding: %v", err1)
		}

		out, err2 := yaml.Marshal(&node)
		if err2 != nil {
			log.Fatalf("Failed to marshal YAML: %v", err2)
		}

		v := viper.New()
		v.SetConfigType("yaml")
		v.ReadConfig(bytes.NewBuffer(out))

		var c Case
		err3 := v.Unmarshal(&c)
		if err3 != nil {
			return nil, err3
		}
		// FIXME: the lines is not correct from now!
		c.FileLines = fileLines
		c.AllKeys = v.AllKeys()
		c.Path = path
		c.Index = count

		cases = append(cases, &c)
	}

	return
}

// ReadLines read the file content, and return the lines
// NOTE: we trans all line to lower case
// should get more info, mut
func ReadLines(path string) (lines map[int]map[int]string, err error) {
	var readFile *os.File
	readFile, err = os.Open(path)
	if err != nil {
		return
	}

	scanner := bufio.NewScanner(readFile)
	scanner.Split(bufio.ScanLines)

	lines = make(map[int]map[int]string, 1)

	index := 1
	lineNo := 0
	for scanner.Scan() {
		lineNo += 1

		line := strings.ToLower(scanner.Text())
		if line == "---" {
			index += 1
		}

		if lines[index] == nil {
			lines[index] = make(map[int]string, 4)
		}

		lines[index][lineNo] = line
	}

	readFile.Close()
	return
}
