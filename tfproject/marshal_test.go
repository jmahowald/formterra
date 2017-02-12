package tfproject

import (
	"bytes"
	"flag"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
	"reflect"
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
		return ioutil.Discard, nil
	}
}

type localFileTemplateIterator struct {
	subdir string
}

var projectAssetsTemplates = localFileTemplateIterator{"assets/project"}

func (l localFileTemplateIterator) getTemplates() (filenames []string, err error) {
	files, err := ioutil.ReadDir(l.subdir)
	if err != nil {
		return
	}
	for _, f := range files {
		filenames = append(filenames, f.Name())
	}
	return
}

func (l localFileTemplateIterator) loadTemplate(name string) ([]byte, error) {
	return ioutil.ReadFile(filepath.Join(l.subdir, name))
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
		{"testwithlists",
			`
Name:listtest
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
			var b bytes.Buffer
			out := projectBufferData{&b}
			processTemplates(out, projectAssetsTemplates, ts)
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

func TestSkelFromYaml(t *testing.T) {
	type args struct {
		data string
	}
	tests := []struct {
		name    string
		args    args
		wantT   TerraformProjectSkeleton
		wantErr bool
	}{{"outputs", args{
		`
name: TestModuleOut
modules:
- name: mod1
  outputs:
  - out1
  - out2
`},
		TerraformProjectSkeleton{TerraformLayer{Name: "TestModuleOut"},
			[]ModuleCall{
				ModuleCall{TerraformModuleDefinition: TerraformModuleDefinition{Name: "mod1", Outputs: []string{"out1", "out2"}}},
			}},
		false},
		{"default values", args{
			`
name: TestDefault
modules:
- name: mod
  vars:
  - var_name: vpc_id
  - var_name: env
    default: testing
  - var_name: tags
    defaults:
    - owner=Josh
    - keep
`},
			TerraformProjectSkeleton{TerraformLayer{Name: "TestDefault"},
				[]ModuleCall{
					ModuleCall{TerraformModuleDefinition: TerraformModuleDefinition{Name: "mod"},
						Variables: []BasicVariableMapping{BasicVariableMapping{VarName: "vpc_id"},
							BasicVariableMapping{VarName: "env", DefaultValue: "testing"},
							BasicVariableMapping{VarName: "tags", DefaultValues: []string{"owner=Josh", "keep"}},
						}},
				},
			},
			false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotT := TerraformProjectSkeleton{}
			err := gotT.UnmarshalYAML([]byte(tt.args.data))
			if (err != nil) != tt.wantErr {
				t.Errorf("SkelFromYaml() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(gotT, tt.wantT) {
				t.Errorf("SkelFromYaml() = %v, want %v", gotT, tt.wantT)
			}
		})
	}
}
