package runner

import (
	"context"
	"os"
	"time"

	"github.com/saucelabs/saucectl/cli/command"
	"github.com/saucelabs/saucectl/cli/config"
)

const (
	logDir           = "/var/log/cont"
	runnerConfigPath = "/home/testrunner/config.yaml"
)

var logFiles = [...]string{
	logDir + "/chrome_browser.log",
	logDir + "/firefox_browser.log",
	logDir + "/supervisord.log",
	logDir + "/video-rec-stderr.log",
	logDir + "/video-rec-stdout.log",
	logDir + "/wait-xvfb.1.log",
	logDir + "/wait-xvfb.2.log",
	logDir + "/wait-xvfb-stdout.log",
	logDir + "/xvfb-tryouts-stderr.log",
	logDir + "/xvfb-tryouts-stdout.log",
	"/home/seluser/videos/video.mp4",
	"/home/seluser/docker.log",
}

// Testrunner describes the test runner interface
type Testrunner interface {
	Context() context.Context
	CLI() *command.SauceCtlCli

	Setup() error
	Run() (int, error)
	Teardown(logDir string) error
}

type baseRunner struct {
	jobConfig    config.JobConfiguration
	runnerConfig config.RunnerConfiguration
	context      context.Context
	cli          *command.SauceCtlCli

	startTime int64
}

func (r baseRunner) Context() context.Context {
	return r.context
}

func (r baseRunner) CLI() *command.SauceCtlCli {
	return r.cli
}

// New creates a new testrunner object
func New(c config.JobConfiguration, cli *command.SauceCtlCli) (Testrunner, error) {
	ctx := context.Background()
	runnerConfig := config.RunnerConfiguration{}

	_, err := os.Stat(runnerConfigPath)
	if os.IsNotExist(err) {
		return localRunner{baseRunner{c, runnerConfig, ctx, cli, makeTimestamp()}, ""}, nil
	}

	return ciRunner{baseRunner{c, runnerConfig, ctx, cli, makeTimestamp()}}, nil
}

func makeTimestamp() int64 {
	return time.Now().UnixNano() / int64(time.Millisecond)
}
