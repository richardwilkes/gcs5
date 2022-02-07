package main

//go:generate go run enumgen.go

import (
	"bytes"
	"fmt"
	"go/format"
	"io"
	"os"
	"path/filepath"
	"strings"
	"text/template"
	"unicode"

	"github.com/richardwilkes/toolbox/log/jot"
	"github.com/richardwilkes/toolbox/txt"
)

const (
	rootDir   = ".."
	genSuffix = "_gen.go"
)

type enumValue struct {
	Name   string
	Key    string
	String string
	Short  string
}

type enumInfo struct {
	Pkg        string
	Name       string
	Desc       string
	StandAlone bool
	Values     []enumValue
}

func main() {
	const (
		enumTmpl = "enum.go.tmpl"
	)
	removeExistingGenFiles()
	processSourceTemplate(enumTmpl, &enumInfo{
		Pkg:        "model/gurps/advantage",
		Name:       "affects",
		Desc:       "describes how an AdvantageModifier affects the point cost",
		StandAlone: true,
		Values: []enumValue{
			{
				Key:    "total",
				String: "to cost",
			},
			{
				Key:    "base_only",
				String: "to base cost only",
				Short:  "(base only)",
			},
			{
				Key:    "levels_only",
				String: "to leveled cost only",
				Short:  "(levels only)",
			},
		},
	})
	processSourceTemplate(enumTmpl, &enumInfo{
		Pkg:        "model/gurps/advantage",
		Name:       "container_type",
		Desc:       "holds the type of an advantage container",
		StandAlone: true,
		Values: []enumValue{
			{
				Key:    "group",
				String: "Group",
			},
			{
				Key:    "meta_trait",
				String: "Meta-Trait",
			},
			{
				Key:    "race",
				String: "Race",
			},
			{
				Key:    "alternative_abilities",
				String: "Alternative Abilities",
			},
		},
	})
	processSourceTemplate(enumTmpl, &enumInfo{
		Pkg:        "model/gurps/advantage",
		Name:       "modifier_cost_type",
		Desc:       "describes how an AdvantageModifier's point cost is applied",
		StandAlone: true,
		Values: []enumValue{
			{
				Key:    "percentage",
				String: "%",
			},
			{
				Key:    "points",
				String: "points",
			},
			{
				Key:    "multiplier",
				String: "×",
			},
		},
	})
	processSourceTemplate(enumTmpl, &enumInfo{
		Pkg:        "model/gurps/attribute",
		Name:       "bonus_limitation",
		Desc:       "holds a limitation for an AttributeBonus",
		StandAlone: true,
		Values: []enumValue{
			{
				Key: "none",
			},
			{
				Key:    "striking_only",
				String: "for striking only",
			},
			{
				Key:    "lifting_only",
				String: "for lifting only",
			},
			{
				Key:    "throwing_only",
				String: "for throwing only",
			},
		},
	})
	processSourceTemplate(enumTmpl, &enumInfo{
		Pkg:        "model/gurps/attribute",
		Name:       "damage_progression",
		Desc:       "controls how Thrust and Swing are calculated",
		StandAlone: true,
		Values: []enumValue{
			{
				Key:    "basic_set",
				String: "Basic Set",
			},
			{
				Key:    "knowing_your_own_strength",
				String: "Knowing Your Own Strength",
				Short:  "Pyramid 3-83, pages 16-19",
			},
			{
				Key:    "no_school_grognard_damage",
				String: "No School Grognard",
				Short:  "https://noschoolgrognard.blogspot.com/2013/04/adjusting-swing-damage-in-dungeon.html",
			},
			{
				Key:    "thrust_equals_swing_minus_2",
				String: "Thrust = Swing-2",
				Short:  "https://github.com/richardwilkes/gcs/issues/97",
			},
			{
				Key:    "swing_equals_thrust_plus_2",
				String: "Swing = Thrust+2",
				Short:  "Houserule originating with Kevin Smyth. See https://gamingballistic.com/2020/12/04/df-eastmarch-boss-fight-and-house-rules/",
			},
			{
				Key:    "phoenix_flame_d3",
				String: "PhoenixFlame d3",
				Short:  "Houserule that use d3s instead of d6s for Damage. See: https://github.com/richardwilkes/gcs/pull/393",
			},
		},
	})
	processSourceTemplate(enumTmpl, &enumInfo{
		Pkg:        "model/gurps/attribute",
		Name:       "threshold_op",
		Desc:       "holds an operation to apply when a pool threshold is hit",
		StandAlone: true,
		Values: []enumValue{
			{
				Key:    "unknown",
				String: "Unknown",
				Short:  "Unknown",
			},
			{
				Key:    "halve_move",
				String: "Halve Move",
				Short:  "Halve Move (round up)",
			},
			{
				Key:    "halve_dodge",
				String: "Halve Dodge",
				Short:  "Halve Dodge (round up)",
			},
			{
				Name:   "HalveST",
				Key:    "halve_st",
				String: "Halve Strength",
				Short:  "Halve Strength (round up; does not affect HP and damage)",
			},
		},
	})
	processSourceTemplate(enumTmpl, &enumInfo{
		Pkg:        "model/gurps/attribute",
		Name:       "type",
		Desc:       "holds the type of an attribute definition",
		StandAlone: true,
		Values: []enumValue{
			{
				Key:    "integer",
				String: "Integer",
			},
			{
				Key:    "decimal",
				String: "Decimal",
			},
			{
				Key:    "pool",
				String: "Pool",
			},
		},
	})
}

func removeExistingGenFiles() {
	root, err := filepath.Abs(rootDir)
	jot.FatalIfErr(err)
	jot.FatalIfErr(filepath.Walk(root, func(path string, info os.FileInfo, err error) error {
		name := info.Name()
		if info.IsDir() {
			if name == ".git" {
				return filepath.SkipDir
			}
		} else {
			if strings.HasSuffix(name, genSuffix) {
				jot.FatalIfErr(os.Remove(path))
			}
		}
		return nil
	}))
}

func processSourceTemplate(tmplName string, info *enumInfo) {
	tmpl, err := template.New(tmplName).Funcs(template.FuncMap{
		"add":          add,
		"emptyIfTrue":  emptyIfTrue,
		"fileLeaf":     filepath.Base,
		"firstToLower": txt.FirstToLower,
		"last":         last,
		"toCamelCase":  txt.ToCamelCase,
		"toIdentifier": toIdentifier,
		"toKey":        toKey,
		"wrapComment":  wrapComment,
	}).ParseFiles(tmplName)
	jot.FatalIfErr(err)
	var buffer bytes.Buffer
	writeGeneratedFromComment(&buffer, tmplName)
	jot.FatalIfErr(tmpl.Execute(&buffer, info))
	var data []byte
	if data, err = format.Source(buffer.Bytes()); err != nil {
		fmt.Println("unable to format source file: " + filepath.Join(info.Pkg, info.Name+genSuffix))
		data = buffer.Bytes()
	}
	dir := filepath.Join(rootDir, info.Pkg)
	jot.FatalIfErr(os.MkdirAll(dir, 0o750))
	jot.FatalIfErr(os.WriteFile(filepath.Join(dir, info.Name+genSuffix), data, 0o640))
}

func writeGeneratedFromComment(w io.Writer, tmplName string) {
	_, err := fmt.Fprintf(w, "// Code generated from \"%s\" - DO NOT EDIT.\n\n", tmplName)
	jot.FatalIfErr(err)
}

func add(a, b int) int {
	return a + b
}

func (e *enumInfo) LocalType() string {
	return txt.FirstToLower(toIdentifier(e.Name)) + "Data"
}

func (e *enumInfo) IDFor(v enumValue) string {
	id := v.Name
	if id == "" {
		id = toIdentifier(v.Key)
	}
	if !e.StandAlone {
		id += toIdentifier(e.Name)
	}
	return id
}

func (e *enumInfo) HasShort() bool {
	for _, one := range e.Values {
		if one.Short != "" {
			return true
		}
	}
	return false
}

func last(in []enumValue) enumValue {
	return in[len(in)-1]
}

func emptyIfTrue(str string, test bool) string {
	if test {
		return ""
	}
	return str
}

func toIdentifier(in string) string {
	var buffer strings.Builder
	useUpper := true
	for i, ch := range in {
		isUpper := ch >= 'A' && ch <= 'Z'
		isLower := ch >= 'a' && ch <= 'z'
		isDigit := ch >= '0' && ch <= '9'
		isAlpha := isUpper || isLower
		if i == 0 && !isAlpha {
			if !isDigit {
				continue
			}
			buffer.WriteString("_")
		}
		if isAlpha {
			if useUpper {
				buffer.WriteRune(unicode.ToUpper(ch))
			} else {
				buffer.WriteRune(unicode.ToLower(ch))
			}
			useUpper = false
		} else {
			if isDigit {
				buffer.WriteRune(ch)
			}
			useUpper = true
		}
	}
	return buffer.String()
}

func toKey(in string) string {
	var buffer strings.Builder
	lastWasUnderscore := false
	lastWasLower := false
	runes := []rune(in)
	for i, ch := range runes {
		isUpper := ch >= 'A' && ch <= 'Z'
		isLower := ch >= 'a' && ch <= 'z'
		isDigit := ch >= '0' && ch <= '9'
		isAlpha := isUpper || isLower
		if buffer.Len() == 0 && !isAlpha {
			if !isDigit {
				continue
			}
		}
		if isAlpha || isDigit {
			if i != 0 && !lastWasUnderscore && isUpper {
				if lastWasLower {
					buffer.WriteRune('_')
				} else if i+1 < len(runes) {
					nextCh := runes[i+1]
					if nextCh >= 'a' && nextCh <= 'z' {
						buffer.WriteRune('_')
					}
				}
			}
			buffer.WriteRune(unicode.ToLower(ch))
			lastWasUnderscore = false
			lastWasLower = isLower
		} else if !lastWasUnderscore {
			buffer.WriteRune('_')
			lastWasUnderscore = true
		}
	}
	return strings.TrimRight(buffer.String(), "_")
}

func wrapComment(in string, cols int) string {
	return txt.Wrap("// ", in, cols)
}
