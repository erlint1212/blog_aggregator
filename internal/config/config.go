package config

import (
    "encoding/json"
    "os"
    "io/ioutil"
)

const configFileName = ".gatorconfig.json"

type Config struct {
    DbUrl               string  `json:"db_url"`
    CurrentUserName     string  `json:"current_user_name"`
}


func Read() (Config, error) {
    file_path, err := getConfigFilePath()
    if err != nil {
        return Config{}, err
    }

    data, err := os.ReadFile(file_path)
    if err != nil {
        return Config{}, err
    }

    var cfg Config
    if err = json.Unmarshal(data, &cfg); err != nil {
		return Config{}, err
	}

    return cfg, nil
}


func (cfg Config) SetUser(user string) error {
    cfg.CurrentUserName = user

    err := write(cfg)
    if err != nil {
        return err
    }

    return nil
}

func getConfigFilePath() (string, error) {
    home_dir, err := os.UserHomeDir()
    if err != nil {
        return  "", err
    }
    file_path := home_dir + "/" + configFileName

    return file_path, nil
}

func write(cfg Config) error {
    cfg_json, err := json.Marshal(cfg)
    if err != nil {
        return err
    }

    file_path, err := getConfigFilePath()
    if err != nil {
        return err
    }


    err = ioutil.WriteFile(file_path, cfg_json, os.ModePerm)
    if err != nil {
        return err
    }

    return nil
}

