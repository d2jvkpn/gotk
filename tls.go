package gotk

import (
	"crypto/tls"

	"github.com/spf13/viper"
)

type TlsConfig struct {
	Enable bool   `mapstructure:"enable"`
	Cert   string `mapstructure:"cert"`
	Key    string `mapstructure:"key"`
}

func NewTlsConfig(vp *viper.Viper, field string) (cert *tls.Certificate, err error) {
	var (
		tlsConfig TlsConfig
		_cert     tls.Certificate
	)

	if err = vp.UnmarshalKey(field, &tlsConfig); err != nil {
		return nil, err
	}

	if !tlsConfig.Enable {
		return nil, nil
	}

	if _cert, err = tls.LoadX509KeyPair(tlsConfig.Cert, tlsConfig.Key); err != nil {
		return nil, err
	}

	return &_cert, nil
}
