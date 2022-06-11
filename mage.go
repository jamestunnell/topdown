//go:build mage

package main

import (
	"github.com/magefile/mage/mg"
	"github.com/magefile/mage/sh"
)

func CLOC() {
	mg.Deps(InstallCLOC)
	sh.RunV("gocloc", ".")
}

func InstallCLOC() {
	sh.Run("go", "install", "github.com/hhatto/gocloc/cmd/gocloc@latest")
}
