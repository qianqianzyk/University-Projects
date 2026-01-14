package log

// Level 日志级别
type Level uint8

// 日志级别常量
const (
	LevelFatal  Level = 0
	LevelPanic  Level = 1
	LevelDpanic Level = 2
	LevelError  Level = 3
	LevelWarn   Level = 4
	LevelInfo   Level = 5
	LevelDebug  Level = 6
)
