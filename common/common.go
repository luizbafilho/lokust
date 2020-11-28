package common

import (
	loadtestsv1beta1 "github.com/luizbafilho/lokust/apis/loadtests/v1beta1"
)

func MakeLabels(test loadtestsv1beta1.LocustTest, component string) map[string]string {
	return map[string]string{
		"lokust-loadtest-name": test.Name,
		"locust-component":     component,
	}
}
