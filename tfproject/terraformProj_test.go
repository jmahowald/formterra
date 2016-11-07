package tfproject

import (
	"bytes"
	"os"
	"path/filepath"
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
	testdir = "./target"
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

// func check(trutc *C)
func (s *MySuite) TestBucket(c *C) {
	c.Skip("Temporarily off to speed up testing")
	req := S3BucketRequest{
		S3BucketID{"testingbucket", "my.test"},
		true,
	}
	layer, exists := req.Create()
	if exists {
		c.Error("Bucket requeset already existed")
	}
	expectedFile := filepath.Join(testdir, "test", "bucket_testingbucket", "s3.tf")
	fileExists(expectedFile, c)
	fileExists(filepath.Join(testdir, "test", "bucket_testingbucket", "Makefile"), c)
	plan, err := layer.PlanCommand()
	check(c, err, "Couldn't get make command")
	plan.Stdout = os.Stdout
	plan.Stderr = os.Stderr
	err = plan.Run()
	check(c, err, "Problems running plan")
}

func (s *MySuite) TestRetreival(c *C) {

	// jsonexample := "examples"

	module := ExternalModule{URI: "./test-fixtures/simpleterraform"}
	projectDef, err := module.Fetch()
	check(c, err, "couldn't get local module")
	log.Debug("Project def:", projectDef)
	c.Assert(projectDef.Name, Equals, "simpleterraform")
	c.Assert(projectDef.RequiredVars, HasLen, 1)
	c.Assert(projectDef.RequiredVars, DeepEquals, []string{"location"})
}

func (s *MySuite) TestModuleClient(c *C) {
	generateModule("./test-fixtures/dcos", "mesos")
	expectedFile := filepath.Join(testdir, "test", "mesos", "module_client.tf")
	fileExists(expectedFile, c)

}

// func (s *MySuite) TestHelloWorld(c *C) {
// 	c.Assert(42, Equals, "42")
// 	c.Assert(io.ErrClosedPipe, ErrorMatches, "io: .*on closed pipe")
// 	c.Check(42, Equals, 42)
//
// }

func setConfig(location string) Config {
	viper.SetConfigFile(location)
	return viperConfig
}
