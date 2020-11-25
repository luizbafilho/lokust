package cmd

import "fmt"

func buildResourceName(testName string, resourceType ...string) string {
	name := fmt.Sprintf("lokust-%s", testName)

	if len(resourceType) > 1 {
		name += "-" + resourceType[0]
	}

	return name
}
