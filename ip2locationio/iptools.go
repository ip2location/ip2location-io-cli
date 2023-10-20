package main

import (
	"encoding/hex"
	"errors"
	"fmt"
	"math"
	"math/big"
	"math/rand"
	"net"
	"regexp"
	"sort"
	"strconv"
	"strings"
)

var maxIPv4Range *big.Int
var maxIPv6Range *big.Int

func RandIP() string {
	return strconv.Itoa(RandNum()) + "." + strconv.Itoa(RandNum()) + "." + strconv.Itoa(RandNum()) + "." + strconv.Itoa(RandNum())
}

func RandNum() int {
	min := 0
	max := 255
	return rand.Intn(max-min) + min
}

// IsIPv4 returns true if the IP address provided is an IPv4.
func IsIPv4(ip string) bool {
	ipaddr := net.ParseIP(ip)

	if ipaddr == nil {
		return false
	}

	v4 := ipaddr.To4()

	if v4 == nil {
		return false
	}

	return true
}

// IPv4ToDecimal returns the IP number for the supplied IPv4 address.
func IPv4ToDecimal(IP string) (*big.Int, error) {
	if !IsIPv4(IP) {
		return nil, errors.New("Not a valid IPv4 address.")
	}

	ipnum := big.NewInt(0)
	ipaddr := net.ParseIP(IP)

	if ipaddr != nil {
		v4 := ipaddr.To4()

		if v4 != nil {
			ipnum.SetBytes(v4)
		}
	}

	return ipnum, nil
}

// IPv6ToDecimal returns the IP number for the supplied IPv6 address.
func IPv6ToDecimal(IP string) (*big.Int, error) {
	if !IsIPv6(IP) {
		return nil, errors.New("Not a valid IPv6 address.")
	}

	ipnum := big.NewInt(0)
	ipaddr := net.ParseIP(IP)

	if ipaddr != nil {
		v6 := ipaddr.To16()

		if v6 != nil {
			ipnum.SetBytes(v6)
		}
	}

	return ipnum, nil
}

// IsIPv6 returns true if the IP address provided is an IPv6.
func IsIPv6(ip string) bool {
	if IsIPv4(ip) {
		return false
	}

	ipaddr := net.ParseIP(ip)

	if ipaddr == nil {
		return false
	}

	v6 := ipaddr.To16()

	if v6 == nil {
		return false
	}

	return true
}

// CIDRToIPv4 returns the IPv4 range for the supplied CIDR.
func CIDRToIPv4(CIDR string) ([]string, error) {
	if strings.Index(CIDR, "/") == -1 {
		return nil, errors.New("Not a valid CIDR.")
	}

	re := regexp.MustCompile(`^[0-9]{1,2}$`)
	arr := strings.Split(CIDR, "/")

	if len(arr) != 2 || !IsIPv4(arr[0]) || !re.MatchString(arr[1]) {
		return nil, errors.New("Not a valid CIDR.")
	}

	ip := arr[0]

	prefix, err := strconv.Atoi(arr[1])
	if err != nil || prefix > 32 {
		return nil, errors.New("Not a valid CIDR.")
	}

	ipstartbn, err := IPv4ToDecimal(ip)
	if err != nil {
		return nil, errors.New("Not a valid CIDR.")
	}
	ipstartlong := ipstartbn.Int64()

	ipstartlong = ipstartlong & (-1 << (32 - prefix))

	bn := big.NewInt(0)
	bn.SetString(strconv.Itoa(int(ipstartlong)), 10)

	ipstart, _ := DecimalToIPv4(bn)

	var total int64 = 1 << (32 - prefix)

	ipendlong := ipstartlong + total - 1

	if ipendlong > 4294967295 {
		ipendlong = 4294967295
	}

	bn.SetString(strconv.Itoa(int(ipendlong)), 10)
	ipend, _ := DecimalToIPv4(bn)

	result := []string{ipstart, ipend}

	return result, nil
}

// CIDRToIPv6 returns the IPv6 range for the supplied CIDR.
func CIDRToIPv6(CIDR string) ([]string, error) {
	if strings.Index(CIDR, "/") == -1 {
		return nil, errors.New("Not a valid CIDR.")
	}

	re := regexp.MustCompile(`^[0-9]{1,3}$`)
	arr := strings.Split(CIDR, "/")

	if len(arr) != 2 || !IsIPv6(arr[0]) || !re.MatchString(arr[1]) {
		return nil, errors.New("Not a valid CIDR.")
	}

	ip := arr[0]

	prefix, err := strconv.Atoi(arr[1])
	if err != nil || prefix > 128 {
		return nil, errors.New("Not a valid CIDR.")
	}

	hexstartaddress, _ := ExpandIPv6(ip)
	hexstartaddress = strings.ReplaceAll(hexstartaddress, ":", "")
	hexendaddress := hexstartaddress

	bits := 128 - prefix
	pos := 31

	for bits > 0 {
		values := []int{4, bits}
		min, _ := minMax(values)
		x, _ := strconv.ParseInt(string(hexendaddress[pos]), 16, 64)
		y := fmt.Sprintf("%x", (x | int64(math.Pow(2, float64(min))-1)))

		hexendaddress = hexendaddress[:pos] + y + hexendaddress[pos+1:]

		bits = bits - 4
		pos = pos - 1
	}

	re2 := regexp.MustCompile(`(.{4})`)
	hexstartaddress = re2.ReplaceAllString(hexstartaddress, "$1:")
	hexstartaddress = strings.TrimSuffix(hexstartaddress, ":")
	hexendaddress = re2.ReplaceAllString(hexendaddress, "$1:")
	hexendaddress = strings.TrimSuffix(hexendaddress, ":")

	result := []string{hexstartaddress, hexendaddress}

	return result, nil
}

// returns the min and max for the array
func minMax(array []int) (int, int) {
	var max int = array[0]
	var min int = array[0]
	for _, value := range array {
		if max < value {
			max = value
		}
		if min > value {
			min = value
		}
	}
	return min, max
}

// ExpandIPv6 returns the expanded form of the supplied IPv6 address.
func ExpandIPv6(IP string) (string, error) {
	if !IsIPv6(IP) {
		return "", errors.New("Not a valid IPv6 address.")
	}

	ipaddr := net.ParseIP(IP)

	ipstr := hex.EncodeToString(ipaddr)
	re := regexp.MustCompile(`(.{4})`)
	ipstr = re.ReplaceAllString(ipstr, "$1:")
	ipstr = strings.TrimSuffix(ipstr, ":")

	return ipstr, nil
}

// DecimalToIPv4 returns the IPv4 address for the supplied IP number.
func DecimalToIPv4(IPNum *big.Int) (string, error) {
	if IPNum.Cmp(big.NewInt(0)) < 0 || IPNum.Cmp(maxIPv4Range) > 0 {
		return "", errors.New("Invalid IP number.")
	}

	buf := make([]byte, 4)
	bytes := IPNum.FillBytes(buf)

	ip := net.IP(bytes)
	return ip.String(), nil
}

// DecimalToIPv6 returns the IPv6 address for the supplied IP number.
func DecimalToIPv6(IPNum *big.Int) (string, error) {
	if IPNum.Cmp(big.NewInt(0)) < 0 || IPNum.Cmp(maxIPv6Range) > 0 {
		return "", errors.New("Invalid IP number.")
	}

	buf := make([]byte, 16)
	bytes := IPNum.FillBytes(buf)

	ip := net.IP(bytes)
	return ip.String(), nil
}

// ListIPv4 returns the list of IP addresses for the supplied IPv4 range.
func ListIPv4(IPFrom string, IPTo string) ([]string, error) {
	if !IsIPv4(IPFrom) || !IsIPv4(IPTo) {
		return nil, errors.New("Not a valid IPv4 address.")
	}

	startipbig, _ := IPv4ToDecimal(IPFrom)
	endipbig, _ := IPv4ToDecimal(IPTo)
	var result []string

	one := big.NewInt(1)
	for endipbig.Cmp(startipbig) >= 0 {
		v4, err := DecimalToIPv4(startipbig)
		if err != nil {
			return result, err
		}
		result = append(result, v4)
		startipbig = startipbig.Add(startipbig, one)
	}

	return result, nil
}

// ListIPv6 returns the list of IP addresses for the supplied IPv6 range.
func ListIPv6(IPFrom string, IPTo string) ([]string, error) {
	if !IsIPv6(IPFrom) || !IsIPv6(IPTo) {
		return nil, errors.New("Not a valid IPv6 address.")
	}

	startipbig, _ := IPv6ToDecimal(IPFrom)
	endipbig, _ := IPv6ToDecimal(IPTo)
	var result []string

	one := big.NewInt(1)
	for endipbig.Cmp(startipbig) >= 0 {
		v6, err := DecimalToIPv6(startipbig)
		if err != nil {
			return result, err
		}
		result = append(result, v6)
		startipbig = startipbig.Add(startipbig, one)
	}

	return result, nil
}

// IPv4ToCIDR returns the CIDR for the supplied IPv4 range.
func IPv4ToCIDR(IPFrom string, IPTo string) ([]string, error) {
	if !IsIPv4(IPFrom) || !IsIPv4(IPTo) {
		return nil, errors.New("Not a valid IPv4 address.")
	}

	startipbig, _ := IPv4ToDecimal(IPFrom)
	endipbig, _ := IPv4ToDecimal(IPTo)
	startip := startipbig.Uint64()
	endip := endipbig.Uint64()
	var result []string
	var maxsize float64
	var maxdiff float64

	for endip >= startip {
		maxsize = 32

		for maxsize > 0 {
			mask := math.Pow(2, 32) - math.Pow(2, 32-(maxsize-1))
			maskbase := startip & uint64(mask)

			if maskbase != startip {
				break
			}

			maxsize = maxsize - 1
		}

		x := math.Log(float64(endip)-float64(startip)+1) / math.Log(2)
		maxdiff = 32 - math.Floor(x)

		if maxsize < maxdiff {
			maxsize = maxdiff
		}

		bn := big.NewInt(0)

		bn.SetString(fmt.Sprintf("%v", startip), 10)

		ip, _ := DecimalToIPv4(bn)
		result = append(result, ip+"/"+fmt.Sprintf("%v", maxsize))
		startip = startip + uint64(math.Pow(2, 32-maxsize))
	}

	return result, nil
}

// IPv6ToCIDR returns the CIDR for the supplied IPv6 range.
func IPv6ToCIDR(IPFrom string, IPTo string) ([]string, error) {
	if !IsIPv6(IPFrom) || !IsIPv6(IPTo) {
		return nil, errors.New("Not a valid IPv6 address.")
	}

	ipfrombin, err := ipToBinary(IPFrom)

	if err != nil {
		return nil, errors.New("Not a valid IPv6 address.")
	}

	iptobin, err := ipToBinary(IPTo)

	if err != nil {
		return nil, errors.New("Not a valid IPv6 address.")
	}

	var result []string

	networksize := 0
	shift := 0
	unpadded := ""
	padded := ""
	networks := make(map[string]int)
	n := 0

	if ipfrombin == iptobin {
		result = append(result, IPFrom+"/128")
		return result, nil
	}

	if ipfrombin > iptobin {
		tmp := ipfrombin
		ipfrombin = iptobin
		iptobin = tmp
	}

	for {
		if string(ipfrombin[len(ipfrombin)-1]) == "1" {
			unpadded = ipfrombin[networksize:128]
			padded = fmt.Sprintf("%-128s", unpadded)      // pad right with spaces
			padded = strings.ReplaceAll(padded, " ", "0") // replace spaces
			networks[padded] = 128 - networksize
			n = strings.LastIndex(ipfrombin, "0")
			if n == 0 {
				ipfrombin = ""
			} else {
				ipfrombin = ipfrombin[0:n]
			}
			ipfrombin = ipfrombin + "1"
			ipfrombin = fmt.Sprintf("%-128s", ipfrombin)        // pad right with spaces
			ipfrombin = strings.ReplaceAll(ipfrombin, " ", "0") // replace spaces
		}

		if string(iptobin[len(iptobin)-1]) == "0" {
			unpadded = iptobin[networksize:128]
			padded = fmt.Sprintf("%-128s", unpadded)      // pad right with spaces
			padded = strings.ReplaceAll(padded, " ", "0") // replace spaces
			networks[padded] = 128 - networksize
			n = strings.LastIndex(iptobin, "1")
			if n == 0 {
				iptobin = ""
			} else {
				iptobin = iptobin[0:n]
			}
			iptobin = iptobin + "0"
			iptobin = fmt.Sprintf("%-128s", iptobin)        // pad right with spaces
			iptobin = strings.ReplaceAll(iptobin, " ", "1") // replace spaces
		}

		if iptobin < ipfrombin {
			// special logic for Go due to lack of do-while
			if ipfrombin >= iptobin {
				break
			}
			continue
		}

		values := []int{strings.LastIndex(ipfrombin, "0"), strings.LastIndex(iptobin, "1")}
		_, max := minMax(values)
		shift = 128 - max
		unpadded = ipfrombin[0 : 128-shift]
		ipfrombin = fmt.Sprintf("%0128s", unpadded)
		unpadded = iptobin[0 : 128-shift]
		iptobin = fmt.Sprintf("%0128s", unpadded)

		networksize = networksize + shift

		if ipfrombin == iptobin {
			unpadded = ipfrombin[networksize:128]
			padded = fmt.Sprintf("%-128s", unpadded)      // pad right with spaces
			padded = strings.ReplaceAll(padded, " ", "0") // replace spaces
			networks[padded] = 128 - networksize
		}

		if ipfrombin >= iptobin {
			break
		}
	}

	keys := make([]string, 0, len(networks))
	for k := range networks {
		keys = append(keys, k)
	}
	sort.Strings(keys)

	for _, k := range keys {
		str, _ := binaryToIP(k)
		result = append(result, str+"/"+fmt.Sprintf("%d", networks[k]))
	}

	return result, nil
}

// converts IPv6 address to binary string representation.
func ipToBinary(ip string) (string, error) {
	if !IsIPv6(ip) {
		return "", errors.New("Not a valid IPv6 address.")
	}

	ipaddr := net.ParseIP(ip)

	binstr := ""
	for i, j := 0, len(ipaddr); i < j; i = i + 1 {
		binstr += fmt.Sprintf("%08b", ipaddr[i])
	}

	return binstr, nil
}

// converts binary string representation to IPv6 address.
func binaryToIP(binstr string) (string, error) {
	re := regexp.MustCompile(`^[01]{128}$`)
	if !re.MatchString(binstr) {
		return "", errors.New("Not a valid binary string.")
	}

	re2 := regexp.MustCompile(`(.{8})`)

	bytes := make([]byte, 16)
	i := 0
	matches := re2.FindAllStringSubmatch(binstr, -1)
	for _, v := range matches {
		x, _ := strconv.ParseUint(v[1], 2, 8)
		bytes[i] = byte(x)
		i = i + 1
	}

	ipaddr := net.IP(bytes)

	return ipaddr.String(), nil
}
