package suman

import (
	"bufio"
	"bytes"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/pkg/errors"

	_sumanUseCase "mlmtool/pkg/usecases/susemanager"
	returnCodes "mlmtool/pkg/util/returnCodes"
)

type CommandRunner func() (*bytes.Buffer, *bytes.Buffer, error)

// DefaultCommandRunner executes a command and returns its stdout and stderr
func execHostname() (*bytes.Buffer, *bytes.Buffer, error) {
	cmd := exec.Command("/bin/hostname", "-f")
	var out, stdErr bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &stdErr
	err := cmd.Run()
	return &out, &stdErr, err
}

// GetCredentials retrieves credentials needed for SUSE Manager Server primary
func GetCredentials(fileName string, runCommands ...CommandRunner) (_sumanUseCase.SumanConfig, error) {
	var sumancfg _sumanUseCase.SumanConfig
	sumancfg.Insecure = true
	fileName = filepath.Clean(fileName)

	file, err := os.Open(fileName)
	if err != nil {
		return sumancfg, errors.New(returnCodes.ErrOpeningFile)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		l := strings.SplitN(scanner.Text(), "=", 2) // Use SplitN to avoid index out of range
		if len(l) < 2 {
			continue // Skip malformed lines
		}
		switch strings.TrimSpace(l[0]) {
		case "username":
			sumancfg.Login = strings.TrimSpace(l[1])
		case "password":
			sumancfg.Password = strings.TrimSpace(l[1])
		}
	}

	runCommand := execHostname
	if len(runCommands) > 0 {
		runCommand = runCommands[0]
	}

	out, stdErr, err := runCommand()
	if err != nil {
		return sumancfg, err
	}
	if stdErr.Len() > 0 {
		return sumancfg, errors.New(stdErr.String())
	}
	if out.Len() == 0 {
		return sumancfg, errors.New("error retrieving the FQDN of the server")
	}
	sumancfg.Host = strings.TrimSpace(out.String())
	// TODO: Clarify trimming for e.g. ending dot, besides trimming for newline/whitespace, was also meant by the more
	// unsafe way of sumancfg.Host[:len(sumancfg.Host)-1], and just remove the next line
	sumancfg.Host = strings.TrimRight(sumancfg.Host, "./")
	return sumancfg, nil
}

// GetCredentialsUyuni retrieves credentials needed for SUSE Manager Server secondary servers
func GetCredentialsUyuni(fileName string) (_sumanUseCase.SumanConfig, error) {
	var sumancfg _sumanUseCase.SumanConfig
	sumancfg.Insecure = true
	fileName = filepath.Clean(fileName)

	file, err := os.Open(fileName)
	if err != nil {
		return sumancfg, errors.New(returnCodes.ErrOpeningFile)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		l := strings.SplitN(scanner.Text(), ":", 2) // Use SplitN to avoid index out of range
		if len(l) < 2 {
			continue // Skip malformed lines
		}
		switch strings.TrimSpace(l[0]) {
		case "user":
			sumancfg.Login = strings.TrimSpace(l[1])
		case "password":
			sumancfg.Password = strings.TrimSpace(l[1])
		case "hubmaster":
			sumancfg.Host = strings.TrimSpace(l[1])
		}
	}
	return sumancfg, nil
}
