package main

import "testing"

var ordertype = map[string][]string{
    "burrito bowl": []string{"bowl", "burrito bowl"},	
	"tacos": []string{"taco"},
	"salad": []string{"salad"},
	"kid's meal": []string{"kid"},
	"sides & drinks": []string{"side", "drink"},
}

func TestSingleMatch(t *testing.T) {
	var tests = []struct{
		q, want string
	}{
		{"burrito", "burrito"}
		{"apple", ""}
	}
	for _, c := range tests{
		got := SingleMatch(c, ordertype)
		if got != c.want {
			t.Errorf("SingleMatch(%q) == %q, want %q", c.q, got, c.want)
		}
	}
}