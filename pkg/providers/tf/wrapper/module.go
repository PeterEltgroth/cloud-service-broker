// Copyright 2018 the Service Broker Project Authors.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//      http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package wrapper

import (
	"fmt"
	"sort"

	"github.com/pivotal/cloud-service-broker/pkg/validation"
	"github.com/hashicorp/hcl2/gohcl"
	"github.com/hashicorp/hcl2/hclparse"
)

// ModuleDefinition represents a module in a Terraform workspace.
type ModuleDefinition struct {
	Name       string
	Definition string
}

var _ (validation.Validatable) = (*ModuleDefinition)(nil)

// Validate checks the validity of the ModuleDefinition struct.
func (module *ModuleDefinition) Validate() (errs *validation.FieldError) {
	return errs.Also(
		validation.ErrIfBlank(module.Name, "Name"),
		validation.ErrIfNotTerraformIdentifier(module.Name, "Name"),
		validation.ErrIfNotHCL(module.Definition, "Definition"),
	)
}

func (module *ModuleDefinition) decode() (terraformModuleHcl, error) {
	defn := terraformModuleHcl{}
	parser := hclparse.NewParser()
	f, parseDiags := parser.ParseHCL([]byte(module.Definition), "")
	if parseDiags.HasErrors() {
		return defn, fmt.Errorf(parseDiags.Error())
	}
	if err := gohcl.DecodeBody(f.Body, nil, &defn ); err != nil {
		return defn, err
	}
	return defn, nil
}

// Inputs gets the input parameter names for the module.
func (module *ModuleDefinition) Inputs() ([]string, error) {
	defn, err := module.decode()

	return sortedKeys(defn.Inputs), err
}

// Outputs gets the output parameter names for the module.
func (module *ModuleDefinition) Outputs() ([]string, error) {
	defn, err := module.decode()

	return sortedKeys(defn.Outputs), err
}

func sortedKeys(m map[string]interface{}) []string {
	var keys []string
	for key, _ := range m {
		keys = append(keys, key)
	}

	sort.Slice(keys, func(i int, j int) bool { return keys[i] < keys[j] })
	return keys
}

// terraformModuleHcl is a struct used for marshaling/unmarshaling details about
// Terraform modules.
//
// See https://www.terraform.io/docs/modules/create.html for their structure.
type terraformModuleHcl struct {
	Inputs  map[string]interface{} `hcl:"variable"`
	Outputs map[string]interface{} `hcl:"output"`
}
