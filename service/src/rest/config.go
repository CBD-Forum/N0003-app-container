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
package rest

import (
	"github.com/op/go-logging"
	"github.com/spf13/viper"
)

var (
	restLogger = logging.MustGetLogger("rest")
 	fabricChaincodeName string
 	fabricPeerAddress string
 	localServerMaxConnectionLimit int
 	localServerAddress string
	iotApiUrl string
)


func InitConfiguration(){
	fabricChaincodeName = viper.GetString("fabric.chaincode.name")
	fabricPeerAddress = viper.GetString("fabric.peer.address")
	localServerAddress = viper.GetString("local.server.address")
	localServerMaxConnectionLimit = viper.GetInt("local.server.max_connection_limit")
	iotApiUrl = viper.GetString("iot.api.url")
	restLogger.Info("Init configuration in rest module - ")
	restLogger.Infof("	fabricChaincodeName 		%s", fabricChaincodeName)
	restLogger.Infof("	fabricPeerAddress   		%s", fabricPeerAddress)
	restLogger.Infof("	localServerAddress 		%s", localServerAddress)
	restLogger.Infof("	localServerMaxConnectionLimit 	%d", localServerMaxConnectionLimit)
	restLogger.Infof("	iotApiUrl			%s", iotApiUrl)
}
