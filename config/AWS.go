package config

type Credentials struct {
	AWSAccessKeyID     string `toml:"aws_access_key_id"`
	AWSSecretAccessKey string `toml:"aws_secret_access_key"`
	AWSSessionToken    string `toml:"aws_session_token"`
	AWSRegion          string `toml:"aws_region"`
	AWSAssumedRole     string `toml:"aws_assumed_role"`
}

type AWSConfig struct {
	Credentials Credentials `toml:"credentials"`
}
