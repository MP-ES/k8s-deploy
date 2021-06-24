package entities_test

import (
	"k8s-deploy/entities"
	"os"
	"testing"

	"github.com/go-test/deep"
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
			&entities.RepositoryRules{
				Name: "repository-all",
				K8sEnvs: []*entities.K8sEnv{{
					Name: "env1"}, {
					Name: "env2"}, {
					Name: "env3"}},
				Images: []*entities.Image{
					{Name: "docker_image_one"},
					{Name: "docker_image_two"}},
				Secrets: []*entities.Secret{
					{Name: "database_user"},
					{Name: "database_password"}},
				ResourcesQuotas: &entities.ResourcesQuotas{
					LimitsCpu: "100m", LimitsMemory: "100Mi"},
				Ingresses: &map[entities.K8sEnv][]*entities.Ingress{
					{Name: "env1"}: {&entities.Ingress{Name: "application.env1.domain.com"}},
					{Name: "env2"}: {&entities.Ingress{Name: "application.env2.domain.com"}},
					{Name: "env3"}: {&entities.Ingress{Name: "application.env3.domain.com"},
						&entities.Ingress{Name: "application.domain.com"}}},
			},
		},
		""},
	{"owner/repository-min",
		&entities.Repository{"repository-min", "https://github.com/owner/repository-min",
			&entities.RepositoryRules{
				Name:            "repository-min",
				K8sEnvs:         []*entities.K8sEnv{{Name: "env1"}},
				Images:          []*entities.Image{{Name: "docker_image"}},
				ResourcesQuotas: &entities.ResourcesQuotas{LimitsCpu: "100m", LimitsMemory: "100Mi"},
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
				t.Errorf("repository error '%s' not equal to expected '%s'", err, test.expectedError)
			}
		} else {
			if diff := deep.Equal(repository, test.expectedRepository); diff != nil {
				t.Errorf("repository not equal to expected")
				t.Error(diff)
			}
		}
	}
}
