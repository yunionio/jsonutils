// Copyright 2019 Yunion
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

package jsonutils

import (
	"testing"
)

func TestYaml(t *testing.T) {
	type SUser struct {
		Name         string
		Passwd       string
		Keys         []string
		LockPassword bool
		FloatValue   float32
	}
	type SFile struct {
		Path    string
		Content string
	}
	type SCallback struct {
		Url string
	}
	type Config struct {
		Users       []SUser
		WriteFiles  []SFile
		Callback    *SCallback
		Runcmd      []string
		Bootcmd     []string
		Packages    []string
		DisableRoot int
		SshPwauth   int
	}
	conf := Config{
		Users: []SUser{
			{
				Name: "root",
				Keys: []string{
					"ssh-rsa AAAAA",
				},
			},
			{
				Name:         "yunion",
				Passwd:       "123456",
				LockPassword: true,
			},
		},
		WriteFiles: []SFile{
			{
				Content: "#\n\n127.0.0.1\tlocalhost\n10.0.0.1\t212222\n\n",
				Path:    "/etc/hosts",
			},
			{
				Content: "gobuild",
				Path:    "/etc/ansible/hosts",
			},
		},
		Callback: &SCallback{
			Url: "https://www.yunion.io/$INSTANCE_ID",
		},
		Runcmd: []string{
			"mkdir -p /var/run/httpd",
		},
	}
	jsonConf := Marshal(&conf)
	yaml := jsonConf.YAMLString()

	t.Logf("\n%s", yaml)

	jsonConf2, err := ParseYAML(yaml)
	if err != nil {
		t.Errorf("%s", err)
	} else {
		yaml2 := jsonConf2.YAMLString()
		t.Logf("\n%s", yaml2)

		if yaml != yaml2 {
			t.Errorf("yaml != yaml2")
		}
	}

}
