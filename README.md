[![Go Report Card](https://goreportcard.com/badge/github.com/ip2location/ip2location-io-cli)](https://goreportcard.com/report/github.com/ip2location/ip2location-io-cli)

IP2Location.io Go CLI
=====================
This Go command line tool enables user to query for an enriched data set, such as country, region, district, city, latitude & longitude, ZIP code, time zone, ASN, ISP, domain, net speed, IDD code, area code, weather station data, MNC, MCC, mobile brand, elevation, usage type, address type, advertisement category and proxy data with an IP address. It supports both IPv4 and IPv6 address lookup.

This program requires an API key to function. You may sign up for a free API key at https://www.ip2location.io/pricing.


Installation
============
```bash
go install github.com/ip2location/ip2location-io-cli/ip2locationio@latest
```



#### Debian/Ubuntu (amd64)

```bash
curl -LO https://github.com/ip2location/ip2location-io-cli/releases/download/v1.0.1/ip2location-io-1.0.1.deb
sudo dpkg -i ip2location-io-1.0.1.deb
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


Example API Response
====================
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
    "is_residential_proxy": false,
    "is_spammer": false,
    "is_scanner": false,
    "is_botnet": false
  }
}
```


LICENCE
=====================
See the LICENSE file.
