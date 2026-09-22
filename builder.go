package main

import (
	"fmt"
	"log"
	"os"
	"os/exec"
	"strings"
)

var PATH string

type Blueprint struct {
	Name   string   `json:"name"`
	Source string   `json:"source"`
	Steps  []string `json:"steps"`
	Needs  []string `json:"needs"`
}

func build(bp Blueprint) {
	PATH = os.Getenv("HOME") + "/.pkgr/"
	fmt.Println(PATH)
	downloadSource(bp)
}

func downloadSource(bp Blueprint) {
	runCommand("wget", bp.Source)
	extractArchive(bp)
}

func runCommand(args ...string) {
	cmd := exec.Command(args[0], args[1:]...)
	cmd.Dir = "build/"
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	err := cmd.Run()
	if err != nil {
		log.Fatal(args[0], ":", err)
	}
}

func extractArchive(bp Blueprint) {
	if strings.Index(bp.Source, ".zip") != -1 {
		runCommand("unzip")
	}
}
