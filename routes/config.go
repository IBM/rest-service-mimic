package routes

type Proxy struct {
	Host               string `json:"host,omitempty"`
	InsecureSkipVerify bool   `json:"skip_insecure_verify,omitempty"`
}

type Response struct {
	Proxy   Proxy                  `json:"proxy,omitempty"`
	Code    int                    `json:"status_code,omitempty"`
	Payload map[string]interface{} `json:"payload,omitempty"`
}

type Route struct {
	Path        string            `json:"path"`
	Methods     []string          `json:"methods"`
	Headers     map[string]string `json:"headers"`
	Response    Response          `json:"response"`
	QueryParams map[string]string `json:"query_params"`
	CacheKey    string            `json:"cache_key"`
}
