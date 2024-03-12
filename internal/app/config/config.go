package config

import (
	"flag"
	"fmt"
	"github.com/caarlos0/env/v6"
	"github.com/denis-oreshkevich/gophkeeper/internal/app/logger"
	"os"
)

const (
	defaultHost = "localhost"

	defaultPort = "8081"
)

var conf Config

func Get() Config {
	return conf
}

func Parse() (*Config, error) {
	flag.StringVar(&conf.ServerAddress, "a", fmt.Sprintf("%s:%s", defaultHost, defaultPort),
		"HTTP server address")
	flag.StringVar(&conf.DataBaseURI, "d",
		"host=localhost port=5433 user=postgres password=postgres dbname=mart sslmode=disable",
		"DataBase URI")

	fileSet := flag.NewFlagSet("file", flag.ExitOnError)

	fileSet.StringVar(&conf.UserLogin, "ul", "", "User login")
	fileSet.StringVar(&conf.UserPassword, "up", "", "User password")

	fileSet.StringVar(&conf.Filename, "f", "", "filename")
	fileSet.StringVar(&conf.FilePath, "fp", "", "Path to file")

	textSet := flag.NewFlagSet("text", flag.ExitOnError)

	textSet.StringVar(&conf.UserLogin, "ul", "", "User login")
	textSet.StringVar(&conf.UserPassword, "up", "", "User password")

	textSet.StringVar(&conf.Text, "t", "", "Text")

	cardSet := flag.NewFlagSet("card", flag.ExitOnError)

	cardSet.StringVar(&conf.UserLogin, "ul", "", "User login")
	cardSet.StringVar(&conf.UserPassword, "up", "", "User password")

	cardSet.StringVar(&conf.CardNum, "n", "", "Number")
	cardSet.StringVar(&conf.CardCVC, "c", "", "CVC")
	cardSet.StringVar(&conf.CardHolderName, "hn", "", "Holder name")

	credSet := flag.NewFlagSet("cred", flag.ExitOnError)

	credSet.StringVar(&conf.UserLogin, "ul", "", "User login")
	credSet.StringVar(&conf.UserPassword, "up", "", "User password")

	credSet.StringVar(&conf.CredentialsLogin, "l", "", "Login")
	credSet.StringVar(&conf.CredentialsPassword, "p", "", "Password")

	flag.BoolVar(&conf.Sync, "sync", false, "Synchronize")

	if len(os.Args) >= 2 {
		switch os.Args[1] {
		case "file":
			fileSet.Parse(os.Args[2:])
			conf.IsFileFlagsParsed = true
		case "text":
			textSet.Parse(os.Args[2:])
			conf.IsTextFlagsParsed = true
		case "card":
			cardSet.Parse(os.Args[2:])
			conf.IsCardFlagsParsed = true
		case "cred":
			credSet.Parse(os.Args[2:])
			conf.IsCredentialsFlagsParsed = true
		default:
			flag.PrintDefaults()
		}
	}

	flag.Parse()
	err := env.Parse(&conf)
	if err != nil {
		return nil, fmt.Errorf("env.Parse: %w", err)
	}

	logger.Log.Info(fmt.Sprintf("initializing Config %+v", conf))

	return &conf, nil
}
