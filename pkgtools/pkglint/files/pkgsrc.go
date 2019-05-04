package pkglint

import (
	"io/ioutil"
	"netbsd.org/pkglint/regex"
	"netbsd.org/pkglint/textproc"
	"os"
	"path/filepath"
	"sort"
	"strings"
)

// Pkgsrc describes a pkgsrc installation.
// In each pkglint run, only a single pkgsrc installation is ever loaded.
// It just doesn't make sense to check multiple pkgsrc installations at once.
type Pkgsrc struct {
	// The top directory (PKGSRCDIR), either absolute or relative to
	// the current working directory.
	topdir string

	// The set of user-defined variables that are added to BUILD_DEFS
	// within the bsd.pkg.mk file.
	buildDefs map[string]bool

	Tools *Tools

	MasterSiteURLToVar map[string]string // "https://github.com/" => "MASTER_SITE_GITHUB"
	MasterSiteVarToURL map[string]string // "MASTER_SITE_GITHUB" => "https://github.com/"

	PkgOptions map[string]string // "x11" => "Provides X11 support"

	suggestedUpdates    []SuggestedUpdate   //
	suggestedWipUpdates []SuggestedUpdate   //
	LastChange          map[string]*Change  //
	listVersions        map[string][]string // See ListVersions

	// Variables that may be overridden by the pkgsrc user.
	// They are typically defined in mk/defaults/mk.conf.
	//
	// Whenever a package uses such a variable, it must add the variable name
	// to BUILD_DEFS.
	UserDefinedVars Scope

	Deprecated map[string]string
	vartypes   VarTypeRegistry
}

func NewPkgsrc(dir string) Pkgsrc {
	return Pkgsrc{
		dir,
		make(map[string]bool),
		NewTools(),
		make(map[string]string),
		make(map[string]string),
		make(map[string]string),
		nil,
		nil,
		make(map[string]*Change),
		make(map[string][]string),
		NewScope(),
		make(map[string]string),
		NewVarTypeRegistry()}
}

func (src *Pkgsrc) loadDefaultBuildDefs() {

	// Some user-defined variables do not influence the binary
	// package at all and therefore do not have to be added to
	// BUILD_DEFS; therefore they are marked as "already added".
	src.addBuildDefs(
		"DISTDIR",
		"FETCH_CMD",
		"FETCH_OUTPUT_ARGS",
		"FETCH_USING",
		"PKGSRC_RUN_TEST")

	// The following variables are used so often that not every
	// package should need to add it to BUILD_DEFS manually.
	src.addBuildDefs(
		"PKGSRC_COMPILER",
		"PKGSRC_USE_SSP",
		"UNPRIVILEGED",
		"USE_CROSS_COMPILE")

	// The following variables are so obscure that they are
	// probably not used in practice.
	src.addBuildDefs(
		"MANINSTALL")

	// The following variables are added to _BUILD_DEFS by the pkgsrc
	// infrastructure and thus don't need to be added by the package again.
	// To regenerate the below list:
	//  grep -hr '^_BUILD_DEFS+=' mk/ | tr ' \t' '\n\n' | sed -e 's,.*=,,' -e '/^_/d' -e '/^$/d' -e 's,.*,"&"\,,' | sort -u
	// TODO: Run the equivalent of the above command at startup.
	src.addBuildDefs(
		"ABI",
		"BUILTIN_PKGS",
		"CFLAGS",
		"CMAKE_ARGS",
		"CONFIGURE_ARGS",
		"CONFIGURE_ENV",
		"CPPFLAGS",
		"FFLAGS",
		"GAMEDATAMODE",
		"GAMEDIRMODE",
		"GAMEMODE",
		"GAMES_GROUP",
		"GAMES_USER",
		"GLIBC_VERSION",
		"INIT_SYSTEM",
		"LDFLAGS",
		"LICENSE",
		"LOCALBASE",
		"MACHINE_ARCH",
		"MACHINE_GNU_ARCH",
		"MULTI",
		"NO_BIN_ON_CDROM",
		"NO_BIN_ON_FTP",
		"NO_SRC_ON_CDROM",
		"NO_SRC_ON_FTP",
		"OBJECT_FMT",
		"OPSYS",
		"OS_VERSION",
		"OSVERSION_SPECIFIC",
		"PKG_HACKS",
		"PKG_OPTIONS",
		"PKG_SYSCONFBASEDIR",
		"PKG_SYSCONFDIR",
		"PKGGNUDIR",
		"PKGINFODIR",
		"PKGMANDIR",
		"PKGPATH",
		"RESTRICTED",
		"USE_ABI_DEPENDS")
}

// LoadInfrastructure reads the pkgsrc infrastructure files to
// extract information like the tools, packages to update,
// user-defined variables.
//
// This work is not done in the constructor to keep the tests
// simple, since setting up a realistic pkgsrc environment requires
// a lot of files.
func (src *Pkgsrc) LoadInfrastructure() {
	src.vartypes.Init(src)
	src.loadMasterSites()
	src.loadPkgOptions()
	src.loadDocChanges()
	src.loadSuggestedUpdates()
	src.loadUserDefinedVars()
	src.loadTools()
	src.initDeprecatedVars()
	src.loadUntypedVars()
	src.loadDefaultBuildDefs()
}

// Latest returns the latest package matching the given pattern.
// It searches the category for subdirectories matching the given
// regular expression, takes the latest of them and replaces its
// name with repl.
//
// Example:
//  Latest("lang", `^php[0-9]+$`, "../../lang/$0")
//      => "../../lang/php72"
func (src *Pkgsrc) Latest(category string, re regex.Pattern, repl string) string {
	versions := src.ListVersions(category, re, repl, true)

	if len(versions) > 0 {
		return versions[len(versions)-1]
	}
	return ""
}

// ListVersions searches the category for subdirectories matching the given
// regular expression, replaces their names with repl and returns a slice
// of them, properly sorted from early to late.
//
// Example:
//  ListVersions("lang", `^php[0-9]+$`, "php-$0")
//      => {"php-53", "php-56", "php-73"}
func (src *Pkgsrc) ListVersions(category string, re regex.Pattern, repl string, errorIfEmpty bool) []string {
	if G.Testing {
		G.Assertf(
			hasPrefix(string(re), "^") && hasSuffix(string(re), "$"),
			"Regular expression %q must be anchored at both ends.", re)
	}

	// TODO: Maybe convert cache key to a struct, to save allocations.
	cacheKey := category + "/" + string(re) + " => " + repl
	if latest, found := src.listVersions[cacheKey]; found {
		return latest
	}

	categoryDir := src.File(category)
	error := func() []string {
		if errorIfEmpty {
			dummyLine.Errorf("Cannot find package versions of %q in %q.", re, categoryDir)
		}
		src.listVersions[cacheKey] = nil
		return nil
	}

	var names []string
	for _, fileInfo := range src.ReadDir(category) {
		name := fileInfo.Name()
		if matches(name, re) {
			names = append(names, name)
		}
	}
	if len(names) == 0 {
		return error()
	}

	// In the pkgsrc directories, the major versions of packages are
	// written without dots, which leads to ambiguities:
	//
	// databases/postgresql: 94 < 95 < 96 < 10 < 11
	// lang/go: 19 < 110 < 111 < 2
	keys := make(map[string]int)
	for _, name := range names {
		if m, pkgbase, versionStr := match2(name, `^(\D+)(\d+)$`); m {
			version := toInt(versionStr, 0)
			if pkgbase == "postgresql" && version < 60 {
				version = 10 * version
			}
			if pkgbase == "go" {
				major := toInt(versionStr[:1], 0)
				minor := toInt(versionStr[1:], 0)
				version = 100*major + minor
			}
			keys[name] = version
		}
	}

	sort.SliceStable(names, func(i, j int) bool {
		if keyI, keyJ := keys[names[i]], keys[names[j]]; keyI != keyJ {
			return keyI < keyJ
		}
		return naturalLess(names[i], names[j])
	})

	var repls = make([]string, len(names), len(names))
	for i, name := range names {
		repls[i] = replaceAll(name, re, repl)
	}

	src.listVersions[cacheKey] = repls
	return repls
}

func (src *Pkgsrc) checkToplevelUnusedLicenses() {
	usedLicenses := G.UsedLicenses
	if usedLicenses == nil {
		return
	}

	licensesDir := src.File("licenses")
	for _, licenseFile := range src.ReadDir("licenses") {
		licenseName := licenseFile.Name()
		if !usedLicenses[licenseName] {
			licensePath := licensesDir + "/" + licenseName
			if fileExists(licensePath) {
				NewLineWhole(licensePath).Warnf("This license seems to be unused.")
			}
		}
	}
}

// loadTools loads the tool definitions from `mk/tools/*`.
func (src *Pkgsrc) loadTools() {
	tools := src.Tools

	toolFiles := []string{"defaults.mk"}
	{
		toc := src.File("mk/tools/bsd.tools.mk")
		mklines := LoadMk(toc, MustSucceed|NotEmpty)
		for _, mkline := range mklines.mklines {
			if mkline.IsInclude() {
				includedFile := mkline.IncludedFile()
				if !contains(includedFile, "/") {
					toolFiles = append(toolFiles, includedFile)
				}
			}
		}
		if len(toolFiles) <= 1 {
			NewLineWhole(toc).Fatalf("Too few tool files.")
		}
	}

	// TODO: parse bsd.prefs.mk and bsd.pkg.mk instead of hardcoding this.
	toolDefs := [...]struct {
		Name     string
		Varname  string
		Validity Validity
	}{
		{"echo", "ECHO", AfterPrefsMk},
		{"echo -n", "ECHO_N", AfterPrefsMk},
		{"false", "FALSE", AtRunTime}, // from bsd.pkg.mk
		{"test", "TEST", AfterPrefsMk},
		{"true", "TRUE", AfterPrefsMk}}

	for _, toolDef := range toolDefs {
		tools.def(toolDef.Name, toolDef.Varname, true, toolDef.Validity, nil)
	}

	for _, basename := range toolFiles {
		mklines := src.LoadMk("mk/tools/"+basename, MustSucceed|NotEmpty)
		mklines.ForEach(func(mkline MkLine) {
			tools.ParseToolLine(mklines, mkline, true, !mklines.indentation.IsConditional())
		})
	}

	for _, relativeName := range [...]string{"mk/bsd.prefs.mk", "mk/bsd.pkg.mk"} {

		mklines := src.LoadMk(relativeName, MustSucceed|NotEmpty)
		mklines.ForEach(func(mkline MkLine) {
			if mkline.IsVarassign() {
				varname := mkline.Varname()
				switch varname {
				case "USE_TOOLS":
					tools.ParseToolLine(mklines, mkline, true, !mklines.indentation.IsConditional())

				case "_BUILD_DEFS":
					// TODO: Compare with src.loadDefaultBuildDefs; is it redundant?
					for _, buildDefsVar := range mkline.Fields() {
						src.addBuildDefs(buildDefsVar)
					}
				}
			}
		})
	}

	if trace.Tracing {
		tools.Trace()
	}
}

// loadUntypedVars scans all pkgsrc infrastructure files in mk/
// to find variable definitions that are not yet covered in
// Pkgsrc.InitVartypes.
//
// Even if pkglint cannot guess the type of each variable,
// at least prevent the "used but not defined" warnings.
func (src *Pkgsrc) loadUntypedVars() {

	// Setting guessed to false prevents the vartype.guessed case in MkLineChecker.CheckVaruse.
	unknownType := Vartype{BtUnknown, NoVartypeOptions, []ACLEntry{{"*", aclpAll}}}

	define := func(varcanon string, mkline MkLine) {
		switch {
		case src.vartypes.DefinedCanon(varcanon):
			// Already defined, can also be a tool.

		case hasPrefix(varcanon, "_"):
			// Variables starting with an underscore are reserved for the
			// infrastructure and are not available for use by packages.

		case contains(varcanon, "$"):
			// Indirect, but not the usual parameterized form. Variables of
			// this form should not be unintentionally visible from outside
			// the infrastructure since they don't follow the pkgsrc naming
			// conventions.

		case hasSuffix(varcanon, "_MK"):
			// Multiple-inclusion guards are internal to the infrastructure.

		default:
			if trace.Tracing {
				trace.Stepf("Untyped variable %q in %s", varcanon, mkline)
			}
			src.vartypes.DefineType(varcanon, &unknownType)
		}
	}

	handleMkFile := func(path string) {
		mklines := LoadMk(path, 0)
		if mklines != nil && len(mklines.mklines) > 0 {
			mklines.collectDefinedVariables()
			mklines.collectUsedVariables()
			for varname, mkline := range mklines.vars.firstDef {
				define(varnameCanon(varname), mkline)
			}
			for varname, mkline := range mklines.vars.used {
				define(varnameCanon(varname), mkline)
			}
		}
	}

	handleFile := func(pathName string, info os.FileInfo, err error) error {
		baseName := info.Name()
		if hasSuffix(baseName, ".mk") || baseName == "mk.conf" {
			handleMkFile(filepath.ToSlash(pathName))
		}
		return nil
	}

	err := filepath.Walk(src.File("mk"), handleFile)
	G.AssertNil(err, "Walk error in pkgsrc infrastructure")
}

func (src *Pkgsrc) parseSuggestedUpdates(lines Lines) []SuggestedUpdate {
	if lines == nil {
		return nil
	}

	var updates []SuggestedUpdate
	state := 0
	for _, line := range lines.Lines {
		text := line.Text

		// TODO: Replace this state transition scheme with explicit code,
		//  hoping that the code will be easier to understand.
		if state == 0 && text == "Suggested package updates" {
			state = 1
		} else if state == 1 && text == "" {
			state = 2
		} else if state == 2 {
			state = 3
		} else if state == 3 && text == "" {
			state = 4
		}

		if state == 3 {
			if m, pkgname, comment := match2(text, `^\to[\t ]([^\t ]+)(?:[\t ]*(.+))?$`); m {
				if m, pkgbase, pkgversion := match2(pkgname, rePkgname); m {
					updates = append(updates, SuggestedUpdate{line.Location, intern(pkgbase), intern(pkgversion), intern(comment)})
				} else {
					line.Warnf("Invalid package name %q.", pkgname)
				}
			} else {
				line.Warnf("Invalid line format %q.", text)
			}
		}
	}
	return updates
}

func (src *Pkgsrc) loadSuggestedUpdates() {
	src.suggestedUpdates = src.parseSuggestedUpdates(Load(src.File("doc/TODO"), MustSucceed))
	src.suggestedWipUpdates = src.parseSuggestedUpdates(Load(src.File("wip/TODO"), NotEmpty))
}

func (src *Pkgsrc) loadDocChangesFromFile(filename string) []*Change {

	warn := !G.Wip

	parseChange := func(line Line) *Change {
		lex := textproc.NewLexer(line.Text)

		space := lex.NextHspace()
		if space == "" {
			return nil
		}

		if space != "\t" {
			if warn {
				line.Warnf("Package changes should be indented using a single tab, not %q.", space)
				line.Explain(
					"To avoid this formatting mistake in the future, just run",
					sprintf("%q", bmake("cce")),
					"after committing the update to the package.")
			}

			return nil
		}

		f := strings.Fields(lex.Rest())
		n := len(f)
		if n != 4 && n != 6 {
			return nil
		}

		action, pkgpath, author, date := f[0], f[1], f[len(f)-2], f[len(f)-1]
		if !hasPrefix(author, "[") || !hasSuffix(date, "]") {
			return nil
		}
		author, date = author[1:], date[:len(date)-1]

		newChange := func(version string) *Change {
			return &Change{
				Location: line.Location,
				Action:   intern(action),
				Pkgpath:  intern(pkgpath),
				Version:  intern(version),
				Author:   intern(author),
				Date:     intern(date),
			}
		}

		switch {
		case action == "Added" && f[2] == "version" && n == 6:
			return newChange(f[3])

		case (action == "Updated" || action == "Downgraded") && f[2] == "to" && n == 6:
			return newChange(f[3])

		case action == "Removed" && (n == 6 && f[2] == "successor" || n == 4):
			return newChange("")

		case (action == "Renamed" || action == "Moved") && f[2] == "to" && n == 6:
			return newChange("")
		}

		line.Warnf("Unknown doc/CHANGES line: %s", line.Text)
		line.Explain(
			"See mk/misc/developer.mk for the rules.")

		return nil
	}

	// Each date in the file should be from the same year as the filename says.
	// This check has been added in 2018.
	// For years earlier than 2018 pkglint doesn't care because it's not a big issue anyway.
	year := ""
	if m, yyyy := match1(filename, `-(\d+)$`); m && yyyy >= "2018" {
		year = yyyy
	}

	infra := false
	lines := Load(filename, MustSucceed|NotEmpty)
	var changes []*Change
	for _, line := range lines.Lines {

		if hasPrefix(line.Text, "\tmk/") {
			infra = true
		}
		if infra {
			if hasSuffix(line.Text, "]") {
				infra = false
			}
			continue
		}

		change := parseChange(line)
		if change == nil {
			continue
		}

		changes = append(changes, change)

		if !warn {
			continue
		}

		if year != "" && change.Date[0:4] != year {
			line.Warnf("Year %s for %s does not match the filename %s.",
				change.Date[0:4], change.Pkgpath, filename)
		}

		if len(changes) >= 2 && year != "" {
			if prev := changes[len(changes)-2]; change.Date < prev.Date {
				line.Warnf("Date %s for %s is earlier than %s in %s.",
					change.Date, change.Pkgpath, prev.Date, line.RefToLocation(prev.Location))
				line.Explain(
					"The entries in doc/CHANGES should be in chronological order, and",
					"all dates are assumed to be in the UTC timezone, to prevent time",
					"warps.",
					"",
					"To fix this, determine which of the involved dates are correct",
					"and which aren't.",
					"",
					"To prevent this kind of mistakes in the future,",
					"make sure that your system time is correct and run",
					sprintf("%q", bmake("cce")),
					"to commit the changes entry.")
			}
		}
	}

	return changes
}

func (src *Pkgsrc) SuggestedUpdates() []SuggestedUpdate {
	if G.Wip {
		return src.suggestedWipUpdates
	} else {
		return src.suggestedUpdates
	}
}

func (src *Pkgsrc) loadDocChanges() {
	docDir := src.File("doc")
	files := src.ReadDir("doc")
	if len(files) == 0 {
		NewLineWhole(docDir).Fatalf("Cannot be read for loading the package changes.")
	}

	var filenames []string
	for _, file := range files {
		filename := file.Name()
		if matches(filename, `^CHANGES-20\d\d$`) && filename >= "CHANGES-2011" { // TODO: Why 2011?
			filenames = append(filenames, filename)
		}
	}

	sort.Strings(filenames)
	src.LastChange = make(map[string]*Change)
	for _, filename := range filenames {
		changes := src.loadDocChangesFromFile(docDir + "/" + filename)
		for _, change := range changes {
			src.LastChange[change.Pkgpath] = change
		}
	}
}

func (src *Pkgsrc) loadUserDefinedVars() {
	mklines := src.LoadMk("mk/defaults/mk.conf", MustSucceed|NotEmpty)

	for _, mkline := range mklines.mklines {
		if mkline.IsVarassign() || mkline.IsCommentedVarassign() {
			src.UserDefinedVars.Define(mkline.Varname(), mkline)
		}
	}
}

func (src *Pkgsrc) initDeprecatedVars() {
	src.Deprecated = map[string]string{
		// December 2003
		"FIX_RPATH": "It has been removed from pkgsrc in 2003.",

		// February 2005
		"LIB_DEPENDS":    "Use DEPENDS instead.",
		"ONLY_FOR_ARCHS": "Use ONLY_FOR_PLATFORM instead.",
		"NOT_FOR_ARCHS":  "Use NOT_FOR_PLATFORM instead.",
		"ONLY_FOR_OPSYS": "Use ONLY_FOR_PLATFORM instead.",
		"NOT_FOR_OPSYS":  "Use NOT_FOR_PLATFORM instead.",

		// May 2005
		"ALL_TARGET":         "Use BUILD_TARGET instead.",
		"DIGEST_FILE":        "Use DISTINFO_FILE instead.",
		"IGNORE":             "Use PKG_FAIL_REASON or PKG_SKIP_REASON instead.",
		"IS_INTERACTIVE":     "Use INTERACTIVE_STAGE instead.",
		"KERBEROS":           "Use the PKG_OPTIONS framework instead.",
		"MASTER_SITE_SUBDIR": "Use some form of MASTER_SITES instead.",
		"MD5_FILE":           "Use DISTINFO_FILE instead.",
		"MIRROR_DISTFILE":    "Use NO_BIN_ON_FTP and/or NO_SRC_ON_FTP instead.",
		"NO_CDROM":           "Use NO_BIN_ON_CDROM and/or NO_SRC_ON_CDROM instead.",
		"NO_PATCH":           "You can just remove it.",
		"NO_WRKSUBDIR":       "Use WRKSRC=${WRKDIR} instead.",
		"PATCH_SITE_SUBDIR":  "Use some form of PATCHES_SITES instead.",
		"PATCH_SUM_FILE":     "Use DISTINFO_FILE instead.",
		"PKG_JVM":            "Use PKG_DEFAULT_JVM instead.",
		"USE_BUILDLINK2":     "You can just remove it.",
		"USE_BUILDLINK3":     "You can just remove it.",
		"USE_CANNA":          "Use the PKG_OPTIONS framework instead.",
		"USE_DB4":            "Use the PKG_OPTIONS framework instead.",
		"USE_DIRS":           "You can just remove it.",
		"USE_ESOUND":         "Use the PKG_OPTIONS framework instead.",
		"USE_GIF":            "Use the PKG_OPTIONS framework instead.",
		"USE_GMAKE":          "Use USE_TOOLS+=gmake instead.",
		"USE_GNU_TOOLS":      "Use USE_TOOLS instead.",
		"USE_IDEA":           "Use the PKG_OPTIONS framework instead.",
		"USE_LIBCRACK":       "Use the PKG_OPTIONS framework instead.",
		"USE_MMX":            "Use the PKG_OPTIONS framework instead.",
		"USE_PKGLIBTOOL":     "Use USE_LIBTOOL instead.",
		"USE_SSL":            "Include \"../../security/openssl/buildlink3.mk\" instead.",

		// July 2005
		"USE_PERL5": "Use USE_TOOLS+=perl or USE_TOOLS+=perl:run instead.",

		// October 2005
		"NO_TOOLS":   "You can just remove it.",
		"NO_WRAPPER": "You can just remove it.",

		// November 2005
		"ALLFILES":       "Use CKSUMFILES instead.",
		"DEPENDS_TARGET": "Use DEPENDS instead.",
		"FETCH_DEPENDS":  "Use DEPENDS instead.",
		"RUN_DEPENDS":    "Use DEPENDS instead.",

		// December 2005
		"USE_CUPS":     "Use the PKG_OPTIONS framework (option cups) instead.",
		"USE_I586":     "Use the PKG_OPTIONS framework (option i586) instead.",
		"USE_INN":      "Use the PKG_OPTIONS framework instead.",
		"USE_OPENLDAP": "Use the PKG_OPTIONS framework (option openldap) instead.",
		"USE_OSS":      "Use the PKG_OPTIONS framework (option oss) instead.",
		"USE_RSAREF2":  "Use the PKG_OPTIONS framework (option rsaref) instead.",
		"USE_SASL":     "Use the PKG_OPTIONS framework (option sasl) instead.",
		"USE_SASL2":    "Use the PKG_OPTIONS framework (option sasl) instead.",
		"USE_SJ3":      "Use the PKG_OPTIONS framework (option sj3) instead.",
		"USE_SOCKS":    "Use the PKG_OPTIONS framework (socks4 and socks5 options) instead.",
		"USE_WNN4":     "Use the PKG_OPTIONS framework (option wnn4) instead.",
		"USE_XFACE":    "Use the PKG_OPTIONS framework instead.",

		// February 2006
		"TOOLS_DEPMETHOD":     "Use the :build or :run modifiers in USE_TOOLS instead.",
		"MANDIR":              "Please use ${PREFIX}/${PKGMANDIR} instead.",
		"DOWNLOADED_DISTFILE": "Use the shell variable $$extract_file instead.",
		"DECOMPRESS_CMD":      "Use EXTRACT_CMD instead.",

		// March 2006
		"INSTALL_EXTRA_TMPL":   "Use INSTALL_TEMPLATE instead.",
		"DEINSTALL_EXTRA_TMPL": "Use DEINSTALL_TEMPLATE instead.",

		// April 2006
		"RECOMMENDED":        "Use ABI_DEPENDS instead.",
		"BUILD_USES_MSGFMT":  "Use USE_TOOLS+=msgfmt instead.",
		"USE_MSGFMT_PLURALS": "Use USE_TOOLS+=msgfmt instead.",

		// May 2006
		"EXTRACT_USING_PAX":       "Use \"EXTRACT_OPTS=-t pax\" instead.",
		"NO_EXTRACT":              "It doesn't exist anymore.",
		"_FETCH_MESSAGE":          "Use FETCH_MESSAGE (different format) instead.",
		"BUILDLINK_DEPENDS.*":     "Use BUILDLINK_API_DEPENDS.* instead.",
		"BUILDLINK_RECOMMENDED.*": "Use BUILDLINK_ABI_DEPENDS.* instead.",
		"SHLIB_HANDLING":          "Use CHECK_SHLIBS_SUPPORTED instead.",
		"USE_RMAN":                "It has been removed.",

		// June 2006
		"DEINSTALL_SRC":      "Use the pkginstall framework instead.",
		"INSTALL_SRC":        "Use the pkginstall framework instead.",
		"DEINSTALL_TEMPLATE": "Use DEINSTALL_TEMPLATES instead.",
		"INSTALL_TEMPLATE":   "Use INSTALL_TEMPLATES instead.",
		"HEADER_TEMPLATE":    "Use HEADER_TEMPLATES instead.",
		"_REPLACE.*":         "Use REPLACE.* instead.",
		"_REPLACE_FILES.*":   "Use REPLACE_FILES.* instead.",
		"MESSAGE":            "Use MESSAGE_SRC instead.",
		"INSTALL_FILE":       "It may only be used internally by pkgsrc.",
		"DEINSTALL_FILE":     "It may only be used internally by pkgsrc.",

		// July 2006
		"USE_DIGEST":           "You can just remove it.",
		"LTCONFIG_OVERRIDE":    "You can just remove it.",
		"USE_GNU_GETTEXT":      "You can just remove it.",
		"BUILD_ENV":            "Use PKGSRC_MAKE_ENV instead.",
		"DYNAMIC_MASTER_SITES": "You can just remove it.",

		// September 2006
		"MAKEFILE": "Use MAKE_FILE instead.",

		// November 2006
		"SKIP_PORTABILITY_CHECK": "Use CHECK_PORTABILITY_SKIP (a list of patterns) instead.",

		// January 2007
		"BUILDLINK_TRANSFORM.*": "Use BUILDLINK_FNAME_TRANSFORM.* instead.",

		// March 2007
		"SCRIPTDIR":       "You can just remove it.",
		"NO_PKG_REGISTER": "You can just remove it.",
		"NO_DEPENDS":      "You can just remove it.",

		// October 2007
		"_PKG_SILENT": "Use RUN (with more error checking) instead.",
		"_PKG_DEBUG":  "Use RUN (with more error checking) instead.",
		"LICENCE":     "Use LICENSE instead.",

		// November 2007
		//USE_NCURSES: Include "../../devel/ncurses/buildlink3.mk" instead.

		// December 2007
		"INSTALLATION_DIRS_FROM_PLIST": "Use AUTO_MKDIRS instead.",

		// April 2009
		"NO_PACKAGE": "It doesn't exist anymore.",
		"NO_MTREE":   "You can just remove it.",

		// July 2012
		"SETGIDGAME": "Use USE_GAMESGROUP instead.",
		"GAMEGRP":    "Use GAMES_GROUP instead.",
		"GAMEOWN":    "Use GAMES_USER instead.",

		// July 2013
		"USE_GNU_READLINE": "Include \"../../devel/readline/buildlink3.mk\" instead.",

		// October 2014
		"SVR4_PKGNAME":           "Just remove it.",
		"PKG_INSTALLATION_TYPES": "Just remove it.",

		// November 2015, commit abccb56
		"EVAL_PREFIX": "All packages are installed in PREFIX now.",

		// January 2016
		"SUBST_POSTCMD.*": "Has been removed, as it seemed unused.",

		// June 2016
		"USE_CROSSBASE": "Has been removed.",
	}
}

// Load loads the file relative to the pkgsrc top directory.
func (src *Pkgsrc) Load(filename string, options LoadOptions) Lines {
	return Load(src.File(filename), options)
}

// LoadMk loads the Makefile relative to the pkgsrc top directory.
func (src *Pkgsrc) LoadMk(filename string, options LoadOptions) MkLines {
	return LoadMk(src.File(filename), options)
}

// ReadDir reads the file listing from the given directory (relative to the pkgsrc root),
// filtering out any ignored files (CVS/*) and empty directories.
func (src *Pkgsrc) ReadDir(dirName string) []os.FileInfo {
	dir := src.File(dirName)
	files, err := ioutil.ReadDir(dir)
	if err != nil {
		return nil
	}

	var relevantFiles []os.FileInfo
	for _, dirent := range files {
		name := dirent.Name()
		if !dirent.IsDir() || !isIgnoredFilename(name) && !isEmptyDir(dir+"/"+name) {
			relevantFiles = append(relevantFiles, dirent)
		}
	}

	return relevantFiles
}

// File resolves a filename relative to the pkgsrc top directory.
//
// Example:
//  NewPkgsrc("/usr/pkgsrc").File("distfiles") => "/usr/pkgsrc/distfiles"
func (src *Pkgsrc) File(relativeName string) string {
	// TODO: Package.File resolves variables, Pkgsrc.File doesn't. They should behave the same.
	return cleanpath(src.topdir + "/" + relativeName)
}

// ToRel returns the path of `filename`, relative to the pkgsrc top directory.
//
// Example:
//  NewPkgsrc("/usr/pkgsrc").ToRel("/usr/pkgsrc/distfiles") => "distfiles"
func (src *Pkgsrc) ToRel(filename string) string {
	return relpath(src.topdir, filename)
}

// IsInfra returns whether the given filename (relative to the pkglint
// working directory) is part of the pkgsrc infrastructure.
func (src *Pkgsrc) IsInfra(filename string) bool {
	rel := src.ToRel(filename)
	return hasPrefix(rel, "mk/") || hasPrefix(rel, "wip/mk/")
}

func (src *Pkgsrc) addBuildDefs(varnames ...string) {
	for _, varname := range varnames {
		src.buildDefs[varname] = true
	}
}

// IsBuildDef returns whether the given variable is automatically added
// to BUILD_DEFS by the pkgsrc infrastructure. In such a case, the
// package doesn't need to add the variable to BUILD_DEFS itself.
func (src *Pkgsrc) IsBuildDef(varname string) bool {
	return src.buildDefs[varname]
}

func (src *Pkgsrc) loadMasterSites() {
	mklines := src.LoadMk("mk/fetch/sites.mk", MustSucceed|NotEmpty)

	for _, mkline := range mklines.mklines {
		if mkline.IsVarassign() {
			varname := mkline.Varname()
			// TODO: Give a plausible reason for the MASTER_SITE_BACKUP exception.
			if hasPrefix(varname, "MASTER_SITE_") && varname != "MASTER_SITE_BACKUP" {
				for _, url := range mkline.ValueFields(mkline.Value()) {
					if matches(url, `^(?:http://|https://|ftp://)`) {
						src.registerMasterSite(varname, url)
					}
				}

				// TODO: register variable type, to avoid redundant definitions in vardefs.go.
			}
		}
	}

	// Explicitly allowed, although not defined in mk/fetch/sites.mk.
	// TODO: Document where this definition comes from and why it is good.
	src.registerMasterSite("MASTER_SITE_LOCAL", "ftp://ftp.NetBSD.org/pub/pkgsrc/distfiles/LOCAL_PORTS/")

	if trace.Tracing {
		trace.Stepf("Loaded %d MASTER_SITE_* URLs.", len(src.MasterSiteURLToVar))
	}
}

func (src *Pkgsrc) registerMasterSite(varname, url string) {
	nameToURL := src.MasterSiteVarToURL
	urlToName := src.MasterSiteURLToVar

	if nameToURL[varname] == "" {
		nameToURL[varname] = url
	}
	urlToName[url] = varname
}

func (src *Pkgsrc) loadPkgOptions() {
	lines := src.Load("mk/defaults/options.description", MustSucceed)

	for _, line := range lines.Lines {
		if m, name, description := match2(line.Text, `^([-0-9a-z_+]+)(?:[\t ]+(.*))?$`); m {
			src.PkgOptions[name] = description
		} else {
			line.Errorf("Invalid line format: %s", line.Text)
		}
	}
}

// VariableType returns the type of the variable
// (possibly guessed based on the variable name),
// or nil if the type cannot even be guessed.
func (src *Pkgsrc) VariableType(mklines MkLines, varname string) (vartype *Vartype) {
	if trace.Tracing {
		defer trace.Call(varname, trace.Result(&vartype))()
	}

	// When scanning mk/** for otherwise unknown variables, their type
	// is set to BtUnknown. These variables must not override the guess
	// based on the variable name.
	vartype = src.vartypes.Canon(varname)
	if vartype != nil && vartype.basicType != BtUnknown {
		return vartype
	}

	if tool := G.ToolByVarname(mklines, varname); tool != nil {
		if trace.Tracing {
			trace.Stepf("Use of tool %+v", tool)
		}
		perms := aclpUse
		if tool.Validity == AfterPrefsMk && mklines.Tools.SeenPrefs {
			perms |= aclpUseLoadtime
		}
		return &Vartype{BtShellCommand, NoVartypeOptions, []ACLEntry{{"*", perms}}}
	}

	if m, toolVarname := match1(varname, `^TOOLS_(.*)`); m {
		if tool := G.ToolByVarname(mklines, toolVarname); tool != nil {
			return &Vartype{BtPathname, NoVartypeOptions, []ACLEntry{{"*", aclpUse}}}
		}
	}

	return src.guessVariableType(varname)
}

func (src *Pkgsrc) guessVariableType(varname string) (vartype *Vartype) {
	allowAll := []ACLEntry{{"*", aclpAll}}
	allowRuntime := []ACLEntry{{"*", aclpAllRuntime}}

	// Guess the data type of the variable based on naming conventions.
	varbase := varnameBase(varname)
	var gtype *Vartype
	switch {
	case hasSuffix(varbase, "DIRS"):
		gtype = &Vartype{BtPathmask, List | Guessed, allowRuntime}
	case hasSuffix(varbase, "DIR") && !hasSuffix(varbase, "DESTDIR"), hasSuffix(varname, "_HOME"):
		// TODO: hasSuffix(varbase, "BASE")
		gtype = &Vartype{BtPathname, Guessed, allowRuntime}
	case hasSuffix(varbase, "FILES"):
		gtype = &Vartype{BtPathmask, List | Guessed, allowRuntime}
	case hasSuffix(varbase, "FILE"):
		gtype = &Vartype{BtPathname, Guessed, allowRuntime}
	case hasSuffix(varbase, "PATH"):
		gtype = &Vartype{BtPathlist, Guessed, allowRuntime}
	case hasSuffix(varbase, "PATHS"):
		gtype = &Vartype{BtPathname, List | Guessed, allowRuntime}
	case hasSuffix(varbase, "_USER"):
		gtype = &Vartype{BtUserGroupName, Guessed, allowAll}
	case hasSuffix(varbase, "_GROUP"):
		gtype = &Vartype{BtUserGroupName, Guessed, allowAll}
	case hasSuffix(varbase, "_ENV"):
		gtype = &Vartype{BtShellWord, List | Guessed, allowRuntime}
	case hasSuffix(varbase, "_CMD"):
		gtype = &Vartype{BtShellCommand, Guessed, allowRuntime}
	case hasSuffix(varbase, "_ARGS"):
		gtype = &Vartype{BtShellWord, List | Guessed, allowRuntime}
	case hasSuffix(varbase, "_CFLAGS"), hasSuffix(varname, "_CPPFLAGS"), hasSuffix(varname, "_CXXFLAGS"):
		gtype = &Vartype{BtCFlag, List | Guessed, allowRuntime}
	case hasSuffix(varname, "_LDFLAGS"):
		gtype = &Vartype{BtLdFlag, List | Guessed, allowRuntime}
	case hasSuffix(varbase, "_MK"):
		// TODO: Add BtGuard for inclusion guards, since these variables may only be checked using defined().
		gtype = &Vartype{BtUnknown, Guessed, allowAll}
	case hasSuffix(varbase, "_SKIP"):
		gtype = &Vartype{BtPathmask, List | Guessed, allowRuntime}
	}

	if gtype == nil {
		vartype = src.vartypes.Canon(varname)
		if vartype != nil {
			return vartype
		}
	}

	if trace.Tracing {
		if gtype != nil {
			trace.Step2("The guessed type of %q is %q.", varname, gtype.String())
		} else {
			trace.Step1("No type definition found for %q.", varname)
		}
	}
	return gtype
}

// Change describes a modification to a single package, from the doc/CHANGES-* files.
type Change struct {
	Location Location
	Action   string
	Pkgpath  string
	Version  string
	Author   string
	Date     string
}

// SuggestedUpdate describes a desired package update, from the doc/TODO file.
type SuggestedUpdate struct {
	Line    Location
	Pkgname string
	Version string
	Comment string
}
