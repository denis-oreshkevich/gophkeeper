package config

import (
	"errors"
	"flag"
	"fmt"
	"github.com/caarlos0/env/v6"
	"github.com/denis-oreshkevich/gophkeeper/internal/shared/logger"
	"os"
	"strconv"
	"strings"
)

const (
	defaultHost = "127.0.0.1"

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
		"host=localhost port=5433 user=postgres password=postgres dbname=keeper sslmode=disable",
		"DataBase URI")

	flag.StringVar(&conf.WorkingDir, "wd", "saved", "Working directory")

	fileSet := flag.NewFlagSet("file", flag.ExitOnError)

	fileSet.StringVar(&conf.WorkingDir, "wd", "saved", "Working directory")
	fileSet.StringVar(&conf.UserLogin, "ul", "", "User login")
	fileSet.StringVar(&conf.UserPassword, "up", "", "User password")

	actionFn := func(s string) error {
		if s == "" {
			return errors.New("empty arg")
		}
		us := strings.ToUpper(s)
		if us != string(ActionGet) && us != string(ActionSave) && us != string(ActionDelete) {
			return fmt.Errorf("%s does not match action", s)
		}
		conf.Action = Action(us)
		return nil
	}

	isNewFn := func(s string) error {
		b, err := strconv.ParseBool(s)
		if err != nil {
			return fmt.Errorf("strconv.ParseBool: %w", err)
		}
		conf.IsNew = b
		return nil
	}

	fileSet.Func("a", "action get, save, delete", actionFn)
	fileSet.Func("in", "is object new", isNewFn)

	fileSet.StringVar(&conf.ID, "id", "", "ID")
	fileSet.StringVar(&conf.Filename, "f", "", "filename")

	textSet := flag.NewFlagSet("text", flag.ExitOnError)

	textSet.StringVar(&conf.WorkingDir, "wd", "saved", "Working directory")
	textSet.StringVar(&conf.UserLogin, "ul", "", "User login")
	textSet.StringVar(&conf.UserPassword, "up", "", "User password")

	textSet.Func("a", "action get, save, delete", actionFn)
	textSet.Func("in", "is object new", isNewFn)

	textSet.StringVar(&conf.ID, "id", "", "ID")
	textSet.StringVar(&conf.Text, "t", "", "Text")

	cardSet := flag.NewFlagSet("card", flag.ExitOnError)

	cardSet.StringVar(&conf.WorkingDir, "wd", "saved", "Working directory")
	cardSet.StringVar(&conf.UserLogin, "ul", "", "User login")
	cardSet.StringVar(&conf.UserPassword, "up", "", "User password")

	cardSet.Func("a", "action get, save, delete", actionFn)
	cardSet.Func("in", "is object new", isNewFn)

	cardSet.StringVar(&conf.ID, "id", "", "ID")
	cardSet.StringVar(&conf.CardNum, "n", "", "Number")
	cardSet.StringVar(&conf.CardCVC, "c", "", "CVC")
	cardSet.StringVar(&conf.CardHolderName, "hn", "", "Holder name")

	credSet := flag.NewFlagSet("cred", flag.ExitOnError)

	credSet.StringVar(&conf.WorkingDir, "wd", "saved", "Working directory")
	credSet.StringVar(&conf.UserLogin, "ul", "", "User login")
	credSet.StringVar(&conf.UserPassword, "up", "", "User password")

	credSet.Func("a", "action get, save, delete", actionFn)
	credSet.Func("in", "is object new", isNewFn)

	credSet.StringVar(&conf.ID, "id", "", "ID")
	credSet.StringVar(&conf.CredentialsLogin, "l", "", "Login")
	credSet.StringVar(&conf.CredentialsPassword, "p", "", "Password")

	syncSet := flag.NewFlagSet("cred", flag.ExitOnError)
	syncSet.StringVar(&conf.WorkingDir, "wd", "saved", "Working directory")
	syncSet.StringVar(&conf.UserLogin, "ul", "", "User login")
	syncSet.StringVar(&conf.UserPassword, "up", "", "User password")

	if len(os.Args) >= 2 {
		switch os.Args[1] {
		case "file":
			err := fileSet.Parse(os.Args[2:])
			if err != nil {
				return nil, fmt.Errorf("fileSet.Parse: %w", err)
			}
			conf.IsFileFlagsParsed = true
		case "text":
			err := textSet.Parse(os.Args[2:])
			if err != nil {
				return nil, fmt.Errorf("textSet.Parse: %w", err)
			}
			conf.IsTextFlagsParsed = true
		case "card":
			err := cardSet.Parse(os.Args[2:])
			if err != nil {
				return nil, fmt.Errorf("cardSet.Parse: %w", err)
			}
			conf.IsCardFlagsParsed = true
		case "cred":
			err := credSet.Parse(os.Args[2:])
			if err != nil {
				return nil, fmt.Errorf("credSet.Parse: %w", err)
			}
			conf.IsCredentialsFlagsParsed = true
		case "sync":
			err := syncSet.Parse(os.Args[2:])
			if err != nil {
				return nil, fmt.Errorf("syncSet.Parse: %w", err)
			}
			conf.IsSync = true
		default:
			flag.PrintDefaults()
			return nil, errors.New("unknown action")
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
