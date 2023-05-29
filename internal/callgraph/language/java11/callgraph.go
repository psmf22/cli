package java

import (
	"embed"
	"io"
	"os"
	"path/filepath"

	"github.com/debricked/cli/internal/callgraph/cgexec"
)

//go:embed embeded/SootWrapper.jar
var jarCallGraph embed.FS

type Callgraph struct {
	cmdFactory       ICmdFactory
	workingDirectory string
	targetClasses    string
	targetDir        string
	ctx              cgexec.IContext
}

func (cg *Callgraph) runCallGraphWithSetup() error {
	jarFile, err := jarCallGraph.Open("embeded/SootWrapper.jar")
	if err != nil {
		return err
	}
	defer jarFile.Close()

	os.TempDir()
	tempDir, err := os.MkdirTemp("", "jar")
	if err != nil {
		return err
	}
	defer os.RemoveAll(tempDir)
	tempJarFile := filepath.Join(tempDir, "SootWrapper.jar")

	jarBytes, err := io.ReadAll(jarFile)
	if err != nil {

		return err
	}

	err = os.WriteFile(tempJarFile, jarBytes, 0600)
	if err != nil {

		return err
	}

	err = cg.runCallGraph(tempJarFile)

	return err
}

func (cg *Callgraph) runCallGraph(callgraphJarPath string) error {
	osCmd, err := cg.cmdFactory.MakeCallGraphGenerationCmd(callgraphJarPath, cg.workingDirectory, cg.targetClasses, cg.targetDir, cg.ctx)
	if err != nil {
		return err
	}

	cmd := cgexec.NewCommand(osCmd)
	err = cgexec.RunCommand(*cmd, cg.ctx)

	return err
}
