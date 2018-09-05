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

import (
	"fmt"
	"path"
	"runtime"
	"strings"

	"github.com/sirupsen/logrus"
)

type callerHook struct {
}

func (callerHook) Levels() []logrus.Level {
	return logrus.AllLevels
}

type caller struct {
	p string // package name
	m string // function name
	l string // line
}

func (callerHook) Fire(entry *logrus.Entry) error {
	c := &caller{
		p: "-",
		m: "-",
		l: "-",
	}

	pc := make([]uintptr, 3, 3)
	cnt := runtime.Callers(7, pc)
	found := false
	for i := 0; i < cnt; i++ {
		fu := runtime.FuncForPC(pc[i] - 1)
		name := fu.Name()
		if !strings.Contains(name, "github.com/sirupsen/logrus") && !found {
			found = true
			continue
		}
		if found {
			file, line := fu.FileLine(pc[i] - 1)
			c.p = path.Base(file)
			c.m = path.Base(name)
			c.l = fmt.Sprintf("%v", line)
			break
		}
	}

	entry.Data["logger.caller"] = c
	return nil
}
