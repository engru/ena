package binding

import (
	"strings"

	"github.com/pkg/errors"

	"github.com/lsytj0413/ena/to"
)

type bindTag struct {
	source       Source
	name         *string
	defaultValue *string
}

func (t *bindTag) Name() string {
	return to.String(t.name)
}

func (t *bindTag) DefaultValue() string {
	return to.String(t.defaultValue)
}

func (t *bindTag) Validate() error {
	if _, ok := availableSource[t.source]; !ok {
		return errors.Errorf("Invalid source %v", t.source)
	}

	switch t.source {
	case SourceHeader, SourcePath, SourceQuery:
		if t.name == nil || (*t.name) == "" {
			return errors.Errorf("Name must not been empty when source is [%v,%v,%v]", SourceHeader, SourcePath, SourceQuery)
		}
	case SourceAuto, SourceBody:
		if t.name != nil || t.defaultValue != nil {
			return errors.Errorf("Name and default must be empty when source is [%v,%v]", SourceAuto, SourceBody)
		}
	}

	return nil
}

func parseBindParameterTag(tag string) (*bindTag, error) {
	results := strings.Split(tag, ",")
	length := len(results)
	if length == 0 || length > 3 {
		return nil, errors.Errorf("Invalid tag %v", tag)
	}

	v := &bindTag{}
	if length >= 1 {
		v.source = Source(strings.TrimSpace(results[0]))
	}
	if length >= 2 {
		v.name = to.StringPtr(strings.TrimSpace(results[1]))
	}
	if length >= 3 {
		index := strings.Split(results[2], "=")
		if len(index) != 2 || index[0] != "default" {
			return nil, errors.Errorf("Invalid default value %v", results[2])
		}

		v.defaultValue = to.StringPtr(strings.TrimSpace(index[1]))
	}

	return v, v.Validate()
}
