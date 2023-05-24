package impls

import (
	"bytes"
	"fmt"

	"github.com/spf13/viper"
)

func LoadYamlConfig(fp, name string) (vp *viper.Viper, err error) {
	vp = viper.New()
	vp.SetConfigName(name)
	vp.SetConfigFile(fp)

	if err = vp.ReadInConfig(); err != nil {
		return nil, fmt.Errorf("ReadInConfig: %w", err)
	}

	return vp, nil
}

func LoadYamlBytes(bts []byte) (vp *viper.Viper, err error) {
	vp = viper.New()
	vp.SetConfigType("yaml")

	if err = vp.ReadConfig(bytes.NewReader(bts)); err != nil {
		return nil, err
	}

	return vp, nil
}

func UnmarshalYamlBytes(bts []byte, obj any) (err error) {
	var vp *viper.Viper

	if vp, err = LoadYamlBytes(bts); err != nil {
		return err
	}

	return vp.Unmarshal(obj)
}

func UnmarshalYamlObjects(fp string, objects map[string]any) (err error) {
	var vp *viper.Viper

	if vp, err = LoadYamlConfig(fp, "..."); err != nil {
		return err
	}

	for k, v := range objects {
		if err = vp.UnmarshalKey(k, v); err != nil {
			return err
		}
	}

	return nil
}
