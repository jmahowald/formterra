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

// S3BucketRequest Used to make  a new bucket in s3
// because of the global namespace, we try to
// enforce having a base fqdn and a name on top of that.
type S3BucketRequest struct {
	Fqdn        string
	BucketName  string
	UnVersioned bool
	CreateUser  bool
}

var bucketRequest tf.S3BucketRequest

func render(bucketRequest tf.S3BucketRequest) {
	bucketRequest.Create()
}

var s3Cmd = cobra.Command{
	Use:   "s3 <command>",
	Short: "Manipulates s3 buckets using terraform",
	Long: `Generates terraform necessary to create an s3 bucket
	this allows us to more finely tune
	`,
	PreRun: func(cmd *cobra.Command, args []string) {
		if bucketRequest.BucketName == "" {
			fail("You must provide a bucket name", cmd)
			return
		}
		if bucketRequest.Fqdn == "" {
			log.Error("You must provide fqdn for the bucket")
			cmd.Usage()
			return
		}
	},

	Run: func(cmd *cobra.Command, args []string) {
		render(bucketRequest)
	},
}

func init() {
	RootCmd.AddCommand(&s3Cmd)

	// Here you will define your flags and configuration settings.
	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// s3Cmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	s3Cmd.Flags().StringVarP(&bucketRequest.BucketName, "bucket", "b", "", "what's the base name for your bucket")
	s3Cmd.Flags().StringVarP(&bucketRequest.Fqdn, "fqdn", "f", "", "is prepended onto your bucket name")
	s3Cmd.Flags().BoolVarP(&bucketRequest.UnVersioned, "unversioned", "u", false, "do you want the bucket unversioned")
	s3Cmd.Flags().BoolVar(&bucketRequest.CreateUser, "createuser", false, "do you want to create a user")
}
