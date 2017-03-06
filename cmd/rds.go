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
	tf "github.com/backpack/formterra/tfproject"
	"github.com/spf13/cobra"
)

var databaseName string
var rdsCommand = cobra.Command{
	Use:   "rds <command>",
	Short: "Creates a terrafomr module that you can use to create a db",
	Long: `Generates terraform necessary to create an rds request
	this allows us to more finely tune
	`,
	PreRun: func(cmd *cobra.Command, args []string) {
		if databaseName == "" {
			fail("You must provide a database name", cmd)
			return
		}
	},

	Run: func(cmd *cobra.Command, args []string) {
		dbRequest := tf.RDSRequest{databaseName}
		dbRequest.Create()
	},
}

func init() {
	RootCmd.AddCommand(&rdsCommand)
	rdsCommand.Flags().StringVar(&databaseName, "database", "", "what name do you want for your rds instance")
}
