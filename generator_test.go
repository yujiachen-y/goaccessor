package main

import "testing"

func TestFillTypeArguments(t *testing.T) {
	type testCase struct {
		param, arg string
		input      string
		output     string
	}
	for _, tc := range []testCase{
		{
			param:  "param",
			arg:    "arg",
			input:  "param param1 param2 _param3 param4_",
			output: "arg param1 param2 _param3 param4_",
		},
		{
			param:  "T",
			arg:    "int",
			input:  "map[T]t",
			output: "map[int]t",
		},
		{
			param:  "T",
			arg:    "int",
			input:  "chan T <-chan T chan<- T chan a <-chan b chan<- c",
			output: "chan int <-chan int chan<- int chan a <-chan b chan<- c",
		},
		{
			param:  "T",
			arg:    "int",
			input:  "T is a type parameter, while T1 is not",
			output: "int is a type parameter, while T1 is not",
		},
		{
			param:  "T",
			arg:    "int",
			input:  "This is T",
			output: "This is int",
		},
		{
			param:  "T",
			arg:    "int",
			input:  "(T)",
			output: "(int)",
		},
		{
			param:  "param",
			arg:    "arg",
			input:  "",
			output: "",
		},
		{
			param:  "param",
			arg:    "param",
			input:  "param should remain param",
			output: "param should remain param",
		},
		{
			param:  "T",
			arg:    "int",
			input:  "func[T]()T is a bad case, the expected result is func()int",
			output: "func[int]()int is a bad case, the expected result is func()int",
		},
	} {
		got := fillTypeArguments(tc.input, tc.param, tc.arg)
		if got != tc.output {
			t.Errorf("got %s, expected %s", got, tc.output)
		}
	}
}
