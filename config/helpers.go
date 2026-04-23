package config

import (
	"fmt"

	"go.yaml.in/yaml/v3"
)

type IntOrString struct {
	IntVal *int
	StrVal *string
}

func (v *IntOrString) UnmarshalYAML(node *yaml.Node) error {
	switch node.Kind {
	case yaml.ScalarNode:
		var i int
		if err := node.Decode(&i); err == nil {
			v.IntVal = &i
			return nil
		}

		var s string
		if err := node.Decode(&s); err != nil {
			return err
		}
		v.StrVal = &s
		return nil
	default:
		return fmt.Errorf("invalid type for IntOrString")
	}
}

func NewIntOrStringForString(s string) IntOrString {
	return IntOrString{StrVal: &s}
}

func NewIntOrStringForInt(i int) IntOrString {
	return IntOrString{IntVal: &i}
}

type FloatOrString struct {
	FloatVal *float64
	StrVal   *string
}

func (v *FloatOrString) UnmarshalYAML(node *yaml.Node) error {
	switch node.Kind {
	case yaml.ScalarNode:
		var i float64
		if err := node.Decode(&i); err == nil {
			v.FloatVal = &i
			return nil
		}

		var s string
		if err := node.Decode(&s); err != nil {
			return err
		}
		v.StrVal = &s
		return nil
	default:
		return fmt.Errorf("invalid type for FloatOrString")
	}
}

func NewFloatOrStringForString(s string) FloatOrString {
	return FloatOrString{StrVal: &s}
}
