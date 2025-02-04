package concourse

const (
	MethodGet   = "get"
	MethodPut   = "put"
	MethodCheck = "check"
)

const (
	EntrypointGet   = "/opt/resource/in"
	EntrypointPut   = "/opt/resource/out"
	EntrypointCheck = "/opt/resource/check"
)

func GetEntrypoint(method string) []string {
	switch method {
	case MethodGet:
		return []string{EntrypointGet}
	case MethodPut:
		return []string{EntrypointPut}
	case MethodCheck:
		return []string{EntrypointCheck}
	}

	return nil
}
