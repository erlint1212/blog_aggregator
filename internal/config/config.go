package config

import (
    "encoding/json"
    "os"
)

const configFileName = ".gatorconfig.json"

type Config struct {
    db_url              string  `json:"db_url"`
    current_user_name   string  `json:"current_user_name"`
}

func (cfg Config) SetUser(user string) {
    cfg.current_user_name = user
}

func Write(cfg Config) error {
    home_dir, err := os.UserHomeDir()
    if err != nil {
        return  err
    }

    cfg_json, err := json.Marshal(cfg)
    if err != nil {
        return err
    }

    jsonFile, err := os.Open(home_dir + "/" + configFileName)
    if err != nil {
        return err
    }
    defer jsonFile.Close()

    cfg_json_byte := []byte(cfg_json)

    _, err = jsonFile.Write(cfg_json_byte)
    if err != nil {
        return err
    }

    return nil
}

func Read() (Config, error) {
    home_dir, err := os.UserHomeDir()
    if err != nil {
        return Config{}, err
    }

    data, err := os.ReadFile(home_dir + "/" + configFileName)
    if err != nil {
        return Config{}, err
    }

    var config Config
    if err = json.Unmarshal(data, &config); err != nil {
		return Config{}, err
	}

    return config, nil
}

