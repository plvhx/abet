package config

import (
    "log"
    "sync"

    "github.com/joho/godotenv"
    "github.com/kelseyhightower/envconfig"
)

var (
    once sync.Once
    conf Config
)

type Config struct {
    AppPort             int    `envconfig:"app_port"`

    DBUser              string `envconfig:"db_user"`
    DBPassword          string `envconfig:"db_password"`
    DBHost              string `envconfig:"db_host"`
    DBPort              string `envconfig:"db_port"`
    DBSchema            string `envconfig:"db_schema"`
    DBMaxIdleConn       int    `envconfig:"db_max_idle_conn"`
    DBMaxOpenDescriptor int    `envconfig:"db_max_open_descriptor"`

    AWSAccessKey        string `envconfig:"aws_access_key"`
    AWSSecretKey        string `envconfig:"aws_secret_key"`
    AWSBucketName       string `envconfig:"aws_bucket_name"`
    AWSEndpoint         string `envconfig:"aws_endpoint"`
    AWSRegion           string `envconfig:"aws_region"`

    SMTPHost            string `envconfig:"smtp_host"`
    SMTPPort            int    `envconfig:"smtp_port"`
    SMTPUser            string `envconfig:"smtp_user"`
    SMTPPassword        string `envconfig:"smtp_password"`
}

func GetConfig() Config {
    once.Do(func() {
        _ = godotenv.Load()

        if err := envconfig.Process("", &conf); err != nil {
            log.Fatal(err)
        }
    })

    return conf
}
