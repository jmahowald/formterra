// Copyright © 2016 NAME HERE <EMAIL ADDRESS>
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package cmd

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

//go:generate go-bindata -pkg cmd -o assets.go assets/

var cfgFile string
var TfDir, EnvName, OwnerName string

var OverWriteFiles bool

//mostly a dev flag that allows me to switch off reading
//templates from embedded assets or lcally
var EmbeddedTemplatesDir string

// RootCmd represents the base command when called without any subcommands
var RootCmd = &cobra.Command{
	Use:   "formterra",
	Short: "A collection of commands to use common terraform patterns",
	Long: `A longer description that spans multiple lines and likely contains
examples and usage of using your application. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	// Uncomment the following line if your bare application
	// has an action associated with it:
	//	Run: func(cmd *cobra.Command, args []string) { },
}

// Execute adds all child commands to the root command sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	if err := RootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(-1)
	}
}

var homeDir = os.Getenv("HOME")

func init() {

	// Here you will define your flags and configuration settings.
	// Cobra supports Persistent Flags, which, if defined here,
	// will be global for your application.

	RootCmd.PersistentFlags().StringVarP(&cfgFile, "config", "c",
		filepath.Join(homeDir, ".formterra", "formterra.yml"), "config file")

	RootCmd.PersistentFlags().StringVarP(&TfDir, "terraform-dir", "d",
		filepath.Join(homeDir, ".formterra", "terraform"), "directory where generated terraform will go")

	RootCmd.PersistentFlags().BoolP("overwrite", "w",
		false, "if there is something already there, do we overwrite when we generate")

	RootCmd.PersistentFlags().StringVarP(&EnvName, "env", "e", "", "what environment name the resource should be tagged with ")
	RootCmd.PersistentFlags().StringVarP(&OwnerName, "owner", "o", "", "what owner should the resources be tagged with")
	cobra.OnInitialize(initConfig)
}

// initConfig reads in config file and ENV variables if set.
func initConfig() {
	if cfgFile != "" { // enable ability to specify config file via flag
		viper.SetConfigFile(cfgFile)
	}
	viper.SetConfigType("yaml")
	// viper.SetConfigName("formterra.yml") // name of config file (without extension)
	// viper.AddConfigPath("$HOME/.formterra")  // adding home directory as first search path

	//TODO figure out what this means
	//viper.AutomaticEnv()          // read in environment variables that match

	viper.BindPFlag("owner", RootCmd.PersistentFlags().Lookup("owner"))
	viper.BindPFlag("terraform-dir", RootCmd.PersistentFlags().Lookup("terraform-dir"))
	viper.BindPFlag("env", RootCmd.PersistentFlags().Lookup("env"))
	viper.BindPFlag("overwrite", RootCmd.PersistentFlags().Lookup("overwrite"))

	// viper.BindPFlag("useViper", RootCmd.PersistentFlags().Lookup("viper"))
	// viper.SetDefault("author", "NAME HERE <EMAIL ADDRESS>")

	// If a config file is found, read it in.
	if err := viper.ReadInConfig(); err == nil {
		fmt.Println("Using config file:", viper.ConfigFileUsed())
	}
}
