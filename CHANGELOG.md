## 0.16.0

BACKWARDS INCOMPATIBILITIES:

- Remove Value() from the parsley.Node interface{}, differentiate between literal and non-literal nodes

## 0.15.0

BACKWARDS INCOMPATIBILITIES:

- Replace `Type() string` on parsley.Node with `Schema() interface{}` to support any user-defined schema
- Allow to set schema for nil parser

IMPROVEMENTS:

- Use Go 1.16
- Move tools dependency file to a subpackage

OTHER:
- Remove CodeCov integration

## 0.14.0

BACKWARDS INCOMPATIBILITIES:

- Remove NL and Whitespaces parsers, these are achievable using the trim parsers
- Change Trim parsers to always skip all whitespaces and return specific whitespace error
- In text.LeftTrim prefer a not found parse error over a whitespace error if the position is the same

IMPROVEMENTS:

- Add struct wrapper for parser functions to allow recursive parser definitions
- Simplify parsley.Error, remove unnecessary msg field (and therefore unnecessary .Error()
- Introduce notfound errors when a parser doesn't match
- Only collect notfound errors in any/choice if they are not in the starting position

## 0.13.4

IMPROVEMENTS:

- Allow WsSpacesForceNl mode in RightTrim

## 0.13.3

IMPROVEMENTS:

- Make the parsley.StaticCheck method public

## 0.13.2

BUGFIXES:

- The Select interpreter should not run static checking on any child nodes, but simply return the index node's type

## 0.13.1

IMPROVEMENTS:

- The Select interpreter should static check the selected node (if the node is checkable)

## 0.13.0

BACKWARDS INCOMPATIBILITIES:

- The parsley.Evaluate function doesn't expect an evaluation context anymore, which was renamed to user context to make it more clear and can be set on parsley.Context
- Move ast.WalkNode to parsley.Walk
- The transform functions don't expect a NodeTransformer object anymore (the node transformation must be set on the interpreter)

IMPROVEMENTS:

 - The evaluation context is called user context now and can be set on parsley.Context
 - Proper recursive static checking (depth-first)
 - Proper recursive node transformation support (breadth-first)
 - The static checking and transformation are disabled by default and can be enabled on parsley.Context
 - Add NodeTransformerRegistry interface

## 0.12.3

IMPROVEMENTS:

- Add NodeTransformer to context to allow AST transformations
- The StaticChecker helper should not parse the input (minor backward incompatibility)

## 0.12.2

BUGFIXES:

- Do not panic in StaticCheck if there is no interpreter

## 0.12.1

IMPROVEMENTS:

- The text types should be constants (e.g. terminal.StringType)

## 0.12.0

IMPROVEMENTS:

- Add parsley.NonTerminalNode interface for nonterminal nodes
- Add parsley.StaticChecker and parsley.StaticCheckable interface for doing static analysis
- Add type to the AST nodes

BACKWARDS INCOMPATIBILITIES:

- The parsley.Interpreter's Eval now expects a parsley.NonTerminalNode instead of a list of nodes
- The ast.NodeList's node methods will run always on the first node only
- The ast.NilNode was renamed to ast.EmptyNode to avoid confusion
- The ast.NewTerminalNode now expects a new valueType parameter
- Most of the text/terminal parsers now will return a custom node type (string node, int node, etc.)

## 0.11.8

IMPROVEMENTS:

- Remove the index from the walk function, make it properly recursive

## 0.11.7

IMPROVEMENTS:

- Make the whole AST walkable, add Walk function to ast.NonTerminalNode

## 0.11.6

IMPROVEMENTS:

- Add keyword handling to parse context

## 0.11.5

IMPROVEMENTS:

- Add TimeDuration parser

## 0.11.4

IMPROVEMENTS:

- Change Sequence.Name() to return *combinator.Sequence
- Add Sequence.Token() method to set result token

## 0.11.3

IMPROVEMENTS:

- Add new whitespace mode: text.WsSpacesForceNl to force new lines

## 0.11.2

IMPROVEMENTS:

- Add text NL parser to match new lines

## 0.11.1

IMPROVEMENTS:

- Allow higher quality errors to be returned by the parsers

## 0.11.0

BACKWARDS INCOMPATIBILITIES:

- Rename parser's ReturnError method to Name
- Rename Seq combinator to SeqOf
- rename Recursive combinator to Seq (the name Recursive doesn't mean anything)

IMPROVEMENTS:

- Add Name() function to combinator.Sequence (same helper as ReturnError)

## 0.10.2

IMPROVEMENTS:

- Add SuppressError combinato

## 0.10.1

IMPROVEMENTS:

- Record possible errors in Context if a combinator returns a result for better error messages

## 0.10.0

BACKWARDS INCOMPATIBILITIES:

- Rename Substring parser to Op

## 0.9.0

BACKWARDS INCOMPATIBILITIES:

- The parser interface was changed to return an error. The error return value was previously removed but since I realised it was a bad design decision to save ~10% in benchmarks.
- The parser interface doesn't contain the Name() method anymore, it wasn't that useful.
- The error field was removed from the parsing context
- The parsing context now carries the file set
- The top level parsley.Parse and parsley.Evaluate now will return with a simple error type which already contains the error position (these are only convenience methods).
- I removed the *OrValue combinators as you can force a single child result with the Single combinator
- Clean up all combinators to only have parser parameters
- Add parser.ReturnError combinator to be able to override a parser's error (if the error's position is the same as the reader's)
- Add ReturnError() helper method to parser.Func

BUGFIXES:

- The string parser won't accept new-line characters (\r, \n) in double-quoted strings anymore

## 0.8.4

BUGFIXES:

- The Single combinator will now return with the original parser's name

## 0.8.3

IMPROVEMENTS:
- Add Go module definition
- Update CircleCI config to use Go 1.11

## 0.8.2

IMPROVEMENTS:
- Add Single combinator: Single will change the result of p if it returns with a non terminal node with only one child. In this case directly the child will be returned.

## 0.8.1

IMPROVEMENTS:

- Add back the SeqTryOrValue, SepByOrValue and SepByOrValue1 combinators as these are required

BUGFIXES:

- Many1 won't match for zero p matches

## 0.8.0

IMPROVEMENTS:

- the library got a huge performance increase by refactoring the String terminal parser. The issue was caused by string concatenation and the excessive use of strconv.UnquoteChar.

BACKWARDS INCOMPATIBILITIES:

- a new parsing context struct was introduced which carries the reader, the result cache, the call count and the best error
- the parser interface was changed to use the parsing context
- the parsley.History interface and the parser.History struct was replaced by parsley.ResultCache (which lives in the context now)
- the terminal.Integer parser will return an int64

## 0.7.0

BACKWARDS INCOMPATIBILITIES:

- extend parsley.Error to have a `Cause() error` method
- rename parsley.NewError() to parsley.NewErrorf()
- add parsley.NewError with error type as input
- parsley.WrapError will store the original cause but update the error message

## 0.6.0

BACKWARDS INCOMPATIBILITIES:

- major refactor of most of the API
- most of the interfaces were moved to the parsley package
- the position handling was rewritten similar to go's token.Pos
- the builders were completely removed
- the reader doesn't handle whitespaces anymore
- the parsers' return value were simplified
- nil node values are not allowed anymore
- most of the combinators API's were simplified to avoid repetition (like name + token)

IMPROVEMENTS:

- the reader became stateless
- most of the tests were rewritten using Ginkgo/Gomega
- common interpreters were added (array, object)
- whitespaces can be handled precisely with new parsers (text.LeftTrim/RightTrim/Trim)
- new error type with position (parsley.Pos)
- new empty node type
- new file and fileset types were introduced to support parsing multiple files better

BUGFIXES:

- the History was using the wrong key when checking left-recursion and wasn't curtailing properly

TODO:

- some of the old tests in the combinator package needs to be rewritten using Gingko (this means we miss a lot of test coverage)

## 0.5.0

BACKWARDS INCOMPATIBILITIES:

- text.NewReader now expects a filename parameter

IMPROVEMENTS:

- Windows-style line endings (\r\n) are automatically replaced to Unix-style line endings (\n) in the text reader.

## 0.4.0

BACKWARDS INCOMPATIBILITIES:

- Typo fix in some methods: peak -> peek

IMPROVEMENTS:

- Move precompiled whitespace regexp to a separate variable in text reader

OTHER:

- Fix example JSON parser + add comparison benchmarks against encoding/json

## 0.3.3

IMPROVEMENTS:

- allow nil pos in reader errors, replace {{err}} placeholder in WrapError
- add filename to position, add new file reader constructor

## 0.3.2

IMPROVEMENTS:

- reader.WrapError() falls back to the cause's error message if no error message was given and the cause is not a reader error

## 0.3.1

IMPROVEMENTS:

- Add SepByOrValue and SepByOrValue1 combinators which will return the value node if only the value parser is matched

## 0.3.0

BACKWARDS INCOMPATIBILITIES:

- Add history object to parser.Parse
- Move Memoize back to the combinators
- Move parsley package to a directory
- Change top-level Parse/Evaluate methods
- Remove parser.Stat, collect call statistics in history

IMPROVEMENTS:

- Add Sentence root parser

## 0.2.3

IMPROVEMENTS:

- reader.WrapError keeps the original error message if empty string is given as message

## 0.2.2

IMPROVEMENTS:

- Add parser.FuncFactory interface

## 0.2.1

BUGFIXES:

- the ast.Node mock wasn't regenerated

IMPROVEMENTS:

- Generate mock for parser.Parser

## 0.2.0

BACKWARDS INCOMPATIBILITIES:

- combinator.Memoize was removed. Use the Memoize method on the history object instead.
- parser.Error, parser.NewError and parser.WrapError was moved to reader
- Interpreter.Eval now returns a reader.Error instead of a general error
- Node.Value now returns a reader.Error instead of a general error
- Terminal.Value now returns a reader.Error instead of a general error
- NonTerminal.Value now returns a reader.Error instead of a general error

## 0.1.5

IMPROVEMENTS:

- Add cause to reader.Error, add parser.WrapError constructor

## 0.1.4

IMPROVEMENTS:

- All combinators will return with a parser.Func type (not the parser.Parser interface)
- Change internal sepBy to a factory
- Define dependencies with Go Dep
- Generate mocks for testing

## 0.1.3

DEPRECATED:

- combinator.Memoize will be removed in version 0.2. Use the Memoize method on the history object instead.

IMPROVEMENTS:

- You don't need to use meaningless names for memoization anymore, but you have to be careful to call Memoize
  for your parsers only once.

CHANGES:

- History.GetParserIndex was removed as only the original combinator.Memoize needed it

## 0.1.2

IMPROVEMENTS:

- Add copyright and license headers to all .go files
- Improve code coverage, add notes for (hopefully) impossible panics
- Add codecov.io integration

## 0.1.1

BUG FIXES:

- IsEOF() in text.Reader was not ignoring whitespaces when ignoreWhitespaces was set to true.

IMPROVEMENTS:

- Add Reset() method to History to allow to reuse the parsers for multiple inputs

## 0.1.0

First release
