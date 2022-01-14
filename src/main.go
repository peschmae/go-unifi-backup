package main

import (
	"fmt"
	"github.com/spf13/viper"
	c "unifi-backup/src/config"
)

func main() {
	// Set the file name of the configurations file
	viper.SetConfigName("config")

	// Set the path to look for the configurations file
	viper.AddConfigPath(".")

	// Enable VIPER to read Environment Variables
	viper.AutomaticEnv()

	viper.SetConfigType("yml")
	var configuration c.Configuration

	if err := viper.ReadInConfig(); err != nil {
		fmt.Printf("Error reading config file, %s", err)
	}

	// Set undefined variables
	viper.SetDefault("backup.keep", "10")
	viper.SetDefault("backup.output_directory", "/tmp/unifi-backup/")

	err := viper.Unmarshal(&configuration)
	if err != nil {
		fmt.Printf("Unable to decode into struct, %v", err)
	}

	// Reading variables using the model
	fmt.Println("Reading variables using the model..")
	fmt.Println("Unifi config")
	fmt.Println("Server is\t", configuration.Unifi.Server)
	fmt.Println("Username is\t\t", configuration.Unifi.Username)
	fmt.Println("Password is\t\t", configuration.Unifi.Password)
	fmt.Println("Version is\t\t", configuration.Unifi.ControllerVersion)
	fmt.Println("Backup config")
	fmt.Println("Output directory is\t\t", configuration.Backup.OutputDirectory)
	fmt.Println("Number of backups to keep is\t\t", configuration.Backup.Keep)

}
