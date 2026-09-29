package main

import (
	"log"
	"os"
	"os/exec"
	"strings"
)

var GLOBALS map[string]string

type Blueprint struct {
	Name   string   `json:"name"`
	Source string   `json:"source"`
	Steps  []string `json:"steps"`
	Needs  []string `json:"needs"`
}

func build(bp Blueprint) {
	GLOBALS = make(map[string]string)
	GLOBALS["$PREFIX"] = os.Getenv("PREFIX") + "/" + bp.Name
	GLOBALS["$NPROC"] = os.Getenv("NPROC")
	cleanBuildDirectory()
	downloadSource(bp)
	extractArchive(bp)
	runSteps(bp)
}

func runSteps(bp Blueprint) {
	files, err := os.ReadDir("build/")
	if err != nil {
		log.Fatal(err)
	}
	for _, step := range bp.Steps {
		parts := strings.Split(evaluateVariable(step), " ")
		cmd := exec.Command(parts[0], parts[1:]...)
		cmd.Dir = "build/" + files[0].Name()
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
		err := cmd.Run()
		if err != nil {
			log.Fatal(err)
		}
	}
}

func evaluateVariable(in string) string {
	for k, v := range GLOBALS {
		if strings.Contains(in, k) {
			return strings.Replace(in, k, v, 1)
		}
	}
	return in
}

func downloadSource(bp Blueprint) {
	cmd := exec.Command("wget", bp.Source)
	cmd.Dir = "build/"
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	err := cmd.Run()
	if err != nil {
		log.Fatal(err)
	}
}

func cleanBuildDirectory() {
	os.RemoveAll("build/")
	os.Mkdir("build/", 0777)
}

func extractArchive(bp Blueprint) {
	if strings.Index(bp.Source, ".zip") != -1 {
		extractArchiveZip(bp)
	}
	if strings.Index(bp.Source, ".tar.") != -1 {
		extractArchiveTar(bp)
	}
}

func extractArchiveZip(bp Blueprint) {
	files, err := os.ReadDir("build/")
	if err != nil {
		log.Fatal(err)
	}
	cmd := exec.Command("unzip", files[0].Name())
	cmd.Dir = "build/"
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	err = cmd.Run()
	if err != nil {
		log.Fatal(err)
	}
	os.Remove("build/" + files[0].Name())
}

func extractArchiveTar(bp Blueprint) {
	files, err := os.ReadDir("build/")
	if err != nil {
		log.Fatal(err)
	}
	cmd := exec.Command("tar", "xf", files[0].Name())
	cmd.Dir = "build/"
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	err = cmd.Run()
	if err != nil {
		log.Fatal(err)
	}
	os.Remove("build/" + files[0].Name())
}
