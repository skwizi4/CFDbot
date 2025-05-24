package config

import (
	"encoding/json"
	"fmt"
	"os"
)

type (
	Config struct {
		CFD   CFD     // api биржи
		Mongo MongoDb // база данных (https://www.mongodb.com)
	}
	CFD struct {
		ApiKey string
		ApiUrl string
	}
	MongoDb struct {
		Uri            string `json:"Uri"`
		DataBaseName   string `json:"DataBaseName"`
		CollectionName string `json:"collection_name"`
	}
)

func ParseConfig(path string) (*Config, error) {
	file, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("error opening config file, %v", err)
	}
	fmt.Println(string(file))
	var Cfg Config
	if err = json.Unmarshal(file, &Cfg); err != nil {
		return nil, fmt.Errorf("error parsing config file, %v", err)
	}
	return &Cfg, nil
}
