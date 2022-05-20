package main

import (
	"flag"
	"fmt"
	"io"
	"os"
	"os/exec"
	"path"
	"regexp"
	"strings"
)

var lineRegex = regexp.MustCompile(`(?m)^\s+(.+?)\s=>\s(.+?)\s\(`)
var diskRegex = regexp.MustCompile(`^(\w):`)

func main() {
	flag.Parse()
	cmd := exec.Command("ldd", flag.Arg(0))

	stdout, err := cmd.Output()
	if err != nil {
		panic(err)
	}

	pathEnv := os.Getenv("PATH")
	pathSplit := strings.Split(pathEnv, ";")

	cleanPath := make([]string, len(pathSplit))
	for i, p := range pathSplit {
		clean := strings.Replace(p, "\\", "/", -1)
		diskMatch := diskRegex.FindAllStringSubmatch(p, -1)
		if len(diskMatch) > 0 {
			clean = "/" + diskMatch[0][1] + clean[2:]
		}
		clean = strings.ToLower(clean)

		println("PATH:", clean)
		cleanPath[i] = clean
	}

	wd, err := os.Getwd()
	if err != nil {
		panic(err)
	}

	matches := lineRegex.FindAllSubmatch(stdout, -1)
	for _, match := range matches {
		println("Found match:", string(match[0]))
		for _, p := range cleanPath {
			if strings.HasPrefix(string(match[2]), p) {
				src := string(match[2])
				src = src[1:2] + ":" + src[2:]
				dst := path.Join(wd, "out", string(match[1]))
				println("Copying", src, "=>", dst)
				if _, err := copyFile(src, dst); err != nil {
					panic(err)
				}
			}
		}
	}
}

func copyFile(src string, dst string) (int64, error) {
	sourceFileStat, err := os.Stat(src)
	if err != nil {
		return 0, err
	}

	if !sourceFileStat.Mode().IsRegular() {
		return 0, fmt.Errorf("%s is not a regular file", src)
	}

	source, err := os.Open(src)
	if err != nil {
		return 0, err
	}
	defer source.Close()

	destination, err := os.Create(dst)
	if err != nil {
		return 0, err
	}
	defer destination.Close()
	nBytes, err := io.Copy(destination, source)
	return nBytes, err
}
