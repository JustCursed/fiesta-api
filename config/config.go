package config

import (
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"github.com/gofiber/fiber/v2/log"
	"github.com/pelletier/go-toml"
	"os"
)

const AuthKey string = "authorized"
const NeedCheckKey string = "needCheck"
const DiscordTokenKey string = "discordToken"
const AccessServersKey string = "accessServers"

var Config struct {
	General struct {
		DevMode bool              `toml:"dev_mode"`
		Address string            `toml:"address"`
		Secret  *ecdsa.PrivateKey `toml:"secret"`
	}
	Database struct {
		Address  string `toml:"address"`
		Username string `toml:"username"`
		Password string `toml:"password"`
	}
	Discord struct {
		Secret      string            `toml:"token"`
		ID          string            `toml:"id"`
		AuthURI     string            `toml:"auth_uri"`
		AccessRoles map[string]string `toml:"access_roles"`
	}
	Log struct {
		AvailableLogTypes []string            `toml:"available_log_types"`
		OptionalLogs      map[string][]string `toml:"optional_logs"`
		Servers           []string            `toml:"servers"`
	}
}

func init() {
	tree, err := toml.LoadFile("config.toml")

	if err == nil { // reading existing config
		err = tree.Unmarshal(&Config)
		Config.General.Secret.Curve = elliptic.P256()

		if err != nil {
			log.Fatalf("failed to unmarshal config: %v", err)
		}
	} else { // create config is not exists
		privateKey, _ := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)

		// default config
		Config.General.Secret = privateKey
		Config.Log.AvailableLogTypes = []string{"chat", "items", "deaths"}
		Config.Log.OptionalLogs = map[string][]string{"htc": {"movable"}}

		data, _ := toml.Marshal(&Config)

		err = os.WriteFile("config.toml", data, 0644)
		if err != nil {
			log.Fatalf("failed to write config: %v", err)
		}
	}
}
