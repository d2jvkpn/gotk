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
	Server        string        `mapstructure:"server"`
	Timeout       time.Duration `mapstructure:"timeout"`
	TlsSkipVerify bool          `mapstructure:"tls_skip_verify"`
	UserAgent     string        `mapstructure:"user_agent"`

	// bot
	ApiKey string `mapstructure:"api_key"`

	// user
	Email    string `mapstructure:"email"`
	Password string `mapstructure:"password"`

	*http.Client
}

func NewClient(fp string, keys ...string) (client *Client, err error) {
	var vp *viper.Viper

	vp = viper.New()
	// vp.SetConfigName(name)
	vp.SetConfigType("yaml")

	vp.SetConfigFile(fp)
	if err = vp.ReadInConfig(); err != nil {
		return nil, fmt.Errorf("read config %q: %w", fp, err)
	}
	vp.SetDefault("timeout", "5s")
	vp.SetDefault("user_agent", "Mozilla/5.0 Gecko/20100101 Firefox/130.0")

	return ClientFromViper(vp, keys...)
}

func ClientFromViper(vp *viper.Viper, keys ...string) (client *Client, err error) {
	client = new(Client)

	if len(keys) > 0 {
		err = vp.UnmarshalKey(keys[0], client)
	} else {
		err = vp.Unmarshal(client)
	}

	if err != nil {
		return nil, fmt.Errorf("unmarshal config: %w", err)
	}

	if client.Server = strings.TrimRight(client.Server, "/"); client.Server == "" {
		return nil, errors.New("empty server")
	}
	if client.Timeout < 0 {
		return nil, fmt.Errorf("invalid timeout: %v", client.Timeout)
	}

	client.Client = new(http.Client)
	client.Client.Timeout = client.Timeout

	if client.TlsSkipVerify {
		client.Client.Transport = &http.Transport{
			TLSClientConfig: &tls.Config{MinVersion: tls.VersionTLS12, InsecureSkipVerify: true},
		}
	}

	return client, nil
}
