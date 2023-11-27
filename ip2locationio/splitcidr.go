package main

import (
	"encoding/binary"
	"errors"
	"math"
	"math/big"
	"net"
	"strconv"
)

func SplitCIDR(cidr string, split string) ([]string, error) {
	var res []string
	ip, _, err := net.ParseCIDR(cidr)
	if err != nil {
		return nil, err
	}

	toSplit, err := strconv.Atoi(split)
	if err != nil {
		return nil, err
	}

	if ip.To4() != nil {
		ipsubnet, err := GetIPv4Subnet(cidr)
		if err != nil {
			return nil, err
		}

		subs, err := SplitCIDRIPv4(ipsubnet, toSplit)
		if err != nil {
			return nil, err
		}

		for _, s := range subs {
			netBitCntStr := strconv.Itoa(int(s.NetBitCnt))
			ipStr := net.IPv4(byte(s.LoIP>>24), byte(s.LoIP>>16), byte(s.LoIP>>8), byte(s.LoIP)).String()

			res = append(res, ipStr+"/"+netBitCntStr)
		}
	} else {
		ipsubnet, err := GetIPv6Subnet(cidr)
		if err != nil {
			return nil, err
		}
		subs, err := SplitCIDRIPv6(ipsubnet, toSplit)
		if err != nil {
			return nil, err
		}
		for _, s := range subs {
			netBitCntStr := strconv.Itoa(int(s.NetBitCnt))
			b := make([]byte, 16)
			binary.BigEndian.PutUint64(b[:8], s.LoIP.Hi)
			binary.BigEndian.PutUint64(b[8:], s.LoIP.Lo)
			ipStr := net.IP(b).String()
			res = append(res, ipStr+"/"+netBitCntStr)
		}
	}

	return res, nil
}

type uint128 struct {
	Hi uint64
	Lo uint64
}

type IPv4Subnet struct {
	NetBitCnt  uint32
	NetMask    uint32
	HostBitCnt uint32
	HostMask   uint32
	LoIP       uint32
	HiIP       uint32
}

type IPv6Subnet struct {
	NetBitCnt  uint32
	NetMask    uint128
	HostBitCnt uint32
	HostMask   uint128
	LoIP       uint128
	HiIP       uint128
}

func NetAndHostMasksIPv4(size uint32) (uint32, uint32) {
	if size > 32 {
		size = 32
	}

	var mask uint32 = 0
	for i := uint32(0); i < size; i++ {
		mask += 1 << (32 - (i + 1))
	}

	return mask, ^mask
}

func NetAndHostMasksIPv6(size uint32) (uint128, uint128) {
	if size > 128 {
		size = 128
	}

	var mask uint128
	if size > 64 {
		mask = uint128{^uint64(0), ^uint64(0) << (128 - size)}
	} else {
		mask = uint128{^uint64(0) << (64 - size), 0}
	}
	maskNot := uint128{^mask.Hi, ^mask.Lo}

	return mask, maskNot
}

func GetIPv4Subnet(cidr string) (IPv4Subnet, error) {
	_, network, err := net.ParseCIDR(cidr)
	if err != nil {
		return IPv4Subnet{}, err
	}

	ones, _ := network.Mask.Size()
	netMask, hostMask := NetAndHostMasksIPv4(uint32(ones))
	start := binary.BigEndian.Uint32(network.IP)
	ipsubnet := IPv4Subnet{
		HostBitCnt: uint32(32 - ones),
		HostMask:   hostMask,
		NetBitCnt:  uint32(ones),
		NetMask:    netMask,
		LoIP:       uint32(start) & netMask,
		HiIP:       (uint32(start) & netMask) | hostMask,
	}

	return ipsubnet, nil
}

func GetIPv6Subnet(cidr string) (IPv6Subnet, error) {
	_, network, err := net.ParseCIDR(cidr)
	if err != nil {
		return IPv6Subnet{}, err
	}

	ones, _ := network.Mask.Size()
	netMask, hostMask := NetAndHostMasksIPv6(uint32(ones))
	starthi := binary.BigEndian.Uint64(network.IP[:8])
	startlo := binary.BigEndian.Uint64(network.IP[8:])
	start := uint128{Hi: starthi, Lo: startlo}
	ip6subnet := IPv6Subnet{
		HostBitCnt: uint32(128 - ones),
		HostMask:   hostMask,
		NetBitCnt:  uint32(ones),
		NetMask:    netMask,
		LoIP:       uint128{Hi: start.Hi & netMask.Hi, Lo: start.Lo & netMask.Lo},
		HiIP:       uint128{Hi: (start.Hi & netMask.Hi) | hostMask.Hi, Lo: (start.Lo & netMask.Lo) | hostMask.Lo},
	}

	return ip6subnet, nil
}

func SplitCIDRIPv4(s IPv4Subnet, split int) ([]IPv4Subnet, error) {
	bitshifts := int(uint32(split) - s.NetBitCnt)
	if bitshifts < 0 || bitshifts > 31 || int(s.NetBitCnt)+bitshifts > 32 {
		return nil, errors.New("Invalid split.")
	}

	hostBits := (32 - s.NetBitCnt) - uint32(bitshifts)
	netMask, hostMask := NetAndHostMasksIPv4(uint32(split))
	ipsubnets := make([]IPv4Subnet, 1<<bitshifts)
	for i := range ipsubnets {
		start := uint32(s.LoIP) + uint32(i*(1<<hostBits))
		ipsubnets[i] = IPv4Subnet{
			HostBitCnt: uint32(32 - split),
			HostMask:   hostMask,
			NetBitCnt:  uint32(split),
			LoIP:       uint32(start) & netMask,
			HiIP:       (uint32(start) & netMask) | hostMask,
		}
	}

	return ipsubnets, nil
}

func SplitCIDRIPv6(s IPv6Subnet, split int) ([]IPv6Subnet, error) {
	bitshifts := int(uint32(split) - s.NetBitCnt)
	if bitshifts < 0 || bitshifts > 128 || int(s.NetBitCnt)+bitshifts > 128 {
		return nil, errors.New("Invalid split. ")
	}

	hostBits := (128 - s.NetBitCnt) - uint32(bitshifts)
	netMask, hostMask := NetAndHostMasksIPv6(uint32(split))
	subnetCount := math.Pow(2, float64(bitshifts))
	subnetCountBig := big.NewFloat(subnetCount)
	hostCount := math.Pow(2, float64(hostBits))
	hostCountbig := big.NewFloat(hostCount)

	var ipsubnets []IPv6Subnet
	for i := big.NewFloat(0); i.Cmp(subnetCountBig) < 0; i.Add(i, big.NewFloat(1)) {
		hostCountMul := new(big.Float)
		hostCountMul.Mul(i, hostCountbig)
		newIP := new(big.Int)
		b := make([]byte, 16)
		binary.BigEndian.PutUint64(b[:8], s.LoIP.Hi)
		binary.BigEndian.PutUint64(b[8:], s.LoIP.Lo)
		newIP.SetBytes(b)
		hostCountAdd := new(big.Int)
		result := new(big.Int)
		hostCountMul.Int(result)
		hostCountAdd.Add(newIP, result)

		ipArr := []byte{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0}
		copy(ipArr[16-len(hostCountAdd.Bytes()):], hostCountAdd.Bytes())
		startIP := uint128{Hi: binary.BigEndian.Uint64(ipArr[:8]), Lo: binary.BigEndian.Uint64(ipArr[8:])}

		subnet := IPv6Subnet{
			HostBitCnt: uint32(128 - split),
			HostMask:   hostMask,
			NetBitCnt:  uint32(split),
			LoIP:       uint128{Hi: startIP.Hi & netMask.Hi, Lo: startIP.Lo & netMask.Lo},
			HiIP:       uint128{Hi: (startIP.Hi & netMask.Hi) | hostMask.Hi, Lo: (startIP.Lo & netMask.Lo) | hostMask.Lo},
		}
		ipsubnets = append(ipsubnets, subnet)
	}

	return ipsubnets, nil
}
