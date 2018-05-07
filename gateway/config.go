package gateway

var config map[string]string

func InitConfig() *map[string]string {
	config = map[string]string{
		"access_log_path" : "gateway/access.log",
		"jwt_auth_secret" : "secret",
	}
	return &config
}

func GetConfig() *map[string]string {
	return &config
}