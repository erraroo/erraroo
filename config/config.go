package config

import "github.com/spf13/viper"

var (
	Env                  string
	LogSql               bool
	MigrationsPath       string
	Port                 int
	Postgres             string
	Redis                string
	SqsLongPollTimeout   int64
	SqsMessagesPerWorker int64
	SqsQueueURL          string
	QueueWorkers         int
	SendGridKey          string
	SendGridUser         string
	MailerBaseURL        string
	TokenSigningKey      []byte
)

func init() {
	viper.SetDefault("Env", "development")
	viper.SetDefault("Port", 3000)
	viper.SetDefault("SqsLongPollTimeout", 20)
	viper.SetDefault("SqsMessagesPerWorker", 1)

	viper.SetConfigName("config")
	viper.AddConfigPath("/etc/erraroo")
	viper.AddConfigPath("/tmp/erraroo")
	viper.AddConfigPath("$HOME/.erraroo")
	viper.ReadInConfig()

	Env = viper.GetString("Env")
	LogSql = viper.GetBool("LogSql")

	// jwt
	TokenSigningKey = []byte(viper.GetString("TokenSigningKey"))

	// Web server
	Port = viper.GetInt("Port")

	// DB
	Postgres = viper.GetString("Postgres")
	Redis = viper.GetString("Redis")

	// Queues
	SqsLongPollTimeout = int64(viper.GetInt("SqsLongPollTimeout"))
	SqsMessagesPerWorker = int64(viper.GetInt("SqsMessagesPerWorker"))
	SqsQueueURL = viper.GetString("SqsQueueURL")
	QueueWorkers = viper.GetInt("QueueWorkers")

	// Email
	SendGridKey = viper.GetString("SendGridKey")
	SendGridUser = viper.GetString("SendGridUser")
	MailerBaseURL = viper.GetString("MailerBaseURL")
}
