package main

import (
	"flag"
	"os"
	"reflect"
	"testing"
)

func Test_getParams(t *testing.T) {
	tests := []struct {
		name    string
		want    paramBlueprint
		wantErr bool
		osArgs  []string
	}{
		{
			"Valid parameters",
			paramBlueprint{"https://github.com/GenesysCloudBlueprints/angular-app-with-genesys-cloud-sdk", "./default_rule.json"},
			false,
			[]string{"cmd", "https://github.com/GenesysCloudBlueprints/angular-app-with-genesys-cloud-sdk", "./default_rule.json"},
		},
		{"No parameters", paramBlueprint{}, true, []string{"cmd"}},
		{
			"Missing rule parameter",
			paramBlueprint{},
			true,
			[]string{"cmd", "https://github.com/GenesysCloudBlueprints/angular-app-with-genesys-cloud-sdk"},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			actualOsArgs := os.Args

			// Restore original os.Args
			defer func() {
				os.Args = actualOsArgs
				flag.CommandLine = flag.NewFlagSet(os.Args[0], flag.ExitOnError)
			}()

			os.Args = tt.osArgs

			got, err := getParams()
			if (err != nil) != tt.wantErr {
				t.Errorf("getFileData() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("getFileData() = %v, want %v", got, tt.want)
			}
		})
	}
}
