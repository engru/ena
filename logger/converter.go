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
	"github.com/sirupsen/logrus"
)

// Converter is interface for user-define extend
type Converter interface {
	Convert(*logrus.Entry) string
}

var (
	defaultConverterMap = make(map[string]Converter)
)

func init() {
	defaultConverterMap["d"] = &dateConverter{timestampFormat: "2006-01-02T15:04:05.000000000Z07:00"}
	defaultConverterMap["level"] = &levelConverter{}
	defaultConverterMap["p"] = &packageConverter{}
	defaultConverterMap["M"] = &methodConverter{}
	defaultConverterMap["L"] = &lineConverter{}
	defaultConverterMap["msg"] = &msgConverter{}
}

type textConverter struct {
	v string
}

func (c *textConverter) Convert(entry *logrus.Entry) string {
	return c.v
}

type dateConverter struct {
	timestampFormat string
}

func (c *dateConverter) Convert(entry *logrus.Entry) string {
	return entry.Time.Format(c.timestampFormat)
}

type levelConverter struct {
}

func (c *levelConverter) Convert(entry *logrus.Entry) string {
	return entry.Level.String()
}

type packageConverter struct {
}

func (c *packageConverter) Convert(entry *logrus.Entry) string {
	d, exists := entry.Data["logger.caller"]
	if !exists {
		return "-"
	}

	return d.(*caller).p
}

type methodConverter struct {
}

func (c *methodConverter) Convert(entry *logrus.Entry) string {
	d, exists := entry.Data["logger.caller"]
	if !exists {
		return "-"
	}

	return d.(*caller).m
}

type lineConverter struct {
}

func (c *lineConverter) Convert(entry *logrus.Entry) string {
	d, exists := entry.Data["logger.caller"]
	if !exists {
		return "-"
	}

	return d.(*caller).l
}

type msgConverter struct{}

func (c *msgConverter) Convert(entry *logrus.Entry) string {
	return entry.Message
}
