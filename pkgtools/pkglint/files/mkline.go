package pkglint

import (
	"fmt"
	"netbsd.org/pkglint/regex"
	"netbsd.org/pkglint/textproc"
	"path"
	"strings"
)

// MkLine is a line from a Makefile fragment.
// There are several types of lines.
// The most common types in pkgsrc are variable assignments,
// shell commands and directives like .if and .for.
type MkLine = *MkLineImpl

type MkLineImpl struct {
	Line
	data interface{} // One of the following mkLine* types
}
type mkLineAssign = *mkLineAssignImpl // See https://github.com/golang/go/issues/28045
type mkLineAssignImpl struct {
	commented         bool   // Whether the whole variable assignment is commented out
	varname           string // e.g. "HOMEPAGE", "SUBST_SED.perl"
	varcanon          string // e.g. "HOMEPAGE", "SUBST_SED.*"
	varparam          string // e.g. "", "perl"
	spaceAfterVarname string
	op                MkOperator //
	valueAlign        string     // The text up to and including the assignment operator, e.g. VARNAME+=\t
	value             string     // The trimmed value
	valueMk           []*MkToken // The value, sent through splitIntoMkWords
	valueMkRest       string     // nonempty in case of parse errors
	fields            []string   // The value, space-separated according to shell quoting rules
	spaceAfterValue   string
	comment           string
}
type mkLineShell struct {
	command string
}
type mkLineComment struct{} // See mkLineAssignImpl.commented for another type of comment line
type mkLineEmpty struct{}
type mkLineDirective = *mkLineDirectiveImpl // See https://github.com/golang/go/issues/28045
type mkLineDirectiveImpl struct {
	indent    string // the space between the leading "." and the directive
	directive string // "if", "else", "for", etc.
	args      string
	comment   string   // mainly interesting for .endif and .endfor
	elseLine  MkLine   // for .if (filled in later)
	cond      MkCond   // for .if and .elif (filled in on first access)
	fields    []string // the arguments for the .for loop (filled in on first access)
}
type mkLineInclude = *mkLineIncludeImpl // See https://github.com/golang/go/issues/28045
type mkLineIncludeImpl struct {
	mustExist       bool     // for .sinclude, nonexistent files are ignored
	sys             bool     // whether the include uses <file.mk> (very rare) instead of "file.mk"
	indent          string   // the space between the leading "." and the directive
	includedFile    string   // the text between the <brackets> or "quotes"
	conditionalVars []string // variables on which this inclusion depends (filled in later, as needed)
}
type mkLineDependency struct {
	targets string
	sources string
}

type MkLineParser struct{}

// Parse parses the text of a Makefile line to see what kind of line
// it is: variable assignment, include, comment, etc.
//
// See devel/bmake/parse.c:/^Parse_File/
func (p MkLineParser) Parse(line Line) *MkLineImpl {
	text := line.Text

	// XXX: This check should be moved somewhere else. NewMkLine should only be concerned with parsing.
	if hasPrefix(text, " ") && line.Basename != "bsd.buildlink3.mk" {
		line.Warnf("Makefile lines should not start with space characters.")
		line.Explain(
			"If this line should be a shell command connected to a target, use a tab character for indentation.",
			"Otherwise remove the leading whitespace.")
	}

	data := p.split(line, text)

	// Check for shell commands first because these cannot have comments
	// at the end of the line.
	if hasPrefix(text, "\t") {
		return p.parseShellcmd(line)
	}

	if mkline := p.parseVarassign(line, data); mkline != nil {
		return mkline
	}
	if mkline := p.parseCommentOrEmpty(line); mkline != nil {
		return mkline
	}
	if mkline := p.parseDirective(line, data); mkline != nil {
		return mkline
	}
	if mkline := p.parseInclude(line); mkline != nil {
		return mkline
	}
	if mkline := p.parseSysinclude(line); mkline != nil {
		return mkline
	}
	if mkline := p.parseDependency(line); mkline != nil {
		return mkline
	}
	if mkline := p.parseMergeConflict(line); mkline != nil {
		return mkline
	}

	// The %q is deliberate here since it shows possible strange characters.
	line.Errorf("Unknown Makefile line format: %q.", text)
	return &MkLineImpl{line, nil}
}

func (p MkLineParser) parseVarassign(line Line, data mkLineSplitResult) MkLine {
	m, a := p.MatchVarassign(line, line.Text, data)
	if !m {
		return nil
	}

	if a.spaceAfterVarname != "" {
		varname := a.varname
		op := a.op
		switch {
		case hasSuffix(varname, "+") && (op == opAssign || op == opAssignAppend):
			break
		case matches(varname, `^[a-z]`) && op == opAssignEval:
			break
		default:
			fix := line.Autofix()
			fix.Notef("Unnecessary space after variable name %q.", varname)
			fix.Replace(varname+a.spaceAfterVarname+op.String(), varname+op.String())
			fix.Apply()
		}
	}

	if a.comment != "" && a.value != "" && a.spaceAfterValue == "" {
		line.Warnf("The # character starts a Makefile comment.")
		line.Explain(
			"In a variable assignment, an unescaped # starts a comment that",
			"continues until the end of the line.",
			"To escape the #, write \\#.",
			"",
			"If this # character intentionally starts a comment,",
			"it should be preceded by a space in order to make it more visible.")
	}

	return &MkLineImpl{line, a}
}

func (p MkLineParser) parseShellcmd(line Line) MkLine {
	return &MkLineImpl{line, mkLineShell{line.Text[1:]}}
}

func (p MkLineParser) parseCommentOrEmpty(line Line) MkLine {
	trimmedText := trimHspace(line.Text)

	if strings.HasPrefix(trimmedText, "#") {
		return &MkLineImpl{line, mkLineComment{}}
	}

	if trimmedText == "" {
		return &MkLineImpl{line, mkLineEmpty{}}
	}

	return nil
}

func (p MkLineParser) parseInclude(line Line) MkLine {
	m, indent, directive, includedFile := MatchMkInclude(line.Text)
	if !m {
		return nil
	}

	return &MkLineImpl{line, &mkLineIncludeImpl{directive == "include", false, indent, includedFile, nil}}
}

func (p MkLineParser) parseSysinclude(line Line) MkLine {
	m, indent, directive, includedFile := match3(line.Text, `^\.([\t ]*)(s?include)[\t ]+<([^>]+)>[\t ]*(?:#.*)?$`)
	if !m {
		return nil
	}

	return &MkLineImpl{line, &mkLineIncludeImpl{directive == "include", true, indent, includedFile, nil}}
}

func (p MkLineParser) parseDependency(line Line) MkLine {
	// XXX: Replace this regular expression with proper parsing.
	// There might be a ${VAR:M*.c} in these variables, which the below regular expression cannot handle.
	m, targets, whitespace, sources := match3(line.Text, `^([^\t :]+(?:[\t ]*[^\t :]+)*)([\t ]*):[\t ]*([^#]*?)(?:[\t ]*#.*)?$`)
	if !m {
		return nil
	}

	if whitespace != "" {
		line.Notef("Space before colon in dependency line.")
	}
	return &MkLineImpl{line, mkLineDependency{targets, sources}}
}

func (p MkLineParser) parseMergeConflict(line Line) MkLine {
	if !matches(line.Text, `^(<<<<<<<|=======|>>>>>>>)`) {
		return nil
	}

	return &MkLineImpl{line, nil}
}

// String returns the filename and line numbers.
func (mkline *MkLineImpl) String() string {
	return sprintf("%s:%s", mkline.Filename, mkline.Linenos())
}

// IsVarassign returns true for variable assignments of the form VAR=value.
//
// See IsCommentedVarassign.
func (mkline *MkLineImpl) IsVarassign() bool {
	data, ok := mkline.data.(mkLineAssign)
	return ok && !data.commented
}

// IsCommentedVarassign returns true for commented-out variable assignments.
// In most cases these are treated as ordinary comments, but in some others
// they are treated like variable assignments, just inactive ones.
func (mkline *MkLineImpl) IsCommentedVarassign() bool {
	data, ok := mkline.data.(mkLineAssign)
	return ok && data.commented
}

// IsShellCommand returns true for tab-indented lines that are assigned to a Make
// target. Example:
//
//  pre-configure:    # IsDependency
//          ${ECHO}   # IsShellCommand
func (mkline *MkLineImpl) IsShellCommand() bool {
	_, ok := mkline.data.(mkLineShell)
	return ok
}

// IsComment returns true for lines that consist entirely of a comment.
func (mkline *MkLineImpl) IsComment() bool {
	_, ok := mkline.data.(mkLineComment)
	return ok || mkline.IsCommentedVarassign()
}

func (mkline *MkLineImpl) IsEmpty() bool {
	_, ok := mkline.data.(mkLineEmpty)
	return ok
}

// IsDirective returns true for conditionals (.if/.elif/.else/.if) or loops (.for/.endfor).
//
// See IsInclude.
func (mkline *MkLineImpl) IsDirective() bool {
	_, ok := mkline.data.(mkLineDirective)
	return ok
}

// IsInclude returns true for lines like: .include "other.mk"
//
// See IsSysinclude for lines like: .include <sys.mk>
func (mkline *MkLineImpl) IsInclude() bool {
	incl, ok := mkline.data.(mkLineInclude)
	return ok && !incl.sys
}

// IsSysinclude returns true for lines like: .include <sys.mk>
//
// See IsInclude for lines like: .include "other.mk"
func (mkline *MkLineImpl) IsSysinclude() bool {
	incl, ok := mkline.data.(mkLineInclude)
	return ok && incl.sys
}

// IsDependency returns true for dependency lines like "target: source".
func (mkline *MkLineImpl) IsDependency() bool {
	_, ok := mkline.data.(mkLineDependency)
	return ok
}

// Varname applies to variable assignments and returns the name
// of the variable that is assigned or appended to.
//
// Example:
//  VARNAME.${param}?=      value   # Varname is "VARNAME.${param}"
func (mkline *MkLineImpl) Varname() string { return mkline.data.(mkLineAssign).varname }

// Varcanon applies to variable assignments and returns the canonicalized variable name for parameterized variables.
// Examples:
//  HOMEPAGE           => "HOMEPAGE"
//  SUBST_SED.anything => "SUBST_SED.*"
//  SUBST_SED.${param} => "SUBST_SED.*"
func (mkline *MkLineImpl) Varcanon() string { return mkline.data.(mkLineAssign).varcanon }

// Varparam applies to variable assignments and returns the parameter for parameterized variables.
// Examples:
//  HOMEPAGE           => ""
//  SUBST_SED.anything => "anything"
//  SUBST_SED.${param} => "${param}"
func (mkline *MkLineImpl) Varparam() string { return mkline.data.(mkLineAssign).varparam }

// Op applies to variable assignments and returns the assignment operator.
func (mkline *MkLineImpl) Op() MkOperator { return mkline.data.(mkLineAssign).op }

// ValueAlign applies to variable assignments and returns all the text
// before the variable value, e.g. "VARNAME+=\t".
func (mkline *MkLineImpl) ValueAlign() string { return mkline.data.(mkLineAssign).valueAlign }

func (mkline *MkLineImpl) Value() string { return mkline.data.(mkLineAssign).value }

// VarassignComment applies to variable assignments and returns the comment.
//
// Example:
//  VAR=value # comment
//
// In the above line, the comment is "# comment".
//
// The leading "#" is included so that pkglint can distinguish between no comment at all and an empty comment.
func (mkline *MkLineImpl) VarassignComment() string { return mkline.data.(mkLineAssign).comment }

func (mkline *MkLineImpl) ShellCommand() string { return mkline.data.(mkLineShell).command }

func (mkline *MkLineImpl) Indent() string {
	if mkline.IsDirective() {
		return mkline.data.(mkLineDirective).indent
	} else {
		return mkline.data.(mkLineInclude).indent
	}
}

// Directive returns the preprocessing directive, like "if", "for", "endfor", etc.
//
// See matchMkDirective.
func (mkline *MkLineImpl) Directive() string { return mkline.data.(mkLineDirective).directive }

// Args returns the arguments from an .if, .ifdef, .ifndef, .elif, .for, .undef.
func (mkline *MkLineImpl) Args() string { return mkline.data.(mkLineDirective).args }

// Cond applies to an .if or .elif line and returns the parsed condition.
//
// If a parse error occurs, it is silently swallowed, returning a
// best-effort part of the condition, or even nil.
func (mkline *MkLineImpl) Cond() MkCond {
	cond := mkline.data.(mkLineDirective).cond
	if cond == nil {
		cond = NewMkParser(mkline.Line, mkline.Args(), true).MkCond()
		mkline.data.(mkLineDirective).cond = cond
	}
	return cond
}

// DirectiveComment is the trailing end-of-line comment, typically at a deeply nested .endif or .endfor.
func (mkline *MkLineImpl) DirectiveComment() string { return mkline.data.(mkLineDirective).comment }

func (mkline *MkLineImpl) HasElseBranch() bool { return mkline.data.(mkLineDirective).elseLine != nil }

func (mkline *MkLineImpl) SetHasElseBranch(elseLine MkLine) {
	data := mkline.data.(mkLineDirective)
	data.elseLine = elseLine
	mkline.data = data
}

func (mkline *MkLineImpl) MustExist() bool { return mkline.data.(mkLineInclude).mustExist }

func (mkline *MkLineImpl) IncludedFile() string { return mkline.data.(mkLineInclude).includedFile }

func (mkline *MkLineImpl) Targets() string { return mkline.data.(mkLineDependency).targets }

func (mkline *MkLineImpl) Sources() string { return mkline.data.(mkLineDependency).sources }

// ConditionalVars applies to .include lines and is a space-separated
// list of those variable names on which the inclusion depends.
// It is initialized later, step by step, when parsing other lines.
func (mkline *MkLineImpl) ConditionalVars() []string {
	return mkline.data.(mkLineInclude).conditionalVars
}
func (mkline *MkLineImpl) SetConditionalVars(varnames []string) {
	include := mkline.data.(mkLineInclude)
	include.conditionalVars = varnames
	mkline.data = include
}

// Tokenize extracts variable uses and other text from the given text.
//
// When used in IsVarassign lines, the given text must have the format
// after stripping the end-of-line comment. Such text is available from
// Value. A shell comment is therefore marked by a simple #, not an escaped
// \# like in Makefiles.
//
// When used in IsShellCommand lines, # does not mark a Makefile comment
// and may thus still appear in the text. Therefore, # marks a shell comment.
//
// Example:
//  input:  ${PREFIX}/bin abc
//  output: [MkToken("${PREFIX}", MkVarUse("PREFIX")), MkToken("/bin abc")]
//
// See ValueTokens, which is the tokenized version of Value.
func (mkline *MkLineImpl) Tokenize(text string, warn bool) []*MkToken {
	if trace.Tracing {
		defer trace.Call(mkline, text)()
	}

	var tokens []*MkToken
	var rest string
	if (mkline.IsVarassign() || mkline.IsCommentedVarassign()) && text == mkline.Value() {
		tokens, rest = mkline.ValueTokens()
	} else {
		p := NewMkParser(mkline.Line, text, true)
		tokens = p.MkTokens()
		rest = p.Rest()
	}

	if warn && rest != "" {
		mkline.Warnf("Internal pkglint error in MkLine.Tokenize at %q.", rest)
	}
	return tokens
}

// ValueSplit splits the given value, taking care of variable references.
// Example:
//
//  ValueSplit("${VAR:Udefault}::${VAR2}two:words", ":")
//  => "${VAR:Udefault}"
//     ""
//     "${VAR2}two"
//     "words"
//
// Note that even though the first word contains a colon, it is not split
// at that point since the colon is inside a variable use.
//
// When several separators are adjacent, this results in empty words in the output.
func (mkline *MkLineImpl) ValueSplit(value string, separator string) []string {
	G.Assertf(separator != "", "Separator must not be empty; use ValueFields to split on whitespace")

	tokens := mkline.Tokenize(value, false)
	var split []string
	cont := false

	out := func(s string) {
		if cont {
			split[len(split)-1] += s
		} else {
			split = append(split, s)
		}
	}

	for _, token := range tokens {
		if token.Varuse != nil {
			out(token.Text)
			cont = true
		} else {
			lexer := textproc.NewLexer(token.Text)
			for !lexer.EOF() {
				if lexer.SkipString(separator) {
					out("")
					cont = false
				}
				idx := strings.Index(lexer.Rest(), separator)
				if idx == -1 {
					idx = len(lexer.Rest())
				}
				if idx > 0 {
					out(lexer.NextString(lexer.Rest()[:idx]))
					cont = true
				}
			}
		}
	}
	return split
}

var notSpace = textproc.Space.Inverse()

// ValueFields splits the given value in the same way as the :M variable
// modifier, taking care of variable references. Example:
//
//  ValueFields("${VAR:Udefault value} ${VAR2}two words;;; 'word three'")
//  => "${VAR:Udefault value}"
//     "${VAR2}two"
//     "words;;;"
//     "'word three'"
//
// Note that even though the first word contains a space, it is not split
// at that point since the space is inside a variable use. Shell tokens
// such as semicolons are also treated as normal characters. Only double
// and single quotes are interpreted.
//
// Compare devel/bmake/files/str.c, function brk_string.
//
// TODO: Compare with brk_string from devel/bmake, especially for backticks.
func (mkline *MkLineImpl) ValueFields(value string) []string {
	if trace.Tracing {
		defer trace.Call(mkline, value)()
	}

	p := NewShTokenizer(mkline.Line, value, false)
	atoms := p.ShAtoms()

	if len(atoms) > 0 && atoms[0].Type == shtSpace {
		atoms = atoms[1:]
	}

	var word strings.Builder
	var words []string
	for _, atom := range atoms {
		if atom.Type == shtSpace && atom.Quoting == shqPlain {
			words = append(words, word.String())
			word.Reset()
		} else {
			word.WriteString(atom.MkText)
		}
	}
	if word.Len() > 0 && atoms[len(atoms)-1].Quoting == shqPlain {
		words = append(words, word.String())
		word.Reset()
	}

	// TODO: Handle parse errors
	word.WriteString(p.parser.Rest())
	rest := word.String()
	_ = rest

	return words
}

func (mkline *MkLineImpl) ValueTokens() ([]*MkToken, string) {
	value := mkline.Value()
	if value == "" {
		return nil, ""
	}

	assign := mkline.data.(mkLineAssign)
	if assign.valueMk != nil || assign.valueMkRest != "" {
		return assign.valueMk, assign.valueMkRest
	}

	// No error checking here since all this has already been done when the
	// whole line was parsed in MkLineParser.Parse.
	p := NewMkParser(nil, value, false)
	assign.valueMk = p.MkTokens()
	assign.valueMkRest = p.Rest()
	return assign.valueMk, assign.valueMkRest
}

// Fields applies to variable assignments and .for loops.
// For variable assignments, it returns the right-hand side, properly split into words.
// For .for loops, it returns all arguments (including variable names), properly split into words.
func (mkline *MkLineImpl) Fields() []string {
	if mkline.IsVarassign() || mkline.IsCommentedVarassign() {
		value := mkline.Value()
		if value == "" {
			return nil
		}

		assign := mkline.data.(mkLineAssign)
		if assign.fields != nil {
			return assign.fields
		}

		assign.fields = mkline.ValueFields(value)
		return assign.fields
	}

	// For .for loops.
	args := mkline.Args()
	if args == "" {
		return nil
	}

	directive := mkline.data.(mkLineDirective)
	if directive.fields != nil {
		return directive.fields
	}

	directive.fields = mkline.ValueFields(args)
	return directive.fields

}

func (mkline *MkLineImpl) WithoutMakeVariables(value string) string {
	var valueNovar strings.Builder
	for _, token := range NewMkParser(nil, value, false).MkTokens() {
		if token.Varuse == nil {
			valueNovar.WriteString(token.Text)
		}
	}
	return valueNovar.String()
}

func (mkline *MkLineImpl) ResolveVarsInRelativePath(relativePath string) string {
	if !contains(relativePath, "$") {
		return cleanpath(relativePath)
	}

	var basedir string
	if G.Pkg != nil {
		basedir = G.Pkg.File(".")
	} else {
		basedir = path.Dir(mkline.Filename)
	}

	tmp := relativePath
	if contains(tmp, "PKGSRCDIR") {
		pkgsrcdir := relpath(basedir, G.Pkgsrc.File("."))

		if G.Testing {
			// Relative pkgsrc paths usually only contain two or three levels.
			// A possible reason for reaching this assertion is a pkglint unit test
			// that uses t.NewMkLines instead of the correct t.SetUpFileMkLines.
			G.Assertf(!contains(pkgsrcdir, "../../../../.."),
				"Relative path %q for %q is too deep below the pkgsrc root %q.",
				pkgsrcdir, basedir, G.Pkgsrc.File("."))
		}
		tmp = strings.Replace(tmp, "${PKGSRCDIR}", pkgsrcdir, -1)
	}

	// Strictly speaking, the .CURDIR should be replaced with the basedir.
	// Depending on whether pkglint is executed with a relative or an absolute
	// path, this would produce diagnostics that "this relative path must not
	// be absolute". Since ${.CURDIR} is usually used in package Makefiles and
	// followed by "../.." anyway, the exact directory doesn't matter.
	tmp = strings.Replace(tmp, "${.CURDIR}", ".", -1)

	// TODO: Add test for exists(${.PARSEDIR}/file).
	// TODO: Add test for evaluating ${.PARSEDIR} in an included package.
	// TODO: Add test for including ${.PARSEDIR}/other.mk.
	// TODO: Add test for evaluating ${.PARSEDIR} in the infrastructure.
	//  This is the only practically relevant use case since the category
	//  directories don't contain any *.mk files that could be included.
	// TODO: Add test that suggests ${.PARSEDIR} in .include to be omitted.
	tmp = strings.Replace(tmp, "${.PARSEDIR}", ".", -1)

	replaceLatest := func(varuse, category string, pattern regex.Pattern, replacement string) {
		if contains(tmp, varuse) {
			latest := G.Pkgsrc.Latest(category, pattern, replacement)
			tmp = strings.Replace(tmp, varuse, latest, -1)
		}
	}

	// These variables are only used in pkgsrc packages, therefore they
	// are replaced with the fixed "../.." regardless of where the text appears.
	replaceLatest("${LUA_PKGSRCDIR}", "lang", `^lua[0-9]+$`, "../../lang/$0")
	replaceLatest("${PHPPKGSRCDIR}", "lang", `^php[0-9]+$`, "../../lang/$0")
	replaceLatest("${PYPKGSRCDIR}", "lang", `^python[0-9]+$`, "../../lang/$0")

	replaceLatest("${PYPACKAGE}", "lang", `^python[0-9]+$`, "$0")
	replaceLatest("${SUSE_DIR_PREFIX}", "emulators", `^(suse[0-9]+)_base$`, "$1")

	if G.Pkg != nil {
		// XXX: Even if these variables are defined indirectly,
		// pkglint should be able to resolve them properly.
		// There is already G.Pkg.Value, maybe that can be used here.
		tmp = strings.Replace(tmp, "${FILESDIR}", G.Pkg.Filesdir, -1)
		tmp = strings.Replace(tmp, "${PKGDIR}", G.Pkg.Pkgdir, -1)
	}

	tmp = cleanpath(tmp)

	if trace.Tracing && relativePath != tmp {
		trace.Step2("resolveVarsInRelativePath: %q => %q", relativePath, tmp)
	}
	return tmp
}

func (mkline *MkLineImpl) ExplainRelativeDirs() {
	mkline.Explain(
		"Directories in the form \"../../category/package\" make it easier to",
		"move a package around in pkgsrc, for example from pkgsrc-wip to the",
		"main pkgsrc repository.")
}

// RefTo returns a reference to another line,
// which can be in the same file or in a different file.
//
// If there is a type mismatch when calling this function, try to add ".line" to
// either the method receiver or the other line.
func (mkline *MkLineImpl) RefTo(other MkLine) string {
	return mkline.Line.RefTo(other.Line)
}

var (
	LowerDash                  = textproc.NewByteSet("a-z---")
	AlnumDot                   = textproc.NewByteSet("A-Za-z0-9_.")
	unescapeMkCommentSafeChars = textproc.NewByteSet("\\#[$").Inverse()
)

// unescapeComment takes a Makefile line, as written in a file, and splits
// it into the main part and the comment.
//
// The comment starts at the first #. Except if it is preceded by an odd number
// of backslashes. Or by an opening bracket.
//
// The main text is returned including leading and trailing whitespace. Any
// escaped # is returned in its unescaped form, that is, \# becomes #.
//
// The comment is returned including the leading "#", if any. If the line has
// no comment, it is an empty string.
func (p MkLineParser) unescapeComment(text string) (main, comment string) {
	var sb strings.Builder

	lexer := textproc.NewLexer(text)

again:
	if plain := lexer.NextBytesSet(unescapeMkCommentSafeChars); plain != "" {
		sb.WriteString(plain)
		goto again
	}

	switch {
	case lexer.SkipByte('$'):
		sb.WriteByte('$')

	case lexer.SkipString("\\#"):
		sb.WriteByte('#')

	case lexer.PeekByte() == '\\' && len(lexer.Rest()) >= 2:
		sb.WriteString(lexer.Rest()[:2])
		lexer.Skip(2)

	case lexer.SkipByte('\\'):
		sb.WriteByte('\\')

	case lexer.SkipString("[#"):
		// See devel/bmake/files/parse.c:/as in modifier/
		sb.WriteString("[#")

	case lexer.SkipByte('['):
		sb.WriteByte('[')

	default:
		main = sb.String()
		if lexer.PeekByte() == '#' {
			return main, lexer.Rest()
		}

		G.Assertf(lexer.EOF(), "unescapeComment(%q): sb = %q, rest = %q", text, main, lexer.Rest())
		return main, ""
	}

	goto again
}

type mkLineSplitResult struct {
	main               string
	tokens             []*MkToken
	spaceBeforeComment string
	hasComment         bool
	comment            string
}

// splitMkLine parses a logical line from a Makefile (that is, after joining
// the lines that end in a backslash) into two parts: the main part and the
// comment.
//
// This applies to all line types except those starting with a tab, which
// contain the shell commands to be associated with make targets. These cannot
// have comments.
func (p MkLineParser) split(line Line, text string) mkLineSplitResult {

	main, comment := p.unescapeComment(text)

	parser := NewMkParser(line, main, line != nil)
	lexer := parser.lexer

	rtrimHspace := func(s string) string {
		end := len(s)
		for end > 0 && isHspace(s[end-1]) {
			end--
		}
		return s[:end]
	}

	parseOther := func() string {
		var sb strings.Builder

		for !lexer.EOF() {
			if lexer.SkipString("$$") {
				sb.WriteString("$$")
				continue
			}

			other := lexer.NextBytesFunc(func(b byte) bool { return b != '$' })
			if other == "" {
				break
			}

			sb.WriteString(other)
		}

		return sb.String()
	}

	var tokens []*MkToken
	for !lexer.EOF() {
		mark := lexer.Mark()

		if varUse := parser.VarUse(); varUse != nil {
			tokens = append(tokens, &MkToken{lexer.Since(mark), varUse})

		} else if other := parseOther(); other != "" {
			tokens = append(tokens, &MkToken{other, nil})

		} else {
			G.Assertf(lexer.SkipByte('$'), "Parse error for %q.", text)
			tokens = append(tokens, &MkToken{"$", nil})
		}
	}

	hasComment := comment != ""
	if hasComment {
		comment = comment[1:]
	}

	G.Assertf(lexer.Rest() == "", "Parse error for %q.", text)

	mainWithSpaces := main
	main = rtrimHspace(main)
	spaceBeforeComment := ifelseStr(true, mainWithSpaces[len(main):], "")
	if spaceBeforeComment != "" && len(tokens) > 0 {
		tokenText := &tokens[len(tokens)-1].Text
		*tokenText = rtrimHspace(*tokenText)
		if *tokenText == "" {
			tokens = tokens[:len(tokens)-1]
		}
	}

	return mkLineSplitResult{main, tokens, spaceBeforeComment, hasComment, comment}
}

func (p MkLineParser) parseDirective(line Line, data mkLineSplitResult) MkLine {
	text := line.Text
	if !hasPrefix(text, ".") {
		return nil
	}

	lexer := textproc.NewLexer(data.main[1:])

	indent := lexer.NextHspace()
	directive := lexer.NextBytesSet(LowerDash)
	switch directive {
	case "if", "else", "elif", "endif",
		"ifdef", "ifndef",
		"for", "endfor", "undef",
		"error", "warning", "info",
		"export", "export-env", "unexport", "unexport-env":
		break
	default:
		// Intentionally not supported are: ifmake ifnmake elifdef elifndef elifmake elifnmake.
		return nil
	}

	lexer.SkipHspace()

	args := lexer.Rest()

	// In .if and .endif lines the space surrounding the comment is irrelevant.
	// Especially for checking that the .endif comment matches the .if condition,
	// it must be trimmed.
	trimmedComment := trimHspace(data.comment)

	return &MkLineImpl{line, &mkLineDirectiveImpl{indent, directive, args, trimmedComment, nil, nil, nil}}
}

// VariableNeedsQuoting determines whether the given variable needs the :Q operator
// in the given context.
//
// This decision depends on many factors, such as whether the type of the context is
// a list of things, whether the variable is a list, whether it can contain only
// safe characters, and so on.
func (mkline *MkLineImpl) VariableNeedsQuoting(mklines MkLines, varuse *MkVarUse, vartype *Vartype, vuc *VarUseContext) (needsQuoting YesNoUnknown) {
	if trace.Tracing {
		defer trace.Call(varuse, vartype, vuc, trace.Result(&needsQuoting))()
	}

	// TODO: Systematically test this function, each and every case, from top to bottom.
	// TODO: Re-check the order of all these if clauses whether it really makes sense.

	vucVartype := vuc.vartype
	if vartype == nil || vucVartype == nil || vartype.basicType == BtUnknown {
		return unknown
	}

	if !vartype.basicType.NeedsQ() {
		if !vartype.List() {
			if vartype.Guessed() {
				return unknown
			}
			return no
		}
		if !vuc.IsWordPart {
			return no
		}
	}

	// A shell word may appear as part of a shell word, for example COMPILER_RPATH_FLAG.
	if vuc.IsWordPart && vuc.quoting == VucQuotPlain {
		if !vartype.List() && vartype.basicType == BtShellWord {
			return no
		}
	}

	// Determine whether the context expects a list of shell words or not.
	wantList := vucVartype.MayBeAppendedTo()
	haveList := vartype.MayBeAppendedTo()
	if trace.Tracing {
		trace.Stepf("wantList=%v, haveList=%v", wantList, haveList)
	}

	// Both of these can be correct, depending on the situation:
	// 1. echo ${PERL5:Q}
	// 2. xargs ${PERL5}
	if !vuc.IsWordPart && wantList && haveList {
		return unknown
	}

	// Pkglint assumes that the tool definitions don't include very
	// special characters, so they can safely be used inside any quotes.
	if tool := G.ToolByVarname(mklines, varuse.varname); tool != nil {
		switch vuc.quoting {
		case VucQuotPlain:
			if !vuc.IsWordPart {
				return no
			}
			// XXX: Should there be a return here? It looks as if it could have been forgotten.
		case VucQuotBackt:
			return no
		case VucQuotDquot, VucQuotSquot:
			return unknown
		}
	}

	// Variables that appear as parts of shell words generally need to be quoted.
	//
	// An exception is in the case of backticks, because the whole backticks expression
	// is parsed as a single shell word by pkglint. (XXX: This comment may be outdated.)
	if vuc.IsWordPart && vucVartype.IsShell() && vuc.quoting != VucQuotBackt {
		return yes
	}

	// SUBST_MESSAGE.perl= Replacing in ${REPLACE_PERL}
	if vucVartype.IsPlainString() {
		return no
	}

	if wantList != haveList {
		if vucVartype.basicType == BtFetchURL && vartype.basicType == BtHomepage {
			return no
		}
		if vucVartype.basicType == BtHomepage && vartype.basicType == BtFetchURL {
			return no // Just for HOMEPAGE=${MASTER_SITE_*:=subdir/}.
		}

		// .for dir in ${PATH:C,:, ,g}
		for _, modifier := range varuse.modifiers {
			if modifier.ChangesWords() {
				return unknown
			}
		}

		return yes
	}

	// Bad: LDADD+= -l${LIBS}
	// Good: LDADD+= ${LIBS:S,^,-l,}
	if wantList {
		return yes
	}

	if trace.Tracing {
		trace.Step1("Don't know whether :Q is needed for %q", varuse.varname)
	}
	return unknown
}

// ForEachUsed calls the action for each variable that is used in the line.
func (mkline *MkLineImpl) ForEachUsed(action func(varUse *MkVarUse, time vucTime)) {

	var searchIn func(text string, time vucTime) // mutually recursive with searchInVarUse

	searchInVarUse := func(varuse *MkVarUse, time vucTime) {
		varname := varuse.varname
		if !varuse.IsExpression() {
			action(varuse, time)
		}
		searchIn(varname, time)
		for _, mod := range varuse.modifiers {
			searchIn(mod.Text, time)
		}
	}

	searchIn = func(text string, time vucTime) {
		if !contains(text, "$") {
			return
		}

		for _, token := range NewMkParser(nil, text, false).MkTokens() {
			if token.Varuse != nil {
				searchInVarUse(token.Varuse, time)
			}
		}
	}

	switch {

	case mkline.IsVarassign():
		searchIn(mkline.Varname(), vucTimeParse)
		searchIn(mkline.Value(), mkline.Op().Time())

	case mkline.IsDirective() && mkline.Directive() == "for":
		searchIn(mkline.Args(), vucTimeParse)

	case mkline.IsDirective() && mkline.Cond() != nil:
		mkline.Cond().Walk(&MkCondCallback{
			VarUse: func(varuse *MkVarUse) {
				searchInVarUse(varuse, vucTimeParse)
			}})

	case mkline.IsShellCommand():
		searchIn(mkline.ShellCommand(), vucTimeRun)

	case mkline.IsDependency():
		searchIn(mkline.Targets(), vucTimeParse)
		searchIn(mkline.Sources(), vucTimeParse)

	case mkline.IsInclude():
		searchIn(mkline.IncludedFile(), vucTimeParse)
	}
}

func (mkline *MkLineImpl) UnquoteShell(str string) string {
	var sb strings.Builder
	n := len(str)

outer:
	for i := 0; i < n; i++ {
		switch str[i] {
		case '"':
			for i++; i < n; i++ {
				switch str[i] {
				case '"':
					continue outer
				case '\\':
					i++
					if i < n {
						sb.WriteByte(str[i])
					}
				default:
					sb.WriteByte(str[i])
				}
			}

		case '\'':
			for i++; i < n && str[i] != '\''; i++ {
				sb.WriteByte(str[i])
			}

		case '\\':
			i++
			if i < n {
				sb.WriteByte(str[i])
			}

		default:
			sb.WriteByte(str[i])
		}
	}

	return sb.String()
}

type MkOperator uint8

const (
	opAssign        MkOperator = iota // =
	opAssignShell                     // !=
	opAssignEval                      // :=
	opAssignAppend                    // +=
	opAssignDefault                   // ?=
	opUseCompare                      // A variable is compared to a value, e.g. in a condition.
	opUseMatch                        // A variable is matched using the :M or :N modifier.
)

func NewMkOperator(op string) MkOperator {
	switch op {
	case "=":
		return opAssign
	case "!=":
		return opAssignShell
	case ":=":
		return opAssignEval
	case "+=":
		return opAssignAppend
	case "?=":
		return opAssignDefault
	}
	panic("Invalid operator: " + op)
}

func (op MkOperator) String() string {
	return [...]string{"=", "!=", ":=", "+=", "?=", "use", "use-loadtime", "use-match"}[op]
}

// Time returns the time at which the right-hand side of the assignment is
// evaluated.
func (op MkOperator) Time() vucTime {
	if op == opAssignShell || op == opAssignEval {
		return vucTimeParse
	}
	return vucTimeRun
}

// VarUseContext defines the context in which a variable is defined
// or used. Whether that is allowed depends on:
//
// * The variable's data type, as defined in vardefs.go.
//
// * When used on the right-hand side of an assigment, the variable can
// represent a list of words, a single word or even only part of a
// word. This distinction decides upon the correct use of the :Q
// operator.
//
// * When used in preprocessing statements like .if or .for, the other
// operands of that statement should fit to the variable and are
// checked against the variable type. For example, comparing OPSYS to
// x86_64 doesn't make sense.
type VarUseContext struct {
	vartype    *Vartype
	time       vucTime
	quoting    VucQuoting
	IsWordPart bool // Example: LOCALBASE=${LOCALBASE}
}

// vucTime is the time at which a variable is used.
//
// See ToolTime, which is the same except that there is no unknown.
type vucTime uint8

const (
	vucTimeUnknown vucTime = iota

	// When Makefiles are loaded, the operators := and != evaluate their
	// right-hand side, as well as the directives .if, .elif and .for.
	// During loading, not all variables are available yet.
	// Variable values are still subject to change, especially lists.
	vucTimeParse

	// All files have been read, all variables can be referenced.
	// Variable values don't change anymore.
	//
	// Well, except for the ::= modifier.
	// But that modifier is usually not used in pkgsrc.
	vucTimeRun
)

func (t vucTime) String() string { return [...]string{"unknown", "parse", "run"}[t] }

// VucQuoting describes in what level of quoting the variable is used.
// Depending on this context, the modifiers :Q or :M can be allowed or not.
//
// The shell tokenizer knows multi-level quoting modes (see ShQuoting),
// but for deciding whether :Q is necessary or not, a single level is enough.
type VucQuoting uint8

const (
	VucQuotUnknown VucQuoting = iota
	VucQuotPlain              // Example: echo LOCALBASE=${LOCALBASE}
	VucQuotDquot              // Example: echo "The version is ${PKGVERSION}."
	VucQuotSquot              // Example: echo 'The version is ${PKGVERSION}.'
	VucQuotBackt              // Example: echo `sed 1q ${WRKSRC}/README`
)

func (q VucQuoting) String() string {
	return [...]string{"unknown", "plain", "dquot", "squot", "backt", "mk-for"}[q]
}

func (vuc *VarUseContext) String() string {
	typename := "no-type"
	if vuc.vartype != nil {
		typename = vuc.vartype.String()
	}
	return sprintf("(%s time:%s quoting:%s wordpart:%v)", typename, vuc.time, vuc.quoting, vuc.IsWordPart)
}

// Indentation remembers the stack of preprocessing directives and their
// indentation. By convention, each directive is indented by 2 spaces.
// An excepting are multiple-inclusion guards, they don't increase the
// indentation.
type Indentation struct {
	levels []indentationLevel
}

func NewIndentation() *Indentation {
	ind := Indentation{}
	ind.Push(nil, 0, "") // Dummy
	return &ind
}

func (ind *Indentation) String() string {
	var s strings.Builder
	for _, level := range ind.levels[1:] {
		_, _ = fmt.Fprintf(&s, " %d", level.depth)
		if len(level.conditionalVars) > 0 {
			_, _ = fmt.Fprintf(&s, " (%s)", strings.Join(level.conditionalVars, " "))
		}
	}
	return "[" + trimHspace(s.String()) + "]"
}

func (ind *Indentation) RememberUsedVariables(cond MkCond) {
	cond.Walk(&MkCondCallback{
		VarUse: func(varuse *MkVarUse) { ind.AddVar(varuse.varname) }})
}

type indentationLevel struct {
	mkline          MkLine   // The line in which the indentation started; the .if/.for
	depth           int      // Number of space characters; always a multiple of 2
	args            string   // The arguments from the .if or .for, or the latest .elif
	conditionalVars []string // Variables on which the current path depends

	// Files whose existence has been checked in an if branch that is
	// related to the current indentation. After a .if exists(fname),
	// pkglint will happily accept .include "fname" in both the then and
	// the else branch. This is ok since the primary job of this file list
	// is to prevent wrong pkglint warnings about missing files.
	checkedFiles []string
}

func (ind *Indentation) Len() int {
	return len(ind.levels)
}

func (ind *Indentation) top() *indentationLevel {
	return &ind.levels[ind.Len()-1]
}

// Depth returns the number of space characters by which the directive
// should be indented.
//
// This is typically two more than the surrounding level, except for
// multiple-inclusion guards.
func (ind *Indentation) Depth(directive string) int {
	switch directive {
	case "if", "elif", "else", "endfor", "endif":
		return ind.levels[imax(0, ind.Len()-2)].depth
	}
	return ind.top().depth
}

func (ind *Indentation) Pop() {
	ind.levels = ind.levels[:ind.Len()-1]
}

func (ind *Indentation) Push(mkline MkLine, indent int, condition string) {
	ind.levels = append(ind.levels, indentationLevel{mkline, indent, condition, nil, nil})
}

// AddVar remembers that the current indentation depends on the given variable,
// most probably because that variable is used in a .if directive.
//
// Variables named *_MK are ignored since they are usually not interesting.
func (ind *Indentation) AddVar(varname string) {
	if hasSuffix(varname, "_MK") {
		return
	}

	vars := &ind.top().conditionalVars
	for _, existingVarname := range *vars {
		if varname == existingVarname {
			return
		}
	}

	*vars = append(*vars, varname)
}

func (ind *Indentation) DependsOn(varname string) bool {
	for _, level := range ind.levels {
		for _, levelVarname := range level.conditionalVars {
			if varname == levelVarname {
				return true
			}
		}
	}
	return false
}

// IsConditional returns whether the current line depends on evaluating
// any variable in an .if or .elif expression or from a .for loop.
//
// Variables named *_MK are excluded since they are usually not interesting.
func (ind *Indentation) IsConditional() bool {
	for _, level := range ind.levels {
		if len(level.conditionalVars) > 0 {
			return true
		}
	}
	return false
}

// Varnames returns the list of all variables that are mentioned in any
// condition or loop surrounding the current line.
//
// Variables named *_MK are excluded since they are usually not interesting.
func (ind *Indentation) Varnames() []string {
	var varnames []string
	for _, level := range ind.levels {
		for _, levelVarname := range level.conditionalVars {
			G.Assertf(
				!hasSuffix(levelVarname, "_MK"),
				"multiple-inclusion guard must be filtered out earlier.")
			varnames = append(varnames, levelVarname)
		}
	}
	return varnames
}

// Args returns the arguments of the innermost .if, .elif or .for.
func (ind *Indentation) Args() string {
	return ind.top().args
}

func (ind *Indentation) AddCheckedFile(filename string) {
	top := ind.top()
	top.checkedFiles = append(top.checkedFiles, filename)
}

func (ind *Indentation) IsCheckedFile(filename string) bool {
	for _, level := range ind.levels {
		for _, levelFilename := range level.checkedFiles {
			if filename == levelFilename {
				return true
			}
		}
	}
	return false
}

func (ind *Indentation) TrackBefore(mkline MkLine) {
	if !mkline.IsDirective() {
		return
	}
	if trace.Tracing {
		trace.Stepf("Indentation before line %s: %s", mkline.Linenos(), ind)
	}

	switch mkline.Directive() {
	case "for", "if", "ifdef", "ifndef":
		ind.Push(mkline, ind.top().depth, mkline.Args())
	}
}

func (ind *Indentation) TrackAfter(mkline MkLine) {
	if !mkline.IsDirective() {
		return
	}

	directive := mkline.Directive()
	args := mkline.Args()

	switch directive {
	case "if":
		cond := mkline.Cond()

		// For multiple-inclusion guards, the indentation stays at the same level.
		guard := cond != nil && cond.Not != nil && hasSuffix(cond.Not.Defined, "_MK")
		if !guard {
			ind.top().depth += 2
		}

	case "for", "ifdef", "ifndef":
		ind.top().depth += 2

	case "elif":
		// Handled here instead of TrackBefore to allow the action to access the previous condition.
		ind.top().args = args

	case "else":
		top := ind.top()
		if top.mkline != nil {
			top.mkline.SetHasElseBranch(mkline)
		}

	case "endfor", "endif":
		if ind.Len() > 1 { // Can only be false in unbalanced files.
			ind.Pop()
		}
	}

	switch directive {
	case "if", "elif":
		cond := mkline.Cond()
		if cond == nil {
			break
		}

		ind.RememberUsedVariables(cond)

		cond.Walk(&MkCondCallback{
			Call: func(name string, arg string) {
				if name == "exists" {
					ind.AddCheckedFile(arg)
				}
			}})
	}

	if trace.Tracing {
		trace.Stepf("Indentation after line %s: %s", mkline.Linenos(), ind)
	}
}

func (ind *Indentation) CheckFinish(filename string) {
	if ind.Len() <= 1 {
		return
	}
	eofLine := NewLineEOF(filename)
	for ind.Len() > 1 {
		openingMkline := ind.levels[ind.Len()-1].mkline
		eofLine.Errorf(".%s from %s must be closed.", openingMkline.Directive(), eofLine.RefTo(openingMkline.Line))
		ind.Pop()
	}
}

// VarbaseBytes contains characters that may be used in the main part of variable names.
// VarparamBytes contains characters that may be used in the parameter part of variable names.
//
// For example, TOOLS_PATH.[ is a valid variable name but [ alone isn't since
// the opening bracket is only allowed in the parameter part of variable names.
//
// This approach differs from the one in devel/bmake/files/parse.c:/^Parse_IsVar,
// but in practice it works equally well. Luckily there aren't many situations
// where a complicated variable name contains unbalanced parentheses or braces,
// which would confuse the devel/bmake parser.
//
// TODO: The allowed characters differ between the basename and the parameter
//  of the variable. The square bracket is only allowed in the parameter part.
var (
	VarbaseBytes  = textproc.NewByteSet("A-Za-z_0-9+---")
	VarparamBytes = textproc.NewByteSet("A-Za-z_0-9#*+---.[")
)

func (p MkLineParser) MatchVarassign(line Line, text string, asdfData mkLineSplitResult) (m bool, assignment mkLineAssign) {

	// A commented variable assignment does not have leading whitespace.
	// Otherwise line 1 of almost every Makefile fragment would need to
	// be scanned for a variable assignment even though it only contains
	// the $NetBSD CVS Id.
	clex := textproc.NewLexer(text)
	commented := clex.SkipByte('#')
	if commented && clex.SkipHspace() || clex.EOF() {
		return false, nil
	}

	withoutLeadingComment := text
	if commented {
		withoutLeadingComment = withoutLeadingComment[1:]
	}

	data := p.split(nil, withoutLeadingComment)

	lexer := NewMkTokensLexer(data.tokens)
	mainStart := lexer.Mark()

	for !commented && lexer.SkipByte(' ') {
	}

	varnameStart := lexer.Mark()
	// TODO: duplicated code in MkParser.Varname
	for lexer.NextBytesSet(VarbaseBytes) != "" || lexer.NextVarUse() != nil {
	}
	if lexer.SkipByte('.') || hasPrefix(data.main, "SITES_") {
		for lexer.NextBytesSet(VarparamBytes) != "" || lexer.NextVarUse() != nil {
		}
	}

	varname := lexer.Since(varnameStart)

	if varname == "" {
		return
	}

	spaceAfterVarname := lexer.NextHspace()

	opStart := lexer.Mark()
	switch lexer.PeekByte() {
	case '!', '+', ':', '?':
		lexer.Skip(1)
	}
	if !lexer.SkipByte('=') {
		return
	}
	op := NewMkOperator(lexer.Since(opStart))

	if hasSuffix(varname, "+") && op == opAssign && spaceAfterVarname == "" {
		varname = varname[:len(varname)-1]
		op = opAssignAppend
	}

	lexer.SkipHspace()

	value := trimHspace(lexer.Rest())
	valueAlign := ifelseStr(commented, "#", "") + lexer.Since(mainStart)
	spaceBeforeComment := data.spaceBeforeComment
	if value == "" {
		valueAlign += spaceBeforeComment
		spaceBeforeComment = ""
	}

	return true, &mkLineAssignImpl{
		commented:         commented,
		varname:           varname,
		varcanon:          varnameCanon(varname),
		varparam:          varnameParam(varname),
		spaceAfterVarname: spaceAfterVarname,
		op:                op,
		valueAlign:        valueAlign,
		value:             value,
		valueMk:           nil, // filled in lazily
		valueMkRest:       "",  // filled in lazily
		fields:            nil, // filled in lazily
		spaceAfterValue:   spaceBeforeComment,
		comment:           ifelseStr(data.hasComment, "#", "") + data.comment,
	}
}

func MatchMkInclude(text string) (m bool, indentation, directive, filename string) {
	lexer := textproc.NewLexer(text)
	if lexer.SkipByte('.') {
		indentation = lexer.NextHspace()
		directive = lexer.NextString("include")
		if directive == "" {
			directive = lexer.NextString("sinclude")
		}
		if directive != "" {
			lexer.NextHspace()
			if lexer.SkipByte('"') {
				// Note: strictly speaking, the full MkVarUse would have to be parsed
				// here. But since these usually don't contain double quotes, it has
				// worked fine up to now.
				filename = lexer.NextBytesFunc(func(c byte) bool { return c != '"' })
				if filename != "" && lexer.SkipByte('"') {
					lexer.NextHspace()
					if lexer.EOF() || lexer.SkipByte('#') {
						m = true
						return
					}
				}
			}
		}
	}
	return false, "", "", ""
}
