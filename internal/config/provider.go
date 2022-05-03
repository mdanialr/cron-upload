package config

// Provider holds data about provider that would be used to upload files.
type Provider struct {
	Name  string `yaml:"name"`  // the name of the provider. support only 'drive'.
	Auth  string `yaml:"auth"`  // the file path for authentication to the provider.
	Token string `yaml:"token"` // a directory where temporary token file for provider 'drive' is stored.
}
