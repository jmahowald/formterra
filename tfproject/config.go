package tfproject

import "github.com/spf13/viper"

// Config basic interface for reading configuration
type Config interface {
	IsSet(key string) bool
	GetString(key string) string
}

type viperconfig struct{}

func (v viperconfig) IsSet(key string) bool {
	return viper.IsSet(key) && viper.GetBool(key)
}

func (v viperconfig) GetString(key string) string {
	return viper.GetString(key)
}

var viperConfig Config = viperconfig{}
