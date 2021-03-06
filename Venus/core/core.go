package venus

import (
  "regexp"
  "./strings"
)

type Matcher struct {
  match string
  replace string
}

func NewMatcher(match string, replace string) *Matcher {
  m := new(Matcher)
  m.match = match
  m.replace = replace
  return m
}

func (m Matcher) Replace(input string) string {
  r:= regexp.MustCompile(m.match)
  return r.ReplaceAllString(input, m.replace)
}

const VARIABLENAME = "[a-zA-Z_]+(?:[a-zA-Z0-9_]*)"
const FUNCTIONNAME = "[a-zA-Z_:.]+(?:[a-zA-Z0-9_:.]*)"
const TYPEKEY = "(num|str|bool|table|tbl)"
const CONDITION = "((?:(?:\\S+(?: |\\t)*?(?:<|>|<=|>=|==|!=)?(?: |\\t)*?\\S*)(?:(?: |\\t)*?(?:&&|\\|\\|)?(?: |\\t)*?)?)+)"

var matchers []Matcher
func InitMatchers() {

  /* BIG THINGS */
  //Typed Variables (remove type)
  matchers = append(matchers, *NewMatcher(TYPEKEY+"\\s*("+VARIABLENAME+")", "$2"))
  //Functions
  matchers = append(matchers, *NewMatcher("fn\\s("+FUNCTIONNAME+")\\s*(\\(.*\\))\\s*"+TYPEKEY+"?\\s*{", "function $1$2"))
  //Anonymous Functions
  matchers = append(matchers, *NewMatcher("fn\\s*(\\(.*\\))\\s*"+TYPEKEY+"?\\s*{", "function $1"))
  //If statements
  matchers = append(matchers, *NewMatcher("if\\s+"+CONDITION+"\\s*{", "if ($1) then"))
  //While loops
  matchers = append(matchers, *NewMatcher("while\\s+"+CONDITION+"\\s*{", "while ($1) do"))
  //Loop loops
  matchers = append(matchers, *NewMatcher("loop\\s*{", "while (true) do"))
  //Numerical For
  matchers = append(matchers, *NewMatcher("for\\s+("+VARIABLENAME+")\\s+=\\s+(.*),\\s*(.*)\\s*{", "for $1 = $2, $3 do"))
  //Indexed For
  matchers = append(matchers, *NewMatcher("for\\s+("+VARIABLENAME+")\\s+in\\s+(.*)\\s+{", "for $1,_ in ipairs($2) do"))
  //Foreach
  matchers = append(matchers, *NewMatcher("foreach\\s+("+VARIABLENAME+")\\s+in\\s+(.*)\\s+{", "for $1,_ in pairs($2) do"))
  //Index Accessor
  matchers = append(matchers, *NewMatcher("([^\\s\\n]*)\\[(.*)&&(.*)\\]", "$1[$2&$1$3]"))
  //Increaser Accessor
  matchers = append(matchers, *NewMatcher("([^\\s\\n]*)\\[\\+\\]", "$1[&$1+1]"))

  //Closing Brackets
  matchers = append(matchers, *NewMatcher("}", "end"))

  //Else's
  matchers = append(matchers, *NewMatcher("else\\s*{", "else"))
  matchers = append(matchers, *NewMatcher("end\\s*else", "else"))

  //Else If
  matchers = append(matchers, *NewMatcher("", ""))

  //([^\\s\\n]*)\\[(.*)##(.*)\\]

  /* SMALL THINGS */
  //@
  matchers = append(matchers, *NewMatcher("@("+VARIABLENAME+")", "self.$1"))
  //@.
  matchers = append(matchers, *NewMatcher("@(\\.|:|\\s|\\))", "self$1"))
  //+= and -=
  matchers = append(matchers, *NewMatcher("("+VARIABLENAME+")\\s*(\\+|-)=", "$1 = $1 $2 "))
  //++ and --
  matchers = append(matchers, *NewMatcher("("+VARIABLENAME+")\\s*(?:(\\+)\\+|(-)-)", "$1 = $1 $2$3 1"))
  //!= means ~=
  matchers = append(matchers, *NewMatcher("!=", "~="))
  //Comments
  matchers = append(matchers, *NewMatcher("#", "--"))
  //Block Comments
  matchers = append(matchers, *NewMatcher("###([\\s\\S]*)###", "--[[ $1 --]]"))
  //& means #
  matchers = append(matchers, *NewMatcher("&", "#"))
}

var postMatchers []Matcher
func InitPostMatchers() {
  //Object Literals
  postMatchers = append(postMatchers, *NewMatcher("<", "{"))
  postMatchers = append(postMatchers, *NewMatcher(">", "}"))
}



func DoMatchers(input string, matches []Matcher) string {
  for _, matcher := range matches {
    input = matcher.Replace(input)
  }
  return input
}

func Translate(input string) (string, error) {

  //Hide Strings (also turns multis into [[]])
  input, strs := vstrings.HideStrings(input)

  //Replace Stuff
  InitMatchers()
  input = DoMatchers(input, matchers)

  //Get strings back
  input, err := vstrings.ShowStrings(input, strs)
  if err != nil {return "", err}

  //Replace Stuff After we get strings back
  InitPostMatchers()
  input = DoMatchers(input, postMatchers)

  //Return
  return input, nil
}
