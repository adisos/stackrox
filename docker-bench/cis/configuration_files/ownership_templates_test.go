package configurationfiles

import (
	"io/ioutil"
	"os"
	"os/user"
	"path/filepath"
	"strconv"
	"testing"

	"bitbucket.org/stack-rox/apollo/pkg/api/generated/api/v1"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func getInt(val string) int {
	v, _ := strconv.Atoi(val)
	return v
}

func setOwnership(file, u, g string) error {
	us, err := user.Lookup(u)
	if err != nil {
		return err
	}
	group, err := user.LookupGroup(g)
	if err != nil {
		return err
	}
	err = os.Chown(file, getInt(group.Gid), getInt(us.Uid))
	return err
}

// Returns file name with these ownership settings
func createTestFileOwnership(u, g string) (string, error) {
	file, err := ioutil.TempFile("", "")
	if err != nil {
		return "", err
	}
	if err := setOwnership(file.Name(), u, g); err != nil {
		return "", err
	}
	return file.Name(), nil
}

func createTestDir() (dir string, fileA string, fileB string, err error) {
	dir, err = ioutil.TempDir("", "")
	if err != nil {
		return
	}
	fileA = filepath.Join(dir, "a.txt")
	if err = ioutil.WriteFile(fileA, []byte("hello world"), 0777); err != nil {
		return
	}
	fileB = filepath.Join(dir, "b.txt")
	if err = ioutil.WriteFile(fileB, []byte("hello world"), 0777); err != nil {
		return
	}
	return
}

// Returns file name with these ownership settings
func createTestDirOwnership(u, g string) (dir string, fileA string, fileB string, err error) {
	dir, a, b, err := createTestDir()
	if err = setOwnership(a, u, g); err != nil {
		return
	}
	if err = setOwnership(b, u, g); err != nil {
		return
	}
	return
}

func TestCompareFileOwnership(t *testing.T) {
	currentUser, err := user.Current()
	require.Nil(t, err)

	// Match the ownership
	expectedResult := v1.BenchmarkTestResult{Result: v1.BenchmarkStatus_PASS}
	f, err := createTestFileOwnership(currentUser.Username, currentUser.Username)
	require.Nil(t, err)
	result := compareFileOwnership(f, currentUser.Username, currentUser.Username)
	assert.Equal(t, expectedResult.Result, result.Result)
	assert.Equal(t, 0, len(result.Notes))

	// Testing for this is hard because of the rules of user:group. The base user requires sudo to chown
	// Compare to a file that doesn't have the right ownership
	// Match the user but not the group
	expectedResult = v1.BenchmarkTestResult{Result: v1.BenchmarkStatus_WARN}
	result = compareFileOwnership("/etc/passwd", currentUser.Username, currentUser.Username)
	result = compareFileOwnership(f, currentUser.Username, "docker")
	assert.Equal(t, expectedResult.Result, result.Result)
	assert.Equal(t, 1, len(result.Notes))
}

func TestFileOwnershipCheck(t *testing.T) {
	// Set up file to check against
	currentUser, err := user.Current()
	require.Nil(t, err)
	expectedResult := v1.BenchmarkTestResult{Result: v1.BenchmarkStatus_PASS}
	f, err := createTestFileOwnership(currentUser.Username, currentUser.Username)
	require.Nil(t, err)

	benchmark := newOwnershipCheck("Test bench", "desc", f, currentUser.Username, currentUser.Username)
	result := benchmark.Run()
	assert.Equal(t, expectedResult.Result, result.Result)
	assert.Equal(t, 0, len(result.Notes))

	// Check empty file
	expectedResult = v1.BenchmarkTestResult{Result: v1.BenchmarkStatus_NOTE}
	benchmark = newOwnershipCheck("Test bench", "desc", "", currentUser.Username, currentUser.Username)
	result = benchmark.Run()
	assert.Equal(t, expectedResult.Result, result.Result)
	assert.Equal(t, 1, len(result.Notes))
}

func TestRecursiveOwnershipCheck(t *testing.T) {
	currentUser, err := user.Current()
	require.Nil(t, err)
	dir, _, _, err := createTestDirOwnership(currentUser.Username, currentUser.Username)
	require.Nil(t, err)

	//expectedResult := v1.BenchmarkTestResult{Result: v1.BenchmarkStatus_PASS}
	expectedResult := v1.BenchmarkTestResult{Result: v1.BenchmarkStatus_PASS}
	benchmark := newRecursiveOwnershipCheck("test bench", "desc", dir, currentUser.Username, currentUser.Username)
	result := benchmark.Run()
	assert.Equal(t, expectedResult.Result, result.Result)
	assert.Equal(t, 0, len(result.Notes))

	// Check empty file
	expectedResult = v1.BenchmarkTestResult{Result: v1.BenchmarkStatus_NOTE}
	benchmark = newRecursiveOwnershipCheck("Test bench", "desc", "", currentUser.Username, currentUser.Username)
	result = benchmark.Run()
	assert.Equal(t, expectedResult.Result, result.Result)
	assert.Equal(t, 1, len(result.Notes))
}
