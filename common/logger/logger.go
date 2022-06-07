package logger

func Info(args ...interface{}) {
	zapLogger.Info(args)
}

func Infof(template string, args ...interface{}) {
	zapLogger.Infof(template, args)
}

func Warn(args ...interface{}) {
	zapLogger.Warn(args)
}
func Warnf(template string, args ...interface{}) {
	zapLogger.Warnf(template, args)
}

func Error(args ...interface{}) {
	zapLogger.Error(args)
}
func Errorf(template string, args ...interface{}) {
	zapLogger.Errorf(template, args)
}
