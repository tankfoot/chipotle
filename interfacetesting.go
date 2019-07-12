package main

import (
	"fmt"
	"encoding/json"

	"github.com/iancoleman/orderedmap"
)

type DataOutput struct {
    Speech string `json:"speech"`
    Entity interface{} `json:"entity"`
}

func main(){
	var p DataOutput
	o := orderedmap.New()
	o.Set("rice", "no rice")
	p.Speech = "hello"
	p.Entity = o
	b, _ := json.Marshal(p)
	fmt.Println(string(b))

}