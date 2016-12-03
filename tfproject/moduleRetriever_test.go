package tfproject

import (
	"bytes"
	"flag"
	"fmt"
	"io/ioutil"
	"path/filepath"
	"reflect"
	"testing"
)

var update = flag.Bool("update", false, "update golden files")

func TestExternalModule_Fetch(t *testing.T) {
	type fields struct {
		Name string
		URI  string
	}
	tests := []struct {
		name    string
		fields  fields
		want    TerraformModuleDefinition
		wantErr bool
	}{
		{
			"simplemodule",
			fields{URI: "./test-fixtures/simpleterraform", Name: "proj"},
			TerraformModuleDefinition{
				Name:          "simpleterraform",
				RequiredVars:  []string{"location"},
				OptionalVars:  []string{"greeting"},
				URI:           "./test-fixtures/simpleterraform",
				LocalLocation: "target/external/simpleterraform",
				Outputs:       []string{"out1", "out2"},
			},
			false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := ExternalModule{
				Name: tt.fields.Name,
				URI:  tt.fields.URI,
			}

			golden := filepath.Join("test-fixtures", tt.name+".golden")
			if *update {
				ioutil.WriteFile(golden, []byte("actual"), 0644)

			}

			expected, _ := ioutil.ReadFile(golden)
			if !bytes.Equal([]byte("actual"), expected) {
				fmt.Printf("placeholder")
			}

			got, err := m.Fetch()
			if (err != nil) != tt.wantErr {
				t.Errorf("ExternalModule.Fetch() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("ExternalModule.Fetch() = %v, want %v", got, tt.want)
			}
		})
	}
}
