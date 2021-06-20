package bootstrap

import (
	"encoding/json"
	"io/ioutil"
	"log"
)

// Configuration model
type Configuration struct {
	AllowedOrigins  []string `json:"allowed_origins"`
	AllowedMethods  []string `json:"allowed_methods"`
	LogFolder       string   `json:"log_folder"`
	Env             string   `json:"env"`
	PortHTTP        uint     `json:"port_http"`
	CertPem         string   `json:"cert_pem"`
	KeyPem          string   `json:"key_pem"`
	PrivateFileSign string   `json:"private_file_sign"`
	PublicFileSign  string   `json:"public_file_sign"`
	Database        Database `json:"database"`
}

func newConfiguration(path string) Configuration {
	file, err := ioutil.ReadFile(path)
	if err != nil {
		log.Fatal(err)
	}

	conf := Configuration{}
	if err := json.Unmarshal(file, &conf); err != nil {
		log.Fatal(err)
	}

	return conf
}

// Database model
type Database struct {
	Engine   string `json:"engine"`
	User     string `json:"user"`
	Password string `json:"password"`
	Server   string `json:"server"`
	Port     uint   `json:"port"`
	Name     string `json:"name"`
	SSLMode  string `json:"ssl_mode"`
}

// DBEngine obtains db engine
func (c Configuration) DBEngine() string { return c.Database.Engine }

// Environment obtains the development enviroment
func (c Configuration) Environment() string { return c.Env }
