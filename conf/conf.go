package conf

import (
	"github.com/Unknwon/goconfig"
)

var Cnf *goconfig.ConfigFile

type CCONF struct {
	MysqlEcology string
	MysqlWallet  string

	Jwt string

	GinIpPort string

	Apiurl_get_all_wallet string
	Apiurl_share_bonus    string
	Apiurl_tfor_info      string
	Apiurl_get_team       string
	Apiurl_user_get_user  string
	Apiurl_auth_get_user  string

	User_tfor   string
	Auth_tfor   string
	Wallet_tfor string

	Real_time_price_api string

	Consul_ip    string
	Consul_port  string
	Service_id   string
	Service_name string

	Schedules       string
	Real_time_price string

	RQ_user_name   string
	RQ_user_passwd string
	RQ_ip_port     string
}

var ConfInfo CCONF

func init() {
	Cnf, _ = goconfig.LoadConfigFile("./conf/app.conf")

	ConfInfo.MysqlEcology, _ = Cnf.GetValue("mysql", "db_ecology")
	ConfInfo.MysqlWallet, _ = Cnf.GetValue("mysql", "db_wallet")

	ConfInfo.Jwt, _ = Cnf.GetValue("jwt", "SignKey")

	ConfInfo.GinIpPort, _ = Cnf.GetValue("gin", "httpport")

	ConfInfo.Apiurl_get_all_wallet, _ = Cnf.GetValue("api", "apiurl_get_all_wallet")
	ConfInfo.Apiurl_share_bonus, _ = Cnf.GetValue("api", "apiurl_share_bonus")
	ConfInfo.Apiurl_tfor_info, _ = Cnf.GetValue("api", "apiurl_tfor_info")
	ConfInfo.Apiurl_get_team, _ = Cnf.GetValue("api", "apiurl_get_team")
	ConfInfo.Apiurl_user_get_user, _ = Cnf.GetValue("api", "apiurl_user_get_user")
	ConfInfo.Apiurl_auth_get_user, _ = Cnf.GetValue("api", "apiurl_auth_get_user")
	ConfInfo.User_tfor, _ = Cnf.GetValue("api", "user_tfor")
	ConfInfo.Auth_tfor, _ = Cnf.GetValue("api", "auth_tfor")
	ConfInfo.Wallet_tfor, _ = Cnf.GetValue("api", "wallet_tfor")
	ConfInfo.Real_time_price_api, _ = Cnf.GetValue("api", "real_time_price_api")
	ConfInfo.Consul_ip, _ = Cnf.GetValue("consul", "consul_ip")
	ConfInfo.Consul_port, _ = Cnf.GetValue("consul", "consul_port")
	ConfInfo.Service_id, _ = Cnf.GetValue("consul", "service_id")
	ConfInfo.Service_name, _ = Cnf.GetValue("consul", "service_name")
	ConfInfo.Schedules, _ = Cnf.GetValue("crontab", "schedules")
	ConfInfo.Real_time_price, _ = Cnf.GetValue("crontab", "real_time_price")
	ConfInfo.RQ_user_name, _ = Cnf.GetValue("rabbit_mq", "user_name")
	ConfInfo.RQ_user_passwd, _ = Cnf.GetValue("crontab", "user_passwd")
	ConfInfo.RQ_ip_port, _ = Cnf.GetValue("crontab", "ip_port")
}
