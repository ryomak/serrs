package serrs

import (
	"fmt"
	"maps"
	"slices"
	"strings"
)

type CustomData interface {
	String() string
	Clone() CustomData
}

type DefaultCustomData map[string]any

func (d DefaultCustomData) String() string {
	bufString := strings.Builder{}
	bufString.WriteString("{")

	dd := make([]string, 0, len(d))
	for k, v := range d {
		dd = append(dd, fmt.Sprintf("%s:%+v", k, v))
	}
	slices.Sort(dd)

	bufString.WriteString(strings.Join(dd, ","))
	bufString.WriteString("}")
	return bufString.String()
}

func (d DefaultCustomData) Clone() CustomData {
	dd := maps.Clone(d)
	return dd
}
