package tfproject

import (
	"bytes"
	"fmt"
	"os"
	"testing"

	log "github.com/Sirupsen/logrus"

	"github.com/spf13/viper"
	// I use this instead of base testing Suite
	// to bring back warm fuzzies of junit

	. "gopkg.in/check.v1"
)

// Hook up gocheck into the "go test" runner.
func Test(t *testing.T) { TestingT(t) }

type MySuite struct{}

var _ = Suite(&MySuite{})
var failure = false

var testdir string

func setTestConfig(config []byte) {
	viper.SetConfigType("yaml") // or viper.SetConfigType("YAML")
	viper.ReadConfig(bytes.NewBuffer(config))

}

func fileExists(path string, c *C) {
	_, err := os.Stat(path)
	c.Assert(err, IsNil, Commentf("file exists at %s", path))
}

func (s *MySuite) SetUpSuite(c *C) {
	testdir = "./testresults"
	log.SetLevel(log.DebugLevel)
	//Cleanup if around from old test
	os.RemoveAll(testdir)
	os.MkdirAll(testdir, 0755)
	viper.Set(TerraformDir, testdir)
	viper.Set("env", "test")

}

func (s *MySuite) TearDownTest(c *C) {

	//If any tests fail, we wanted to mark that so we don't clean up and can examine
	failure = true
}

func (s *MySuite) TearDownSuite(c *C) {
	if !failure {
		os.RemoveAll(testdir)
	} else {
		c.Log("Failures, please examine:", testdir)
	}
}

func check(c *C, err error, msg ...string) {
	if err != nil {
		c.Error(msg, err)
	}
}

var projectStruct = `
name: TestSource
modules:
- name: mod1
  outputs:
  - out1
  - out2
  orig_uri: http://testlocation
  module_vars:
  - module_name: mod2
    mappings:
    - source_var_name: mod2_out
      var_name: foo
    - var_name: bar
  -  module_name: mod3
     mappings:
     - var_name: value3
       source_var_name: mod3var
  remote_source_vars:
  - source_name: vpc_layer
    mappings:
    - source_var_name: remote1_out
      var_name: foo
    - var_name: bar2
      source_var_name: mod3var
    config:
      bucket: mytestingbucket
      key: mytestingbucketkey
      region: us-east-1
  vars:
  - source_var_name: var1_in
    var_name: var1_out
  - var_name: bar3
    default: mydefault
`

func (s *MySuite) TestModuleMarshalling(c *C) {
	var proj TerraformProjectSkeleton

	moduleCall := ModuleCall{

		TerraformModuleDefinition: TerraformModuleDefinition{
			Name:    "mod1",
			URI:     "http://testlocation",
			Outputs: []string{"out1", "out2"},
		},
		ModuleVariables: []FromModuleMappings{
			FromModuleMappings{
				"mod2",
				[]BasicVariableMapping{
					BasicVariableMapping{"foo", "mod2_out", "", ""},
					BasicVariableMapping{VarName: "bar"},
				},
			},
			FromModuleMappings{
				"mod3",
				[]BasicVariableMapping{
					BasicVariableMapping{"value3", "mod3var", "", ""},
				},
			},
		},
		RemoteVariables: []FromRemoteMappings{
			FromRemoteMappings{
				RemoteSourceName: "vpc_layer",
				Mappings: []BasicVariableMapping{
					BasicVariableMapping{"foo", "remote1_out", "", ""},
					BasicVariableMapping{"bar2", "mod3var", "", ""},
				},
				Config: map[string]string{
					"bucket": "mytestingbucket",
					"key":    "mytestingbucketkey",
					"region": "us-east-1",
				},
			},
		},
		Variables: []BasicVariableMapping{
			BasicVariableMapping{"var1_out", "var1_in", "", ""},
			BasicVariableMapping{VarName: "bar3", DefaultValue: "mydefault"},
		},
	}

	expectedProj := TerraformProjectSkeleton{
		TerraformLayer{Name: "TestSource"},
		[]ModuleCall{moduleCall},
	}

	yamlOut, _ := expectedProj.MarshalYAML()

	fmt.Print(string(yamlOut))

	err := proj.UnmarshalYAML([]byte(projectStruct))
	c.Assert(err, IsNil)
	c.Assert(proj, NotNil)
	c.Assert(proj, DeepEquals, expectedProj)

	vars := proj.GetAllVars()
	for _, variable := range vars {
		c.Assert(variable.VarName, Not(Equals), "")
	}
}

func (s *MySuite) TestProjectGeneration(c *C) {
	var proj TerraformProjectSkeleton
	err := proj.UnmarshalYAML([]byte(projectStruct))
	err = proj.GenerateSkeleton()
	log.Info("result is ", err)
	c.Assert(err, IsNil)
}

func setConfig(location string) Config {
	viper.SetConfigFile(location)
	return viperConfig
}
