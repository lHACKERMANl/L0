package config

import (
	"log"
	"os"

	"gopkg.in/yaml.v3"
)

type LoginData struct {
	Services struct {
		Postgres struct {
			Image         string `yaml:"image"`
			ContainerName string `yaml:"container_name"`
			Environment   struct {
				PostgresDB       string `yaml:"POSTGRES_DB"`
				PostgresUser     string `yaml:"POSTGRES_USER"`
				PostgresPassword string `yaml:"POSTGRES_PASSWORD"`
			} `yaml:"environment"`
			Ports []string `yaml :"ports"`
		} `yaml:"postgres"`
		Nats struct {
			DomainName    string `yaml:"domainname"`
			Image         string `yaml:"image"`
			ContainerName string `yaml:"container_name"`
			Ports         string `yaml:"porst"`
		} `yaml: "nats"`
	} `yaml:"services"`
}

func GetDataFromDockerCompose(path string) (LoginData, error) {
	dockerFile, err := os.Open(path)
	if err != nil {
		log.Fatal(err)
	}
	defer dockerFile.Close()

	var config LoginData
	decoder := yaml.NewDecoder(dockerFile)
	err = decoder.Decode(&config)
	if err != nil {
		log.Fatal(err)
		return config, err
	}

	return config, nil
}
