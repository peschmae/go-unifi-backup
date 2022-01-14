package main

import (
	"bytes"
	"crypto/tls"
	"encoding/json"
	"fmt"
	"github.com/spf13/pflag"
	"github.com/spf13/viper"
	"io"
	"io/ioutil"
	"net/http"
	"net/http/cookiejar"
	"os"
	"time"
	c "unifi-backup/src/config"
)

var client http.Client

func init() {
	// this is rather minimalist and can be insecure. But we are only accessing one server so we're okay
	jar, err := cookiejar.New(nil)
	if err != nil {
		fmt.Printf("Got error while creating cookie jar %s", err.Error())
	}
	// the Unifi controller uses a self-signed certificate by default, so we skip the verification
	tr := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}
	client = http.Client{
		Jar:       jar,
		Transport: tr,
	}
}

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
		fmt.Println("")
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

	if debug {
		fmt.Printf("Trying to login to Unifi Controller at %v\n", configuration.Unifi.Server)
	}

	loginValues := map[string]string{"username": configuration.Unifi.Username, "password": configuration.Unifi.Password}
	jsonData, err := json.Marshal(loginValues)
	if err != nil {
		fmt.Printf("Error while unmarshaling: %v", err)
	}
	res, err := client.Post(configuration.Unifi.Server+"/api/login", "application/json", bytes.NewBuffer(jsonData))
	if err != nil {
		fmt.Printf("Error while making login request: %v", err)
	}
	if res.StatusCode != 200 {
		fmt.Println("Couldn't login to the unifi controller")
		if debug {
			body, _ := ioutil.ReadAll(res.Body)
			bodyString := string(body)
			fmt.Println(bodyString)
		}
		os.Exit(1)
	} else if debug {
		fmt.Printf("Successfully logged in")
	}

	if debug {
		fmt.Println("Trying to trigger backup creation")
	}
	backupValues := map[string]interface{}{"cmd": "async-backup", "days": 0}
	jsonData, err = json.Marshal(backupValues)

	if err != nil {
		fmt.Printf("Error while unmarshaling: %v", err)
	}
	res, err = client.Post(configuration.Unifi.Server+"/api/s/default/cmd/system", "application/json", bytes.NewBuffer(jsonData))
	if err != nil {
		fmt.Printf("Error while making login request: %v", err)
	}
	if res.StatusCode != 200 {
		fmt.Printf("Couldn't trigger creation of backup")
		os.Exit(1)
	} else if debug {
		fmt.Println("Triggered backup creation")
	}

	if debug {
		fmt.Println("Downloading file")
	}

	// Get the data
	res, err = client.Get(configuration.Unifi.Server + "/dl/backup/" + configuration.Unifi.ControllerVersion + ".unf")
	if err != nil {
		fmt.Println(err)
	}
	defer res.Body.Close()

	var filepath = configuration.Backup.OutputDirectory + time.Now().Format("2006-01-02-15-04") + ".unf"
	// Create the file
	out, err := os.Create(filepath)
	if err != nil {
		fmt.Println(err)
	}
	defer out.Close()

	// Write the body to file
	_, err = io.Copy(out, res.Body)

	fmt.Printf("Created backup at %v\n", filepath)

}
