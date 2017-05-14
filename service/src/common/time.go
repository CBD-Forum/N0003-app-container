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
	"time"
)

const (
	DATE_TIME_FORMAT string = "2006-01-02 15:04:05"
)

// GetCurrentTime returns current time as string
func GetCurrentTime() string {
	return time.Now().Format(DATE_TIME_FORMAT)
}

// AddTime add currentTime string with duration configured in the config.yml, returns the result time as string
func AddTime(currentTime string, duration time.Duration) string {
	tm, _ := time.Parse(DATE_TIME_FORMAT, currentTime)
	return tm.Add(duration).Format(DATE_TIME_FORMAT)


}

// IsTimeBefore compares timea and timeb string, returns true if timea is before timeb, else return false
func IsTimeBefore(timea, timeb string) bool {
	a, _ := time.Parse(DATE_TIME_FORMAT, timea)
	b, _ := time.Parse(DATE_TIME_FORMAT, timeb)
	return a.Before(b)
}
