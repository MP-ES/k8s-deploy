module k8s-deploy

go 1.16

require (
	github.com/go-openapi/jsonpointer v0.19.5 // indirect
	github.com/go-openapi/jsonreference v0.19.5 // indirect
	github.com/go-openapi/swag v0.19.14 // indirect
	github.com/go-test/deep v1.0.7
	github.com/golang/protobuf v1.5.2 // indirect
	github.com/google/go-github/v35 v35.2.0
	github.com/hashicorp/go-multierror v1.1.1
	github.com/joho/godotenv v1.3.0
	github.com/mikefarah/yq/v4 v4.11.1
	github.com/sethvargo/go-githubactions v0.4.0
	golang.org/x/crypto v0.1.0 // indirect
	golang.org/x/oauth2 v0.0.0-20210514164344-f6687ab2804c
	gopkg.in/op/go-logging.v1 v1.0.0-20160211212156-b2cb9fa56473
	gopkg.in/yaml.v2 v2.4.0
	sigs.k8s.io/kustomize/api v0.8.11
	sigs.k8s.io/kustomize/kyaml v0.11.0
)
