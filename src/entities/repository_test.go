package entities_test

import (
	"k8s-deploy/entities"
	"os"
	"reflect"
	"testing"

	"github.com/gdexlab/go-render/render"
)

type repositoryTest struct {
	githubRepository   string
	expectedRepository *entities.Repository
	expectedError      string
}

var repositoryTests = [...]repositoryTest{
	{"", nil, "couldn't get the repository"},
	{"wrong-string", nil, "repository name format different from expected"},
	{"owner/repository-all",
		&entities.Repository{"repository-all", "https://github.com/owner/repository-all",
			&struct {
				Name            string   "yaml:\"name\""
				K8sEnvs         []string "yaml:\"k8s-envs,flow\""
				Images          []string "yaml:\"images,flow\""
				Secrets         []string "yaml:\"secrets,flow\""
				ResourcesQuotas *struct {
					LimitsCpu    string "yaml:\"limits.cpu\""
					LimitsMemory string "yaml:\"limits.memory\""
				} "yaml:\"resources-quotas\""
				RequestsIngresses *map[string][]string "yaml:\"requests-ingresses\""
			}{"repository-all",
				[]string{"env1", "env2", "env3"},
				[]string{"docker_image_one", "docker_image_two"},
				[]string{"database_user", "database_password"},
				&struct {
					LimitsCpu    string "yaml:\"limits.cpu\""
					LimitsMemory string "yaml:\"limits.memory\""
				}{"100m", "100Mi"},
				&map[string][]string{
					"env1": {"application.env1.domain.com"},
					"env2": {"application.env2.domain.com"},
					"env3": {"application.env3.domain.com", "application.domain.com"},
				},
			},
		},
		""},
	{"owner/repository-min",
		&entities.Repository{"repository-min", "https://github.com/owner/repository-min",
			&struct {
				Name            string   "yaml:\"name\""
				K8sEnvs         []string "yaml:\"k8s-envs,flow\""
				Images          []string "yaml:\"images,flow\""
				Secrets         []string "yaml:\"secrets,flow\""
				ResourcesQuotas *struct {
					LimitsCpu    string "yaml:\"limits.cpu\""
					LimitsMemory string "yaml:\"limits.memory\""
				} "yaml:\"resources-quotas\""
				RequestsIngresses *map[string][]string "yaml:\"requests-ingresses\""
			}{"repository-min",
				[]string{"env1"},
				[]string{"docker_image"},
				nil,
				&struct {
					LimitsCpu    string "yaml:\"limits.cpu\""
					LimitsMemory string "yaml:\"limits.memory\""
				}{"100m", "100Mi"},
				nil,
			},
		},
		""},
}

func TestGetRepository(t *testing.T) {
	for _, test := range repositoryTests {
		orig := os.Getenv("GITHUB_REPOSITORY")
		os.Setenv("GITHUB_REPOSITORY", test.githubRepository)
		t.Cleanup(func() { os.Setenv("GITHUB_REPOSITORY", orig) })

		repository, err := entities.GetRepository(&entities.GitOpsRepository{Owner: "MP-ES", Repository: "k8s-deploy", PathSchemas: "testdata"})

		if err != nil {
			if test.expectedError == "" || err.Error() != test.expectedError {
				t.Errorf("repository error %s not equal to expected %s", err, test.expectedError)
			}
		} else {
			if !reflect.DeepEqual(repository, test.expectedRepository) {
				t.Errorf("repository\n%s\nnot equal to expected\n%s", render.Render(repository), render.Render(test.expectedRepository))
			}
		}
	}
}
