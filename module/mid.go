package module

import (
	"fmt"
	"net"
	"strconv"
	"strings"

	"github.com/l-dandelion/webcrawler/errors"
)

type MID string

var (
	DefaultSNGen = NewSNGenertor(1, 0)
	midTemplate  = "%s%d|%s"
)

func GenMID(mtype Type, sn uint64, maddr net.Addr) (MID, error) {
	if !LegalType(mtype) {
		errMsg := fmt.Sprintf("illegal module type: %s", mtype)
		return "", errors.NewIllegalParameterError(errMsg)
	}
	letter := legalTypeLetterMap[mtype]
	var midStr string
	if maddr == nil {
		midStr = fmt.Sprintf(midTemplate, letter, sn, "")
		midStr = midStr[:len(midStr)-1]
	} else {
		midStr = fmt.Sprintf(midTemplate, letter, sn, maddr.String())
	}
	return MID(midStr), nil
}

func LegalMID(mid MID) bool {
	if _, err := SplitMID(mid); err == nil {
		return true
	}
	return false
}

func SplitMID(mid MID) ([]string, error) {
	var letter, snStr, addr string
	midStr := string(mid)
	if len(midStr) <= 1 {
		return nil, errors.NewIllegalParameterError("insufficient MID")
	}
	letter = midStr[:1]
	if !LegalLetter(letter) {
		return nil, errors.NewIllegalParameterError(fmt.Sprintf("illegal module type letter: %s", letter))
	}
	snAndAddr := midStr[1:]
	index := strings.LastIndex(snAndAddr, "|")
	if index < 0 {
		snStr = snAndAddr
		if !legalSN(snStr) {
			return nil, errors.NewIllegalParameterError(fmt.Sprintf("illegal module SN: %s", snStr))
		}
	} else {
		snStr = snAndAddr[:index]
		if !legalSN(snStr) {
			return nil, errors.NewIllegalParameterError(fmt.Sprintf("illegal module SN: %s", snStr))
		}
		addr = snAndAddr[index+1:]
		index = strings.LastIndex(addr, ":")
		if index <= 0 {
			return nil, errors.NewIllegalParameterError(fmt.Sprintf("illegal module address: %s", addr))
		}
		ipStr := addr[:index]
		if ip := net.ParseIP(ipStr); ip == nil {
			return nil, errors.NewIllegalParameterError(fmt.Sprintf("illegal module IP: %s", ipStr))
		}
		portStr := addr[index+1:]
		if _, err := strconv.ParseUint(portStr, 10, 64); err != nil {
			return nil, errors.NewIllegalParameterError(fmt.Sprintf("illegal module port: %s", portStr))
		}
	}
	return []string{letter, snStr, addr}, nil
}

func legalSN(snStr string) bool {
	_, err := strconv.ParseUint(snStr, 10, 64)
	return err == nil
}
