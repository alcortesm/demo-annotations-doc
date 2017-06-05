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
		path:  "$GOPATH/src/github.com/alcortesm/demo-annotations-doc/java.go",
		skip:  10,
		rules: javaRules,
	},
	bash: localRules{
		path:  "$GOPATH/src/github.com/alcortesm/demo-annotations-doc/bash.go",
		skip:  14,
		rules: bashRules,
	},
}

func main() {
	lang := parseArgs()
	if err := report(lang); err != nil {
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

func report(l lang) error {
	a, ok := knownLocalRules[l]
	if !ok {
		return fmt.Errorf("unknown language %s", l)
	}

	text, err := original(a)
	if err != nil {
		return err
	}
	fmt.Println(text)

	text = doc(a.rules)
	fmt.Println(text)

	return nil
}

func doc(r *ann.Rule) string {
	asMarkdown := (*docgen.RulesAsMarkdown)(r)
	return asMarkdown.String()
}

func original(l localRules) (string, error) {
	path := os.Expand(l.path, os.Getenv)
	raw, err := ioutil.ReadFile(path)
	if err != nil {
		return "", err
	}

	return tail(string(raw), l.skip), nil
}

func tail(s string, n int) string {
	lines := strings.Split(s, "\n")
	lines = lines[n:]
	return strings.Join(lines, "\n")
}
