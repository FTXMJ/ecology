package logs

import (
	"github.com/lestrrat-go/file-rotatelogs"
	"github.com/rifflock/lfshook"
	log "github.com/sirupsen/logrus"
	"os"
	"time"
)

var Log *log.Logger

func init() {
	var file *os.File
	file, er := os.Open("./logs/logtemp/log" + time.Now().Format("2006") + "-" + time.Now().Format("01") + "-" + time.Now().Format("02"))
	if er != nil && os.IsNotExist(er) {
		file, _ = os.Create("./logs/logtemp/log" + time.Now().Format("2006") + "-" + time.Now().Format("01") + "-" + time.Now().Format("02"))
		defer file.Close()
	}
	// 实例化
	Log = log.New()
	// 设置输出
	//Log.Out = file
	log.SetFormatter(&log.JSONFormatter{
		// PrettyPrint: true,//格式化json
		TimestampFormat: "2006-01-02 15:04:05", //时间格式化
	})
	// 设置日志级别
	Log.SetLevel(log.InfoLevel)
	Log.SetOutput(os.Stdout)
	Log.SetReportCaller(true)
	//logFilePath := config.Log_FILE_PATH
	//logFileName := config.LOG_FILE_NAME
	//// 日志文件
	//fileName := path.Join(logFilePath, logFileName)
	// 写入文件
	// 设置 rotatelogs
	logWriter, err := rotatelogs.New(
		// 分割后的文件名称
		"./logs/logtemp/log"+".%Y-%m-%d.log",
		// 生成软链，指向最新日志文件
		// 设置最大保存时间(7天)
		rotatelogs.WithMaxAge(7*24*time.Hour),
		// 设置日志切割时间间隔(1天)
		rotatelogs.WithRotationTime(24*time.Hour),
	)
	if err != nil {
		Log.Errorf("config local file system for logger error: %v", err)
	}
	writeMap := lfshook.WriterMap{
		log.InfoLevel:  logWriter,
		log.FatalLevel: logWriter,
		log.DebugLevel: logWriter,
		log.WarnLevel:  logWriter,
		log.ErrorLevel: logWriter,
		log.PanicLevel: logWriter,
	}
	lfHook := lfshook.NewHook(writeMap, &log.JSONFormatter{
		TimestampFormat: "2006-01-02 15:04:05",
	})
	// 新增 Hook
	Log.AddHook(lfHook)
}
