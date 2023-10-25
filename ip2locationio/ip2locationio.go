package main

import (
	"flag"
	"fmt"
	"math/big"
	"os"
	"reflect"
	"strings"
)

var outputFormat string
var apiKey string
var myLanguage string
var myIP string
var filterFields string

const version string = "1.0.1"
const programName string = "IP2Location.io Command Line"

var showVer bool = false

func init() {
	maxIPv4Range = big.NewInt(4294967295)
	maxIPv6Range = big.NewInt(0)
	maxIPv6Range.SetString("340282366920938463463374607431768211455", 10)

	// read config for API key if exist
	LoadConfig()
}

func main() {
	flag.StringVar(&outputFormat, "o", "json", "Output format: json | pretty")
	flag.StringVar(&apiKey, "k", "", "API key: Get your API key from https://ip2location.io")
	flag.StringVar(&myLanguage, "l", "", "Language: ar | cs | da | de | en | es | et | fi | fr | ga | it | ja | ko | ms | nl | pt | ru | sv | tr | vi | zh-cn | zh-tw")
	flag.StringVar(&filterFields, "f", "", `Filter fields: Field names separted by comma. E.g., "country_code,city_name,continent.name"`)
	flag.BoolVar(&showVer, "v", false, "Show version")

	flag.Usage = func() {
		PrintUsage()
	}
	flag.Parse()

	if showVer {
		PrintVersion()
		return
	}

	if apiKey == "" {
		apiKey = config.APIKey
	}

	var arg = flag.Arg(0)

	if arg == "config" {
		UpdateAPIKey(flag.Arg(1))
		return
	} else if arg == "randip" {
		PrintRandIP()
		return
	} else if arg == "cidr2list" {
		PrintCIDR2List(flag.Arg(1))
		return
	} else if arg == "range2list" {
		PrintRange2List(flag.Arg(1), flag.Arg(2))
		return
	} else if arg == "cidr2range" {
		PrintCIDR2Range(flag.Arg(1))
		return
	} else if arg == "range2cidr" {
		PrintRange2CIDR(flag.Arg(1), flag.Arg(2))
		return
	} else if len(arg) == 0 {
		myIP = MyPublicIP()
	} else if !IsIPv4(arg) && !IsIPv6(arg) {
		fmt.Println("Not a valid IP address.")
		return
	} else {
		myIP = arg
	}

	filterFields = strings.TrimSpace(filterFields)

	if filterFields != "" {
		PrintFiltered()
	} else {
		PrintNormal()
	}
}

func PrintVersion() {
	fmt.Printf("%s Version: %s\n", programName, version)
}

func PrintRandIP() {
	fmt.Printf("%s\n", RandIP())
}

func PrintCIDR2List(cidr string) {
	res, err := CIDRToIPv4(cidr)

	if err != nil {
		res, err := CIDRToIPv6(cidr)

		if err != nil {
			fmt.Println(err)
		} else {
			res, err := ListIPv6(res[0], res[1])
			if err != nil {
				fmt.Println(err)
			} else {
				for _, element := range res {
					fmt.Println(element)
				}
			}
		}
	} else {
		res, err := ListIPv4(res[0], res[1])
		if err != nil {
			fmt.Println(err)
		} else {
			for _, element := range res {
				fmt.Println(element)
			}
		}
	}
}

func PrintRange2List(fromIP string, toIP string) {
	res, err := ListIPv4(fromIP, toIP)

	if err != nil {
		res, err := ListIPv6(fromIP, toIP)

		if err != nil {
			fmt.Println("Invalid IP addresses.")
		} else {
			for _, element := range res {
				fmt.Println(element)
			}
		}
	} else {
		for _, element := range res {
			fmt.Println(element)
		}
	}
}

func PrintCIDR2Range(cidr string) {
	res, err := CIDRToIPv4(cidr)

	if err != nil {
		res, err := CIDRToIPv6(cidr)

		if err != nil {
			fmt.Println(err)
		} else {
			fmt.Printf("%s-%s\n", res[0], res[1])
		}
	} else {
		fmt.Printf("%s-%s\n", res[0], res[1])
	}
}

func PrintRange2CIDR(fromIP string, toIP string) {
	res, err := IPv4ToCIDR(fromIP, toIP)

	if err != nil {
		res, err := IPv6ToCIDR(fromIP, toIP)

		if err != nil {
			fmt.Println("Invalid IP addresses.")
		} else {
			for _, element := range res {
				fmt.Println(element)
			}
		}
	} else {
		for _, element := range res {
			fmt.Println(element)
		}
	}
}

func PrintFiltered() {
	ipl, err := LookUpMap(myIP, myLanguage)

	if err != nil {
		fmt.Println(err)
	} else {
		var field string
		fields := strings.Split(filterFields, ",")
		for i := 0; i < len(fields); i++ {
			field = strings.TrimSpace(fields[i])
			fmt.Print(field)
			if i+1 < len(fields) {
				fmt.Print(",")
			}
		}
		fmt.Println("")
		for i := 0; i < len(fields); i++ {
			field = strings.TrimSpace(fields[i])
			subfields := strings.Split(field, ".")

			// traverse the nested map
			var subfield string
			iplsub := ipl
			for j := 0; j < len(subfields); j++ {
				subfield = subfields[j]

				if v, exists := iplsub[subfield]; exists {
					if v == nil {
						break
					}
					if j+1 == len(subfields) { // end of the traversal
						switch t := reflect.TypeOf(v).Kind(); t {
						case reflect.String:
							v2 := v.(string)
							v2 = strings.ReplaceAll(v2, `"`, `\"`)
							fmt.Printf(`"%s"`, v2)
						case reflect.Float64:
							v2 := v.(float64) // all numbers are converted to float
							if subfield == "latitude" || subfield == "longitude" {
								fmt.Print(v2) // maintain as float
							} else {
								fmt.Print(int(v2))
							}
						case reflect.Slice:
							fmt.Printf("%v", v)
						case reflect.Bool:
							v2 := v.(bool)
							fmt.Print(v2)
						default:
							fmt.Print("")
						}
					} else { // still need to drill down the map
						iplsub = iplsub[subfield].(map[string]interface{})
					}
				} else {
					break
				}
			}
			if i+1 < len(fields) {
				fmt.Print(",")
			}
		}
		fmt.Println("")
	}
}

func PrintNormal() {
	json, err := LookUpJSON(myIP, myLanguage)

	if err != nil {
		fmt.Println(err)
	} else {
		if outputFormat == "json" {
			fmt.Printf("%s\n", json)
			return
		}
		pretty, err := PrettyString(json)

		if err != nil {
			fmt.Println(err)
		} else {
			fmt.Println(pretty)
		}
	}
}

func PrintUsage() {
	fmt.Printf("%s Version %s\n", programName, version)
	var usage string = `
To query IP geolocation:

  Usage: EXE [OPTION]... <IP ADDRESS>

    -v                   Display the version and exit

    -h                   Print this help

    -k                   Specify the IP2Location.io API key
                         Get your API key from https://www.ip2location.io

    -l                   Specify the translaction language, only supported in Plus and Security plans
                         Valid values: ar | cs | da | de | en | es | et | fi | fr | ga | it | ja | ko | ms | nl | pt | ru | sv | tr | vi | zh-cn | zh-tw

    -o                   Specify the output format
                         Valid values: json (default) | pretty

    -f                   Filter the result fields
                         Field names separated by comma and using period for nested field
                         E.g. country_name,region_code,continent.name,country.translation.value

To store the API key

  Usage: EXE config <API KEY>


Other functions:

To generate random IPv4 address

  Usage: EXE randip

To convert CIDR to range

  Usage: EXE cidr2range <CIDR>

To convert range to CIDR

  Usage: EXE range2cidr <START IP> <END IP>

To list out the IPs in a CIDR

  Usage: EXE cidr2list <CIDR>

To list out the IPs in a range

  Usage: EXE range2list <START IP> <END IP>
`

	usage = strings.ReplaceAll(usage, "EXE", os.Args[0])
	fmt.Println(usage)
}
