package entities

import (
	"fmt"
	"k8s-deploy/utils"
	"regexp"
	"strings"

	"github.com/hashicorp/go-multierror"
)

const regexSearchImageLine = `^( )*(- )*image: .*$`
const regexSearchImageNameInLine = `^.*/%s/(?P<image>.*)$`

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
	var globalErr *multierror.Error

	imageLines, err := utils.SearchPatternInFileLineByLine(appDeployPath, regexSearchImageLine)
	if err != nil {
		globalErr = multierror.Append(globalErr, err)
	} else {
		regex := regexp.MustCompile(fmt.Sprintf(regexSearchImageNameInLine, repoRules.Name))
		for _, line := range imageLines {
			image := sanitizeImageLine(line, regex)
			if image != nil && !repoRules.IsImageEnabled(*image) {
				globalErr = multierror.Append(globalErr,
					fmt.Errorf("image '%s' is not enabled in repository '%s'. Check the GitOps repository",
						*image, repoRules.Name))
			}
		}
	}

	return globalErr.ErrorOrNil()
}

func GetImagesTagReplace(appDeployPath string, repoName string, eventSHA string) (map[string]string, error) {
	var globalErr *multierror.Error
	imagesReplace := map[string]string{}

	imageLines, err := utils.SearchPatternInFileLineByLine(appDeployPath, regexSearchImageLine)
	if err != nil {
		globalErr = multierror.Append(globalErr, err)
	} else {
		regex := regexp.MustCompile(fmt.Sprintf(regexSearchImageNameInLine, repoName))
		for _, line := range imageLines {
			image := getApplicationImage(line, regex)
			if image != nil {
				imagesReplace[*image] = eventSHA
			}
		}
	}

	if globalErr == nil {
		return imagesReplace, nil
	}
	return nil, globalErr.ErrorOrNil()
}

func getApplicationImage(imageLine string, regex *regexp.Regexp) *string {
	slice := strings.Split(imageLine, ":")

	if len(slice) > 1 {
		result := regex.FindStringSubmatch(slice[1])
		if len(result) > 1 {
			image := strings.TrimSpace(slice[1])
			return &image
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
