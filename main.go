package main

import (
	"encoding/json"
	"fmt"
	"os"
	"strings"

	"github.com/getkin/kin-openapi/openapi3"
	"github.com/spf13/cobra"
)

const (
	componentsPrefix = "#/components/schemas/"
)

var (
	filePath   = "openapi.yml"
	respCodes  = []string{}
	mimeType   = ""
	outputFile string
)

func convert(cmd *cobra.Command, args []string) {
	codes := map[string]struct{}{}
	for _, code := range respCodes {
		codes[code] = struct{}{}
	}
	isAllCodes := len(codes) == 0

	swg, err := openapi3.NewSwaggerLoader().LoadSwaggerFromFile(filePath)
	if err != nil {
		panic(err)
	}
	for _, path := range swg.Paths {
		for method, op := range path.Operations() {
			for code, resp := range op.Responses {
				if _, ok := codes[code]; !isAllCodes && !ok {
					continue
				}
				respComponentName := strings.Title(op.OperationID) + strings.Title(method) + code + "Resp"
				if resp.Ref != "" || resp.Value == nil {
					continue
				}
				mt, ok := resp.Value.Content[mimeType]
				if !ok || mt.Schema == nil || mt.Schema.Ref != "" {
					continue
				}
				swg.Components.Schemas[respComponentName] = &openapi3.SchemaRef{
					Value: mt.Schema.Value,
				}
				mt.Schema.Ref = componentsPrefix + respComponentName
			}
		}
	}
	b, _ := json.MarshalIndent(swg, "", "  ")
	if err := os.WriteFile(outputFile, b, os.FileMode(0644)); err != nil {
		panic(err)
	}
}

func init() {
	cmd := &cobra.Command{
		Use: "run [--in=./a.yaml] [--respcodes=200,400] [--mimetype='application/json'] [--out=a.json]",
		Run: convert,
	}
	cmd.PersistentFlags().StringVar(&filePath, "in", "./openapi.yaml", "Input yaml file path\ndefault: ./openapi.yaml")
	cmd.PersistentFlags().StringSliceVar(&respCodes, "respcodes", nil, "Anonymous responses with these codes will be converted\ndefault: all code\nexample: 200,400")
	cmd.PersistentFlags().StringVar(&mimeType, "mimetype", "application/json", "Anonymous responses with this mimeType will be converted\ndefault: application/json")
	cmd.PersistentFlags().StringVar(&outputFile, "out", "a.json", "")

	rootCmd.AddCommand(cmd)
}

var rootCmd = &cobra.Command{
	Use:   "openapi-resp-convert",
	Short: `This program will convert the anonymous response (root schema has no ref to components) to ref with components; input is a yaml file and output as a JSON.`,
}

func main() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
