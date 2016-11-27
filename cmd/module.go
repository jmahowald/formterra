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
	"io"
	"log"
	"os"

	"io/ioutil"

	"github.com/spf13/cobra"
)

// moduleCmd represents the module command
var moduleCmd = &cobra.Command{
	Use:   "module",
	Short: "Interacts with terraform modules, bridging variables",
	Long:  ``,
	// Run: func(cmd *cobra.Command, args []string) {
	// 	// TODO: Work your own magic here
	// 	fmt.Println("module called")
	// },
}

var output string
var input string

func getReader() (io.Reader, error) {
	var reader io.Reader

	if input == "" || input == "-" {
		reader = os.Stdin
	} else {
		reader, err := os.Open(input)
		if err != nil {
			log.Fatalf("Could not open %s for reading", input)
		}
		return reader, err
	}
	return reader, nil
}

func getWriter() (io.Writer, error) {
	if output == "" || output == "-" {
		return os.Stdout, nil
	} else {
		writer, err := os.Create(output)
		if err != nil {
			log.Fatalf("Could not open %s for writing", output)
			return writer, err
		}
		return writer, nil
	}
}

func read() []byte {
	var reader io.Reader
	reader, err := getReader()
	if err != nil {
		log.Fatalf("Could not open %s for reading", input)
	}
	data, err := ioutil.ReadAll(reader)
	if err != nil {
		log.Fatalf("Could not read from %s", reader)
	}
	return data
}

func write(data []byte) {
	var writer io.Writer
	writer, err := getWriter()
	if err != nil {
		log.Fatalf("Could not open %s for writing", output)
	}
	//defer writer.Close()
	_, err = fmt.Fprint(writer, string(data))
	if err != nil {
		log.Fatalf("Problem writing to %s", writer)
	}
}

func init() {
	moduleCmd.PersistentFlags().StringVarP(&output, "output", "o", "-", "Output use - to output to stdout (default)")
	moduleCmd.PersistentFlags().StringVarP(&input, "input", "i", "-", "Input use - to use stdin(default)")

	RootCmd.AddCommand(moduleCmd)

}
