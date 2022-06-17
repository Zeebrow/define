package main

import (
	"fmt"
	"net/url"
)

type ProgramInfo struct {
	ProgramName string
	Version     string
	BuildDate   string
	CommitHash  string
	OS          string
	GoVer       string
}

func (p *ProgramInfo) GetInfo() string {
	return fmt.Sprintf("%s version %s (%s) built %s", p.ProgramName, p.Version, p.CommitHash, p.BuildDate)
}

func (pi *ProgramInfo) NewBugReport(title string, errorMsg string) {
	label := url.QueryEscape("Auto-generated")
	body := fmt.Sprintf("Version: %s\nCommit Hash: %s\nBuild Date: %s\nOS: %s\nGo Version: %s\nError Message: %s\n",
		ProgInfo.Version, ProgInfo.CommitHash, ProgInfo.BuildDate, ProgInfo.OS, ProgInfo.GoVer, errorMsg)

	fmt.Println("Oops! Something went wrong. If you would like to let me know it happened, click the link below to open a new issue:")
	fmt.Println()
	fmt.Printf("https://github.com/Zeebrow/define/issues/new?title=%s&labels=%s&body=%s\n\n",
		url.QueryEscape(title), url.QueryEscape(label), url.QueryEscape(body))
}
