package utils

import (
	"io"
	"log"
	"os"
)

var LOG_MAP map[string]*log.Logger

func LogInto() {

	// 创建不同的日志文件
	gatewayfile, err := os.OpenFile("gateway.log", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		log.Fatal("无法创建日志文件:", err)
	}
	//defer gatewayfile.Close()

	Servicefile, err := os.OpenFile("service.log", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		log.Fatal("无法创建日志文件:", err)
	}
	//defer Servicefile.Close()

	// 创建不同的日志输出对象
	Gaylog := log.New(io.MultiWriter(os.Stdout, gatewayfile), "[GatewayLog]", log.Ldate|log.Ltime)
	Syslog := log.New(io.MultiWriter(os.Stdout, Servicefile), "[ServiceLog]", log.Ldate|log.Ltime)
	// 将日志同时输出到控制台和日志文件
	// 初始化 map
	LOG_MAP = make(map[string]*log.Logger)
	LOG_MAP["Gaylog"] = Gaylog
	LOG_MAP["Syslog"] = Syslog
}

func logFmt(logname string, logtype string, uuid string, message string) {
	go func() {
		slog, ok := LOG_MAP[logname]
		if ok {
			slog.Println(" " + logtype + " [" + uuid + "] " + message)
		}
	}()
}

func GayINFO(uuid string, message string) {
	logFmt("Gaylog", "[INFO]", uuid, message)
}

func GayDEBUG(uuid string, message string) {
	logFmt("Gaylog", "[DEBUG]", uuid, message)
}

func GayERROR(uuid string, message string) {
	logFmt("Gaylog", "[ERROR]", uuid, message)
}

func SysINFO(uuid string, message string) {
	// // 设置日志前缀和日志标志
	// gatewaylog.SetPrefix("[INFO]")
	logFmt("Syslog", "[INFO]", uuid, message)
}

func SysDEBUG(uuid string, message string) {
	logFmt("Syslog", "[DEBUG]", uuid, message)
}

func SysERROR(uuid string, message string) {
	logFmt("Syslog", "[ERROR]", uuid, message)
}

// 注意log.Fatal会调用os.Exit(1)退出程序
func GayEXIT(uuid string, message string) {
	logFmt("Gaylog", "[EXIT]", uuid, message)
}

// 注意log.Fatal会调用os.Exit(1)退出程序
func SysEXIT(uuid string, message string) {
	logFmt("Syslog", "[EXIT]", uuid, message)
}
