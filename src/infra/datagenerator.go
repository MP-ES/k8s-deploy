package infra

func GenerateKustomizationData(kEnv string) interface{} {
	data := make(map[string]interface{})

	data["Namespace"] = kEnv

	return data
}
