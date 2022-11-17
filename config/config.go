package config

import (
	"os"
	"path"
	"runtime"
	"strings"

	"github.com/tkanos/gonfig"
)

type Configuration struct {
	GCPProjectId      string
	GCPSubscriptionId string
}

var ConfigurationMap *Configuration

func InitConfig() {
	if ConfigurationMap == nil {
		ConfigurationMap = new(Configuration)
		err := gonfig.GetConf(getFileName(), ConfigurationMap)
		if err != nil {
			os.Exit(500)
		}
	}
}

func getFileName() string {
	fileName := strings.Join([]string{"config.", "json"}, "")
	_, dirPath, _, _ := runtime.Caller(0)
	dirName := path.Dir(dirPath)
	return path.Join(dirName, fileName)
}
