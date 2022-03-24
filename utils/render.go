package utils

import (
	"fmt"

	"github.com/PrinceMerluza/devcenter-content-linter/transform_data"
	"github.com/tidwall/pretty"
)

func Render(data string) {
	if transform_data.TemplateFile != "" {
		mp := transform_data.ConvertJsonToMap(data)
		res := transform_data.ProcessTemplateFile(mp)
		fmt.Println(res)
		return
	}
	if transform_data.TemplateStr != "" {
		mp := transform_data.ConvertJsonToMap(data)
		res := transform_data.ProcessTemplateStr(mp)
		fmt.Println(res)
		return
	}
	result := pretty.Pretty([]byte(data))
	fmt.Printf("%s", result)
}
