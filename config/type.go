package config

import "fmt"

type Config struct {
	DB     DBConfig     `yaml:"db"`
	Secret SecretConfig `yaml:"secret"`
}

type DBConfig struct {
	Host     string `yaml:"host"`
	Port     int    `yaml:"port"`
	User     string `yaml:"user"`
	Password string `yaml:"password"`
	Database string `yaml:"database"`
}

type SecretConfig struct {
	RsaPrivatePem string `yaml:"rsa_private_pem"`
	RsaPublicPem  string `yaml:"rsa_public_pem"`
}

func (db *DBConfig) ToJdbcUrl() string {
	if db == nil {
		return ""
	}

	return fmt.Sprintf("postgres://%s:%s@%s:%d/%s?sslmode=disable", db.User, db.Password, db.Host, db.Port, db.Database)
}
