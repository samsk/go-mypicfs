package ruleFS;
////
// Regexp rules
//

import (
//	"fmt"
//	"path/filepath"
	"regexp"

	picfs "github.com/samsk/go-mypicfs/fs"

//	"github.com/hanwen/go-fuse/fuse"
//	"github.com/hanwen/go-fuse/fuse/pathfs"
//	"github.com/hanwen/go-fuse/fuse/nodefs"
);

type RegexpRuleData struct {
	Ident int
	StringSubmatch []string
}

type RegexpRuleContext struct {
	data RegexpRuleData
}

type regexpRule struct {
	fs picfs.ContextFS
	regex *regexp.Regexp
	ident int
}

// RegexpRuleMatch definition
type RegexpRuleMatch struct {
	rules []regexpRule
	def regexpRule
}

// RegexpRuleMatch implementation
func NewRegexpRules() (RegexpRuleMatch) {
	var _ picfs.RuleFSMatch = (*RegexpRuleMatch)(nil)
	return RegexpRuleMatch {};
}

func (match *RegexpRuleMatch) Match(regexMatch string, identIn int, fsIn picfs.ContextFS) {
	match.rules = append(match.rules, regexpRule {
		fs: fsIn,
		regex: regexp.MustCompile(regexMatch),
		ident: identIn,
	})
}

func (match *RegexpRuleMatch) MatchDefault(identIn int, fsIn picfs.ContextFS) {
	match.def = regexpRule {
		fs: fsIn,
		ident: identIn,
	}
}


func (match *RegexpRuleMatch) MatchPath(name string) (picfs.ContextFS, picfs.Context) {
	var res []string
	var data RegexpRuleData

	// TODO: what about cache
	data.Ident = -1
	data.StringSubmatch = nil

//	fmt.Printf(">> MatchPath(%s)\n", name);
	for _, element := range match.rules {
		res = element.regex.FindStringSubmatch(name)
		//fmt.Printf(">> MatchPath(%s) = %v => %d\n", name, res, element.ident);
		if (res != nil) {
			data.Ident = element.ident
			data.StringSubmatch = res

			return element.fs, data
		}
	}

	if (match.def.fs != nil) {
		return match.def.fs, data
	}

	return nil, nil;
}
