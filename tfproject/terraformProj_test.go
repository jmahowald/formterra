package tfproject

import (
	"bytes"
	"os"
	"path/filepath"
	"testing"

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
		// os.RemoveAll(testdir)
	} else {
		c.Log("Failures, please examine:", testdir)
	}
}

func (s *MySuite) TestBucket(c *C) {
	req := S3BucketRequest{
		BucketName:  "testbucket",
		UnVersioned: true,
		Fqdn:        "my.test",
	}

	tf := TerraformS3(req)
	// Hmm, this smells like I should be capturing the request in the object
	tf.Write(req)

	expectedFile := filepath.Join(testdir, "test", "bucket", "s3.tf")
	fileExists(expectedFile, c)
	fileExists(filepath.Join(testdir, "test", "bucket", "Makefile"), c)

	// tf := cmd.

	//
	// req := S3BucketRequest(Org)
	//
	// tf :-= TerraformS3(S3BucketRequest{
	//   Buck
	// })
	// func Render(bucketRequest S3BucketRequest) {
	//
	// 	layer := tf.TerraformLayer{Name: "bucket"}
	// 	s3Proj := tf.PredefinedTerraformProjects{
	// 		TerraformLayer: layer,
	// 		Templates:      []string{"s3.tf"},
	// 	}
	// 	s3Proj.Write(bucketRequest)
	// }

}

// any approach to require this configuration into your program.
var yamlExample = []byte(`

Hacker: true
name: steve
hobbies:
- skateboarding
- go
clothing:
  jacket: leather
  trousers: denim
age: 35
eyes : brown
beard: true
`)

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
