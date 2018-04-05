// Copyright (c) 2017 Opsidian Ltd.
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

package terminal_test

import (
	"fmt"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/ginkgo/extensions/table"
	. "github.com/onsi/gomega"
	"github.com/opsidian/parsley/ast"
	"github.com/opsidian/parsley/data"
	"github.com/opsidian/parsley/parsley"
	"github.com/opsidian/parsley/text"
	"github.com/opsidian/parsley/text/terminal"
)

var _ = Describe("String", func() {

	Context("when backquotes are allowed", func() {

		var p = terminal.String(true)

		DescribeTable("should match",
			func(input string, startPos int, value interface{}, nodePos parsley.Pos, endPos int) {
				r := text.NewReader(text.NewFile("textfile", []byte(input)))
				curtailingParsers, res, err := p.Parse(nil, data.EmptyIntMap, r, startPos)
				Expect(curtailingParsers).To(Equal(data.EmptyIntSet))
				Expect(err).ToNot(HaveOccurred())
				node := res[0].(*ast.TerminalNode)
				Expect(node.Token()).To(Equal("STRING"))
				Expect(node.Value(nil)).To(Equal(value))
				Expect(node.Pos()).To(Equal(nodePos))
				Expect(node.ReaderPos()).To(Equal(endPos))
			},
			Entry(`"foo" beginning`, `"foo"`, 0, "foo", parsley.Pos(1), 5),
			Entry(`"foo" middle`, `--- "foo" ---`, 4, "foo", parsley.Pos(5), 9),
			Entry(`"foo" end`, `--- "foo"`, 4, "foo", parsley.Pos(5), 9),
			Entry(`""`, `""`, 0, "", parsley.Pos(1), 2),
			Entry(`"a"`, `"a"`, 0, "a", parsley.Pos(1), 3),
			Entry(`" a "`, `" a "`, 0, " a ", parsley.Pos(1), 5), // Should keep whitespaces in the string
			Entry(`"a" "b"`, `"a" "b"`, 0, "a", parsley.Pos(1), 3),
			Entry(`"abcd"`, `"abcd"`, 0, "abcd", parsley.Pos(1), 6),
			Entry(`"'"`, `"'"`, 0, "'", parsley.Pos(1), 3),
			Entry(`"\a\b\f\n\r\t\v"`, `"\a\b\f\n\r\t\v"`, 0, "\a\b\f\n\r\t\v", parsley.Pos(1), 16),
			Entry(`"\x67"`, `"\x67"`, 0, "\x67", parsley.Pos(1), 6),
			Entry(`"\uAB12"`, `"\uAB12"`, 0, "\uAB12", parsley.Pos(1), 8),
			Entry(`"\U0001F355"`, `"\U0001F355"`, 0, "\U0001F355", parsley.Pos(1), 12),
			Entry("``", "``", 0, ``, parsley.Pos(1), 2),
			Entry("`a`", "`a`", 0, `a`, parsley.Pos(1), 3),
			Entry("` a `", "` a `", 0, ` a `, parsley.Pos(1), 5), // Should keep whitespaces in the string
			Entry("`a` `b`", "`a` `b`", 0, `a`, parsley.Pos(1), 3),
			Entry("`abcd`", "`abcd`", 0, `abcd`, parsley.Pos(1), 6),
			Entry("`'`", "`'`", 0, `'`, parsley.Pos(1), 3),
			Entry("`"+`\a\b\f\n\r\t\v`+"`", "`"+`\a\b\f\n\r\t\v`+"`", 0, `\a\b\f\n\r\t\v`, parsley.Pos(1), 16),
			Entry("`"+`\x67`+"`", "`"+`\x67`+"`", 0, `\x67`, parsley.Pos(1), 6),
			Entry("`"+`\uAB12`+"`", "`"+`\uAB12`+"`", 0, `\uAB12`, parsley.Pos(1), 8),
			Entry("`"+`\U0001F355`+"`", "`"+`\U0001F355`+"`", 0, `\U0001F355`, parsley.Pos(1), 12),
		)

		DescribeTable("should not match",
			func(input string, startPos int, errPos parsley.Pos) {
				r := text.NewReader(text.NewFile("textfile", []byte(input)))
				curtailingParsers, res, err := p.Parse(nil, data.EmptyIntMap, r, startPos)
				Expect(curtailingParsers).To(Equal(data.EmptyIntSet))
				Expect(err).To(MatchError("was expecting string literal"))
				Expect(err.Pos()).To(Equal(errPos))
				Expect(res).To(BeNil())
			},
			Entry("empty", ``, 0, parsley.Pos(1)),
			Entry("pos test", `--- x`, 4, parsley.Pos(5)),
			Entry(`a`, `a`, 0, parsley.Pos(1)),
			Entry(`''`, `''`, 0, parsley.Pos(1)),
			Entry(`'a'`, `'a'`, 0, parsley.Pos(1)),
		)

		DescribeTable("unfinished string literal",
			func(input string) {
				r := text.NewReader(text.NewFile("textfile", []byte(input)))
				curtailingParsers, res, err := p.Parse(nil, data.EmptyIntMap, r, 0)
				Expect(curtailingParsers).To(Equal(data.EmptyIntSet))
				Expect(err).To(MatchError(fmt.Sprintf("was expecting '%s'", string(input[0]))))
				Expect(err.Pos()).To(Equal(parsley.Pos(5)))
				Expect(res).To(BeNil())
			},
			Entry("`foo", "`foo"),
			Entry(`"foo`, `"foo`),
		)
	})

	Context("when backquotes are not allowed", func() {

		var p = terminal.String(false)

		It("should match double-quoted strings", func() {
			r := text.NewReader(text.NewFile("textfile", []byte(`"foo"`)))
			_, res, err := p.Parse(nil, data.EmptyIntMap, r, 0)
			Expect(err).ToNot(HaveOccurred())
			Expect(res).ToNot(BeEmpty())
		})

		It("should not match backquoted strings", func() {
			r := text.NewReader(text.NewFile("textfile", []byte("`foo`")))
			_, res, err := p.Parse(nil, data.EmptyIntMap, r, 0)
			Expect(err).To(HaveOccurred())
			Expect(res).To(BeEmpty())
		})
	})

})
