package ucenter

var config map[string]string

func InitConfig() *map[string]string {
	config = map[string]string{
		"access_log_path" : "gateway/access.log",
		"jwt_auth_secret" : "secret",
		"pid_path" : "gateway/pid",
		"server_port" : ":4002",
	}
	return &config
}

func GetConfig() *map[string]string {
	return &config
}