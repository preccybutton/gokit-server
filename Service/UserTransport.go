package Service

import (
	"context"
	"encoding/json"
	"github.com/gorilla/mux"
	"github.com/pkg/errors"
	"net/http"
	"strconv"
	"xuexi/util"
)

func DecodeUserRequest(c context.Context, r *http.Request) (interface{}, error){
	vars := mux.Vars(r)
	if uid, ok := vars["uid"];ok{
		uid, _ := strconv.Atoi(uid)
		return UserRequest{
			uid,
			r.Method,
			r.URL.Query().Get("token"),	//token是get参数，没有会返回空的参数
		}, nil
	}
	return nil, errors.New("no")
}

func EncodeUserResponse(ctx context.Context, w http.ResponseWriter, response interface{}) error{
	w.Header().Set("Content-type", "application/json")
	return json.NewEncoder(w).Encode(response)
}

func MyErrEncoder(_ context.Context, err error, w http.ResponseWriter){
	contentType, body := "text/plain; charset=utf-8", []byte(err.Error())
	w.Header().Set("content-type", contentType)
	if myerr, ok :=err.(*util.MyError); ok{	//如果实现err的对象是MyError，则根据他的code设置
		w.WriteHeader(myerr.Code)
	}else{
		w.WriteHeader(404)	//否则默认设置404
	} //发送status code
	w.Write(body)

}