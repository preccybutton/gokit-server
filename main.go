package main

import (
	"flag"
	"fmt"
	"github.com/go-kit/kit/endpoint"
	kitlog "github.com/go-kit/kit/log"
	httptransport "github.com/go-kit/kit/transport/http"
	"github.com/gorilla/mux"
	"golang.org/x/time/rate"
	"log"
	"net/http"
	"os"
	"os/signal"
	"strconv"
	"syscall"
	"xuexi/Service"
	"xuexi/util"
)

func main()  {

	name := flag.String("name", "", "服务名称")
	port := flag.Int("p", 0, "服务端口")
	flag.Parse()
	if *name == ""{
		log.Fatal("请指定服务名")
	}
	if *port == 0{
		log.Fatal("请指定端口")
	}
	util.SetServiceNameAndPort(*name, *port)

	var logger kitlog.Logger
	{
		logger = kitlog.NewLogfmtLogger(os.Stdout)
		logger = kitlog.WithPrefix(logger, "mykit", "1.0")
		logger = kitlog.With(logger, "time", kitlog.DefaultTimestampUTC)
		logger = kitlog.With(logger, "caller", kitlog.DefaultCaller)
	}
	user := &Service.UserService{}
	limit := rate.NewLimiter(1,5)
	var endp endpoint.Endpoint
	{
		endp = Service.UserServiceLogMiddleware(logger)(Service.GenUserEndPoint(user))	//生成一个日志中间件，同一生成日志
		endp = Service.RateLimit(limit)(endp)  //生成一个限流中间件中间件，然后把endpoint作为参数传入进去
		endp = Service.CheckTokenMiddleware()(endp) //生成一个token验证中间件
	}


	//自定义option函数，里面可以填入不同情况的option函数，这里是自定义错误处理函数，不使用默认的Encoder
	options := []httptransport.ServerOption{
		httptransport.ServerErrorEncoder(Service.MyErrEncoder),
	}

	serverHandler := httptransport.NewServer(endp, Service.DecodeUserRequest, Service.EncodeUserResponse, options...)

	//增加handler用于获取用户token
	accessService := &Service.AccessService{}
	accessServiceEndpoint := Service.AccessEndpoint(accessService)
	accessHandler := httptransport.NewServer(accessServiceEndpoint, Service.DecodeAccessRequest, Service.EncodeAccessResponse)

	r := mux.NewRouter()
	{
		r.Methods("POST").Path("/access-token").Handler(accessHandler)
		r.Methods("GET", "DELETE").Path(`/user/{uid:\d+}`).Handler(serverHandler)
		r.Methods("GET").Path("/health").HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
			writer.Header().Set("Content-type", "application/json")
			writer.Write([]byte(`{"stataus":"ok"}`))	//健康检查
		})
	}
	errChan := make(chan error)
	go func() {
		util.RegService()
		err := http.ListenAndServe("10.61.2.202:"+strconv.Itoa(*port), r)
		if err != nil{
			errChan <- err
		}
	}()
	var sig_c chan os.Signal
	go func() {
		sig_c = make(chan  os.Signal)
		signal.Notify(sig_c, syscall.SIGINT, syscall.SIGTERM)
		errChan <- fmt.Errorf("%s", <-sig_c)
	}()

	getErr := <- errChan
	fmt.Println(getErr)
	util.Unregservice()
}
