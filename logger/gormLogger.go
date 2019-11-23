package logger

func NewGormLogger(logger Logger) GormLoggerImpl {
	return GormLoggerImpl{log: logger}
}

type GormLoggerImpl struct {
	log Logger
}

func (l GormLoggerImpl) Print(v ...interface{}) {
	l.log.Debugf("", v)
}
