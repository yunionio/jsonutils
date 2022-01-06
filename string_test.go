package jsonutils

import "testing"

func TestJSONFloat_String(t *testing.T) {
	type fields struct {
		JSONValue JSONValue
		data      float64
		bit       int
	}
	tests := []struct {
		name   string
		fields fields
		want   string
	}{
		{
			name: "float32",
			fields: fields{
				data: float64(0.9),
				bit:  32,
			},
			want: "0.9",
		},
		{
			name: "float64",
			fields: fields{
				data: float64(0.9),
				bit:  64,
			},
			want: "0.9",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			this := &JSONFloat{
				JSONValue: tt.fields.JSONValue,
				data:      tt.fields.data,
				bit:       tt.fields.bit,
			}
			if got := this.String(); got != tt.want {
				t.Errorf("JSONFloat.String() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestConfigNull(t *testing.T) {
	config := NewDict()
	config.Add(Marshal(nil), "config")
	t.Logf("config: %s", config.String())
	conf2, err := ParseString(config.String())
	if err != nil {
		t.Fatalf("ParseString error %s", err)
	}
	t.Logf("config2: %s", conf2.String())
	conf := conf2.(*JSONDict)
	if !conf.Contains("config") {
		conf.Add(NewDict(), "config")
	}
	if !conf.Contains("config", "default") {
		err = conf.Add(NewDict(), "config", "default")
		if err != nil {
			t.Fatalf("add config default fail %s", err)
		}
	}
	syncConf := map[string]string{
		"api_server": "https://127.0.0.1",
	}
	for k, v := range syncConf {
		if _, ok := conf.GetString("config", "default", k); ok == nil {
			continue
		} else {
			conf.Add(NewString(v), "config", "default", k)
		}
	}
	t.Logf("configs: %s", conf)
}
