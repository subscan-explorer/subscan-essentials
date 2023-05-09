package util

import "testing"

func TestMapToStruct(t *testing.T) {
	type Birthday struct {
		Year  int
		Month int
		Day   int
	}
	type Person struct {
		Name     string
		Age      int
		Birthday Birthday
	}
	input := map[string]interface{}{
		"Name": "John",
		"Age":  30,
		"Birthday": map[string]interface{}{
			"Year":  1990,
			"Month": 1,
			"Day":   1,
		},
	}
	ret, err := MapInterfaceAsStruct[Person](input)
	if err != nil {
		t.Error(err)
	}
	want := Person{
		Name: "John",
		Age:  30,
		Birthday: Birthday{
			Year:  1990,
			Month: 1,
			Day:   1,
		},
	}
	if ret != want {
		t.Errorf("got %v, want %v", ret, want)
	}
}

func TestMapToStruct2(t *testing.T) {
	type Birthday struct {
		Year  int `json:"year"`
		Month int `json:"month_name"`
		Day   int `json:"day"`
	}
	type Person struct {
		Name     string
		Age      int
		Birthday Birthday
	}
	input := map[string]interface{}{
		"Name": "John",
		"Age":  30,
		"Birthday": map[string]interface{}{
			"year":       1990,
			"month_name": 1,
			"day":        1,
		},
	}
	ret, err := MapInterfaceAsStruct[Person](input)
	if err != nil {
		t.Error(err)
	}
	want := Person{
		Name: "John",
		Age:  30,
		Birthday: Birthday{
			Year:  1990,
			Month: 1,
			Day:   1,
		},
	}
	if ret != want {
		t.Errorf("got %v, want %v", ret, want)
	}
}
