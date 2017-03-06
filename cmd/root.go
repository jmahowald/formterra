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

	log "github.com/Sirupsen/logrus"

	"github.com/backpack/formterra/core"
	tf "github.com/backpack/formterra/tfproject"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// "github.com/prometheus/common/log"

var cfgFile string
var debugLogging bool

// RootCmd represents the base command when called without any subcommands
var RootCmd = &cobra.Command{
	Use:   "formterra",
	Short: "A collection of commands to use common terraform patterns",

	// Set logging for all commands
	PersistentPreRun: func(cmd *cobra.Command, args []string) {
		if debugLogging {
			log.SetLevel(log.DebugLevel)
		}
	},

	// Uncomment the following line if your bare application
	// has an action associated with it:
	//	Run: func(cmd *cobra.Command, args []string) { },
}

var versionCommand = &cobra.Command{
	Use: "version",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Printf("Version: %s BuildTime:%s", core.Version, core.BuildTime)
	},
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
	RootCmd.PersistentFlags().BoolVar(&debugLogging, "debug", false, "turn on debug")

	RootCmd.PersistentFlags().StringP(tf.TerraformDir, "d", ".", "directory where generated terraform will go")

	//TODO should use a contant
	RootCmd.PersistentFlags().BoolP(tf.Overwrite, "w",
		false, "if there is something already there, do we overwrite when we generate")

	RootCmd.PersistentFlags().StringP(tf.Env, "e", "", "what environment name the resource should be tagged with ")
	cobra.OnInitialize(initConfig)
	RootCmd.AddCommand(versionCommand)

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

	viper.BindPFlag(tf.Owner, RootCmd.PersistentFlags().Lookup(tf.Owner))
	viper.BindPFlag(tf.TerraformDir, RootCmd.PersistentFlags().Lookup(tf.TerraformDir))
	viper.BindPFlag(tf.Env, RootCmd.PersistentFlags().Lookup(tf.Env))
	viper.BindPFlag(tf.Overwrite, RootCmd.PersistentFlags().Lookup(tf.Overwrite))

	// viper.BindPFlag("useViper", RootCmd.PersistentFlags().Lookup("viper"))
	// viper.SetDefault("author", "NAME HERE <EMAIL ADDRESS>")

	// If a config file is found, read it in.
	if err := viper.ReadInConfig(); err == nil {
		log.Info("Using config file:", viper.ConfigFileUsed())
	}
}

func fail(msg string, cmd *cobra.Command) {
	fmt.Println(msg)
	cmd.Usage()
	os.Exit(-1)
}
