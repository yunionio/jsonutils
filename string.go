package jsonutils

import (
	"bytes"
	"fmt"
)

func (this *JSONString) String() string {
	return quoteString(this.data)
}

func (this *JSONValue) String() string {
	return "null"
}

func (this *JSONInt) String() string {
	return fmt.Sprintf("%d", this.data)
}

func (this *JSONFloat) String() string {
	return fmt.Sprintf("%f", this.data)
}

func (this *JSONBool) String() string {
	if this.data {
		return "true"
	} else {
		return "false"
	}
}

func (this *JSONDict) String() string {
	var buffer bytes.Buffer
	buffer.WriteByte('{')
	var idx = 0
	for _, k := range this.SortedKeys() {
		v := this.data[k]
		if idx > 0 {
			buffer.WriteString(",")
		}
		buffer.WriteString(quoteString(k))
		buffer.WriteByte(':')
		buffer.WriteString(v.String())
		idx++
	}
	buffer.WriteByte('}')
	return buffer.String()
}

func (this *JSONArray) String() string {
	var buffer bytes.Buffer
	buffer.WriteByte('[')
	for idx, v := range this.data {
		if idx > 0 {
			buffer.WriteString(",")
		}
		buffer.WriteString(v.String())
	}
	buffer.WriteByte(']')
	return buffer.String()
}
