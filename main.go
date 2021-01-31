package main

import (
	"flag"
	peacemakr_go_sdk "github.com/peacemakr-io/peacemakr-go-sdk/pkg"
	"github.com/spf13/viper"
	"io/ioutil"
	"log"
	"os"
	"strings"
)

type PeacemakrConfig struct {
	Verbose               bool
	Host                  string
	ApiKey                string
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
			Verbose:               false,
			Host:                  "https://api.peacemakr.io",
			PersisterFileLocation: "/tmp/.peacemakr",
			ClientName:            "peacemakr-cli",
			ApiKey:                viper.GetString("ApiKey"),
		}

		if configuration.Verbose {
			log.Printf("Config:\n Verbose: %v\n Host: %v\n Persister file location: %v\n Client Name: %v\n", configuration.Verbose, configuration.Host, configuration.PersisterFileLocation, configuration.ClientName)
		}
		return &configuration
	}

	err := viper.Unmarshal(&configuration)
	if err != nil {
		log.Fatalf("unable to read config, %v", err)
	}

	if configuration.Verbose {
		log.Printf("Config:\n Verbose: %v\n Host: %v\n Persister file location: %v\n Client Name: %v\n", configuration.Verbose, configuration.Host, configuration.PersisterFileLocation, configuration.ClientName)
	}

	return &configuration
}

func encryptOrFailCommon(sdk peacemakr_go_sdk.PeacemakrSDK, useDomain string, inputData []byte, to *os.File) {
	var encryptedData []byte
	var err error
	if useDomain != "" {
		encryptedData, err = sdk.EncryptInDomain(inputData, useDomain)
	} else {
		encryptedData, err = sdk.Encrypt(inputData)
	}

	if err != nil {
		log.Fatalf("failed to encrypted due to error %v", err)
	}

	_, err = to.Write(encryptedData)
	if err != nil {
		log.Fatalf("failed to write encrypted data due to error %v", err)
	}
}

func encryptOrFail(sdk peacemakr_go_sdk.PeacemakrSDK, inputData []byte, to *os.File) {
	encryptOrFailCommon(sdk, "", inputData, to)
}

func encryptOrFailInDomain(sdk peacemakr_go_sdk.PeacemakrSDK, useDomain string, inputData []byte, to *os.File) {
	if useDomain == "" {
		log.Fatalf("Attemping to encrypt using a specific use domain, yet the use domain was not provided")
	}
	encryptOrFailCommon(sdk, useDomain, inputData, to)
}

func decryptOrFail(sdk peacemakr_go_sdk.PeacemakrSDK, inputData []byte, to *os.File) {
	decryptedData, err := sdk.Decrypt(inputData)
	if err != nil {
		log.Fatalf("failed to decrypt due to error %v", err)
	}

	_, err = to.Write(decryptedData)
	if err != nil {
		log.Fatalf("failed to write decrypted data due to error %v", err)
	}
}

func validatePeacemakrCiphertext(sdk peacemakr_go_sdk.PeacemakrSDK, inputData []byte) {
	isPeacemakrCiphertext := sdk.IsPeacemakrCiphertext(inputData)
	if isPeacemakrCiphertext {
		log.Println("Is a Peacemakr ciphertext")
		// Exit successfully
		os.Exit(0)
	} else {
		// Exit error
		log.Fatalf("Is not a Peacemakr ciphertext")
	}
}

func signOrFail(sdk peacemakr_go_sdk.PeacemakrSDK, inputData []byte, to *os.File) {
	signedBlob, err := sdk.SignOnly(inputData)
	if err != nil {
		log.Fatalf("failed to sign due to error %v", err)
	}

	_, err = to.Write(signedBlob)
	if err != nil {
		log.Fatalf("failed to write signedBlob due to error %v", err)
	}
}

func verifyOrFail(sdk peacemakr_go_sdk.PeacemakrSDK, inputData []byte, to *os.File) {
	verifiedMessage, err := sdk.VerifyOnly(inputData)
	if err != nil {
		log.Fatalf("failed to verify due to error %v", err)
	}

	_, err = to.Write(verifiedMessage)
	if err != nil {
		log.Fatalf("failed to write verifiedMessage due to error %v", err)
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

func isFlagPassed(name string) bool {
	found := false
	flag.Visit(func(f *flag.Flag) {
		if f.Name == name {
			found = true
		}
	})
	return found
}

func main() {
	customConfig := flag.String("config", "peacemakr.yml", "custom config file e.g. (peacemakr.yml)")
	inputFileName := flag.String("inputFileName", "", "inputFile to encrypt/decrypt")
	outputFileName := flag.String("outputFileName", "", "outputFile to encrypt/decrypt")
	shouldEncrypt := flag.Bool("encrypt", false, "Should the application encrypt the message")
	shouldDecrypt := flag.Bool("decrypt", false, "Should the application decrypt the ciphertext")
	shouldSign := flag.Bool("signOnly", false, "Should the application sign the message")
	shouldVerify := flag.Bool("verifyOnly", false, "Should the application verify the input blob")
	useDomain := flag.String("domain", "", "A specific use domain to encrypt; `-domain=DOMAIN_NAME`")
	shouldValidateCiphertext := flag.Bool("is-peacemakr-ciphertext", false, "Should the application "+
		"validate whether the ciphertext is a Peacemakr ciphertext or not")

	flag.Parse()

	config := LoadConfigs(*customConfig)

	if config.ApiKey == "" {
		log.Fatal("Must provide an API key!")
	}

	if shouldEncrypt == nil && shouldDecrypt == nil && shouldValidateCiphertext == nil && shouldSign == nil && shouldVerify == nil {
		log.Fatal("Must specify either encrypt OR decrypt OR is-peacemakr-ciphertext")
	}

	if shouldEncrypt != nil && shouldDecrypt != nil && shouldValidateCiphertext != nil && shouldSign != nil && shouldVerify != nil &&
		*shouldEncrypt && *shouldDecrypt && *shouldValidateCiphertext && *shouldSign && *shouldVerify {
		log.Fatal("Must not attempt multiple functions simultaneously")
	}

	// if encrypting in a specific use domain, ensure the use domain value was set
	domainFlagIsSet := isFlagPassed("domain")
	if domainFlagIsSet && useDomain != nil && *useDomain == "" && shouldEncrypt != nil && *shouldEncrypt {
		log.Fatal("An attempt was made to use a specific use domain to encrypt, yet the use domain name was not provided in " +
			"cli, consider -domain=DOMAIN_NAME")
	}

	if config.Verbose {
		log.Println("Finish parsing flag and config")
		log.Printf("Input Filename: %s, Output Filename: %s", *inputFileName, *outputFileName)
	}

	if config.Verbose {
		log.Println("Setting up SDK...")
	}

	if _, err := os.Stat(config.PersisterFileLocation); os.IsNotExist(err) {
		if err := os.MkdirAll(config.PersisterFileLocation, os.ModePerm); err != nil {
			log.Fatalf("Unable to create persister directory: %v", err)
		}
	}

	logger := log.New(os.Stderr, "Peacemakr CLI", log.LstdFlags)
	if !config.Verbose {
		logger.SetOutput(ioutil.Discard)
	}
	sdk, err := peacemakr_go_sdk.GetPeacemakrSDK(
		config.ApiKey,
		config.ClientName,
		&config.Host,
		GetDiskPersister(config.PersisterFileLocation),
		logger)

	if err != nil {
		log.Fatalf("Failed to create peacemakr sdk due to %v", err)
	}

	inputFile, err := loadInputFile(*inputFileName)
	if err != nil {
		log.Fatalf("Error loading input file", err)
	}

	if inputFile == nil {
		log.Fatalf("Error input file pointer is null")
	}

	outputFile, err := loadOutputFile(*outputFileName)
	if err != nil {
		log.Fatalf("Error loading output file", err)
	}

	if outputFile == nil {
		log.Fatalf("Error output file pointer is null")
	}

	if inputFile == outputFile {
		log.Fatalf("Error in-place operation is not supproted (input and output file pointer are the same)")
	}

	if config.Verbose {
		log.Printf("registering client")
	}

	registerOrFail(sdk)

	data, err := ioutil.ReadAll(inputFile)
	if err != nil {
		log.Fatalf("failed to read stdin due to error %v", err)
	}

	if shouldEncrypt != nil && *shouldEncrypt {
		if config.Verbose {
			log.Println("Encrypting")
		}
		if domainFlagIsSet && useDomain != nil && *useDomain != "" {
			if config.Verbose {
				log.Println("A use domain has been set to: " + *useDomain)
			}
			encryptOrFailInDomain(sdk, *useDomain, data, outputFile)
		} else {
			encryptOrFail(sdk, data, outputFile)
		}
	} else if shouldDecrypt != nil && *shouldDecrypt {
		if config.Verbose {
			log.Println("Decrypting")
		}

		decryptOrFail(sdk, data, outputFile)
	} else if shouldValidateCiphertext != nil && *shouldValidateCiphertext {
		if config.Verbose {
			log.Println("Validating the ciphertext is a Peacemakr ciphertext")
		}

		validatePeacemakrCiphertext(sdk, data)
	} else if shouldSign != nil && *shouldSign {
		if config.Verbose {
			log.Println("Signing the message")
		}

		signOrFail(sdk, data, outputFile)
	} else if shouldVerify != nil && *shouldVerify {
		if config.Verbose {
			log.Println("Verifying the message")
		}

		verifyOrFail(sdk, data, outputFile)
	}
}
