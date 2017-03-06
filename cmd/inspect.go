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
	"github.com/ghodss/yaml"
	tf "github.com/backpack/formterra/tfproject"
	"github.com/spf13/cobra"
)

var uris []string
var name string
var generate bool

// inspectCmd represents the inspect command
var inspectCmd = &cobra.Command{
	Use:   "inspect",
	Short: "downloads terraform modules and examines them, producing a yaml that you can subsequently use to generate a skeleton",
	Long:  `Use this to generate yaml.  Then edit and create clients for the nodes.`,
	Run: func(cmd *cobra.Command, args []string) {
		if len(uris) < 1 {
			log.Error("You must supply a uri of a terraform module")
			cmd.Usage()
			return
		}

		if name == "" {
			log.Error("You must suply a name")
			cmd.Usage()
			return
		}

		// For each uri, fetch and get the module definition to add to our list
		modDefs := make([]tf.TerraformModuleDefinition, 0, len(uris))
		for _, uri := range uris {
			def, err := tf.ExternalModule{URI: uri}.Fetch()
			if err != nil {
				log.Errorf("could not fetch module %s:%s", uri, err)
				return
			}
			modDefs = append(modDefs, def)
		}

		skel := tf.CreateSkeleton(modDefs, name)
		data, err := yaml.Marshal(skel)
		if err != nil {
			log.Errorf("Error marshalling modules %s:%s", skel, err)
			return
		}
		write(data)

		if generate {
			skel.GenerateSkeleton()
		}
	},
}

func init() {
	moduleCmd.AddCommand(inspectCmd)
	inspectCmd.Flags().StringVar(&name, "name", "", "name for the resulting project skeleton (required)")
	inspectCmd.Flags().StringSliceVarP(&uris, "uri", "u", []string{}, "uri to be inspected (can specify multiple)")
	inspectCmd.Flags().BoolVarP(&generate, "gen", "g", false, "generate a skeleton as well")

	moduleCmd.AddCommand(genTerraformCommand)

}

// genTerraformCommand represents the command to generate a skeleton based off a project defintion
var genTerraformCommand = &cobra.Command{
	Use: "gen",

	Short: "Generates terraform to call modules from a skeleton",
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
