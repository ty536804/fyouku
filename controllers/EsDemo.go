package controllers

import (
	"encoding/json"
	"fmt"
	"fyoukuApi/Services/Es"
	"github.com/astaxie/beego"
)

type EsDemoController struct {
	beego.Controller
}

func (e *EsDemoController) Add() {
	body := map[string]interface{}{
		"id":    1,
		"title": "张三",
	}
	Es.EsAdd("fyouku_demo", "user-1", body)
	e.Ctx.WriteString("add")
}

func (e *EsDemoController) Edit() {
	body := map[string]interface{}{
		"id":    1,
		"title": "李四",
	}
	Es.EsAdd("fyouku_demo", "user-1", body)
	e.Ctx.WriteString("edit")
}

func (e *EsDemoController) Delete() {
	Es.EsDelete("fyouku_demo", "user-1")
	e.Ctx.WriteString("delete")
}

// @router /es/search [*]
func (e *EsDemoController) Search() {
	sort := []map[string]string{
		map[string]string{"id": "desc"}}
	query := map[string]interface{}{
		"bool": map[string]interface{}{
			"must": map[string]interface{}{
				"term": map[string]interface{}{
					"id": 1,
				},
			},
		},
	}

	res := Es.EsSearch("fyouku_demo", query, 0, 10, sort)
	var resData []ResData
	for _, v := range res.Hits {
		var data ResData
		err := json.Unmarshal([]byte(v.Source), &data)
		if err != nil {
			fmt.Println(err)
		}
		resData = append(resData, data)
	}
	e.Ctx.WriteString("search")
}

type ResData struct {
	Title string
	Id    int
}
