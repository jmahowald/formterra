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
	"os"

	"io/ioutil"

	log "github.com/Sirupsen/logrus"
	tf "github.com/jmahowald/formterra/tfproject"
	"github.com/spf13/cobra"
)

// TfClientRequest
type TfClientRequest struct {
	Uri               string
	ExistingTfVarPath string
	CreateTfVars      bool
}

var externalModuleDef = false
var projectConfig string
var projectName string
var clientRequest TfClientRequest

func genClient(clientRequest TfClientRequest) {

	module := tf.ExternalModule{URI: clientRequest.Uri}
	projectDef, err := module.Fetch()
	if err != nil {
		log.Fatalf("error fetching %s %v", clientRequest.Uri, err)
	}
	log.Debugf("Project definition vars %s", projectDef.RequiredVars)
}

// clientCmd represents the client command
var clientCmd = &cobra.Command{
	Use: "client",

	Short: "Generates terraform to call modules",
	Long: `Generates terraform necessary to interact with
	existing terraform modules
	`,

	PreRun: func(cmd *cobra.Command, args []string) {
		if projectConfig == "" {
			fail("You must provide a configuration file for the terraform skeleto to generate", cmd)
			return
		}
		// if clientRequest.Uri == "" {
		// 	fail("You must provide the uri of the terraform project", cmd)
		// 	return
		// }
	},

	Run: func(cmd *cobra.Command, args []string) {

		configFileInfo, err := os.Stat(projectConfig)
		if err != nil {
			log.Warnf("Could not find file at %s", projectConfig)
			return
		}
		if projectName == "" {
			projectName = configFileInfo.Name()
		}
		bytes, err := ioutil.ReadFile(projectConfig)
		if err != nil {
			log.Fatal("Could not read contents of %s", projectConfig)
		}
		var skeleton = tf.TerraformProjectSkeleton{}
		err = skeleton.UnmarshalYAML(bytes)
		if err != nil {
			log.Fatal("error parsing %s:%s", projectConfig, err)
		}

		skeleton.Name = projectName
		skeleton.GenerateSkeleton()
	},
}

func init() {
	moduleCmd.AddCommand(clientCmd)
	clientCmd.Flags().StringVar(&projectConfig, "config", "", "Points to a terraform skeleton config")
	clientCmd.Flags().StringVar(&projectName, "name", "n", "project name (defaults to last name of uri)")
}
