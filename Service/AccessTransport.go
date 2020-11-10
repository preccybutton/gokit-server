package Service

import (
	"context"
	"encoding/json"
	"github.com/pkg/errors"
	"github.com/tidwall/gjson"
	"io/ioutil"
	"net/http"
)

func DecodeAccessRequest(c context.Context, r *http.Request) (interface{}, error){
	body, _ := ioutil.ReadAll(r.Body)
	result := gjson.Parse(string(body))	//解析json结构体，gjson专门用于解析json，支持多种匹配
	if result.IsObject(){
		username := result.Get("username")
		userpass := result.Get("userpass")
		return AccessRequest{username.String(), userpass.String(), r.Method}, nil
	}
	return nil, errors.New("参数错误")
}

func EncodeAccessResponse(ctx context.Context, w http.ResponseWriter, response interface{}) error{
	w.Header().Set("Content-type", "application/json")
	return json.NewEncoder(w).Encode(response)
}