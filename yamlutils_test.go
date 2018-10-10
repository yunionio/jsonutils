package jsonutils

import (
	"testing"
)

func TestIndentLines(t *testing.T) {
	cases := []struct {
		in       []string
		in_array bool
		out      []string
	}{
		{[]string{"abc", "def"}, true, []string{"- abc", "  def"}},
		{[]string{"abc", "def"}, false, []string{"  abc", "  def"}},
	}
	for _, c := range cases {
		got := indentLines(c.in, c.in_array)
		if got[0] != c.out[0] {
			t.Errorf("IndentLines(%s, %v) = %s != %s", c.in, c.in_array, got, c.out)
		}
	}
}

func TestYaml(t *testing.T) {
	type SUser struct {
		Name string
		Passwd string
		Keys []string
	}
	type SFile struct {
		Path string
		Content string
	}
	type SCallback struct {
		Url string
	}
	type Config struct {
		Users []SUser
		WriteFiles []SFile
		Callback *SCallback
		Runcmd []string
		Bootcmd []string
		Packages []string
	}
	conf := Config{
		Users: []SUser {
			{
				Name: "root",
				Keys: []string {
					"ssh-rsa AAAAA",
				},
			},
			{
				Name: "yunion",
				Passwd: "123456",
			},
		},
		WriteFiles: []SFile {
			{
				Content: "127.0.0.1\tlocalhost\n",
				Path: "/etc/hosts",
			},
			{
				Content: "gobuild",
				Path: "/etc/ansible/hosts",
			},
		},
		Callback: &SCallback{
			Url: "https://www.yunion.io/$INSTANCE_ID",
		},
		Runcmd: []string {
			"mkdir -p /var/run/httpd",
		},
	}
	jsonConf := Marshal(&conf)
	t.Logf("\n%s", jsonConf.YAMLString())
}
