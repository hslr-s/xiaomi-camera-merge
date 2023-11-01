package log

import (
	"log"
	"os"
)

type Log struct {
	Path string
}

func NewLog(path string) *Log {
	return &Log{
		Path: path,
	}
}

func (l *Log) WriteContent(content ...any) error {
	// 创建或打开一个日志文件
	logFile, err := os.OpenFile(l.Path, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		return err
	}
	defer logFile.Close()

	// 设置日志输出目标为文件
	log.SetOutput(logFile)

	// 输出日志消息
	// log.Println("这是一条日志消息")
	// log.Printf("这是一条带格式的日志消息：%s", "参数")

	// 如果需要，你也可以设置日志的格式和前缀
	log.SetFlags(log.Ldate | log.Ltime) // 设置日期和时间
	// log.SetPrefix("[MyApp] ")           // 设置前缀

	log.Println(content...)
	return nil
}
