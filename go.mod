module github.com/peacemakr-io/peacemakr-cli

go 1.12

require (
	github.com/kennygrant/sanitize v1.2.4
	github.com/patrickmn/go-cache v2.1.0+incompatible
	// poiting to a specific commit: `go get github.com/peacemakr-io/peacemakr-go-sdk@77efbe3bd32f05b738397473ced76dc31bf08c99`
	github.com/peacemakr-io/peacemakr-go-sdk v0.0.11-0.20201217045855-77efbe3bd32f
	github.com/spf13/viper v1.7.1
)
