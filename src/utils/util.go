package utils

import (
	"bytes"
	"encoding/base64"
	"io"

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
