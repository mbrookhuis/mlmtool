package suman

import (
	"bytes"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"mlmtool/pkg/usecases/susemanager"
)

func mockCommandRunner() (*bytes.Buffer, *bytes.Buffer, error) {
	out := bytes.NewBufferString("mocked.hostname\n")
	return out, bytes.NewBuffer(nil), nil
}

func createTempFile(t *testing.T, name string, fileContent string) *os.File {
	tmpFile, err := os.CreateTemp("", name)
	if err != nil {
		t.Fatalf("Cannot create temporary file: %s", err)
	}
	if _, err = tmpFile.Write([]byte(fileContent)); err != nil {
		t.Fatalf("Cannot write to temporary file: %s", err)
	}
	if err = tmpFile.Close(); err != nil {
		t.Fatalf("Cannot close temporary file: %s", err)
	}

	return tmpFile
}

func assertExpected(t *testing.T, actual *susemanager.SumanConfig, expected *susemanager.SumanConfig) {
	assert.Equal(t, expected.Login, actual.Login)
	assert.Equal(t, expected.Password, actual.Password)
	assert.Equal(t, expected.Insecure, actual.Insecure)
	assert.Equal(t, expected.Host, actual.Host)
}

func TestGetCredentials(t *testing.T) {
	fileContent := "username=testuser\npassword=testpass"
	tmpFile := createTempFile(t, "testcreds", fileContent)
	defer os.Remove(tmpFile.Name())
	expected := susemanager.SumanConfig{Login: "testuser", Password: "testpass", Insecure: true, Host: "mocked.hostname"}
	actual, err := GetCredentials(tmpFile.Name(), mockCommandRunner)
	require.NoError(t, err)
	assertExpected(t, &actual, &expected)
}

func TestGetCredentialsUyuni(t *testing.T) {
	fileContent := "user: testuser\npassword: testpass\nhubmaster: mocked.hostname"
	tmpFile := createTempFile(t, "testcredsuyuni", fileContent)
	defer os.Remove(tmpFile.Name())
	expected := susemanager.SumanConfig{Login: "testuser", Password: "testpass", Insecure: true, Host: "mocked.hostname"}
	actual, err := GetCredentialsUyuni(tmpFile.Name())
	require.NoError(t, err)
	assertExpected(t, &actual, &expected)
}
