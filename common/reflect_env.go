// Copyright (c) 2018 soren yang
//
// Licensed under the MIT License
// you may not use this file except in complicance with the License.
// You may obtain a copy of the License at
//
//     https://opensource.org/licenses/MIT
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package common

import (
	"errors"
	"fmt"
	"os"
	"reflect"
	"regexp"
)

// Env will fill the value element with enviroment value
func Env(value interface{}) error {
	resultv := reflect.ValueOf(value)
	if resultv.Kind() != reflect.Ptr || (resultv.Elem().Kind() != reflect.Struct && resultv.Elem().Kind() != reflect.String) {
		return errors.New("value argument must be a struct/string address")
	}

	reg, err := regexp.Compile(`^\${(.+)}$`)
	if err != nil {
		return err
	}
	if resultv.Elem().Kind() == reflect.String {
		env := reg.FindStringSubmatch(resultv.Elem().String())
		if len(env) == 2 {
			fmt.Println(os.Getenv(env[1]))
			resultv.Elem().SetString(os.Getenv(env[1]))
		}
	}

	return nil
}
