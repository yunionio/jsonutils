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
