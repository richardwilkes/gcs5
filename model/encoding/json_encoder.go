/*
 * Copyright ©1998-2022 by Richard A. Wilkes. All rights reserved.
 *
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, version 2.0. If a copy of the MPL was not distributed with
 * this file, You can obtain one at http://mozilla.org/MPL/2.0/.
 *
 * This Source Code Form is "Incompatible With Secondary Licenses", as
 * defined by the Mozilla Public License, version 2.0.
 */

package encoding

import (
	"bufio"
	"io"
	"os"
	"path/filepath"
	"strings"
	"unicode/utf8"

	"github.com/richardwilkes/toolbox/errs"
	"github.com/richardwilkes/toolbox/xio/fs/safe"
	"github.com/richardwilkes/toolbox/xmath/fixed"
)

// JSONEncoder holds a simple JSON encoder.
type JSONEncoder struct {
	writer     *bufio.Writer
	indent     string
	depth      int
	needComma  bool
	needIndent bool
}

// SaveJSON writes JSON-encoded data to the file path, creating any intermediate directories required. The function 'f'
// is expected to produce the actual data. If 'format' is true, the JSON will be pretty-formatted for human consumption.
func SaveJSON(filePath string, format bool, f func(*JSONEncoder)) error {
	if err := os.MkdirAll(filepath.Dir(filePath), 0o750); err != nil {
		return errs.NewWithCause(filePath, err)
	}
	if err := safe.WriteFileWithMode(filePath, func(w io.Writer) error {
		indent := ""
		if format {
			indent = "  "
		}
		encoder := &JSONEncoder{
			writer: bufio.NewWriter(w),
			indent: indent,
		}
		f(encoder)
		return encoder.Close()
	}, 0o640); err != nil {
		return errs.NewWithCause(filePath, err)
	}
	return nil
}

// StartObject starts a JSON object.
func (w *JSONEncoder) StartObject() {
	w.start('{')
}

// EndObject ends a JSON object.
func (w *JSONEncoder) EndObject() {
	w.end('}')
}

// StartArray starts a JSON array.
func (w *JSONEncoder) StartArray() {
	w.start('[')
}

// EndArray ends a JSON array.
func (w *JSONEncoder) EndArray() {
	w.end(']')
}

// Key emits the key.
func (w *JSONEncoder) Key(key string) {
	if w.needComma {
		w.doComma()
	} else if w.indent != "" {
		w.doIndent()
	}
	w.writeQuotedString(key)
	w.writer.WriteByte(':')
	if w.indent != "" {
		w.writer.WriteByte(' ')
	}
}

// Bool emits a boolean value.
func (w *JSONEncoder) Bool(value bool) {
	w.commaIfNeeded()
	if value {
		w.writer.WriteString("true")
	} else {
		w.writer.WriteString("false")
	}
}

// KeyedBool emits a key and boolean value. If 'omitFalse' is true, then only true values will be emitted.
func (w *JSONEncoder) KeyedBool(key string, value, omitFalse bool) {
	if omitFalse && !value {
		return
	}
	w.Key(key)
	if value {
		w.writer.WriteString("true")
	} else {
		w.writer.WriteString("false")
	}
	w.needComma = true
}

// Number emits a numeric value.
func (w *JSONEncoder) Number(value fixed.F64d4) {
	w.commaIfNeeded()
	w.writer.WriteString(value.String())
}

// KeyedNumber emits a key and numeric value.
func (w *JSONEncoder) KeyedNumber(key string, value fixed.F64d4, omitZero bool) {
	if omitZero && value == 0 {
		return
	}
	w.Key(key)
	w.writer.WriteString(value.String())
	w.needComma = true
}

// String emits a string value.
func (w *JSONEncoder) String(value string) {
	w.commaIfNeeded()
	w.writeQuotedString(value)
}

// KeyedString emits a key and string value. If 'omitEmpty' is true, then the value will only be emitted if it isn't
// empty. If 'trimFirst' is also true, then the value will be run through strings.TrimSpace() first.
func (w *JSONEncoder) KeyedString(key, value string, omitEmpty, trimFirst bool) {
	if omitEmpty {
		if trimFirst {
			value = strings.TrimSpace(value)
		}
		if value == "" {
			return
		}
	}
	w.Key(key)
	w.writeQuotedString(value)
	w.needComma = true
}

// Close the encoder, flushing any remaining data.
func (w *JSONEncoder) Close() error {
	if w.indent != "" {
		w.writer.WriteByte('\n')
	}
	return errs.Wrap(w.writer.Flush())
}

func (w *JSONEncoder) start(ch byte) {
	if w.needComma {
		w.doComma()
	} else if w.needIndent {
		w.doIndent()
	}
	w.writer.WriteByte(ch)
	if w.indent != "" {
		w.writer.WriteByte('\n')
		w.depth++
		w.needIndent = true
	}
}

func (w *JSONEncoder) end(ch byte) {
	if w.indent != "" {
		w.writer.WriteByte('\n')
		w.depth--
		w.doIndent()
	}
	w.writer.WriteByte(ch)
	w.needComma = true
}

func (w *JSONEncoder) doIndent() {
	for i := 0; i < w.depth; i++ {
		w.writer.WriteString(w.indent)
	}
	w.needIndent = false
}

func (w *JSONEncoder) doComma() {
	w.needComma = false
	w.writer.WriteByte(',')
	if w.indent != "" {
		w.writer.WriteByte('\n')
		w.doIndent()
	}
}

func (w *JSONEncoder) commaIfNeeded() {
	if w.needComma {
		w.writer.WriteByte(',')
		if w.indent != "" {
			w.writer.WriteByte('\n')
			w.doIndent()
		}
	} else {
		if w.indent != "" {
			w.doIndent()
		}
		w.needComma = true
	}
}

func (w *JSONEncoder) writeQuotedString(s string) {
	w.writer.WriteByte('"')
	for _, ch := range s {
		if ch < 32 {
			switch ch {
			case '\\':
				w.writer.WriteString(`\\`)
			case '"':
				w.writer.WriteString(`\"`)
			case '\b':
				w.writer.WriteString(`\b`)
			case '\t':
				w.writer.WriteString(`\t`)
			case '\n':
				w.writer.WriteString(`\n`)
			case '\r':
				w.writer.WriteString(`\r`)
			case '\f':
				w.writer.WriteString(`\f`)
			default:
				w.writer.WriteString(`\u00`)
				const hex = "0123456789abcdef"
				w.writer.WriteByte(hex[ch>>4])
				w.writer.WriteByte(hex[ch&0xf])
			}
		} else {
			switch ch {
			case '&':
				w.writer.WriteString(`\u0026`) // escape for javascript
			case '<':
				w.writer.WriteString(`\u003c`) // escape for javascript
			case '>':
				w.writer.WriteString(`\u003e`) // escape for javascript
			case '\u2028': // Line separator; escape for javascript
				w.writer.WriteString(`\u2028`)
			case '\u2029': // Paragraph separator; escape for javascript
				w.writer.WriteString(`\u2029`)
			case utf8.RuneError:
				w.writer.WriteString(`\ufffd`)
			default:
				w.writer.WriteRune(ch)
			}
		}
	}
	w.writer.WriteByte('"')
}
