package utils

import (
	"bytes"
	"errors"
	"fmt"
	"os"
	"os/exec"
)

var CmdHandler CmdExec

func runMergeCmd(cmd *exec.Cmd, paths []string, mergeFilePath string) error {
	var stderr bytes.Buffer
	cmd.Stderr = &stderr
	err := cmd.Run()
	if err != nil {
		return fmt.Errorf("%s\n%s", err, stderr.String())
	}

	if mergeFilePath != "" {
		os.Remove(mergeFilePath) // nolint
	}
	// remove parts
	for _, path := range paths {
		os.Remove(path) // nolint
	}
	return nil
}

func runMergeCmdAndroid(cmd string, paths []string, mergeFilePath string) error {
	result := CmdHandler.RunCmd(cmd)
	if result != "ok" {
		return errors.New("ffmpeg error:" + result)
	}

	if mergeFilePath != "" {
		os.Remove(mergeFilePath) // nolint
	}
	// remove parts
	for _, path := range paths {
		os.Remove(path) // nolint
	}
	return nil
}

// MergeFilesWithSameExtension merges files that have the same extension into one.
// Can also handle merging audio and video.
func MergeFilesWithSameExtension(paths []string, mergedFilePath string) error {
	cmds := []string{
		"-y",
	}
	for _, path := range paths {
		cmds = append(cmds, "-i", path)
	}
	cmds = append(cmds, "-c:v", "copy", "-c:a", "copy", mergedFilePath)

	if CmdHandler != nil {
		fullCmd := ""
		for _, item := range cmds {
			fullCmd = fullCmd + " " + item
		}

		return runMergeCmdAndroid(fullCmd, paths, "")
	}

	return runMergeCmd(exec.Command("ffmpeg", cmds...), paths, "")
}

// MergeToMP4 merges video parts to an MP4 file.
func MergeToMP4(paths []string, mergedFilePath string, filename string) error {
	mergeFilePath := filename + ".txt" // merge list file should be in the current directory

	// write ffmpeg input file list
	mergeFile, _ := os.Create(mergeFilePath)
	for _, path := range paths {
		mergeFile.Write([]byte(fmt.Sprintf("file '%s'\n", path))) // nolint
	}
	mergeFile.Close() // nolint

	cmd := exec.Command(
		"ffmpeg", "-y", "-f", "concat", "-safe", "0",
		"-i", mergeFilePath, "-c", "copy", "-bsf:a", "aac_adtstoasc", mergedFilePath,
	)

	if CmdHandler != nil {
		fullCmd := "-y -f concat -safe 0 -i " + mergeFilePath + " -c copy -bsf:a aac_adtstoasc " + mergedFilePath
		return runMergeCmdAndroid(fullCmd, paths, mergeFilePath)
	}

	return runMergeCmd(cmd, paths, mergeFilePath)
}

type CmdExec interface {
	RunCmd(cmd string) string
}
