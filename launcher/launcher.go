package launcher

import (
	"github.com/screepers/screeps-launcher/v1/install"
	"log"
	"os"
	"os/exec"
	// "path/filepath"
	// "strings"
)

type Launcher struct {
	config    *Config
	needsInit bool
}

func (l *Launcher) Prepare() {
	c := NewConfig()
	_, err := c.GetConfig()
	if err != nil {
		log.Fatalf("Error loading config: %v", err)
	}
	l.config = c
	l.needsInit = false
	if _, err := os.Stat("package.json"); os.IsNotExist(err) {
		l.needsInit = true
	}
}

func (l *Launcher) Upgrade() error {
	os.Remove("yarn.lock")
	return l.Apply()
}

func (l *Launcher) Apply() error {
	var err error
	if _, err := os.Stat(install.GetNodePath()); os.IsNotExist(err) {
		log.Print("Installing Node")
		err = install.Node("Carbon")
		if err != nil {
			return err
		}
		// This requires an admin prompt, need to figure out howto prompt the user.
		// if runtime.GOOS == "windows" {
		// 	log.Print("Installing windows-build-tools (This may take a while)")
		// 	err = install.WindowsBuildTools()
		// 	if err != nil {
		// 		return err
		// 	}
		// }
	}
	if _, err := os.Stat("deps/yarn/bin/yarn"); os.IsNotExist(err) {
		log.Print("Installing Yarn")
		err = install.Yarn()
		if err != nil {
			return err
		}
	}
	log.Print("Writing package.json")
	err = writePackage(l.config)
	if err != nil {
		return err
	}
	log.Print("Running yarn")
	err = runYarn()
	if err != nil {
		return err
	}
	if l.needsInit {
		log.Print("Initializing server")
		initServer(l.config)
	}
	log.Print("Writing mods.json")
	err = writeMods(l.config)
	if err != nil {
		return err
	}
	return nil
}

func (l *Launcher) Start() error {
	err := l.Apply()
	if err != nil {
		return err
	}
	log.Print("Starting Server")
	runServer(l.config)
	return nil
}

func (l *Launcher) Cli() error {
	log.Print("Starting CLI")
	runCli(l.config)
	return nil
}

func cmdExists(cmdName string) bool {
	_, err := exec.LookPath(cmdName)
	return err == nil
}

func runYarn() error {
	cmd := exec.Command(install.GetNodePath(), install.GetYarnPath())
	//newPath := filepath.SplitList(os.Getenv("PATH"))
	//cwd, _ := os.Getwd()
	//nodePath := filepath.Dir(install.GetNodePath())
	//newPath = append([]string{filepath.Join(cwd, nodePath)}, newPath...)
	//cmd.Env = append(cmd.Env, "PATH="+strings.Join(newPath, string(filepath.ListSeparator)))
	//log.Print(cmd.Env)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	err := cmd.Run()
	return err
}