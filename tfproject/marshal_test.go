package tfproject

import (
	"bytes"
	"flag"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

// In the particular tests below, we're only going to concern ourselves
// with the rendering of the acutal terraform in project.tf
// the contents of the makefile, manifest, tfvars example are not interesting
type projectBufferData struct {
	buff *bytes.Buffer
}

func (b projectBufferData) getWriter(name string) (io.Writer, error) {
	if strings.Index(name, "project.tf") != -1 {
		return b.buff, nil
	} else {
		//Ignoring output for file and giving new buffer
		return &bytes.Buffer{}, nil
	}
}

var update = flag.Bool("update", false, "update golden files")

// Most of the logic of the terraform skeletion generation produces
// a terraform file that reflects the data.  This tests various cases
// against a golden file.  If you run the test with a -update flag
// you can simply capture the output for future runs
func TestTerraformProjectSkeleton_UnmarshalYAMLAndGenerate(t *testing.T) {
	tests := []struct {
		name    string
		input   string
		wantErr bool
	}{{"simpleWithDefaultValue",
		`
Name: simpletest
Modules:
- Name: simple
  local_path: test-fixtures/simple
  vars:
  - var_name: location
    default: world
`,
		false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ts := &TerraformProjectSkeleton{}
			if err := ts.UnmarshalYAML([]byte(tt.input)); (err != nil) != tt.wantErr {
				t.Errorf("TerraformProjectSkeleton.UnmarshalYAML() error = %v, wantErr %v", err, tt.wantErr)
			}
			tplHandler := AssetTemplates{"project"}
			var b bytes.Buffer
			out := projectBufferData{&b}
			processTemplates(out, tplHandler, ts)
			golden := filepath.Join("test-fixtures", tt.name, "golden.tf")
			if *update {
				os.MkdirAll(filepath.Join("test-fixtures", tt.name), 0755)
				ioutil.WriteFile(golden, out.buff.Bytes(), 0644)
				t.Errorf("Updated golden master for: with %s -- %s", golden, out.buff.String())
			}
			expected, _ := ioutil.ReadFile(golden)
			if out.buff.String() != string(expected) {
				t.Errorf("TerraformProjectSkeleton.MarshalYAML() = %s, want %s", out.buff.String(), expected)
			}

		})
	}
}
