// Copyright [2016] [Cuiting Shi ]
// 
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
// 
// http://www.apache.org/licenses/LICENSE-2.0
// 
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.
// 
package main

import (
	"fmt"
	"strings"
	"github.com/op/go-logging"
	"cbdforum/app-container/service/src/rest"
	"github.com/spf13/viper"
	"cbdforum/app-container/service/src/common"
)

var logger *logging.Logger = logging.MustGetLogger("Console")

const (
	envRoot string = "container"
	configFileName string = "configuration"
	configFilePath string = "$GOPATH/src/cbdforum/app-container/service/src"
)


func main() {
	SetUpConfiguration(envRoot)
	rest.StartContainerRESTServer()
}


func SetUpConfiguration(envRoot string){
	// For environment variables
	viper.SetEnvPrefix(envRoot)
	viper.AutomaticEnv()
	replacer := strings.NewReplacer(".", "_")
	viper.SetEnvKeyReplacer(replacer)

	// Set the configuration file
	viper.SetConfigName(configFileName)
	viper.AddConfigPath(configFilePath)
	viper.AddConfigPath(".")

	if err := viper.ReadInConfig(); err != nil {
		panic(fmt.Errorf("Fatal Error when reading config file %s in %s: %s\n", configFileName, configFilePath, err))
	}

	common.InitConfig()
	rest.InitConfiguration()
}

