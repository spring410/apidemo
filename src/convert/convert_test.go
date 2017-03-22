package convert

import (
	"fmt"
	"testing"
)

func Test_StringToInt(t *testing.T) {
	fmt.Println("start to test..")

	//int, on 32bit system,it is int32, on 64bit system, it is int64.
	var v int
	var e error

	v, e = StringToInt("-1")
	if e != nil {
		t.Error("Error converting the string to int.")
	} else {
		if -1 != v {
			t.Error("Error")
		}
	}

	v, e = StringToInt("-10000")
	if e != nil {
		t.Error("Error converting the string to int.")
	} else {
		if -10000 != v {
			t.Error("Error")
		}
	}

	v, e = StringToInt("0")
	if e != nil {
		t.Error("Error converting the string to int.")
	} else {
		if 0 != v {
			t.Error("Error")
		}
	}

	v, e = StringToInt("123456789")
	if e != nil {
		t.Error("Error converting the string to int.")
	} else {
		if 123456789 != v {
			t.Error("Error")
		}
	}
}

func Test_StringToIntBase16(t *testing.T) {

	var v int64
	var e error
	var s string
	s = "A"

	v, e = StringToIntBase16(s)
	if e != nil {
		t.Error("Error converting the string to int.")
	} else {
		if 0xA != v {
			fmt.Println(s, "?", v)
			t.Error("Error")
		}
	}

	s = "10"
	v, e = StringToIntBase16(s)
	if e != nil {
		t.Error("Error converting the string to int.")
	} else {
		if 0x10 != v {
			fmt.Println(s, "?", v)
			t.Error("Error")
		}
	}

	s = "11"
	v, e = StringToIntBase16(s)
	if e != nil {
		t.Error("Error converting the string to int.")
	} else {
		if 0x11 != v {
			fmt.Println(s, "?", v)
			t.Error("Error")
		}
	}

	s = "F"
	v, e = StringToIntBase16(s)
	fmt.Println(v)
	if e != nil {
		t.Error("Error converting the string to int.")
	} else {
		if 0xF != v {
			fmt.Println(s, "?", v)
			t.Error("Error")
		}
	}

}

func Test_Int64ToString(t *testing.T) {
	var str string
	var v int64
	v = -100
	str = Int64ToString(v)
	fmt.Println(str)
}

func Test_IntToString(t *testing.T) {
	var str string
	var v int
	v = -10090
	str = IntToString(v)
	fmt.Println(str)
}
