package gtpv2_test

import (
	"fmt"
	"testing"

	"github.com/blorticus/gtpv2"
)

var badYamlDefintions = []string{
	"lkaj;ion",
	`---
Gtpv2Pdus:
	- Name: malprop
	  Type: Yoodle
`,
}

func TestInvalidValuesForReadYamlTemplateFromString(t *testing.T) {
	for definitionNumber, definitionString := range badYamlDefintions {
		if _, err := gtpv2.ReadYamlTemplateFromString(definitionString); err == nil {
			t.Errorf("[ReadYamlTemplateFromString] On badYamlDefinition at index (%d) expected error, but received none", definitionNumber)
		}
	}
}

func errorIfValidTemplateReadFails(testname string, template *gtpv2.Template, errorFromRead error) error {
	if errorFromRead != nil {
		return fmt.Errorf("%s, expected no error, got = (%s)", testname, errorFromRead.Error())
	}

	if template == nil {
		return fmt.Errorf("%s, expected template != nil; template == nil", testname)
	}

	return nil
}

var validYamlDefinitions = []string{
	`---
Gtpv2Pdus:
    - Name: csr
      Type: CreateSessionRequest
`,
}

func TestReadYamlTemplateFromString(t *testing.T) {
	template, err := gtpv2.ReadYamlTemplateFromString("")

	if err = errorIfValidTemplateReadFails("[ReadYamlTemplateFromString] on empty string", template, err); err != nil {
		t.Fatal(err)
	}

	template, err = gtpv2.ReadYamlTemplateFromString(validYamlDefinitions[0])

	if err = errorIfValidTemplateReadFails("[ReadYamlTemplateFromString] validYamlDefinition[0]", template, err); err != nil {
		t.Fatal(err)
	}

}
