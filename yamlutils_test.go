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
