package handlers

import (
	"encoding/json"
	"net/http"
)

/*
 |    Function: GoldPrice
 |      Author: Mr.Sancho
 |        Date: 2016-04-17
 |   Arguments:
 |      Return:
 | Description: 获取金币价格列表
 |
*/
func GoldPrice(r *http.Request) (int, string) {
	jsonRlt, _ := json.Marshal(gGoldBeansPrices)
	return 200, string(jsonRlt)
}
