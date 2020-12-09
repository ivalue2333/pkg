package httpxroute

import (
	"errors"
	"reflect"
)

var (
	ErrMustValid = errors.New("method must be valid")
	ErrMustFunc  = errors.New("method must be func")

	ErrMustHasTwoParam = errors.New("method must has two param")
	ErrMustHasTwoOut   = errors.New("method must has two out")

	ErrMustPtr           = errors.New("param must be ptr")
	ErrMustPointToStruct = errors.New("param must point to struct")

	ErrMustError = errors.New("method ret must be error")

	replyErrorType = reflect.TypeOf((*error)(nil)).Elem()
)

func CheckMethod(method interface{}) (mV reflect.Value, reqT, respT reflect.Type, err error) {
	mV = reflect.ValueOf(method)
	if !mV.IsValid() {
		err = ErrMustValid
		return
	}

	mT := mV.Type()
	if mT.Kind() != reflect.Func {
		err = ErrMustFunc
		return
	}

	if mT.NumIn() != 2 {
		err = ErrMustHasTwoParam
		return
	}
	if mT.NumOut() != 2 {
		err = ErrMustHasTwoOut
		return
	}

	// in param
	reqT = mT.In(1)
	if reqT.Kind() != reflect.Ptr {
		err = ErrMustPtr
		return
	}
	if reqT.Elem().Kind() != reflect.Struct {
		err = ErrMustPointToStruct
		return
	}
	reqT = reqT.Elem()

	// return data
	respT = mT.Out(0)
	if respT.Kind() != reflect.Ptr {
		err = ErrMustPtr
		return
	}
	if !(respT.Elem().Kind() == reflect.Struct || respT.Elem().Kind() == reflect.Interface) {
		err = ErrMustPointToStruct
		return
	}

	retT := mT.Out(1)
	if retT != replyErrorType {
		err = ErrMustError
		return
	}

	return mV, reqT, respT, err
}
