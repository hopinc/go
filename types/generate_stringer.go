//go:build ignore
// +build ignore

package main

import (
	"bytes"
	"os"
	"strings"
	"text/template"
)

// Defines the types that the string function is generated for.
var types = []string{
	"LeapUnavailableEvent", "LeapAvailableEvent", "LeapChannelStateUpdateEvent",
	"LeapPipeRoomAvailableEvent", "LeapPipeRoomUpdateEvent", "IgniteDeploymentPatchOpts",
	"Container", "ContainerMetadata", "DeploymentConfig", "DeploymentConfigPartial",
	"VolumeDefinition", "VolumeFormat", "RolloutState", "DeploymentRollout",
	"Build", "BuildMethod", "BuildMetadata", "BuildState", "IgniteGatewayUpdateOpts",
	"HealthCheck", "HealthCheckCreateOpts", "HealthCheckState",
}

const stringTemplate = `// String returns the string representation of this value. This function is auto-generated.
func (x {{.TypeData}}) String() string {
	return stringifyValue(x)
}`

const fileStart = `package types

// Code generated by generate_stringer.go; DO NOT EDIT.

//go:generate go run generate_stringer.go

`

func generateFile(typeData []string) string {
	fileContents := fileStart
	fileData := []string{}
	for _, t := range typeData {
		tpl, err := template.New("tpl").Parse(stringTemplate)
		if err != nil {
			panic(err)
		}
		buf := &bytes.Buffer{}
		err = tpl.Execute(buf, map[string]string{"TypeData": t})
		if err != nil {
			panic(err)
		}
		fileData = append(fileData, buf.String())
	}
	fileContents += strings.Join(fileData, "\n\n")
	return fileContents
}

func main() {
	// Generate the file.
	fileContents := generateFile(types)

	// Write the file.
	err := os.WriteFile("stringer_gen.go", []byte(fileContents+"\n"), 0644)
	if err != nil {
		panic(err)
	}
}
