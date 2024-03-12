package config

type Config struct {
	ServerAddress string `env:"RUN_ADDRESS"`
	DataBaseURI   string `env:"DATABASE_URI"`

	Sync bool `env:"SYNC"`

	UserLogin    string `env:"USER_LOGIN"`
	UserPassword string `env:"USER_PASSWORD"`

	Filename string
	FilePath string

	Text string

	CardNum        string
	CardCVC        string
	CardHolderName string

	CredentialsLogin    string
	CredentialsPassword string

	IsFileFlagsParsed        bool
	IsTextFlagsParsed        bool
	IsCardFlagsParsed        bool
	IsCredentialsFlagsParsed bool
}
