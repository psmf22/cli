package pip

import (
	"os/exec"
	"path/filepath"
)

type ICmdFactory interface {
	MakeCreateVenvCmd(file string) (*exec.Cmd, error)
	MakeActivateVenvCmd(file string) (*exec.Cmd, error)
	MakeInstallCmd(file string) (*exec.Cmd, error)
	MakeCatCmd(file string) (*exec.Cmd, error)
	MakeListCmd() (*exec.Cmd, error)
	MakeShowCmd(list []string) (*exec.Cmd, error)
}

type CmdFactory struct{}

func (_ CmdFactory) MakeCreateVenvCmd(file string) (*exec.Cmd, error) {
	path, err := exec.LookPath("python")

	fpath := filepath.Join(filepath.Dir(file), filepath.Base(file)+".venv")

	return &exec.Cmd{
		Path: path,
		Args: []string{"python", "-m", "venv", fpath, "--clear"},
	}, err
}

func (_ CmdFactory) MakeActivateVenvCmd(file string) (*exec.Cmd, error) {
	path, err := exec.LookPath("bash")

	fpath := filepath.Join(filepath.Dir(file), filepath.Base(file)+".venv", "bin", "activate")

	return &exec.Cmd{
		Path: path,
		Args: []string{"bash", "-c", "source " + fpath},
	}, err
}

func (_ CmdFactory) MakeInstallCmd(file string) (*exec.Cmd, error) {
	path, err := exec.LookPath("pip")

	return &exec.Cmd{
		Path: path,
		Args: []string{"pip", "install", "-r", file},
	}, err
}

func (_ CmdFactory) MakeCatCmd(file string) (*exec.Cmd, error) {
	path, err := exec.LookPath("cat")

	return &exec.Cmd{
		Path: path,
		Args: []string{"cat", file},
	}, err
}

func (_ CmdFactory) MakeListCmd() (*exec.Cmd, error) {
	path, err := exec.LookPath("pip")

	return &exec.Cmd{
		Path: path,
		Args: []string{"pip", "list"},
	}, err
}

func (_ CmdFactory) MakeShowCmd(list []string) (*exec.Cmd, error) {
	path, err := exec.LookPath("pip")

	args := []string{"pip", "show"}
	args = append(args, list...)

	return &exec.Cmd{
		Path: path,
		Args: args,
	}, err
}