package ginx

import (
	"crypto/tls"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
)

type HttpConfig struct {
	Path string `mapstructure:"path"`
	Cors string `mapstructure:"cors"`
	Tls  bool   `mapstructure:"tls"`
	Cert string `mapstructure:"cert"`
	Key  string `mapstructure:"key"`
}

func (self *HttpConfig) FromFile(fp string) (err error) {
	vp := new(viper.Viper)

	vp.SetConfigType("yaml")
	vp.SetConfigName("")
	vp.SetConfigFile(fp)

	if err = vp.Unmarshal(self); err != nil {
		return err
	}

	return nil
}

func (self *HttpConfig) SetServer(server *http.Server) (err error) {
	var cert tls.Certificate

	if self.Tls {
		if cert, err = tls.LoadX509KeyPair(self.Cert, self.Key); err != nil {
			return err
		}

		server.TLSConfig = &tls.Config{
			Certificates: []tls.Certificate{cert},
		}
	}

	return nil
}

func (self *HttpConfig) SetEngine(engine *gin.Engine) {
	if self.Cors != "" {
		engine.Use(Cors(self.Cors))
	}

	router := &engine.RouterGroup
	*router = *(router.Group(self.Path))
}
