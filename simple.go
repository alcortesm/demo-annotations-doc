package main

import (
	. "github.com/bblfsh/sdk/uast"
	. "github.com/bblfsh/sdk/uast/ann"

	"github.com/bblfsh/bash-driver/driver/normalizer/intellij"
	"gopkg.in/src-d/go-errors.v0"
)

var (
	ErrAnError = errors.NewKind("a simple error")
)

var simpleRules = On(Any).Self(
	On(Not(intellij.File)).Error(ErrAnError.New()),
	On(intellij.File).Roles(File).Descendants(
		On(intellij.Comment).Roles(Comment),
	),
)
