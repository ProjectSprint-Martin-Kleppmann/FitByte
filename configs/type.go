package configs

type Config struct {
	App    App          `mapstructure:"app" validate:"required"`
	DB     Database     `mapstructure:"database" validate:"required"`
	Secret SecretConfig `mapstructure:"secret" validate:"required"`
	Minio  MinioConfig  `mapstructure:"minio" validate:"required"`
}

type App struct {
	Port string `mapstructure:"port" validate:"required"`
}

type Database struct {
	Host     string `mapstructure:"host" validate:"required"`
	Port     string `mapstructure:"port" validate:"required"`
	Username string `mapstructure:"username" validate:"required"`
	Password string `mapstructure:"password" validate:"required"`
	DBName   string `mapstructure:"dbname" validate:"required"`
}

type SecretConfig struct {
	JWTSecret string `mapstructure:"jwt_secret" validate:"required"`
}

type MinioConfig struct {
	Endpoint        string `mapstructure:"endpoint" validate:"required"`
	AccessKeyID     string `mapstructure:"access_key_id" validate:"required"`
	SecretAccessKey string `mapstructure:"secret_access_key" validate:"required"`
	UseSSL          bool   `mapstructure:"use_ssl"`
	Bucket          string `mapstructure:"bucket" validate:"required"`
}
