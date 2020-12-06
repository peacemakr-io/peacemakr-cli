
package main

import (
	"flag"
	peacemakr_go_sdk "github.com/peacemakr-io/peacemakr-go-sdk/pkg"
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
	var configuration PeacemakrConfig

	viper.SetConfigFile(configName)

	// Also permit environment overrides.
	viper.SetEnvPrefix("PEACEMAKR")

	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))
	viper.BindEnv("ApiKey")
	viper.AutomaticEnv() // Bind to all configs, overriding config from env when in both file and env var.

    // If no config was found, we use default values
	if err := viper.MergeInConfig(); err != nil {
		configuration = PeacemakrConfig{
				Verbose: false,
				Host: "https://api.peacemakr.io",
				PersisterFileLocation: "/tmp/.peacemakr",
				ClientName: "peacemakr-cli",
				ApiKey: viper.GetString("ApiKey"),
		}

		if configuration.Verbose {
			log.Printf("Config:\n Verbose: %v\n Host: %v\n Persister file location: %v\n Client Name: %v\n",  configuration.Verbose, configuration.Host, configuration.PersisterFileLocation,  configuration.ClientName)
		}
		return &configuration
	}

	err := viper.Unmarshal(&configuration)
	if err != nil {
		log.Fatalf("unable to read config, %v", err)
	}

	if configuration.Verbose {
		log.Printf("Config:\n Verbose: %v\n Host: %v\n Persister file location: %v\n Client Name: %v\n",  configuration.Verbose, configuration.Host, configuration.PersisterFileLocation,  configuration.ClientName)
	}

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

func canonicalAction(action *string) string {
	if action == nil {
		log.Fatalf("failed to provide an action")
	}

	actionStr := strings.ToLower(*action)

	if actionStr != "encrypt" && actionStr != "decrypt" {
		log.Fatalf("unkonwn action: ", *action)
	}

	return actionStr
}

func loadInputFile(inputFileName string) (*os.File, error) {
	var inputFile *os.File
	var err error
	if inputFileName == "" {
		inputFile = os.Stdin 
	} else {
		inputFile, err = os.Open(inputFileName)
		if err != nil {
			log.Printf("Error opening the file %v", err)
			return nil, err
		}
	}
	return inputFile, nil
}

func loadOutputFile(outputFileName string) (*os.File, error) {
	var outputFile *os.File
	var err error
	if outputFileName == "" {
		outputFile = os.Stdout
	} else {
		outputFile, err = os.OpenFile(outputFileName, os.O_CREATE|os.O_WRONLY, os.ModePerm)
		if err != nil {
			log.Printf("Error opening the file %v", err)
			return nil, err
		}
	}
	return outputFile, nil
}

type CustomLogger struct{}
func (l *CustomLogger) Printf(format string, args ...interface{}) {
	log.Printf(format, args...)
}
func main() {
	action := flag.String("action", "encrypt", "action= encrypt|decrypt")
	customConfig := flag.String("config", "peacemakr.yml", "custom config file e.g. (peacemakr.yml)")
	inputFileName := flag.String("inputFileName", "", "inputFile to encrypt/decrypt")
	outputFileName := flag.String("outputFileName", "", "outputFile to encrypt/decrypt")
	flag.Parse()

	actionStr := canonicalAction(action)

	config := LoadConfigs(*customConfig)

	if config.Verbose {
		log.Println("Finish parsing flag and config")
		log.Printf("inputfilename: %s, OutputFilename: %s", *inputFileName, *outputFileName)
	}

	if config.Verbose {
		log.Println("Setting up SDK...")
	}

	if _, err := os.Stat(config.PersisterFileLocation); os.IsNotExist(err) {
		os.MkdirAll(config.PersisterFileLocation, os.ModePerm)
	}

	sdk, err := peacemakr_go_sdk.GetPeacemakrSDK(
		config.ApiKey,
		config.ClientName,
		&config.Host,
		GetDiskPersister(config.PersisterFileLocation),
		log.New(os.Stdout, "MyProjectCrypto", log.LstdFlags))


	if err != nil {
		log.Fatalf("Failed to create peacemakr sdk due to %v", err)
	}


	inputFile, err := loadInputFile(*inputFileName)
	if err != nil {
		log.Fatalf("Error loading input file", err)
	}
	outputFile, err := loadOutputFile(*outputFileName)
	if err != nil {
		log.Fatalf("Error loading output file", err)
	}

	if config.Verbose {
		log.Printf("registering client")
		registerOrFail(sdk)
	}

	if actionStr == "encrypt" {
		if config.Verbose {
			log.Println("In encrypting")
		}
		for err = sdk.Register(); err != nil; {
			log.Println("Encrypting client, failed to register, trying again...")
			time.Sleep(time.Duration(1) * time.Second)
		}

		if config.Verbose {
			log.Println("Encrypting")
		}

		encryptOrFail(sdk, inputFile, outputFile)
	} else if actionStr == "decrypt" {
		if config.Verbose {
			log.Println("In decrypting")
		}
		for err = sdk.Register(); err != nil; {
			log.Println("Decrypting client, failed to register, trying again...")
			time.Sleep(time.Duration(1) * time.Second)
		}
		decryptOrFail(sdk, inputFile, outputFile)
	}
}
