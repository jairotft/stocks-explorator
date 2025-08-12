package exceptions

type AppException struct {
	Code   int    `json:"code"`
	Detail string `json:"detail"`
	Stack  string `json:"stack,omitempty"`
}
