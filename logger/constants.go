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

package logger

var (
	logrusPackageID = "github.com/sirupsen/logrus"

	convertKeyDate    = "d"
	convertKeyLevel   = "level"
	convertKeyPackage = "P"
	convertKeyMethod  = "M"
	convertKeyFile    = "F"
	convertKeyLine    = "L"
	convertKeyMsg     = "msg"

	tokenTypeConverter = "c"
	tokenTypeToken     = "t"

	loggerCallerKeyName = "ena.logger.caller"
)
