package logs

import (
	"github.com/astaxie/beego/logs"
)

//"filename" : "searchApi.log", // 文件名
//"maxlines" : 1000,       // 最大行
//"maxsize"  : 10240,       // 最大Size
var jsonConfig = `{
        "filename" : "logtemp/ecology.log", 
        "maxlines" : 100000,        
        "maxsize"  : 10240,       
        "maxdays":7,
		 "color":true,
        "separate":["emergency", "alert", "critical", "error", "warning", "notice", "info", "debug"]
    }`

var Log *logs.BeeLogger

func init() {
	//logs.NewLogger()
	Log = logs.Async(1000) // 创建一个日志记录器，参数为缓冲区的大小
	Log.SetLogger(logs.AdapterMultiFile, jsonConfig)
	//beego.BeeLogger.DelLogger("console")    //不输出到控制台
	Log.SetLevel(logs.LevelDebug) // 设置日志写入缓冲区的等级
	Log.EnableFuncCallDepth(true) // 输出log时能显示输出文件名和行号（非必须）
	Log.SetLogFuncCallDepth(2)
	// Log.Async(1000)        // message 的长度，单位是字节，这里设置了1000
}

/*func init() {
	//Log = logs.NewLogger(10000) // 创建一个日志记录器，参数为缓冲区的大小
	//// 设置配置文件
	//jsonConfig := `{
    //    "filename" : "/logtemp/ecology.log", // 文件名
    //    "maxlines" : 1000,       // 最大行
    //    "maxsize"  : 10240       // 最大Size
    //}`
	//Log.SetLogger("file", jsonConfig) // 设置日志记录方式：本地文件记录
	//Log.SetLevel(logs.LevelDebug)     // 设置日志写入缓冲区的等级
	//Log.EnableFuncCallDepth(true)     // 输出log时能显示输出文件名和行号（非必须）
	//
	//Log.Emergency("Emergency")
	//Log.Alert("Alert")
	//Log.Critical("Critical")
	//Log.Error("Error")
	//Log.Warning("Warning")
	//Log.Notice("Notice")
	//Log.Informational("Informational")
	//Log.Debug("Debug")
	//
	//Log.Flush() // 将日志从缓冲区读出，写入到文件


	config := make(map[string]interface{})
	config["filename"] = "/logtemp/logcollect.log"
	config["level"] = logs.LevelDebug

	configStr, err := json.Marshal(config)
	if err != nil {
		fmt.Println("marshal failed,err:", err)
		return
	}
	logs.SetLogger(logs.AdapterFile, string(configStr))
	logs.Debug("this is a test,my name is %s", "stu01")
	logs.Trace("this is a trace,my name is %s", "stu02")
	logs.Warn("this is a warn,my name is %s", "stu03")
}*/
