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
	"fmt"
	"crypto/sha1"
	"crypto/md5"
	"time"
)


// SHA1string returns the hash of string s, using the SHA1 algorithm
func SHA1string(s string)(string){
	digestbyte := sha1.Sum([]byte(s))
	return fmt.Sprintf("%x", digestbyte)
}

// MD5string returns the hash of string s, using the MD5 algorithm
func MD5string(s string)(string){
	digestbyte := md5.Sum([]byte(s))
	return fmt.Sprintf("%x", digestbyte)
}

// ComputeSessionToken returns session token
func ComputeSessionToken(userid, sessionid, password string)string {
	return SHA1string(userid + sessionid + password)
}

// GenSessionExpireTime returns session expire time
func GenSessionExpireTime() string {
	return AddTime(GetCurrentTime(), localServerSessionDuration)
}


// GetLocalSessionDuration returns the session duration
func GetLocalSessionDuration() time.Duration {
	return localServerSessionDuration
}
