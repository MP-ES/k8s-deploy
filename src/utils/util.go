package utils

import (
	"bufio"
	"bytes"
	"encoding/base64"
	"io"
	"os"
	"regexp"
	"strings"

	"gopkg.in/yaml.v2"
)

func UnmarshalSingleYamlKeyFromMultifile(base64FileContent *string, out interface{}) error {
	in, err := base64.StdEncoding.DecodeString(*base64FileContent)
	if err != nil {
		return err
	}

	r := bytes.NewReader(in)
	decoder := yaml.NewDecoder(r)
	for {
		if err := decoder.Decode(out); err != nil {
			// Break when there are no more documents to decode
			if err != io.EOF {
				return err
			}
			break
		}
	}
	return nil
}

func SearchPatternInFileLineByLine(fileName string, pattern string) ([]string, error) {
	file, err := os.Open(fileName)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	regex, err := regexp.Compile(pattern)
	if err != nil {
		return nil, err
	}

	matchList := []string{}

	fileScanner := bufio.NewScanner(file)
	for fileScanner.Scan() {
		line := fileScanner.Text()
		if regex.MatchString(line) {
			matchList = append(matchList, strings.TrimSpace(line))
		}
	}

	if err := fileScanner.Err(); err != nil {
		return nil, err
	}

	return SliceRemoveDuplicateElements(matchList), nil
}
