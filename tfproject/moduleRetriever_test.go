package tfproject

import (
	"reflect"
	"testing"
)

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
				localLocation: "target/external/simpleterraform",
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
