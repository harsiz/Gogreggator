package config

import (
	"encoding/json"
	"os"
)

const ConfigFilePath = "/root/.gatorconfig.json"

func Read() (Config, error) {
	data, err := os.ReadFile(ConfigFilePath)
	if err != nil {
		return Config{}, err
	}
	var c Config
	err = json.Unmarshal(data, &c)
	if err != nil {
		return Config{}, err
	}

	return c, nil
}

func Write(c Config) (error) {
	jsonStr, err := json.Marshal(c)
	if err != nil {
		return err
	}

	err = os.WriteFile(ConfigFilePath, jsonStr, 0664)
	if err != nil {
		return err
	}
	return nil
}

func (c Config) SetUser(name string) error {
	c.CurrentUserName = name
	if err := Write(c); err != nil {
		return err
	}
	return nil
}

