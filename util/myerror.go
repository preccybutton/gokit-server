package util

//自定义error生成器
type MyError struct {
	Code int
	Message string
}

func NewMyError(code int, message string) error{
	return &MyError{Code: code, Message: message}
}

func (this *MyError) Error() string{
	return this.Message
}
