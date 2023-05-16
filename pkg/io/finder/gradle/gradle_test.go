package gradle

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"testing"

	writerTestdata "github.com/debricked/cli/pkg/io/writer/testdata"

	"github.com/stretchr/testify/assert"
)

func TestNewGradleSetup(t *testing.T) {

	gs := NewGradleSetup()
	assert.NotNil(t, gs)
}

func TestErrors(t *testing.T) {

	walkError := SetupWalkError{message: "test"}
	assert.Equal(t, "test", walkError.Error())

	scriptError := SetupScriptError{message: "test"}
	assert.Equal(t, "test", scriptError.Error())

	subprojectError := SetupSubprojectError{message: "test"}
	assert.Equal(t, "test", subprojectError.Error())

}

func TestSetupFilePathMappings(t *testing.T) {
	gs := NewGradleSetup()
	files := []string{filepath.Join("testdata", "project", "build.gradle")}
	gs.setupFilePathMappings(files)

	assert.Len(t, gs.GradlewMap, 1)
	assert.Len(t, gs.SettingsMap, 1)
}

func TestSetupFilePathMappingsNoFiles(t *testing.T) {
	gs := NewGradleSetup()
	gs.setupFilePathMappings([]string{})

	assert.Len(t, gs.GradlewMap, 0)
	assert.Len(t, gs.SettingsMap, 0)
}

func TestSetupFilePathMappingsNoGradlew(t *testing.T) {
	gs := NewGradleSetup()
	files := []string{filepath.Join("testdata", "project", "subproject", "build.gradle")}
	gs.setupFilePathMappings(files)

	assert.Len(t, gs.GradlewMap, 0)
	assert.Len(t, gs.SettingsMap, 0)
}

func TestSetupGradleProjectMappings(t *testing.T) {
	gs := NewGradleSetup()
	gs.CmdFactory = &mockCmdFactory{}

	gs.SettingsMap = map[string]string{
		filepath.Join("testdata", "project"): filepath.Join("testdata", "project", "settings.gradle"),
	}
	gs.SubProjectMap = map[string]string{}
	err := gs.setupGradleProjectMappings()
	// assert GradleSetupSubprojectError
	assert.NotNil(t, err)

	assert.Len(t, gs.GradleProjects, 1)
}

type mockCmdFactory struct {
	createFile bool
}

func (m *mockCmdFactory) MakeFindSubGraphCmd(workingDirectory string, _ string, _ string) (*exec.Cmd, error) {
	if m.createFile {
		fileName := filepath.Join(workingDirectory, multiProjectFilename)
		content := []byte(workingDirectory)
		file, err := os.Create(fileName)
		if err != nil {

			return nil, err
		}
		defer file.Close()
		_, err = file.Write(content)
		if err != nil {

			return nil, err
		}
	}
	// if windows use dir
	if runtime.GOOS == "windows" {
		// gradlewOsName = "gradlew.bat"
		return exec.Command("dir"), nil
	}

	return exec.Command("ls"), nil
}

func TestSetupSubProjectPathsNoFileCreated(t *testing.T) {
	gs := NewGradleSetup()
	gs.CmdFactory = &mockCmdFactory{createFile: false}

	absPath, _ := filepath.Abs(filepath.Join("testdata", "project"))
	gradleProject := Project{Dir: absPath, Gradlew: filepath.Join("testdata", "project", "gradlew")}
	err := gs.setupSubProjectPaths(gradleProject)
	fmt.Println(err)
	assert.NotNil(t, err)
	assert.Len(t, gs.SubProjectMap, 0)
}

func TestSetupSubProjectPaths(t *testing.T) {
	gs := NewGradleSetup()
	gs.CmdFactory = &mockCmdFactory{createFile: true}

	absPath, _ := filepath.Abs(filepath.Join("testdata", "project"))
	gradleProject := Project{Dir: absPath, Gradlew: filepath.Join("testdata", "project", "gradlew")}
	err := gs.setupSubProjectPaths(gradleProject)
	assert.Nil(t, err)
	assert.Len(t, gs.SubProjectMap, 1)

	absPath, _ = filepath.Abs(filepath.Join("testdata", "project", "subproject"))
	gradleProject = Project{Dir: absPath, Gradlew: filepath.Join("testdata", "project", "gradlew")}
	err = gs.setupSubProjectPaths(gradleProject)
	assert.Nil(t, err)
	assert.Len(t, gs.SubProjectMap, 2)
}

func TestSetupSubProjectPathsError(t *testing.T) {
	gs := NewGradleSetup()

	absPath, _ := filepath.Abs(filepath.Join("testdata", "project"))
	gradleProject := Project{Dir: absPath, Gradlew: filepath.Join("testdata", "project", "gradlew")}
	err := gs.setupSubProjectPaths(gradleProject)

	assert.NotNil(t, err)
}

func TestGetGradleW(t *testing.T) {
	gs := NewGradleSetup()

	gs.GradlewMap = map[string]string{
		filepath.Join("testdata", "project"): filepath.Join("testdata", "project", "gradlew"),
	}

	gradlew := gs.GetGradleW(filepath.Join("testdata", "project", "subproject"))

	assert.Equal(t, filepath.Join("testdata", "project", "gradlew"), gradlew)

	gradlew = gs.GetGradleW(filepath.Join("testdata", "project"))

	assert.Equal(t, filepath.Join("testdata", "project", "gradlew"), gradlew)
}

type mockInitScriptHandler struct {
	writeInitFileErr error
}

func (_ mockInitScriptHandler) ReadInitFile() ([]byte, error) {
	return gradleInitScript.ReadFile("gradle-init/gradle-init-script.groovy")
}

func (i mockInitScriptHandler) WriteInitFile() error {
	return i.writeInitFileErr
}

type mockFileHandler struct {
	setupWalkErr error
}

func (f mockFileHandler) Find(_ []string) (map[string]string, map[string]string, error) {
	return nil, nil, f.setupWalkErr
}

func TestConfigureErrors(t *testing.T) {
	gs := NewGradleSetup()
	gs.Writer = &writerTestdata.FileWriterMock{}
	err := gs.Configure([]string{"testdata/project"})
	assert.NotNil(t, err)

	gs.MetaFileFinder = mockFileHandler{setupWalkErr: SetupScriptError{message: "mock error"}}
	err = gs.Configure([]string{"testdata/project"})
	assert.Equal(t, "mock error", err.Error())

	gs.InitScriptHandler = mockInitScriptHandler{writeInitFileErr: SetupScriptError{message: "write-init-file-err"}}
	err = gs.Configure([]string{"testdata/project"})
	assert.Equal(t, "write-init-file-err", err.Error())
}

func TestConfigure(t *testing.T) {
	gs := NewGradleSetup()
	gs.Writer = &writerTestdata.FileWriterMock{}
	gs.MetaFileFinder = mockFileHandler{setupWalkErr: nil}
	gs.InitScriptHandler = mockInitScriptHandler{writeInitFileErr: nil}

	err := gs.Configure([]string{"testdata/project"})
	assert.NoError(t, err)
}