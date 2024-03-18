package config

type Action string

const (
	ActionGet Action = "GET"

	ActionSave Action = "SAVE"

	ActionDelete Action = "DELETE"
)

type Config struct {
	ServerAddress string `env:"RUN_ADDRESS"`
	DataBaseURI   string `env:"DATABASE_URI"`

	UserLogin    string `env:"USER_LOGIN"`
	UserPassword string `env:"USER_PASSWORD"`

	WorkingDir string `env:"WORKING_DIR"`

	ID    string
	IsNew bool

	Filename string

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
	IsSync                   bool

	Action Action
}
