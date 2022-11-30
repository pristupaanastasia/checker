package config

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/pristupaanastasia/checker/logger"
	"io/ioutil"
	"os"
	"strconv"
	"time"
)

type Config struct {
	StartParser string `json:"start-parser"`
	Goroutine   int    `json:"goroutine"`
	Url         string `json:"url"`
	IdTelegram  int64  `json:"id-telegram"`
	TokenTg     string `json:"token-tg"`
	Port        string `json:"port"`
	log         logger.Log
}

func NewConfig(log logger.Log) (*Config, error) {
	var config *Config

	config = &Config{log: log}
	err := config.Parse()
	return config, err
}

func (c *Config) Update(ctx context.Context) {
	ticker := time.Tick(30 * time.Second)
	for {
		select {
		case <-ticker:
			err := c.Parse()
			if err != nil {
				c.log.Error(err)
			}
		case <-ctx.Done():
			return
		}

	}
}
func (c *Config) Parse() error {
	c.ParseEnv()
	c.ParseJson()
	if c.IdTelegram == 0 || c.TokenTg == "" ||
		c.Port == "" || c.Url == "" || c.Goroutine == 0 || c.StartParser == "" {
		return errors.New(fmt.Sprintf("Some config options are empty %+v", c))
	}
	return nil
}
func (c *Config) ParseEnv() {
	if c.Goroutine == 0 {
		c.Goroutine, _ = strconv.Atoi(os.Getenv("GOROUTINE"))
	}
	if c.Url == "" {
		c.Url = os.Getenv("URL")
	}
	if c.IdTelegram == 0 {
		id, _ := strconv.Atoi(os.Getenv("ID_TELEGRAM"))
		c.IdTelegram = int64(id)
	}
	if c.TokenTg == "" {
		c.TokenTg = os.Getenv("TOKEN_TELEGRAM")
	}
	if c.Port == "" {
		c.Port = os.Getenv("PORT")
	}

	if c.StartParser == "" {
		c.StartParser = os.Getenv("START_PARSER")
	}
}

func (c *Config) ParseJson() {
	payload := Config{}

	f, err := os.Open("config.json")
	if err != nil {
		return
	}
	buf, err := ioutil.ReadAll(f)
	if err != nil {
		return
	}
	if err = json.Unmarshal(buf, &payload); err != nil {
		return
	}
	if payload.Url != "" {
		c.Url = payload.Url
	}
	if payload.StartParser != "" {
		c.StartParser = payload.StartParser
	}
	if payload.Goroutine != 0 {
		c.Goroutine = payload.Goroutine
	}
	if payload.Port != "" {
		c.Port = payload.Port
	}
	if payload.IdTelegram != 0 {
		c.IdTelegram = payload.IdTelegram
	}
	if payload.TokenTg != "" {
		c.TokenTg = payload.TokenTg
	}
}
