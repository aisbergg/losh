// Copyright 2022 André Lehmann
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//      http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package errors

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/aisbergg/go-errors/pkg/errors"
	"github.com/gookit/color"
)

// FormatColorfulCLIMessage formats an error message for CLI output.
func FormatColorfulCLIMessage(err error) string {
	// colors
	boldStyle := color.Bold
	boldRedStyle := color.Style{color.Bold, color.Red}
	boldBlueStyle := color.Style{color.Bold, color.Blue}

	// build the error message
	msgBuilder := &strings.Builder{}
	upkErr := errors.Unpack(err, true)
	msgBuilder.WriteString(boldRedStyle.Render("⚠️  Error Occurred ⚠️\n"))

	for _, upkElm := range upkErr {
		msgBuilder.WriteString(fmt.Sprintf("\n%s\n", boldRedStyle.Render(upkElm.Msg)))
		for _, frame := range upkElm.PartialStack {
			loc := boldStyle.Render(frame.File) + ":" + strconv.Itoa(frame.Line)
			msgBuilder.WriteString(fmt.Sprintf("at %s in %s\n", loc, boldBlueStyle.Render(frame.Name)))
		}
	}

	return msgBuilder.String()
}
