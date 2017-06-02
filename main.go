package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"strings"

	"github.com/bblfsh/sdk/docgen"
	"github.com/bblfsh/sdk/uast/ann"
)

type localRules struct {
	path  string // relative to $GOPATH
	skip  int    // how many lines to skip from the file for pretty printing it
	rules *ann.Rule
}

type lang string

const (
	java lang = "java"
	bash lang = "bash"
)

var knownLocalRules = map[lang]localRules{
	java: localRules{
		path:  "/src/github.com/alcortesm/demo-annotations-doc/java.go",
		skip:  11,
		rules: javaRules,
	},
	bash: localRules{
		path:  "/src/github.com/alcortesm/demo-annotations-doc/bash.go",
		skip:  14,
		rules: bashRules,
	},
}

func main() {
	lang := parseArgs()
	localRules, ok := knownLocalRules[lang]
	if !ok {
		fmt.Fprintln(os.Stderr, "unknown language %s", lang)
	}
	/*
		if err := report(localRules); err != nil {
			fmt.Fprintln(os.Stderr, err)
			os.Exit(1)
		}
	*/
	if err := printTable(localRules); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

func parseArgs() lang {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "invalid number of command arguments")
		usage()
		os.Exit(1)
	}
	return lang(os.Args[1])
}

func usage() {
	fmt.Fprintln(os.Stderr, "usage:")
	fmt.Fprintf(os.Stderr, "\t%s [java|bash]\n", os.Args[0])
}

func report(lr localRules) error {
	desc := fmt.Sprint(lr.rules)
	descSplit := strings.Split(desc, "\n")
	path := os.Expand(lr.path, os.Getenv)
	raw, err := ioutil.ReadFile(os.Getenv("GOPATH") + path)
	if err != nil {
		return err
	}
	rawSplit := strings.Split(string(raw), "\n")
	rawSplit = rawSplit[lr.skip:]
	printSideBySide(rawSplit, descSplit)
	return nil
}

func printSideBySide(a, b []string) {
	a = tabsToSpaces(a)
	b = tabsToSpaces(b)

	maxNLines := len(a)
	if len(b) > maxNLines {
		maxNLines = len(b)
	}

	maxLineLenA := maxLineLen(a)
	format := fmt.Sprintf("%%-%ds | %%s\n", maxLineLenA)

	var aLine, bLine string
	for i := 0; i < maxNLines; i++ {
		if i >= len(a) {
			aLine = ""
		} else {
			aLine = a[i]
		}
		if i >= len(b) {
			bLine = ""
		} else {
			bLine = b[i]
		}

		fmt.Printf(format, aLine, bLine)
	}
}

func tabsToSpaces(s []string) []string {
	ret := make([]string, len(s))
	for i, l := range s {
		ret[i] = strings.Replace(l, "\t", "  ", -1)
	}
	return ret
}

func maxLineLen(a []string) int {
	m := 0
	var s string
	for _, l := range a {
		s = fmt.Sprint("%s", l)
		if len(s) > m {
			m = len(s)
		}
	}
	return m
}

func printTable(lr localRules) error {
	asMarkdown := docgen.RulesAsMarkdown(*lr.rules)
	text, err := asMarkdown.MarshalText()
	if err != nil {
		return fmt.Errorf("printTable: %s", err)
	}
	fmt.Println(string(text))
	return nil
}
