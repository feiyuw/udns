package utils

import (
	"bufio"
	"encoding/binary"
	"net"
	"os"
	"strconv"
	"strings"

	"udns/logger"
)

const (
	file  = "/proc/net/route"
	line  = 1    // line containing the gateway addr. (first line: 0)
	sep   = "\t" // field separator
	field = 2    // field containing hex gateway address (first field: 0)
)

// GetDefaultGateway is used to fetch default gateway IP address
func GetDefaultGateway() string {
	file, err := os.Open(file)
	if err != nil {
		return ""
	}
	defer func() {
		if err = file.Close(); err != nil {
			logger.Error("gateway", err)
		}
	}()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		// jump to line containing the agteway address
		for i := 0; i < line; i++ {
			scanner.Scan()
		}

		// get field containing gateway address
		tokens := strings.Split(scanner.Text(), sep)
		gatewayHex := "0x" + tokens[field]

		// cast hex address to uint32
		d, _ := strconv.ParseInt(gatewayHex, 0, 64)
		d32 := uint32(d)

		// make net.IP address from uint32
		ipd32 := make(net.IP, 4)
		binary.LittleEndian.PutUint32(ipd32, d32)

		// format net.IP to dotted ipV4 string
		return net.IP(ipd32).String()
	}
	return ""
}
