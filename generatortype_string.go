// Code generated by "stringer -type GeneratorType"; DO NOT EDIT.

package main

import "strconv"

func _() {
	// An "invalid array index" compiler error signifies that the constant values have changed.
	// Re-run the stringer command to generate them again.
	var x [1]struct{}
	_ = x[GeneratorTypeUnknown-0]
	_ = x[GeneratorTypeVariable-1]
	_ = x[GeneratorTypeStructure-2]
	_ = x[GeneratorTypeField-3]
}

const _GeneratorType_name = "GeneratorTypeUnknownGeneratorTypeVariableGeneratorTypeStructureGeneratorTypeField"

var _GeneratorType_index = [...]uint8{0, 20, 41, 63, 81}

func (i GeneratorType) String() string {
	if i < 0 || i >= GeneratorType(len(_GeneratorType_index)-1) {
		return "GeneratorType(" + strconv.FormatInt(int64(i), 10) + ")"
	}
	return _GeneratorType_name[_GeneratorType_index[i]:_GeneratorType_index[i+1]]
}