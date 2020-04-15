package runner

import (
	"context"
	"io/ioutil"
	"os"
	"os/exec"
	"path"
	"path/filepath"

	"github.com/saucelabs/saucectl/cli/command"
	"github.com/saucelabs/saucectl/cli/config"
)

type ciRunner struct {
	baseRunner
}

func newCIRunner(c config.JobConfiguration, cli *command.SauceCtlCli) (ciRunner, error) {
	runner := ciRunner{}

	// read runner config file
	rc, err := config.NewRunnerConfiguration(runnerConfigPath)
	if err != nil {
		return runner, err
	}

	runner.cli = cli
	runner.context = context.Background()
	runner.jobConfig = c
	runner.startTime = makeTimestamp()
	runner.runnerConfig = rc
	return runner, nil
}

func (r ciRunner) Setup() error {
	// copy files from repository into target dir
	for _, pattern := range r.jobConfig.Files {
		matches, err := filepath.Glob(pattern)
		if err != nil {
			continue
		}

		for _, file := range matches {
			r.cli.Logger.Info().Msg("Copy file " + file + " to " + r.runnerConfig.TargetDir)
			if err := copyFile(file, r.runnerConfig.TargetDir); err != nil {
				return err
			}
		}
	}
	return nil
}

func (r ciRunner) Run() (int, error) {
	cmd := exec.Command(r.runnerConfig.ExecCommand[0], r.runnerConfig.ExecCommand[1])

	cmd.Stdout = r.cli.Out()
	cmd.Stderr = r.cli.Out()
	cmd.Dir = r.runnerConfig.RootDir
	err := cmd.Run()

	if err != nil {
		return 1, err
	}
	return 0, nil
}

func (r ciRunner) Teardown(logDir string) error {
	if logDir != "" {
		return nil
	}

	for _, containerSrcPath := range logFiles {
		file := filepath.Base(containerSrcPath)
		dstPath := filepath.Join(logDir, file)
		if err := copyFile(containerSrcPath, dstPath); err != nil {
			continue
		}
	}

	return nil
}

func copyFile(src string, targetDir string) error {
	pwd, err := os.Getwd()
	if err != nil {
		return err
	}

	srcFile := src
	if !path.IsAbs(srcFile) {
		srcFile = filepath.Join(pwd, src)
	}

	input, err := ioutil.ReadFile(srcFile)
	if err != nil {
		return err
	}

	err = ioutil.WriteFile(targetDir+filepath.Base(srcFile), input, 0644)
	if err != nil {
		return err
	}

	return nil
}
