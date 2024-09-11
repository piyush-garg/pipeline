/*
Copyright 2022 The Tekton Authors

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package git

import (
	"context"
	"reflect"
	"strings"

	"github.com/tektoncd/pipeline/pkg/resolution/resolver/framework"
)

const (
	// DefaultTimeoutKey is the configuration field name for controlling
	// the maximum duration of a resolution request for a file from git.
	DefaultTimeoutKey = "fetch-timeout"

	// DefaultURLKey is the configuration field name for controlling
	// the git url to fetch the remote resource from.
	DefaultURLKey = "default-url"

	// DefaultRevisionKey is the configuration field name for controlling
	// the revision to fetch the remote resource from.
	DefaultRevisionKey = "default-revision"

	// DefaultOrgKey is the configuration field name for setting a default organization when using the SCM API.
	DefaultOrgKey = "default-org"

	// ServerURLKey is the config map key for the SCM provider URL
	ServerURLKey = "server-url"
	// SCMTypeKey is the config map key for the SCM provider type
	SCMTypeKey = "scm-type"
	// APISecretNameKey is the config map key for the token secret's name
	APISecretNameKey = "api-token-secret-name"
	// APISecretKeyKey is the config map key for the containing the token within the token secret
	APISecretKeyKey = "api-token-secret-key"
	// APISecretNamespaceKey is the config map key for the token secret's namespace
	APISecretNamespaceKey = "api-token-secret-namespace"
)

type GitResolverConfig struct {
	ScmTokens map[string]ScmInfo
}

type ScmInfo struct {
	Timeout            string `json:"fetch-timeout"`
	URL                string `json:"default-url"`
	Revision           string `json:"default-revision"`
	Org                string `json:"default-org"`
	ServerURL          string `json:"server-url"`
	SCMType            string `json:"scm-type"`
	APISecretName      string `json:"api-token-secret-name"`
	APISecretKey       string `json:"api-token-secret-namespace"`
	APISecretNamespace string `json:"api-secret-namespace"`
}

func GetGitConfig(ctx context.Context) GitResolverConfig {
	var scmInfo interface{} = &ScmInfo{}
	structType := reflect.TypeOf(scmInfo).Elem()
	gitResolverConfig := GitResolverConfig{ScmTokens: map[string]ScmInfo{}}
	conf := framework.GetResolverConfigFromContext(ctx)
	for k, v := range conf {
		if k == "identifier" {
			_, ok := gitResolverConfig.ScmTokens[v]
			if ok {
				continue
			}
			gitResolverConfig.ScmTokens[k] = ScmInfo{}
		}
		key := strings.Split(k, ".")
		if len(key) >= 3 && key[0] == "provider" {
			keyValue := strings.Join(key[2:], ".")
			_, ok := gitResolverConfig.ScmTokens[key[1]]
			if !ok {
				gitResolverConfig.ScmTokens[key[1]] = ScmInfo{}
			}
			for i := 0; i < structType.NumField(); i++ {
				field := structType.Field(i)
				fieldName := field.Name
				jsonTag := field.Tag.Get("json")
				if keyValue == jsonTag {
					tokenDetails := gitResolverConfig.ScmTokens[key[1]]
					var scm interface{} = &tokenDetails
					structValue := reflect.ValueOf(scm).Elem()
					structValue.FieldByName(fieldName).SetString(v)
					gitResolverConfig.ScmTokens[key[1]] = structValue.Interface().(ScmInfo)
				}
			}
		}
		if key[0] != "provider" {
			_, ok := gitResolverConfig.ScmTokens["default"]
			if !ok {
				gitResolverConfig.ScmTokens["default"] = ScmInfo{}
			}
			for i := 0; i < structType.NumField(); i++ {
				field := structType.Field(i)
				fieldName := field.Name
				jsonTag := field.Tag.Get("json")
				if k == jsonTag {
					tokenDetails := gitResolverConfig.ScmTokens["default"]
					var scm interface{} = &tokenDetails
					structValue := reflect.ValueOf(scm).Elem()
					structValue.FieldByName(fieldName).SetString(v)
					gitResolverConfig.ScmTokens["default"] = structValue.Interface().(ScmInfo)
				}
			}
		}
	}
	return gitResolverConfig
}
