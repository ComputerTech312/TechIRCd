package main

import (
	"encoding/json"
	"fmt"
	"os"
	"time"
)

type Config struct {
	Server struct {
		Name        string `json:"name"`
		Network     string `json:"network"`
		Description string `json:"description"`
		Version     string `json:"version"`
		AdminInfo   string `json:"admin_info"`
		Listen      struct {
			Host      string `json:"host"`
			Port      int    `json:"port"`
			SSLPort   int    `json:"ssl_port"`
			EnableSSL bool   `json:"enable_ssl"`
		} `json:"listen"`
		SSL struct {
			CertFile   string `json:"cert_file"`
			KeyFile    string `json:"key_file"`
			RequireSSL bool   `json:"require_ssl"`
		} `json:"ssl"`
	} `json:"server"`

	Limits struct {
		MaxClients          int `json:"max_clients"`
		MaxChannels         int `json:"max_channels"`
		MaxChannelUsers     int `json:"max_channel_users"`
		MaxNickLength       int `json:"max_nick_length"`
		MaxChannelLength    int `json:"max_channel_length"`
		MaxTopicLength      int `json:"max_topic_length"`
		MaxKickLength       int `json:"max_kick_length"`
		MaxAwayLength       int `json:"max_away_length"`
		PingTimeout         int `json:"ping_timeout"`
		RegistrationTimeout int `json:"registration_timeout"`
		FloodLines          int `json:"flood_lines"`
		FloodSeconds        int `json:"flood_seconds"`
	} `json:"limits"`

	Features struct {
		EnableOper     bool   `json:"enable_oper"`
		EnableServices bool   `json:"enable_services"`
		EnableModes    bool   `json:"enable_modes"`
		EnableCTCP     bool   `json:"enable_ctcp"`
		EnableDCC      bool   `json:"enable_dcc"`
		CaseMapping    string `json:"case_mapping"`
	} `json:"features"`

	Channels struct {
		DefaultModes  string   `json:"default_modes"`
		AutoJoin      []string `json:"auto_join"`
		AdminChannels []string `json:"admin_channels"`
		Modes         struct {
			BanListSize    int `json:"ban_list_size"`
			ExceptListSize int `json:"except_list_size"`
			InviteListSize int `json:"invite_list_size"`
		} `json:"modes"`
	} `json:"channels"`

	Opers []struct {
		Name     string   `json:"name"`
		Password string   `json:"password"`
		Host     string   `json:"host"`
		Class    string   `json:"class"`
		Flags    []string `json:"flags"`
	} `json:"opers"`

	MOTD []string `json:"motd"`

	Logging struct {
		Level      string `json:"level"`
		File       string `json:"file"`
		MaxSize    int    `json:"max_size"`
		MaxBackups int    `json:"max_backups"`
		MaxAge     int    `json:"max_age"`
	} `json:"logging"`
}

func LoadConfig(filename string) (*Config, error) {
	data, err := os.ReadFile(filename)
	if err != nil {
		return nil, fmt.Errorf("failed to read config file: %v", err)
	}

	var config Config
	if err := json.Unmarshal(data, &config); err != nil {
		return nil, fmt.Errorf("failed to parse config file: %v", err)
	}

	return &config, nil
}

func (c *Config) PingTimeoutDuration() time.Duration {
	return time.Duration(c.Limits.PingTimeout) * time.Second
}

func (c *Config) RegistrationTimeoutDuration() time.Duration {
	return time.Duration(c.Limits.RegistrationTimeout) * time.Second
}

func DefaultConfig() *Config {
	return &Config{
		Server: struct {
			Name        string `json:"name"`
			Network     string `json:"network"`
			Description string `json:"description"`
			Version     string `json:"version"`
			AdminInfo   string `json:"admin_info"`
			Listen      struct {
				Host      string `json:"host"`
				Port      int    `json:"port"`
				SSLPort   int    `json:"ssl_port"`
				EnableSSL bool   `json:"enable_ssl"`
			} `json:"listen"`
			SSL struct {
				CertFile   string `json:"cert_file"`
				KeyFile    string `json:"key_file"`
				RequireSSL bool   `json:"require_ssl"`
			} `json:"ssl"`
		}{
			Name:        "TechIRCd",
			Network:     "TechNet",
			Description: "A modern IRC server written in Go",
			Version:     "1.0.0",
			AdminInfo:   "admin@example.com",
			Listen: struct {
				Host      string `json:"host"`
				Port      int    `json:"port"`
				SSLPort   int    `json:"ssl_port"`
				EnableSSL bool   `json:"enable_ssl"`
			}{
				Host:      "localhost",
				Port:      6667,
				SSLPort:   6697,
				EnableSSL: false,
			},
			SSL: struct {
				CertFile   string `json:"cert_file"`
				KeyFile    string `json:"key_file"`
				RequireSSL bool   `json:"require_ssl"`
			}{
				CertFile:   "server.crt",
				KeyFile:    "server.key",
				RequireSSL: false,
			},
		},
		Limits: struct {
			MaxClients          int `json:"max_clients"`
			MaxChannels         int `json:"max_channels"`
			MaxChannelUsers     int `json:"max_channel_users"`
			MaxNickLength       int `json:"max_nick_length"`
			MaxChannelLength    int `json:"max_channel_length"`
			MaxTopicLength      int `json:"max_topic_length"`
			MaxKickLength       int `json:"max_kick_length"`
			MaxAwayLength       int `json:"max_away_length"`
			PingTimeout         int `json:"ping_timeout"`
			RegistrationTimeout int `json:"registration_timeout"`
			FloodLines          int `json:"flood_lines"`
			FloodSeconds        int `json:"flood_seconds"`
		}{
			MaxClients:          1000,
			MaxChannels:         100,
			MaxChannelUsers:     500,
			MaxNickLength:       30,
			MaxChannelLength:    50,
			MaxTopicLength:      307,
			MaxKickLength:       307,
			MaxAwayLength:       307,
			PingTimeout:         300,
			RegistrationTimeout: 60,
			FloodLines:          20,
			FloodSeconds:        10,
		},
		Features: struct {
			EnableOper     bool   `json:"enable_oper"`
			EnableServices bool   `json:"enable_services"`
			EnableModes    bool   `json:"enable_modes"`
			EnableCTCP     bool   `json:"enable_ctcp"`
			EnableDCC      bool   `json:"enable_dcc"`
			CaseMapping    string `json:"case_mapping"`
		}{
			EnableOper:     true,
			EnableServices: false,
			EnableModes:    true,
			EnableCTCP:     true,
			EnableDCC:      false,
			CaseMapping:    "rfc1459",
		},
		Channels: struct {
			DefaultModes  string   `json:"default_modes"`
			AutoJoin      []string `json:"auto_join"`
			AdminChannels []string `json:"admin_channels"`
			Modes         struct {
				BanListSize    int `json:"ban_list_size"`
				ExceptListSize int `json:"except_list_size"`
				InviteListSize int `json:"invite_list_size"`
			} `json:"modes"`
		}{
			DefaultModes:  "+nt",
			AutoJoin:      []string{"#general"},
			AdminChannels: []string{"#admin"},
			Modes: struct {
				BanListSize    int `json:"ban_list_size"`
				ExceptListSize int `json:"except_list_size"`
				InviteListSize int `json:"invite_list_size"`
			}{
				BanListSize:    100,
				ExceptListSize: 100,
				InviteListSize: 100,
			},
		},
		Opers: []struct {
			Name     string   `json:"name"`
			Password string   `json:"password"`
			Host     string   `json:"host"`
			Class    string   `json:"class"`
			Flags    []string `json:"flags"`
		}{
			{
				Name:     "admin",
				Password: "changeme",
				Host:     "*@localhost",
				Class:    "admin",
				Flags:    []string{"global_kill", "remote", "connect", "squit"},
			},
		},
		MOTD: []string{
			"Welcome to TechIRCd!",
			"A modern IRC server written in Go",
			"Enjoy your stay on TechNet!",
		},
		Logging: struct {
			Level      string `json:"level"`
			File       string `json:"file"`
			MaxSize    int    `json:"max_size"`
			MaxBackups int    `json:"max_backups"`
			MaxAge     int    `json:"max_age"`
		}{
			Level:      "info",
			File:       "techircd.log",
			MaxSize:    100,
			MaxBackups: 3,
			MaxAge:     28,
		},
	}
}

func SaveConfig(config *Config, filename string) error {
	data, err := json.MarshalIndent(config, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal config: %v", err)
	}

	if err := os.WriteFile(filename, data, 0644); err != nil {
		return fmt.Errorf("failed to write config file: %v", err)
	}

	return nil
}
