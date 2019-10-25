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

package art

// tokenizerState is the phase of tokenizer
type tokenizerState = int

const (
	// literalState when parse the literal
	literalState = 1

	// formatModifierState when parse the format modifier
	formatModifierState = 2

	// keyworkState when parse the converter keywork
	keywordState = 3

	// optionState when parse the keywork option
	optionState = 4
)

const (
	// percentChar of %, is the begin of keywork
	percentChar = '%'

	// leftCurly of {, is the begin of option
	leftCurly = '{'

	// rightCurly of }, is the end of option
	rightCurly = '}'

	// escapeChar of \
	escapeChar = '\\'
)
