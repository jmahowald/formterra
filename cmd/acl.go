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
	tf "github.com/jmahowald/formterra/tfproject"
	"github.com/spf13/cobra"
)

var aclRequest tf.NetworkACLAllowRequest

var aclCmd = cobra.Command{
	Use:   "acl -cidr <x> -port y ",
	Short: "Turns off all but necessary ports",
	Long: `Generates terraform necessary to turn off 
    ports except the specified set to the specified set of cidrs
    For multiple cidrs just specify multipe -cidr.  Similar for ports
	`,
	PreRun: func(cmd *cobra.Command, args []string) {
		if len(aclRequest.Cidrs) == 0 {
			fail("You must supply at least 1 cidr", cmd)
			return
		}
		if len(aclRequest.Ports) == 0 {
			fail("You must supply at least 1 port", cmd)
			return
		}

	},

	Run: func(cmd *cobra.Command, args []string) {
		aclRequest.Create()
	},
}

func init() {
	RootCmd.AddCommand(&aclCmd)

	// Here you will define your flags and configuration settings.
	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// s3Cmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	aclCmd.Flags().StringSliceVarP(&aclRequest.Cidrs, "cidr", "i", []string{}, "cidr.  Can have multiple cidrs")
	aclCmd.Flags().StringSliceVarP(&aclRequest.Ports, "port", "p", []string{}, "port.  Can have multiple flags")
	aclCmd.Flags().BoolVarP(&aclRequest.TCP, "tcp", "t", true, "for tcp")
	aclCmd.Flags().BoolVarP(&aclRequest.UDP, "udp", "u", false, "for udp")

}
