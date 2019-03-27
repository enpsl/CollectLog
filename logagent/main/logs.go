package main

import (
	"encoding/json"
	"fmt"
	"github.com/astaxie/beego/logs"
	"strings"
)

func convertLevel(level string) int {
	level = strings.ToLower(level)
	switch level {
	case "debug":
		return logs.LevelDebug
	case "warn":
		return logs.LevelWarn
	case "trace":
		return logs.LevelTrace
	case "notice":
		return logs.LevelNotice
	default:
		return logs.LevelDebug
	}
}

func initLogger() (err error) {
	config := make(map[string]interface{})
	config["filename"] = appConfig.logPath
	config["level"] = convertLevel(appConfig.logLevel)
	configStr, err := json.Marshal(config)
	if err != nil {
		fmt.Println("marshal failed, err:", err)
		return
	}
	logs.SetLogger(logs.AdapterFile, string(configStr))
	return
}
