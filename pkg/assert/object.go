/* MIT License
 * Copyright (c) 2012-2020 Mat Ryer, Tyler Bunnell and contributors.
 *
 * Permission is hereby granted, free of charge, to any person obtaining a copy
 * of this software and associated documentation files (the "Software"), to deal
 * in the Software without restriction, including without limitation the rights
 * to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
 * copies of the Software, and to permit persons to whom the Software is
 * furnished to do so, subject to the following conditions:
 *
 * The above copyright notice and this permission notice shall be included in all
 * copies or substantial portions of the Software.
 *
 * THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
 * IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
 * FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
 * AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
 * LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
 * OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
 * SOFTWARE.
 */

/* NOTE: copied from https://github.com/stretchr/testify/assert/assertion_compare.go and modified
 *  The original versions of the files are MIT licensed
 */

package assert

import (
	"fmt"

	"github.com/wklken/httptest/pkg/util"

	"github.com/stretchr/testify/assert"
)

func Equal(actual interface{}, expected interface{}) (bool, string) {
	equal := assert.ObjectsAreEqual(actual, expected)
	if !equal {
		// TODO: truncate the middle, keep the begin and end
		actualStr := util.TruncateString(fmt.Sprintf("%v", actual), 100)

		actualValue := prettyLine(actualStr)

		// not equal, maybe is the type wrong
		return false, fmt.Sprintf("not equal, expected=`%v`(%T), actual=`%v`(%T)",
			expected, expected, actualValue, actual)
	} else {
		return true, "OK"
	}

}
