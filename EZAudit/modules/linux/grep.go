package linux

import (
	"Just-Go-IT/EZAudit/artifact"
	"Just-Go-IT/EZAudit/global"
	"Just-Go-IT/EZAudit/interact"
	"Just-Go-IT/EZAudit/registry"
	"errors"
	"strconv"
)

func init() {
	registry.Register("grep", &grep{}, false, registry.Linux)
	registry.Register("egrep", &grep{}, false, registry.Linux)
	registry.Register("fgrep", &grep{}, false, registry.Linux)
	registry.Register("rgrep", &grep{}, false, registry.Linux)
	registry.Register("pgrep", &grep{}, false, registry.Linux)
}

type grep struct {
	command       string
	admin         bool
	searchPattern string
	patternSyntax struct {
		basicRegexp  bool // -G
		extendRegexp bool // -E
		fixedStrings bool // -F
		perlRegexp   bool // -P
	}
	matchingControl struct {
		patterns     string // -e
		patternFile  string // -f
		ignoreCase   bool   // -i
		noIgnoreCase bool   // --no-ignore-case
		invertMatch  bool   // -v
		wordRegexp   bool   // -w
		lineRegexp   bool   // -x
	}
	outputControl struct {
		count               bool // -c
		filesWithoutMatches bool // -L
		filesWithMatches    bool // -l
		maxCount            int  // -m
		onlyMatching        bool // -o
		recursive           bool // -r
	}
	target string
}

func (g grep) New(p map[string]interface{}) (global.Module, error) {

	ok := false

	for key := range p {
		switch key {
		case "target":
			g.target, ok = p["target"].(string)
			if !ok {
				return nil, errors.New("the key \"target\" is set for the module. The value must be a \"string\"")
			}
		case "searchPattern":
			g.searchPattern, ok = p["searchPattern"].(string)
			if !ok {
				return nil, errors.New("the key \"searchPattern\" is set for the module. The value must be a \"string\"")
			}
		case "admin":
			g.admin, ok = p["admin"].(bool)
			if !ok {
				return nil, errors.New("the key \"admin\" is set but the value, which should be a bool could not be parsed")
			}
		case "patternSyntax":
			patternSyntax, ok := p["patternSyntax"].(map[string]interface{})
			if !ok {
				return nil, errors.New("the object \"patternSyntax\" could not be parsed")
			}
			for patternKey := range patternSyntax {
				switch patternKey {
				case "basicRegexp":
					g.patternSyntax.basicRegexp, ok = patternSyntax["basicRegexp"].(bool)
					if !ok {
						return nil, errors.New("the key \"basicRegexp\" is set but the value, which should be a \"bool\" could not be parsed")
					}
				case "extendRegexp":
					g.patternSyntax.extendRegexp, ok = patternSyntax["extendRegexp"].(bool)
					if !ok {
						return nil, errors.New("the key \"extendRegexp\" is set but the value, which should be a \"bool\" could not be parsed")
					}
				case "fixedStrings":
					g.patternSyntax.fixedStrings, ok = patternSyntax["fixedStrings"].(bool)
					if !ok {
						return nil, errors.New("the key \"fixedStrings\" is set but the value, which should be a \"bool\" could not be parsed")
					}
				case "perlRegexp":
					g.patternSyntax.perlRegexp, ok = patternSyntax["perlRegexp"].(bool)
					if !ok {
						return nil, errors.New("the key \"perlRegexp\" is set but the value, which should be a \"bool\" could not be parsed")
					}
				default:
					if patternKey != "basicRegexp" && patternKey != "extendRegexp" && patternKey != "fixedStrings" && patternKey != "perlRegexp" {
						return nil, errors.New("there is no key called: \"" + patternKey + "\" in the object \"patternSyntax\" in the module grep")
					}
				}
			}
		case "matchingControl":
			matchingControl, ok := p["matchingControl"].(map[string]interface{})
			if !ok {
				return nil, errors.New("the object \"matchingControl\" could not be parsed")
			}
			for matchingControlKey := range matchingControl {
				switch matchingControlKey {
				case "patterns":
					g.matchingControl.patterns, ok = matchingControl["patterns"].(string)
					if !ok {
						return nil, errors.New("the key \"patterns\" is set but the value, which should be a \"string\" could not be parsed")
					}
				case "patternFile":
					g.matchingControl.patternFile, ok = matchingControl["patternFile"].(string)
					if !ok {
						return nil, errors.New("the key \"patternFile\" is set but the value, which should be a \"string\" could not be parsed")
					}
				case "ignoreCase":
					g.matchingControl.ignoreCase, ok = matchingControl["ignoreCase"].(bool)
					if !ok {
						return nil, errors.New("the key \"ignoreCase\" is set but the value, which should be a \"bool\" could not be parsed")
					}
				case "noIgnoreCase":
					g.matchingControl.noIgnoreCase, ok = matchingControl["noIgnoreCase"].(bool)
					if !ok {
						return nil, errors.New("the key \"noIgnoreCase\" is set but the value, which should be a \"bool\" could not be parsed")
					}
				case "invertMatch":
					g.matchingControl.invertMatch, ok = matchingControl["invertMatch"].(bool)
					if !ok {
						return nil, errors.New("the key \"invertMatch\" is set but the value, which should be a \"bool\" could not be parsed")
					}
				case "wordRegexp":
					g.matchingControl.wordRegexp, ok = p["wordRegexp"].(bool)
					if !ok {
						return nil, errors.New("the key \"wordRegexp\" is set but the value, which should be a \"bool\" could not be parsed")
					}
				case "lineRegexp":
					g.matchingControl.lineRegexp, ok = p["lineRegexp"].(bool)
					if !ok {
						return nil, errors.New("the key \"lineRegexp\" is set but the value, which should be a \"bool\" could not be parsed")
					}
				default:
					if matchingControlKey != "patterns" && matchingControlKey != "patternFile" && matchingControlKey != "ignoreCase" && matchingControlKey != "noIgnoreCase" && matchingControlKey != "invertMatch" && matchingControlKey != "wordRegexp" {
						return nil, errors.New("there is no key called: \"" + matchingControlKey + "\" in the object \"matchingControl\" in the module grep")
					}
				}
			}
		case "outputControl":
			outputControl, ok := p["outputControl"].(map[string]interface{})
			if !ok {
				return nil, errors.New("the object \"outputControl\" could not be parsed")
			}
			for outputControlKey := range outputControl {
				switch outputControlKey {
				case "count":
					g.outputControl.count, ok = outputControl["count"].(bool)
					if !ok {
						return nil, errors.New("the key \"count\" is set but the value, which should be a \"bool\" could not be parsed")
					}
				case "filesWithoutMatches":
					g.outputControl.filesWithoutMatches, ok = outputControl["filesWithoutMatches"].(bool)
					if !ok {
						return nil, errors.New("the key \"filesWithoutMatches\" is set but the value, which should be a \"bool\" could not be parsed")
					}
				case "filesWithMatches":
					g.outputControl.filesWithMatches, ok = outputControl["filesWithMatches"].(bool)
					if !ok {
						return nil, errors.New("the key \"filesWithMatches\" is set but the value, which should be a \"bool\" could not be parsed")
					}
				case "maxCount":
					g.outputControl.maxCount, ok = outputControl["maxCount"].(int)
					if !ok {
						return nil, errors.New("the key \"noIgnoreCase\" is set but the value, which should be a \"int\" could not be parsed")
					}
				case "onlyMatching":
					g.outputControl.onlyMatching, ok = outputControl["onlyMatching"].(bool)
					if !ok {
						return nil, errors.New("the key \"onlyMatching\" is set but the value, which should be a \"bool\" could not be parsed")
					}
				case "recursive":
					g.outputControl.recursive, ok = outputControl["recursive"].(bool)
					if !ok {
						return nil, errors.New("the key \"recursive\" is set but the value, which should be a \"bool\" could not be parsed")
					}
				default:
					if outputControlKey != "count" && outputControlKey != "filesWithoutMatches" && outputControlKey != "filesWithMatches" && outputControlKey != "maxCount" && outputControlKey != "onlyMatching" && outputControlKey != "recursive" {
						return nil, errors.New("there is no key called: \"" + outputControlKey + "\" in the object \"outputControl\" in the module grep")
					}
				}
			}
		default:
			if key != "target" && key != "admin" && key != "patternSyntax" && key != "matchingControl" && key != "outputControl" && key != "searchPattern" {
				return nil, errors.New("there is no key called: \"" + key + "\" in the the module grep")
			}
		}
	}
	if g.matchingControl.ignoreCase && g.matchingControl.noIgnoreCase {
		return nil, errors.New("the keys  \"ignoreCase\" and \" noIgnoreCase\" can not be set at the same time")
	}

	if g.matchingControl.wordRegexp && g.matchingControl.lineRegexp {
		return nil, errors.New("the keys  \"wordRegexp\" and \" lineRegexp\" can not be set at the same time")
	}
	return &g, nil
}

func (g *grep) Execute(s *global.Step) (output string, err error) {
	g.command = ""
	if g.admin {
		g.command += "sudo "
	}

	// Which Module
	if s.ModuleName == "egrep" {
		g.command += "egrep"
		g.patternSyntax.extendRegexp = false
	} else if s.ModuleName == "fgrep" {
		g.command += "fgrep"
		g.patternSyntax.fixedStrings = false
	} else if s.ModuleName == "rgrep" {
		g.command += "rgrep"
		g.outputControl.recursive = false
	} else if s.ModuleName == "pgrep" {
		g.command += "pgrep"
		g.patternSyntax.perlRegexp = false
	} else {
		g.command += "grep"
	}
	// PatternSyntax
	if g.patternSyntax.extendRegexp {
		g.command += " -E"
	}
	if g.patternSyntax.fixedStrings {
		g.command += " -F"
	}
	if g.patternSyntax.basicRegexp {
		g.command += " -G"
	}
	if g.patternSyntax.perlRegexp {
		g.command += " -P"
	}
	// Matching Control
	if g.matchingControl.patterns != "" {
		g.command += " -e \"" + g.matchingControl.patterns + "\""
	}
	if g.matchingControl.patternFile != "" {
		g.command += " -f " + g.matchingControl.patternFile
	}

	if g.matchingControl.ignoreCase {
		g.command += " -i"
	}
	if g.matchingControl.noIgnoreCase {
		g.command += " --no-ignore-case"
	}
	if g.matchingControl.invertMatch {
		g.command += " -v"
	}
	if g.matchingControl.wordRegexp {
		g.command += " -w"
	}
	if g.matchingControl.lineRegexp {
		g.command += " -x"
	}

	// Output Control
	if g.outputControl.count {
		g.command += " -c"
	}
	if g.outputControl.recursive {
		g.command += " -r"
	}
	if g.outputControl.filesWithoutMatches {
		g.command += " -L"
	}
	if g.outputControl.filesWithMatches {
		g.command += " -l"
	}
	if g.outputControl.onlyMatching {
		g.command += " -o"
	}
	if g.outputControl.maxCount != 0 {
		g.command += " -m " + strconv.Itoa(g.outputControl.maxCount)
	}
	if g.searchPattern != "" {
		g.command += " " + g.searchPattern
	}

	g.command += " " + g.target

	output, err = interact.ShellPipe(g.command, s)
	if err != nil {
		if err.Error() != "exit status 1" && output != "" {
			return
		}
		err = nil
	}

	if g.target != "" {
		artifact.SaveFile(g.target, *s)
	}

	artifact.SaveString(output, *s)

	return
}
