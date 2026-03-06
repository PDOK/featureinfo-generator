package main

import (
	"encoding/json"
	"os"
	"testing"

	"github.com/pdok/featureinfo-generator/pkg/types"
)

func Test_writeHTMLFileV2(t *testing.T) {
	data, err := os.ReadFile("../test/scheme_v2.json")
	if err != nil {
		t.Errorf("File not found.")
	}

	var scheme types.Scheme
	err = json.Unmarshal(data, &scheme)
	if err != nil {
		t.Errorf("Unmarshal error.")
	}
	destFolder := "../test/html/v2/"
	filename := "all"
	type args struct {
		scheme     types.Scheme
		destFolder *string
		filename   *string
	}
	tests := []struct {
		name string
		args args
		want int
	}{
		{
			name: "Test 1 layers",
			args: args{
				scheme:     scheme,
				destFolder: &destFolder,
				filename:   &filename,
			},
			want: 1,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := generateHTMLTemplateForAll(tt.args.scheme, tt.args.destFolder, tt.args.filename); got != tt.want {
				t.Errorf("generateHtmlTemplateForAll() = %v, want %v", got, tt.want)
			}
		})
	}
}
