package main

import (
	"log"
	"slices"
	"testing"

	"github.com/gospider007/gson"
)

func TestFind(t *testing.T) {
	txt := `
	{
		"aa2":11,
		"bb":{
			"aa":22
		}
	}
	`
	jsonData, err := gson.Decode(txt)
	if err != nil {
		t.Error(err)
	}
	log.Print(jsonData.Find("aa"))
	if jsonData.Find("aa").Int() != 22 {
		t.Error("find aa error")
	}
}
func TestFinds(t *testing.T) {
	txt := `
	{
		"aa":11,
		"bb":{
			"aa":22,
			"bb":{
				"aa":22
			}
		}
	}
	`
	jsonData, err := gson.Decode(txt)
	if err != nil {
		t.Error(err)
	}
	rs := []string{}
	for _, v := range jsonData.Finds("aa") {
		rs = append(rs, v.String())
	}
	log.Print(rs)
	if !slices.Equal(rs, []string{"11", "22", "22"}) {
		t.Error("find aa error")
	}
}
