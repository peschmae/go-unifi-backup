package main

import (
	"fmt"
	"github.com/spf13/pflag"
	"github.com/spf13/viper"
	"os"
	c "unifi-backup/src/config"
)

func main() {

	// Set the file name of the configurations file
	viper.SetConfigName("config")
	// Set the path to look for the configurations file
	viper.AddConfigPath(".")
	viper.SetConfigType("yml")
	var configuration c.Configuration
	var debug bool

	pflag.BoolVar(&debug, "debug", false, "Output debug logs")
	pflag.Parse()

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

	if debug {
		// Reading variables using the model
		fmt.Println("Unifi config")
		fmt.Println("Server:\t\t", configuration.Unifi.Server)
		fmt.Println("Username:\t", configuration.Unifi.Username)
		fmt.Println("Password:\t", configuration.Unifi.Password)
		fmt.Println("Version:\t", configuration.Unifi.ControllerVersion)
		fmt.Println("\nBackup config")
		fmt.Println("Output directory:\t\t", configuration.Backup.OutputDirectory)
		fmt.Println("Number of backups to keep:\t", configuration.Backup.Keep)
	}

	if _, err := os.Stat(configuration.Backup.OutputDirectory); err != nil {
		if os.IsNotExist(err) {
			// path/to/whatever does not exist
			err := os.MkdirAll(configuration.Backup.OutputDirectory, os.ModePerm)
			if err != nil {
				fmt.Printf("Couldn't create directory: %v\n", configuration.Backup.OutputDirectory)
				fmt.Println(err)
				os.Exit(1)
			}
		}
		if os.IsPermission(err) {
			fmt.Printf("Access denied to: %v\n", configuration.Backup.OutputDirectory)
			if debug {
				fmt.Println(err)
			}
			os.Exit(1)
		}
	}

}
