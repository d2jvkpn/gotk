package vocechat

import (
	"crypto/tls"
	"errors"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/spf13/viper"
)

type Client struct {
	Server   string `mapstructure:"server"`
	Username string `mapstructure:"username"`

	// Account
	Email    string `mapstructure:"email"`
	Password string `mapstructure:"password"`

	// Bot
	ApiKey string `mapstructure:"api_key"`

	SendToId      uint32 `mapstructure:"send_to_id"`
	SendToName    string `mapstructure:"send_to_name"`
	SendToKind    string `mapstructure:"send_to_kind"` // user or group
	Timeout       int64  `mapstructure:"timeout"`      // seconds
	TlsSkipVerify bool   `mapstructure:"tls_skip_verify"`

	cli *http.Client
}

func NewClient(fp string, key string) (client *Client, err error) {
	conf := viper.New()
	// conf.SetConfigName(name)
	conf.SetConfigType("yaml")
	conf.SetDefault(fmt.Sprintf("%s.timeout", key), 5)

	conf.SetConfigFile(fp)
	if err = conf.ReadInConfig(); err != nil {
		return nil, fmt.Errorf("read config %q: %w", fp, err)
	}

	client = new(Client)
	if err = conf.UnmarshalKey(key, client); err != nil {
		return nil, fmt.Errorf("unmarshal config: %w", err)
	}

	if client.Server = strings.TrimRight(client.Server, "/"); client.Server == "" {
		return nil, errors.New("empty server")
	}

	if client.SendToId == 0 {
		return nil, errors.New("invalid send_to_id")
	}
	if client.SendToKind != "user" && client.SendToKind != "group" {
		return nil, fmt.Errorf("invalid send_to_kind: %s", client.SendToKind)
	}

	if client.Timeout < 0 {
		return nil, fmt.Errorf("invalid timeout: %v", client.Timeout)
	}

	client.cli = new(http.Client)
	client.cli.Timeout = time.Duration(client.Timeout) * time.Second

	if client.TlsSkipVerify {
		client.cli.Transport = &http.Transport{
			TLSClientConfig: &tls.Config{MinVersion: tls.VersionTLS12, InsecureSkipVerify: true},
		}
	}

	return client, nil
}
