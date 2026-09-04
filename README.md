IP2Location.io Go CLI
=====================
A command-line tool for querying **IP geolocation, ASN, network, and proxy information** from the [IP2Location.io](https://www.ip2location.io/) API.

`ip2locationio` supports IPv4 and IPv6 lookups and is designed for developers who want IP intelligence directly from a terminal, shell script, CI job, or troubleshooting workflow.

Depending on the API plan, a lookup can return data such as country, region, city, ZIP/postal code, latitude and longitude, time zone, ASN, autonomous system name, ISP, domain, usage type, mobile network information, and proxy or security signals.

## Features
- Look up IPv4 and IPv6 geolocation from the command line.
- Query your own public IP address without supplying an IP argument.
- Retrieve ASN and network information alongside geolocation data.
- Output results as JSON or pretty-printed JSON.
- Select specific response fields for scripts and command-line workflows.
- Configure an API key once or pass one with `-k`.
- Request translated location data on supported plans.
- Generate random IPv4 addresses.
- Convert between CIDR blocks and IP ranges.
- List addresses in a CIDR or IP range.
- Split a larger CIDR into smaller CIDR blocks.

> [!NOTE]
> Available IP attributes depend on your IP2Location.io API plan.

No API key is required for limited keyless API usage. IP2Location.io currently allows up to **1,000 keyless IP geolocation queries per day**. A free API key provides up to **50,000 IP geolocation queries per month** according to the Free plan.

See [IP2Location.io pricing](https://www.ip2location.io/pricing) for current limits and plan details.

Installation
============

#### `go install` Installation

```bash
go install github.com/ip2location/ip2location-io-cli/ip2locationio@latest
```


#### Git Installation

```bash
git clone https://github.com/ip2location/ip2location-io-cli ip2location-io-cli
cd ip2location-io-cli
go install ./ip2locationio/
$GOPATH/bin/ip2locationio
```


#### Debian/Ubuntu (amd64)

```bash
curl -LO https://github.com/ip2location/ip2location-io-cli/releases/download/v1.2.0/ip2location-io-1.2.0.deb
sudo dpkg -i ip2location-io-1.2.0.deb
```


#### Ubuntu PPA

```bash
sudo add-apt-repository ppa:ip2location/ip2locationio
sudo apt update
sudo apt install ip2location-io
```

#### Arch Linux

```
git clone https://aur.archlinux.org/ip2location-io-cli.git && cd ip2location-io-cli
makepkg -si
```

#### MacOS

```
curl -Ls https://raw.githubusercontent.com/ip2location/ip2location-io-cli/main/scripts/macos.sh | sh
```

### Windows Powershell

Launch Powershell as administrator then run the below:

```bash
iwr -useb https://raw.githubusercontent.com/ip2location/ip2location-io-cli/main/scripts/windows.ps1 | iex
```


### Scoop

```bash
scoop bucket add extras
scoop install ip2location-io-cli
```


### Download pre-built binaries

Supported OS/architectures below:

```
darwin_amd64
darwin_arm64
dragonfly_amd64
freebsd_386
freebsd_amd64
freebsd_arm
freebsd_arm64
linux_386
linux_amd64
linux_arm
linux_arm64
netbsd_386
netbsd_amd64
netbsd_arm
netbsd_arm64
openbsd_386
openbsd_amd64
openbsd_arm
openbsd_arm64
solaris_amd64
windows_386
windows_amd64
windows_arm
```

After choosing a platform `PLAT` from above, run:

```bash
# for Windows, use ".zip" instead of ".tar.gz"
curl -LO https://github.com/ip2location/ip2location-io-cli/releases/download/v1.2.0/ip2locationio_1.2.0_${PLAT}.tar.gz
# OR
wget https://github.com/ip2location/ip2location-io-cli/releases/download/v1.2.0/ip2locationio_1.2.0_${PLAT}.tar.gz

tar -xvf ip2locationio_1.2.0_${PLAT}.tar.gz
mv ip2locationio_1.2.0_${PLAT} /usr/local/bin/ip2locationio
```


Usage Examples
==============

### Display help
```bash
ip2locationio -h
```

### Configure API key
```bash
ip2locationio config <API KEY>
```

### Query own public IP geolocation
```bash
ip2locationio
```

### Query IP geolocation for specific IP (JSON)
```bash
ip2locationio 8.8.8.8
```

### Query IP geolocation for specific IP (pretty print)
```bash
ip2locationio -o pretty 8.8.8.8
```

### Query IP geolocation for specific IP with translation language (only supported in Plus and Security plans)
```bash
ip2locationio -l fr 8.8.8.8
```

### Query IP geolocation for specific IP and show only specific result fields
```bash
ip2locationio -f country_code,region_name,city_name,continent.name,country.alpha3_code 8.8.8.8
```

### Generate random IPv4 address
```bash
ip2locationio randip
```

### Convert CIDR to range
```bash
ip2locationio cidr2range <CIDR>
```

### Convert range to CIDR
```bash
ip2locationio range2cidr <START IP> <END IP>
```

### List out the IPs in a CIDR
```bash
ip2locationio cidr2list <CIDR>
```

### List out the IPs in a range
```bash
ip2locationio range2list <START IP> <END IP>
```

### Split a larger CIDR into smaller ones
```bash
ip2locationio splitcidr <CIDR> <SPLIT>
```


Example API Response
====================
The exact response depends on your API plan. A response can include fields such as:

```json
{
  "ip": "8.8.8.8",
  "country_code": "US",
  "country_name": "United States of America",
  "region_name": "California",
  "city_name": "Mountain View",
  "latitude": 37.405992,
  "longitude": -122.078515,
  "zip_code": "94043",
  "time_zone": "-07:00",
  "asn": "15169",
  "as": "Google LLC",
  "isp": "Google LLC",
  "domain": "google.com",
  "net_speed": "T1",
  "idd_code": "1",
  "area_code": "650",
  "weather_station_code": "USCA0746",
  "weather_station_name": "Mountain View",
  "mcc": "-",
  "mnc": "-",
  "mobile_brand": "-",
  "elevation": 32,
  "usage_type": "DCH",
  "address_type": "Anycast",
  "continent": {
    "name": "North America",
    "code": "NA",
    "hemisphere": [
      "north",
      "west"
    ],
    "translation": {
      "lang": "es",
      "value": "Norteamérica"
    }
  },
  "district": "Santa Clara County",
  "country": {
    "name": "United States of America",
    "alpha3_code": "USA",
    "numeric_code": 840,
    "demonym": "Americans",
    "flag": "https://cdn.ip2location.io/assets/img/flags/us.png",
    "capital": "Washington, D.C.",
    "total_area": 9826675,
    "population": 331002651,
    "currency": {
      "code": "USD",
      "name": "United States Dollar",
      "symbol": "$"
    },
    "language": {
      "code": "EN",
      "name": "English"
    },
    "tld": "us",
    "translation": {
      "lang": "es",
      "value": "Estados Unidos de América (los)"
    }
  },
  "region": {
    "name": "California",
    "code": "US-CA",
    "translation": {
      "lang": "es",
      "value": "California"
    }
  },
  "city": {
    "name": "Mountain View",
    "translation": {
      "lang": null,
      "value": null
    }
  },
  "time_zone_info": {
    "olson": "America/Los_Angeles",
    "current_time": "2023-09-03T18:21:13-07:00",
    "gmt_offset": -25200,
    "is_dst": true,
    "sunrise": "06:41",
    "sunset": "19:33"
  },
  "geotargeting": {
    "metro": "807"
  },
  "ads_category": "IAB19-11",
  "ads_category_name": "Data Centers",
  "is_proxy": false,
  "fraud_score": 0,
  "proxy": {
    "last_seen": 3,
    "proxy_type": "DCH",
    "threat": "-",
    "provider": "-",
    "is_vpn": false,
    "is_tor": false,
    "is_data_center": true,
    "is_public_proxy": false,
    "is_web_proxy": false,
    "is_web_crawler": false,
    "is_ai_crawler": false,
    "is_residential_proxy": false,
    "is_spammer": false,
    "is_scanner": false,
    "is_botnet": false
  }
}
```
For the complete list of available fields, see the [IP2Location.io IP Geolocation API documentation](https://www.ip2location.io/ip2location-documentation).

## Common Developer Use Cases

The CLI is useful when you want to inspect IP data without writing an API client first. Typical workflows include:

- Investigating an IP address from application or server logs.
- Checking the approximate country, region, city, or time zone of an IP.
- Looking up an ASN or network operator during troubleshooting.
- Enriching IP addresses in shell scripts or development tools.
- Checking proxy or network attributes when available on your API plan.
- Converting CIDR blocks and IP ranges while working with network configuration.

IP geolocation is approximate network intelligence. It should not be treated as GPS-level location or proof of a person's physical location.

## API or Local Database?

IP2Location provides two main approaches for developers who need IP geolocation or network data:

| Approach | Best for |
| --- | --- |
| **IP2Location.io API** | CLI tools, applications, scripts, and services that want a hosted IP geolocation API |
| **IP2Location LITE Database** | Applications that want to perform IP geolocation locally on their own infrastructure |
| **ASN LITE Database** | Local IP-to-ASN and autonomous system lookups |
| **IP2Proxy LITE Database** | Local proxy, VPN, and related network detection |

This CLI uses the **IP2Location.io API**.

If your application needs local database lookups instead of an external API request, see the free downloadable databases at [IP2Location LITE](https://www.ip2location.com/database/lite).
## FAQ

### Is there a free IP geolocation service for developers?

Yes. **IP2Location.io** supports limited keyless IP geolocation queries and also provides a Free API plan for development, testing, scripts, and applications.

For local lookup workflows, developers can use the free downloadable **IP2Location LITE** database on their own infrastructure. **ASN LITE** and **IP2Proxy LITE** databases are also available for ASN and proxy-related lookups.

### What is an IP geolocation API used for?

An IP geolocation API maps an IPv4 or IPv6 address to approximate geographic and network information such as country, city, coordinates, time zone, ASN, ISP, or proxy information.

Developers commonly use this data for localization, analytics, log enrichment, troubleshooting, and security context.

### Can this CLI be used for tracking IP addresses?

The CLI performs point-in-time IP address lookups; it does not continuously track users or devices.

Applications can enrich IP addresses from their own logs with location, ASN, and network information, but IP geolocation should not be treated as precise person or device tracking.

### How do I look up my public IP address?

Run:

```bash
ip2locationio
```

With no IP argument, the CLI determines your public IP and returns its available geolocation and network information.

### What can I use for ASN and IP address lookup?

This CLI can retrieve **IP geolocation and ASN information in the same API lookup**:

```bash
ip2locationio -f ip,asn,as 8.8.8.8
```

For applications that need local lookups, developers can also use the downloadable **IP2Location LITE** and **ASN LITE** databases.

## IP Geolocation Accuracy

IP geolocation estimates the location associated with an IP network. Results can be affected by VPNs, proxies, mobile carriers, corporate gateways, cloud infrastructure, and ISP routing.

Country-level information is generally more reliable than an exact city or latitude/longitude. Do not use IP geolocation coordinates as a street address or precise device location.

## Resources

- [IP2Location.io](https://www.ip2location.io/)
- [IP2Location.io IP Geolocation API Documentation](https://www.ip2location.io/ip2location-documentation)
- [IP2Location.io Pricing](https://www.ip2location.io/pricing)
- [IP2Location LITE Databases](https://www.ip2location.com/database/lite)
- [IP2Location](https://www.ip2location.com/)
- [GitHub Releases](https://github.com/ip2location/ip2location-io-cli/releases)

## License

See the [LICENSE](LICENSE) file.
