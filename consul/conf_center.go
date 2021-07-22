package consul

import (
	"time"

	"github.com/spf13/viper"
	_ "github.com/spf13/viper/remote"
)

type ConfClient interface {
	ParseConfig(confName string) (*viper.Viper, error)
}

type confClient struct {
	Host     string
	ConfType string
}

func NewConfClient(host string) ConfClient {
	return &confClient{Host: host, ConfType: "json"}
}

func (c *confClient) ParseConfig(confName string) (*viper.Viper, error) {
	var runtime_viper = viper.New()

	runtime_viper.AddRemoteProvider("consul", c.Host, confName)
	runtime_viper.SetConfigType(c.ConfType)
	err := runtime_viper.ReadRemoteConfig()
	if err != nil {
		return nil, err
	}

	go func() {
		for {
			time.Sleep(time.Second * 10)
			err := runtime_viper.WatchRemoteConfig()
			if err != nil {
				continue				
			}
		}
	}()

	return runtime_viper, nil
}
