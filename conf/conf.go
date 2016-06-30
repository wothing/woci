/**
 * Copyright 2015-2016, Wothing Co., Ltd.
 * All rights reserved.
 *
 * Created by Elvizlai on 2016/04/23 10:44
 */

package conf

import (
	"encoding/json"
	"io/ioutil"
	"os"
	"fmt"
	"runtime"
	"strings"

	"github.com/pborman/uuid"
	"github.com/wothing/log"
)

type config struct {
	TRACER     string
	Concurrent int
	Env        map[string]interface{}
	Initial    []string
	Before     []string
	Modules    []Module
	Test       []string
	After      []string
}

type Module struct {
	Name  string
	Build string
	Start string
	Clean string
}

var Config = &config{Concurrent: runtime.NumCPU(), TRACER: uuid.New()[:8]}

func ParseConfig(configFile string) {
	data, err := ioutil.ReadFile(configFile)
	if err != nil {
		log.Fatal("read '%s' error: %v", configFile, err)
	}
	log.Debugf("load '%s' succeed", configFile)

	err = json.Unmarshal(data, Config)
	if err != nil {
		log.Fatalf("unmarshal '%s' error: %v", configFile, err)
	}

	dataStr := string(data)
	for i, j := 0, len(Config.Env); i < j; i++ {
		for k, v := range Config.Env {
			dataStr = strings.Replace(dataStr, "$"+k, fmt.Sprint(v), -1)
		}
	}

	json.Unmarshal([]byte(dataStr), Config)
}

func GenUUID() {
	Config.TRACER = uuid.New()[:8]
	err := ioutil.WriteFile("/TRACER", []byte(Config.TRACER), os.ModeTemporary)
	if err != nil {
		log.Fatal(err)
	}
}

func RestoreUUID() {
	data, err := ioutil.ReadFile("/TRACER")
	if err != nil {
		log.Fatal("can not do this, read file '/tracer' failed")
	}
	Config.TRACER = string(data)
}
