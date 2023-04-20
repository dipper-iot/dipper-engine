package arithmetic

type Option struct {
	Operators   map[string]string `json:"operators"`
	NextError   string            `json:"next_error"`
	NextSuccess string            `json:"next_success"`
	Debug       bool              `json:"debug"`
}
