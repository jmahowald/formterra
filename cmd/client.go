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

// clientCmd represents the client command
var clientCmd = &cobra.Command{
	Use: "client",

	Short: "Generates terraform to call modules",
	Long: `Generates terraform necessary to interact with
	existing terraform modules
	`,

	Run: func(cmd *cobra.Command, args []string) {
		input := read()
		var skeleton = tf.TerraformProjectSkeleton{}
		err := skeleton.UnmarshalYAML(input)
		if err != nil {
			log.Fatal("error parsing %s:%s", projectConfig, err)
		}
		skeleton.GenerateSkeleton()
	},
}

func init() {
	moduleCmd.AddCommand(clientCmd)
	clientCmd.Flags().StringVarP(&projectName, "name", "n", "", "project name")
}
