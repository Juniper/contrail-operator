// Copyright 2018 The Operator-SDK Authors
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package e2e

import (
	"os"
	"testing"

	f "github.com/operator-framework/operator-sdk/pkg/test"
)

var buildTag string
var cemRelease string

func TestMain(m *testing.M) {
	os.Setenv("TEST_NAMESPACE", "contrail")
	scmRevision := getEnv("BUILD_SCM_REVISION", "latest")
	scmBranch := getEnv("BUILD_SCM_BRANCH", "master")
	buildTag = scmBranch + "." + scmRevision
	cemRelease = getEnv("CEM_RELEASE", "master-latest")
	f.MainEntry(m)
}

func getEnv(key, fallback string) string {
	if value, ok := os.LookupEnv(key); ok {
		return value
	}
	return fallback
}
