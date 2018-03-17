package module

import (
	"fmt"
	"net"
	"strconv"

	"github.com/l-dandelion/webcrawler/errors"
)

type mAddr struct {
	network string
	address string
}

func (maddr *mAddr) Network() string {
	return maddr.network
}

func (maddr *mAddr) String() string {
	return maddr.address
}

func NewAddr(network, ip string, port uint64) (net.Addr, error) {
	if network != "http" && network != "https" {
		errMsg := fmt.Sprintf("illegal network for module address: %s", network)
		return nil, errors.NewIllegalParameterError(errMsg)
	}
	if parseIP := net.ParseIP(ip); parseIP == nil {
		errMsg := fmt.Sprintf("illegal IP for module address: %s", ip)
		return nil, errors.NewIllegalParameterError(errMsg)
	}
	return &mAddr{
		network: network,
		address: ip + ":" + strconv.Itoa(int(port)),
	}, nil
}
