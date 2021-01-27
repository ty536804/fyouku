package Es

import (
	"encoding/json"
	"fmt"
	"github.com/astaxie/beego/httplib"
)

var esUrl string

func init() {
	esUrl = "http://127.0.0.1:9200"
}

type ReqSearchData struct {
	Hits HisData `json:"hits"`
}

type HisData struct {
	Total TotalData    `json:"total"`
	Hits  []HisTwoData `json:"hits"`
}

type HisTwoData struct {
	Source json.RawMessage `json:"source"`
}

type TotalData struct {
	Value    int
	Relation string
}

// @Title EsSearch
// @Description 搜索功能
// @Param indexName string 索引
// @Param query map[string]interface{} 查询参数
// @Param form int
// @Param size int
// @Param sort []map[string]string 排序
func EsSearch(indexName string, query map[string]interface{}, form int, size int, sort []map[string]string) HisData {
	searchQuery := map[string]interface{}{
		"query": query,
		"form":  form,
		"size":  size,
		"sort":  sort,
	}
	req := httplib.Post(esUrl + indexName + "/_search")
	req.JSONBody(searchQuery)
	str, err := req.String()
	if err != nil {
		fmt.Println(err)
	}
	var stb ReqSearchData
	err = json.Unmarshal([]byte(str), &stb)
	return stb.Hits
}

// @Title EsAdd
// @Description 添加ES
// @Param indexName string 					索引
// @Param Id 		string					ID
// @Param Body 		map[string]interface{} 	内容
func EsAdd(indexName, Id string, body map[string]interface{}) bool {
	req := httplib.Post(esUrl + indexName + "/_doc/" + Id)
	req.JSONBody(body)
	_, err := req.String()
	if err != nil {
		fmt.Println("添加失败:", err)
		return false
	}
	return true
}

// @Title EsEdit
// @Description 修改Es
// @Param indexName string 					索引
// @Param Id 		string					ID
// @Param Body 		map[string]interface{} 	内容
func EsEdit(indexName, Id string, body map[string]interface{}) bool {
	bodyData := map[string]interface{}{
		"doc": body,
	}
	req := httplib.Post(esUrl + indexName + "/_doc/" + Id + "/_update")
	req.JSONBody(bodyData)
	_, err := req.String()
	if err != nil {
		fmt.Printf("ES编辑失败:%s,ID是：%s\n", err, Id)
		return false
	}
	return true
}

// @Title EsDelete
// @Description 删除Es
// @Param indexName string 					索引
// @Param Id 		string					ID
func EsDelete(indexName, Id string) bool {
	req := httplib.Delete(esUrl + indexName + "/_doc/" + Id)
	_, err := req.String()
	if err != nil {
		fmt.Println("ES删除数据失败:", err)
		return false
	}
	return true
}
