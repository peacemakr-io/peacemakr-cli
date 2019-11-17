
package main

import (
	"flag"
	peacemakr_go_sdk "github.com/peacemakr-io/peacemakr-go-sdk/pkg"
	"github.com/peacemakr-io/peacemakr-go-sdk/pkg/utils"
	"github.com/spf13/viper"
	"io/ioutil"
	"log"
	"os"
	"strings"
	"time"
)

type PeacemakrConfig struct {
	Verbose               bool
	Host                  string
	ApiKey	              string
	PersisterFileLocation string
	ClientName            string
}

func LoadConfigs(configName string) *PeacemakrConfig {
	viper.SetConfigName(configName)
	viper.AddConfigPath(".")
	viper.SetConfigType("yml")

	// Also permit environment overrides.
	viper.SetEnvPrefix("PEACEMAKR_")

	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))

	viper.AutomaticEnv() // Bind to all configs, overriding config from env when in both file and env var.

	var configuration PeacemakrConfig

	if err := viper.MergeInConfig(); err != nil {
		log.Fatalf("Error reading config, %v", err)
	}

	err := viper.Unmarshal(&configuration)
	if err != nil {
		log.Fatalf("unable to read config, %v", err)
	}
	log.Printf("Successfully read in config")

	log.Println("Config: ", configuration)

	return &configuration
}

func encryptOrFail(sdk peacemakr_go_sdk.PeacemakrSDK, from, to *os.File) {
	if from == nil {
		log.Fatalf("missing 'from' in encryption")
	}

	if to == nil {
		log.Fatalf("missing 'to' in encryption")
	}

	if from == to {
		log.Fatalf("in-place encryption is not supproted (from and to are the same)")
	}

	data, err := ioutil.ReadAll(from)
	if err != nil {
		log.Fatalf("failed to read stdin due to error %v", err)
	}


	encryptedData, err := sdk.Encrypt(data)
	if err != nil {
		log.Fatalf("failed to encrypted due to error %v", err)
	}

	_, err = to.Write(encryptedData)
	if err != nil {
		log.Fatalf("failed to write encrypted data due to error %v", err)
	}
}

func decryptOrFail(sdk peacemakr_go_sdk.PeacemakrSDK, from, to *os.File) {
	if from == nil {
		log.Fatalf("missing 'from' in decryption")
	}

	if to == nil {
		log.Fatalf("missing 'to' in decryption")
	}

	if from == to {
		log.Fatalf("in-place decryption is not supproted (from and to are the same)")
	}

	data, err := ioutil.ReadAll(from)
	if err != nil {
		log.Fatalf("failed to read stdin due to error %v", err)
	}


	decryptedData, err := sdk.Decrypt(data)
	if err != nil {
		log.Fatalf("failed to decrypt due to error %v", err)
	}

	_, err = to.Write(decryptedData)
	if err != nil {
		log.Fatalf("failed to write decrypted data due to error %v", err)
	}
}

func registerOrFail(sdk peacemakr_go_sdk.PeacemakrSDK) {
	err := sdk.Register()
	if err != nil {
		log.Fatalf(" failed to register due to %v", err)
	}
}

type CustomLogger struct{}
func (l *CustomLogger) Printf(format string, args ...interface{}) {
	log.Printf(format, args...)
}
func main() {
	config := LoadConfigs("peacemakr")

	persister := utils.GetDiskPersister(config.PersisterFileLocation)

	if config.Verbose {
		log.Println("Setting up SDK...")
	}

	sdk, err := peacemakr_go_sdk.GetPeacemakrSDK(config.ApiKey, config.ClientName, &config.Host, persister, &CustomLogger{}, false)
	if err != nil {
		log.Fatalf("Failed to create peacemakr sdk due to %v", err)
	}

	action := flag.String("action", "encrypt", "action= encrypt|decrypt")
	flag.Parse()

	if action == nil {
		log.Fatalf("Failed to provide an action")
	}

	actionStr := strings.ToLower(*action)

	if actionStr == "encrypt" {
		registerOrFail(sdk)
		encryptOrFail(sdk, os.Stdin, os.Stdout)
	} else if actionStr == "decrypt" {
		registerOrFail(sdk)
		decryptOrFail(sdk, os.Stdin, os.Stdout)
	} else {
		log.Fatalf("Unknown action specified %s", *action)
	}
}
