package sftp_test

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"

	"github.com/restic/restic/backend"
	"github.com/restic/restic/backend/sftp"
	"github.com/restic/restic/backend/test"

	. "github.com/restic/restic/test"
)

var tempBackendDir string

//go:generate go run ../test/generate_backend_tests.go

func createTempdir() error {
	if tempBackendDir != "" {
		return nil
	}

	tempdir, err := ioutil.TempDir("", "restic-local-test-")
	if err != nil {
		return err
	}

	fmt.Printf("created new test backend at %v\n", tempdir)
	tempBackendDir = tempdir
	return nil
}

func init() {
	sftpserver := ""

	for _, dir := range strings.Split(TestSFTPPath, ":") {
		testpath := filepath.Join(dir, "sftp-server")
		fd, err := os.Open(testpath)
		fd.Close()
		if !os.IsNotExist(err) {
			sftpserver = testpath
			break
		}
	}

	if sftpserver == "" {
		SkipMessage = "sftp server binary not found, skipping tests"
		return
	}

	test.CreateFn = func() (backend.Backend, error) {
		err := createTempdir()
		if err != nil {
			return nil, err
		}

		return sftp.Create(tempBackendDir, sftpserver)
	}

	test.OpenFn = func() (backend.Backend, error) {
		err := createTempdir()
		if err != nil {
			return nil, err
		}
		return sftp.Open(tempBackendDir, sftpserver)
	}

	test.CleanupFn = func() error {
		if tempBackendDir == "" {
			return nil
		}

		fmt.Printf("removing test backend at %v\n", tempBackendDir)
		err := os.RemoveAll(tempBackendDir)
		tempBackendDir = ""
		return err
	}
}
