package config

import "github.com/spf13/viper"

var (
	Env                string
	MigrationsPath     string
	Port               int
	Postgres           string
	Public             string
	SessionAuthKey     string
	SessionCryptKey    string
	AwsAccessKeyID     string
	AwsSecretAccessKey string
	AwsRegion          string
	SqsQueueURL        string
	QueueWorkers       int
	SendGridKey        string
	SendGridUser       string
	MailerBaseURL      string
	TokenSigningKey    []byte
)

func init() {

	viper.SetDefault("Env", "development")
	viper.SetDefault("Port", 3000)
	viper.SetConfigName("config")
	viper.AddConfigPath("/etc/erraroo")
	viper.AddConfigPath("/tmp/erraroo")
	viper.AddConfigPath("$HOME/.erraroo")
	viper.ReadInConfig()

	Env = viper.GetString("Env")
	Port = viper.GetInt("Port")
	Postgres = viper.GetString("Postgres")
	Public = viper.GetString("Public")
	SessionAuthKey = viper.GetString("SessionAuthKey")
	SessionCryptKey = viper.GetString("SessionCryptKey")
	AwsAccessKeyID = viper.GetString("AwsAccessKeyID")
	AwsSecretAccessKey = viper.GetString("AwsSecretAccessKey")
	AwsRegion = viper.GetString("AwsRegion")
	SqsQueueURL = viper.GetString("SqsQueueURL")
	QueueWorkers = viper.GetInt("QueueWorkers")
	SendGridKey = viper.GetString("SendGridKey")
	SendGridUser = viper.GetString("SendGridUser")
	MailerBaseURL = viper.GetString("MailerBaseURL")
	TokenSigningKey = []byte(viper.GetString("SessionCryptKey"))
}
