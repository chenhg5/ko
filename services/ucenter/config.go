package ucenter

var config map[string]string

func InitConfig() *map[string]string {
	config = map[string]string{
		"access_log_path" : "services/ucenter/access.log",
		"jwt_auth_secret" : "secret",
		"pid_path" : "ucenter/pid",
		"server_port" : ":4002",

		"database_ip":"",
		"database_port":"",
		"database_user":"",
		"database_pwd":"",
		"database_name":"",
	}
	return &config
}

func GetConfig() *map[string]string {
	return &config
}


type Connections struct {
	DATABASE_IP       string
	DATABASE_PORT     string
	DATABASE_USERNAME string
	DATABASE_PASSWORD string
	DATABASE_NAME     string
}

func GetCons() map[string]*Connections {

	return map[string]*Connections{
		"official_account" : &Connections{
			DATABASE_USERNAME: "",
			DATABASE_IP: "",
			DATABASE_PASSWORD:  "",
			DATABASE_NAME: "",
			DATABASE_PORT: "",
		},
	}
}