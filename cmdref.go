package cmdref

import (
	"os"
	"path/filepath"
)

type Action int

const (
	Create Action = iota
	Update
	Remove
	View
	Exit
)

const CmdDirName = "cmdref"
const CmdFileName = "cmdref.json"

func CmdFilePath() (string, error) {
	d, err := os.UserConfigDir()
	if err != nil {
		return "", err
	}
	cmdStoreDir := filepath.Join(d, CmdDirName)
	err = CreateDirIfNotExists(cmdStoreDir)
	if err != nil {
		return "", err
	}
	return filepath.Join(cmdStoreDir, CmdFileName), nil
}

func CreateDirIfNotExists(dpath string) error {
	if stat, err := os.Stat(dpath); err == nil && stat.IsDir() {
		return nil // directory exists
	}

	//0755 Commonly used on web servers. The owner can read, write, execute. Everyone else can read and execute but not modify the file.
	//
	//0777 Everyone can read write and execute. On a web server, it is not advisable to use ‘777’ permission for your files and folders, as it allows anyone to add malicious code to your server.
	//
	//0644 Only the owner can read and write. Everyone else can only read. No one can execute the file.
	//
	//0655 Only the owner can read and write, but not execute the file. Everyone else can read and execute, but cannot modify the file.
	err := os.MkdirAll(dpath, 0777)
	if err != nil {
		return err
	}
	return nil
}
