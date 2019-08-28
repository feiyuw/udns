package utils

import (
	"os/exec"
	"strings"
)

func GetDefaultGateway() string {
	cmd := exec.Command("/sbin/route", "-n", "get", "0.0.0.0")
	output, err := cmd.CombinedOutput()
	if err != nil {
		return ""
	}
	for _, line := range strings.Split(string(output), "\n") {
		fields := strings.Fields(line)
		if len(fields) >= 2 && fields[0] == "gateway:" {
			return fields[1]
		}
	}
	return ""
}
