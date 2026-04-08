package logger

import "go.uber.org/zap"

func NewProduction() (ILogger, error) {
	l, err := zap.NewProduction()
	if err != nil {
		return nil, err
	}

	return NewZapLogger(l), nil
}

//func NewDevelopment() (ILogger, error) {
//	l, err := zap.NewDevelopment()
//	if err != nil {
//		return nil, err
//	}
//
//	return NewZapLogger(l), nil
//}
