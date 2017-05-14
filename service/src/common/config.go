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
package common

import (
	"github.com/spf13/viper"
	"github.com/op/go-logging"
	"time"
)

var (
	logger *logging.Logger = logging.MustGetLogger("common")
	databaseName string
	databaseDsn string
	localServerSessionDuration time.Duration
)

func InitConfig(){
	databaseName = viper.GetString("database.name")
	databaseDsn = viper.GetString("database.dsn")
	localServerSessionDuration = viper.GetDuration("local.server.session.duration")
	logger.Info("Init configuration in common module - ")
	logger.Infof("	databaseName 			%s", databaseName)
	logger.Infof("	databaseDsn   			%s", databaseDsn)
	logger.Infof("	localServerSessionDuration 	%v", localServerSessionDuration)
}
