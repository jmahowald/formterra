package tfproject

import (
	"flag"
	"os"
	"path/filepath"
	"testing"

	log "github.com/Sirupsen/logrus"
	"github.com/spf13/viper"
)

func TestMain(m *testing.M) {
	flag.Parse()

	testdir := "./target"
	log.SetLevel(log.DebugLevel)
	//Cleanup if around from old test
	os.RemoveAll(testdir)
	os.MkdirAll(testdir, 0755)

	viper.Set(TerraformDir, testdir)

	// if !testing.Short() {
	//     setupDatabase()
	// }
	result := m.Run()

	// Everything passed, cleanup after ourselves.  Otherwise leave around for inspection
	if result == 0 {
		// os.RemoveAll(testdir)
	}
	os.Exit(result)
}

func TestS3BucketCreate(t *testing.T) {

	req := S3BucketRequest{
		S3BucketID{"my.test", "testingbucket"},
		true,
		true,
		CorsConfig{},
	}
	layer, _ := req.Create()
	dir, _ := layer.dir()
	expectedFile := filepath.Join(dir, "s3.tf")
	_, err := os.Stat(expectedFile)
	if err != nil {
		t.Errorf("Didn't create expected file:%s", expectedFile)
	}
	_, err = layer.PlanCommand()
	if err != nil {
		t.Errorf("Didn't create a makeable file")
	}
}

// func TestS3BucketRequest_Create(t *testing.T) {
// 	type fields struct {
// 		S3BucketID  S3BucketID
// 		UnVersioned bool
// 	}
// 	tests := []struct {
// 		name   string
// 		fields fields
// 		want   TerraformLayer
// 		want1  bool
// 	}{
// 		{
// 			"projecreate",
// 			fields{S3BucketID{"my.test", "testingbucket"}, false},
// 			TerraformLayer{"testingbucket"},
// 			true,
// 		},
// 	}

// 	// 	req := S3BucketRequest{
// 	// 	S3BucketID{"testingbucket", "my.test"},
// 	// 	true,
// 	// }
// 	// layer, _ := req.Create()
// 	// expectedFile := filepath.Join(testdir, "test", "bucket_my.test", "s3.tf")
// 	// fileExists(expectedFile, c)
// 	// fileExists(filepath.Join(testdir, "test", "bucket_my.test", "Makefile"), c)
// 	// _, err := layer.PlanCommand()
// 	// check(c, err, "Couldn't get make command")

// 	for _, tt := range tests {
// 		t.Run(tt.name, func(t *testing.T) {
// 			s := S3BucketRequest{
// 				S3BucketID:  tt.fields.S3BucketID,
// 				UnVersioned: tt.fields.UnVersioned,
// 			}
// 			got, got1 := s.Create()

// 			if !reflect.DeepEqual(got, tt.want) {
// 				t.Errorf("S3BucketRequest.Create() got = %v, want %v", got, tt.want)
// 			}
// 			if got1 != tt.want1 {
// 				t.Errorf("S3BucketRequest.Create() got1 = %v, want %v", got1, tt.want1)
// 			}
// 		})
// 	}
// }
