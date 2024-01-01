package tmux

import (
	"errors"
	"fmt"
	"os"
	"os/exec"
	"strings"
	"syscall"
)

func execCmd(args []string) error {
	tmux, err := exec.LookPath("tmux")
	if err != nil {
		return err
	}
	args = append([]string{tmux}, args...)
	return syscall.Exec(tmux, args, os.Environ())
}

func GetSessions() ([]string, error) {
	cmd := exec.Command("tmux", "ls")
	cmdOutput, err := cmd.Output()
	if err != nil {
		return nil, err
	}
	lines := strings.Split(string(cmdOutput), "\n")

	sessions := make([]string, 0)
	for _, line := range lines {
		if line == "" {
			continue
		}
		parts := strings.Split(line, ":")
		if len(parts) < 2 {
			return nil, errors.New("Invalid tmux session line: " + line)
		}
		sessions = append(sessions, parts[0])
	}
	return sessions, nil
}

func GetCurrentSession() (string, error) {
	cmd := exec.Command("tmux", "display-message", "-p", "#S")
	cmdOutput, err := cmd.Output()
	if err != nil {
		return "", err
	}
	return strings.TrimSpace(string(cmdOutput)), nil
}

func CreateSession(basePath, session string) error {
	cmd := exec.Command(
		"tmux", "new",
		"-s",
		session,
		"-d", "-c",
		fmt.Sprintf("%s/%s", basePath, session),
	)
	return cmd.Run()
}

func SwitchToSession(session string) error {
	return execCmd([]string{"switch", "-t", session})
}
