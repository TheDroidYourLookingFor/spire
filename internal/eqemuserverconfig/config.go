package eqemuserverconfig

import (
	"encoding/json"
	"github.com/Akkadius/spire/internal/logger"
	"github.com/Akkadius/spire/internal/pathmgmt"
	"os"
	"path/filepath"
	"time"
)

// Config is the struct type
type Config struct {
	logger   *logger.AppLogger
	pathmgmt *pathmgmt.PathManagement
	config   EQEmuConfigJson
}

// NewConfig creates a new Config struct
func NewConfig(logger *logger.AppLogger, pathmgmt *pathmgmt.PathManagement) *Config {
	return &Config{
		logger:   logger,
		pathmgmt: pathmgmt,
	}
}

// DatabaseConfig is the struct that represents the database configuration in eqemu_config.json
type DatabaseConfig struct {
	Db       string `json:"db,omitempty"`
	Host     string `json:"host,omitempty"`
	Port     string `json:"port,omitempty"`
	Username string `json:"username,omitempty"`
	Password string `json:"password,omitempty"`
}

type WebAdminLauncherConfig struct {
	RunSharedMemory             bool   `json:"runSharedMemory"`
	RunLoginserver              bool   `json:"runLoginserver"`
	RunQueryServ                bool   `json:"runQueryServ"`
	IsRunning                   bool   `json:"isRunning"`
	MinZoneProcesses            int    `json:"minZoneProcesses,omitempty"`
	StaticZones                 string `json:"staticZones,omitempty"`
	UpdateOpcodesOnStart        bool   `json:"updateOpcodesOnStart"`
	DeleteLogFilesOlderThanDays int    `json:"deleteLogFilesOlderThanDays"`
}

// EQEmuConfigJson is the struct that represents the eqemu_config.json file
type EQEmuConfigJson struct {
	Server struct {
		Discord *struct {
			Channelid    string `json:"channelid"`
			Itemurl      string `json:"itemurl"`
			Refreshrate  string `json:"refreshrate"`
			Clientid     string `json:"clientid"`
			Clientsecret string `json:"clientsecret"`
			Serverid     string `json:"serverid"`
			Username     string `json:"username"`
		} `json:"discord,omitempty"`
		Zones struct {
			Defaultstatus string `json:"defaultstatus"`
			Ports         struct {
				Low  string `json:"low"`
				High string `json:"high"`
			} `json:"ports"`
		} `json:"zones,omitempty"`
		Qsdatabase *DatabaseConfig `json:"qsdatabase,omitempty"`
		// deprecated
		Chatserver *struct {
			Port string `json:"port"`
			Host string `json:"host"`
		} `json:"chatserver,omitempty"`
		// deprecated
		Mailserver *struct {
			Host string `json:"host"`
			Port string `json:"port"`
		} `json:"mailserver,omitempty"`
		// replaces chatserver and mailserver
		Ucs *struct {
			Host string `json:"host"`
			Port string `json:"port"`
		} `json:"ucs,omitempty"`
		World struct {
			AutoDatabaseUpdates *bool `json:"auto_database_updates,omitempty"`
			Locked              bool  `json:"locked,omitempty"`
			API                 struct {
				Enabled bool `json:"enabled"`
			} `json:"api,omitempty"`
			Address      string `json:"address,omitempty"`
			Localaddress string `json:"localaddress,omitempty"`
			Loginserver1 *struct {
				Port     string `json:"port"`
				Account  string `json:"account"`
				Password string `json:"password"`
				Host     string `json:"host"`
				Legacy   string `json:"legacy,omitempty"`
			} `json:"loginserver1,omitempty"`
			Loginserver2 *struct {
				Account  string `json:"account"`
				Password string `json:"password"`
				Legacy   string `json:"legacy,omitempty"`
				Host     string `json:"host"`
				Port     string `json:"port"`
			} `json:"loginserver2,omitempty"`
			TCP struct {
				IP   string `json:"ip"`
				Port string `json:"port"`
			} `json:"tcp,omitempty"`
			Telnet struct {
				IP      string `json:"ip"`
				Port    string `json:"port"`
				Enabled string `json:"enabled"`
			} `json:"telnet,omitempty"`
			Key          string `json:"key"`
			Shortname    string `json:"shortname"`
			Longname     string `json:"longname"`
			Loginserver3 *struct {
				Account  string `json:"account"`
				Password string `json:"password"`
				Legacy   string `json:"legacy,omitempty"`
				Host     string `json:"host"`
				Port     string `json:"port"`
			} `json:"loginserver3,omitempty"`
		} `json:"world"`
		Database        *DatabaseConfig `json:"database"`
		ContentDatabase *DatabaseConfig `json:"content_database,omitempty"`
		Files           struct {
			Opcodes     string `json:"opcodes"`
			MailOpcodes string `json:"mail_opcodes"`
		} `json:"files"`
		Directories struct {
			Patches string `json:"patches"`
			Opcodes string `json:"opcodes"`
		} `json:"directories"`
	} `json:"server"`
	WebAdmin *struct { // Occulus
		Discord *struct {
			CrashLogWebhook string `json:"crash_log_webhook,omitempty"`
		} `json:"discord,omitempty"`
		Application struct {
			Key   string `json:"key,omitempty"`
			Admin struct {
				Password string `json:"password,omitempty"`
			} `json:"admin,omitempty"`
		} `json:"application,omitempty"`
		Launcher *WebAdminLauncherConfig `json:"launcher,omitempty"`
		Quests   struct {
			HotReload bool `json:"hotReload"`
		} `json:"quests"`
		ServerCodePath string `json:"serverCodePath,omitempty"`
	} `json:"web-admin,omitempty"`
	Spire struct {
		EncryptionKey string `json:"encryption_key,omitempty"`
		HttpPort      int    `json:"http_port,omitempty"`
		LauncherStart bool   `json:"launcher_start"` // starts server launcher
	} `json:"spire,omitempty"`
}

var cachedConfig *EQEmuConfigJson
var lastModifiedTime time.Time

// Get returns the eqemu config json
// If the file has been modified since the last read, it will re-read the file
// We utilize a cache to prevent reading the file on every call
func (e *Config) Get() EQEmuConfigJson {
	configFile := e.pathmgmt.GetEQEmuServerConfigFilePath()
	if len(configFile) > 0 {
		stat, _ := os.Stat(configFile)
		if stat.ModTime().After(lastModifiedTime) || lastModifiedTime.IsZero() {
			e.logger.Debug().Any("path", configFile).Msg("Reading eqemu config file")
			body, err := os.ReadFile(e.pathmgmt.GetEQEmuServerConfigFilePath())
			if err != nil {
				e.logger.Fatal().Err(err).Any("path", configFile).Msg("unable to read eqemu config file")
			}

			config := EQEmuConfigJson{}
			err = json.Unmarshal(body, &config)
			if err != nil {
				e.logger.Fatal().Err(err).Any("path", configFile).Msg("unable to unmarshal eqemu config file")
			}

			lastModifiedTime = stat.ModTime()
			cachedConfig = &config

			return config
		} else if cachedConfig != nil {
			return *cachedConfig
		}
	}

	return EQEmuConfigJson{}
}

// Exists will return true if the eqemu_config.json file exists
func (e *Config) Exists() bool {
	return len(e.pathmgmt.GetEQEmuServerConfigFilePath()) > 0
}

// Save will save the config to the eqemu_config.json file
func (e *Config) Save(c EQEmuConfigJson) error {
	if c.WebAdmin != nil && c.WebAdmin.Discord != nil {
		if len(c.WebAdmin.Discord.CrashLogWebhook) == 0 {
			c.WebAdmin.Discord = nil
		}
		if c.WebAdmin.Launcher != nil {
			if c.WebAdmin.Launcher.MinZoneProcesses == 0 {
				c.WebAdmin.Launcher.MinZoneProcesses = 10
			}
		}
	}

	file, err := json.MarshalIndent(c, "", " ")
	if err != nil {
		return err
	}

	path := filepath.Join(e.pathmgmt.GetEQEmuServerPath(), "eqemu_config.json")
	err = os.WriteFile(path, file, 0755)
	if err != nil {
		return err
	}

	return nil
}
