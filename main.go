package main

import (
	"bufio"
	"flag"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"strings"
)

var (
	targetDirectory string
	filePtn         string
)

type Profile struct {
	Args []string
	Ext  string
}

var profileMap = map[string]Profile{
	"_d.wav": {Args: []string{"-c:a", "aac"}, Ext: "_s.m4a"},
}

func getProfile(profile string) (Profile, error) {
	for k, y := range profileMap {
		if k == profile {
			return y, nil
		}
	}
	return Profile{}, fmt.Errorf("ERROR")
}

func init() {
	flag.StringVar(&targetDirectory, "target", "", "")
	flag.StringVar(&filePtn, "pattern", "", "")
}

func main() {
	flag.Parse()
	fPtn := regexp.MustCompile(fmt.Sprintf("%s", filePtn))
	files, err := os.ReadDir(targetDirectory)
	if err != nil {
		panic(err)
	}

	for _, entry := range files {
		name := entry.Name()

		//create a log file
		ffmpegLog, err := os.Create(filepath.Join(targetDirectory, fmt.Sprintf("%s.log", name)))
		if err != nil {
			panic(err)
		}
		defer ffmpegLog.Close()
		writer := bufio.NewWriter(ffmpegLog)

		if fPtn.MatchString(name) {
			inputFile := filepath.Join(targetDirectory, name)
			profile, err := getProfile(filePtn)
			if err != nil {
				panic(err)
			}
			cmds := []string{"-i", inputFile}
			cmds = append(cmds, profile.Args...)
			outputFile := filepath.Join(targetDirectory, strings.ReplaceAll(name, filePtn, profile.Ext))
			cmds = append(cmds, outputFile)

			fmt.Printf("Transcoding: %s\n", inputFile)
			cmd := exec.Command("ffmpeg", cmds...)
			cmd.Stdout = os.Stdout
			cmd.Stderr = writer
			if err := cmd.Run(); err != nil {
				panic(err)
			}
			writer.Flush()
			fmt.Printf("transcoded: %s\n", outputFile)
		}
	}
}
