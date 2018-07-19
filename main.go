package main

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

const (
	// ExitOK is 0
	ExitOK = 0
	// ExitError is 1
	ExitError = 1
)

func usage() {
	msg := `Edit the line selected by grep`
	fmt.Println(msg)
}

func grep(cfg *config, arg string) ([]byte, error) {
	var buf bytes.Buffer
	var err error
	c1 := exec.Command("memo", "grep", arg)
	c2 := exec.Command("sed", "-e", "s@"+cfg.MemoDir+"/@@")
	c2.Stdin, err = c1.StdoutPipe()
	if err != nil {
		return nil, err
	}
	c2.Stdout = &buf
	if err = c2.Start(); err != nil {
		return nil, err
	}
	if err = c1.Run(); err != nil {
		return nil, err
	}
	if err = c2.Wait(); err != nil {
		return nil, err
	}

	return buf.Bytes(), nil
}

func filter(cfg *config, out []byte) (*bytes.Buffer, error) {
	var buf bytes.Buffer
	var cmd *exec.Cmd
	if cfg.SelectCmd == "fzf" {
		// TODO: Extract this setting to the setting file?
		option := "--multi --cycle --bind=ctrl-u:half-page-up,ctrl-d:half-page-down"
		cmd = exec.Command(cfg.SelectCmd, (strings.Split(option, " "))[0:]...)
	} else {
		cmd = exec.Command(cfg.SelectCmd)
	}
	cmd.Stderr = os.Stderr
	cmd.Stdout = &buf
	cmd.Stdin = bytes.NewReader(out)

	if err := cmd.Run(); err != nil {
		// If the file is not selected, then it exit 0
		if len(buf.String()) == 0 {
			// os.Exit(ExitOK)
			return nil, nil
		}
		return nil, err
	}

	return &buf, nil
}

func edit(cfg *config, buf *bytes.Buffer) error {
	line := strings.TrimSpace(buf.String())
	filename := strings.TrimSpace(strings.Split(line, ":")[0])
	lineno := strings.TrimSpace(strings.Split(line, ":")[1])

	cmd := exec.Command(cfg.Editor, "+"+lineno, filepath.Join(cfg.MemoDir, filename))
	cmd.Stderr = os.Stderr
	cmd.Stdout = os.Stdout
	cmd.Stdin = os.Stdin
	return cmd.Run()
}

func cmd(arg string) error {
	var cfg config
	if err := cfg.load(); err != nil {
		return err
	}

	out, err := grep(&cfg, arg)
	if err != nil {
		return err
	}

	buf, err := filter(&cfg, out)
	if err != nil {
		return err
	}
	if buf == nil {
		// If the file is not selected(buf==nil), then it exit 0
		return err
	}

	return edit(&cfg, buf)
}

func returnCode(err error) int {
	if err != nil {
		fmt.Fprintf(os.Stderr, "%s: %v\n", os.Args[0], err)
		return ExitError
	}
	return ExitOK
}

func main() {
	arg := "."
	if len(os.Args) > 1 {
		if os.Args[1] == "-usage" {
			usage()
			os.Exit(ExitOK)
		}
		arg = os.Args[1]
	}
	os.Exit(returnCode(cmd(arg)))
}
