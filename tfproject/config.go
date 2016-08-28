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

// By default use viper, but it's overritable
var viperConfig Config = viperconfig{}

// Overwrite flag setting in config on whether we default overwrite files if they
// exist
const Overwrite = "overwrite"

// TerraformDir Property Name to overwrite where to store the generated terraform
const TerraformDir = "terraform-dir"

// Env All terraform projects belong to an environment
const Env = "env"

// Owner used to tag instances
const Owner = "owner"
