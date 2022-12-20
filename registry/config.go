package registry

type Config interface {
	GetInt(string) int
	GetString(string) string
	GetStringSlice(string) []string
	Get(string) any
}
