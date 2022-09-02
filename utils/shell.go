package utils

import (
	"bytes"
	"os/exec"
)

const (
	shellBin = "bash"
)

var (
	TrojanConf TrojanConfig
)

type TrojanConfig struct {
	BinFile string `yaml:"bin-file"`
	ApiAddr string `yaml:"api-addr"`
	ApiPort string `yaml:"api-port"`
}

func Shell(command string) (error, string, string) {
	var stdout bytes.Buffer
	var stderr bytes.Buffer
	cmd := exec.Command(shellBin, "-c", command)
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr
	err := cmd.Run()
	return err, stdout.String(), stderr.String()
}

func RunTrojanCLI(command string) (error, string, string) {
	return Shell(TrojanConf.BinFile + " -api-addr " + TrojanConf.ApiAddr + ":" + TrojanConf.ApiPort +
		" -api " + command)
}
