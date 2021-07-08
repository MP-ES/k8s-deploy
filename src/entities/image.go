package entities

import (
	"fmt"
	"k8s-deploy/utils"
	"regexp"
	"strings"
)

type Image struct {
	Name string
}

func (i *Image) String() string {
	return i.Name
}

func (i *Image) UnmarshalYAML(unmarshal func(interface{}) error) error {
	var output string
	if err := unmarshal(&output); err != nil {
		return err
	}
	i.Name = output
	return nil
}

func ValidateImagesFromAppDeploy(appDeployPath string, repoRules *RepositoryRules) error {
	imageLines, err := utils.SearchPatternInFileLineByLine(appDeployPath, "^( )*image: .*$")
	if err != nil {
		return err
	}

	regex := regexp.MustCompile(fmt.Sprintf(`^.*/%s/(?P<image>.*)$`, repoRules.Name))
	for _, line := range imageLines {
		image := sanitizeImageLine(line, regex)
		if image != nil {
			// validar imagem
			fmt.Printf("%s", *image)
		}
	}

	return nil
}

func sanitizeImageLine(imageLine string, regex *regexp.Regexp) *string {
	slice := strings.Split(imageLine, ":")

	if len(slice) > 1 {
		result := regex.FindStringSubmatch(slice[1])
		if len(result) > 1 {
			return &result[1]
		}
	}

	return nil
}
