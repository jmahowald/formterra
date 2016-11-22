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
	"strings"

	log "github.com/Sirupsen/logrus"
	"github.com/jmahowald/formterra/tfproject"
	"github.com/spf13/cobra"
)

var name string
var clientFile string

// inspectCmd represents the inspect command
var inspectCmd = &cobra.Command{
	Use:   "inspect <moduleuri>",
	Short: "downloads a terraform module and tells you what variables are required and optional",
	Long:  `Use this to generate yaml.  Then edit and create clients for the nodes.`,
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) != 1 {
			log.Error("You must supply a uri of a terraform module")
			cmd.Usage()
			return
		}
		uri := args[0]
		// url, err := url.Parse(uri)
		// if err != nil {
		// 	log.Errorf("Invalid uri %s :%s", uri, err)
		// 	return
		// }
		if name == "" {
			urlParts := strings.Split(uri, "/")
			name = urlParts[len(urlParts)-1]
		}
		def, err := tfproject.ExternalModule{Name: name, URI: uri}.Fetch()
		if err != nil {
			log.Error("could not fetch module:", err)
			return
		}

		bytes, err := def.MarshalYAML()
		if err != nil {
			log.Errorf("Error marshalling project %v:%s", def, err)
		}

		//TODO this should be able to be a file and stdout should just be a flag
		fmt.Print(string(bytes))
	},
}

func init() {
	moduleCmd.AddCommand(inspectCmd)
	inspectCmd.Flags().StringVar(&name, "name", "", "name of the module")
	inspectCmd.Flags().StringVar(&name, "name", "", "name of the module")

}
