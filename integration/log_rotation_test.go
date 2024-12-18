//go:build !windows
// +build !windows

package integration

import (
	"bufio"
	"net/http"
	"os"
	"strings"
	"syscall"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
	"github.com/traefik/traefik/v2/integration/try"
	"github.com/traefik/traefik/v2/pkg/log"
)

const (
	traefikTestLogFileRotated       = traefikTestLogFile + ".rotated"
	traefikTestAccessLogFileRotated = traefikTestAccessLogFile + ".rotated"
)

// Log rotation integration test suite.
type LogRotationSuite struct{ BaseSuite }

func TestLogRotationSuite(t *testing.T) {
	suite.Run(t, new(LogRotationSuite))
}

func (s *LogRotationSuite) SetupSuite() {
	s.BaseSuite.SetupSuite()

	os.Remove(traefikTestAccessLogFile)
	os.Remove(traefikTestLogFile)
	os.Remove(traefikTestAccessLogFileRotated)

	s.createComposeProject("access_log")
	s.composeUp()
}

func (s *LogRotationSuite) TearDownSuite() {
	s.BaseSuite.TearDownSuite()

	generatedFiles := []string{
		traefikTestLogFile,
		traefikTestLogFileRotated,
		traefikTestAccessLogFile,
		traefikTestAccessLogFileRotated,
	}

	for _, filename := range generatedFiles {
		if err := os.Remove(filename); err != nil {
			log.WithoutContext().Warning(err)
		}
	}
}

func (s *LogRotationSuite) TestAccessLogRotation() {
	// Start Traefik
	cmd, _ := s.cmdTraefik(withConfigFile("fixtures/access_log_config.toml"))
	defer s.displayTraefikLogFile(traefikTestLogFile)

	// Verify Traefik started ok
	s.verifyEmptyErrorLog("traefik.log")

	s.waitForTraefik("server1")

	// Make some requests
	req, err := http.NewRequest(http.MethodGet, "http://127.0.0.1:8000/", nil)
	require.NoError(s.T(), err)
	req.Host = "frontend1.docker.local"

	err = try.Request(req, 500*time.Millisecond, try.StatusCodeIs(http.StatusOK), try.HasBody())
	require.NoError(s.T(), err)

	// Rename access log
	err = os.Rename(traefikTestAccessLogFile, traefikTestAccessLogFileRotated)
	require.NoError(s.T(), err)

	// in the midst of the requests, issue SIGUSR1 signal to server process
	err = cmd.Process.Signal(syscall.SIGUSR1)
	require.NoError(s.T(), err)

	// continue issuing requests
	err = try.Request(req, 500*time.Millisecond, try.StatusCodeIs(http.StatusOK), try.HasBody())
	require.NoError(s.T(), err)
	err = try.Request(req, 500*time.Millisecond, try.StatusCodeIs(http.StatusOK), try.HasBody())
	require.NoError(s.T(), err)

	// Verify access.log.rotated output as expected
	s.logAccessLogFile(traefikTestAccessLogFileRotated)
	lineCount := s.verifyLogLines(traefikTestAccessLogFileRotated, 0, true)
	assert.GreaterOrEqual(s.T(), lineCount, 1)

	// make sure that the access log file is at least created before we do assertions on it
	err = try.Do(1*time.Second, func() error {
		_, err := os.Stat(traefikTestAccessLogFile)
		return err
	})
	assert.NoError(s.T(), err, "access log file was not created in time")

	// Verify access.log output as expected
	s.logAccessLogFile(traefikTestAccessLogFile)
	lineCount = s.verifyLogLines(traefikTestAccessLogFile, lineCount, true)
	assert.Equal(s.T(), 3, lineCount)

	s.verifyEmptyErrorLog(traefikTestLogFile)
}

func (s *LogRotationSuite) TestTraefikLogRotation() {
	// Start Traefik
	cmd := s.traefikCmd(withConfigFile("fixtures/traefik_log_config.toml"))

	s.waitForTraefik("server1")

	// Rename traefik log
	err := os.Rename(traefikTestLogFile, traefikTestLogFileRotated)
	require.NoError(s.T(), err)

	// issue SIGUSR1 signal to server process
	err = cmd.Process.Signal(syscall.SIGUSR1)
	require.NoError(s.T(), err)

	err = cmd.Process.Signal(syscall.SIGTERM)
	require.NoError(s.T(), err)

	// Allow time for switch to be processed
	err = try.Do(3*time.Second, func() error {
		_, err = os.Stat(traefikTestLogFile)
		return err
	})
	require.NoError(s.T(), err)

	// we have at least 6 lines in traefik.log.rotated
	lineCount := s.verifyLogLines(traefikTestLogFileRotated, 0, false)

	// GreaterOrEqualThan used to ensure test doesn't break
	// If more log entries are output on startup
	assert.GreaterOrEqual(s.T(), lineCount, 5)

	// Verify traefik.log output as expected
	lineCount = s.verifyLogLines(traefikTestLogFile, lineCount, false)
	assert.GreaterOrEqual(s.T(), lineCount, 7)
}

func (s *LogRotationSuite) logAccessLogFile(fileName string) {
	output, err := os.ReadFile(fileName)
	require.NoError(s.T(), err)
	log.WithoutContext().Infof("Contents of file %s\n%s", fileName, string(output))
}

func (s *LogRotationSuite) verifyEmptyErrorLog(name string) {
	err := try.Do(5*time.Second, func() error {
		traefikLog, e2 := os.ReadFile(name)
		if e2 != nil {
			return e2
		}
		assert.Empty(s.T(), string(traefikLog))

		return nil
	})
	require.NoError(s.T(), err)
}

func (s *LogRotationSuite) verifyLogLines(fileName string, countInit int, accessLog bool) int {
	rotated, err := os.Open(fileName)
	require.NoError(s.T(), err)
	rotatedLog := bufio.NewScanner(rotated)
	count := countInit
	for rotatedLog.Scan() {
		line := rotatedLog.Text()
		if accessLog {
			if len(line) > 0 {
				if !strings.Contains(line, "/api/rawdata") {
					s.CheckAccessLogFormat(line, count)
					count++
				}
			}
		} else {
			count++
		}
	}

	return count
}
