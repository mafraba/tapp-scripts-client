package main

import (
	"crypto/tls"
	"encoding/xml"
	"io/ioutil"
	"net/http"
	"os"
	"log"
)

type TappConfig struct {
	XMLName     xml.Name `xml:"tapp"`
	ApiEndpoint string   `xml:"server,attr"`
	LogFile     string   `xml:"log_file,attr"`
	LogLevel    string   `xml:"log_level,attr"`
	Certificate Cert     `xml:"ssl"`
}

type Cert struct {
	Cert string `xml:"cert,attr"`
	Key  string `xml:"key,attr"`
	Ca   string `xml:"server_ca,attr"`
}

func openTappConfiguration(fileLocation string) (config TappConfig) {
	xmlFile, err := os.Open(fileLocation)
	if err != nil {
		log.Println("Error opening file:", err)
		return
	}
	defer xmlFile.Close()
	b, _ := ioutil.ReadAll(xmlFile)
	// var config TappConfig
	xml.Unmarshal(b, &config)
	return config
}

func createClient(config TappConfig) (client *http.Client) {
	/**
	 * Loads Clients Certificates and creates and 509KeyPair
	 */
	cert, err := tls.LoadX509KeyPair(config.Certificate.Cert, config.Certificate.Key)
	if err != nil {
		log.Fatalln(err)
	}

	/**
	 * Creates a client with specific transport configurations
	 */
	transport := &http.Transport{
		TLSClientConfig: &tls.Config{Certificates: []tls.Certificate{cert}, InsecureSkipVerify: true},
	}
	client = &http.Client{Transport: transport}
	return client
}