package core

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/open-lms-test-functionality/schemas"
)

var Config schemas.ProjectConfiguration

func ReadEnvFile() {
	environment := os.Getenv("ENVIRONMENT")
	if environment == "production" {
		fmt.Println("production environment", os.Getenv("DATA"))
		err := json.Unmarshal([]byte(os.Getenv("DATA")), &Config)
		if err != nil {
			fmt.Println("Error while reading env data into struct :: ", err)
			os.Exit(1)
		}
	} else {
		currentWorkingDirectory, err := os.Getwd()
		if err != nil {
			fmt.Println("Error while getting current working directory :: ", err)
			os.Exit(1)
		}
		fmt.Println(currentWorkingDirectory)
		fileName := currentWorkingDirectory + "/.env"

		file, err := os.ReadFile(fileName)
		if err != nil {
			fmt.Println("Error while reading env file :: ", err)
			os.Exit(1)
		}

		err = json.Unmarshal(file, &Config)
		if err != nil {
			fmt.Println("Error while converting env file data into struct :: ", err)
			os.Exit(1)
		}
	}
	// currentWorkingDirectory, err := os.Getwd()
	// if err != nil {
	// 	fmt.Println("Error while getting current working directory :: ", err)
	// 	os.Exit(1)
	// }
	// fmt.Println(currentWorkingDirectory)
	// fileName := currentWorkingDirectory + "/.env"

	// file, err := os.ReadFile(fileName)
	// if err != nil {
	// 	fmt.Println("Error while reading env file :: ", err)
	// 	os.Exit(1)
	// }

	// err = json.Unmarshal(file, &Config)
	// if err != nil {
	// 	fmt.Println("Error while converting env file data into struct :: ", err)
	// 	os.Exit(1)
	// }

	Config.DBString = fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=%s",
		Config.DBConfig.DBUser,
		Config.DBConfig.DBPassword,
		Config.DBConfig.DBHost,
		Config.DBConfig.DBPort,
		Config.DBConfig.DBName,
		Config.DBConfig.DBSSLMode,
	)

	fmt.Println(Config)

}
