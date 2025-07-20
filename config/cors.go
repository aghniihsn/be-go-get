package config

var allowedOrigins = []string{
	"*",
}

func GetAllowedOrigins() []string {
	return allowedOrigins
}
