module k8s-deploy

go 1.15

require (
	github.com/gdexlab/go-render v1.0.1
	github.com/go-test/deep v1.0.7
	github.com/google/go-github v17.0.0+incompatible
	github.com/google/go-github/v35 v35.2.0
	github.com/hashicorp/go-multierror v1.1.1
	github.com/joho/godotenv v1.3.0
	github.com/sethvargo/go-githubactions v0.4.0
	golang.org/x/oauth2 v0.0.0-20210514164344-f6687ab2804c
	gopkg.in/yaml.v2 v2.4.0
	sigs.k8s.io/kustomize/api v0.8.10
	sigs.k8s.io/yaml v1.2.0
)
