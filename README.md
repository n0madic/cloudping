# cloudping

This util performs a ping to measure the network latency from your location to the various cloud data centers around the world.

### Supported providers

- Alibaba Cloud
- Amazon Web Services
- Baidu Cloud
- Digital Ocean
- Google Cloud Platform
- Hetzner Cloud
- Huawei Cloud
- IBM Cloud
- Linode
- Microsoft Azure
- Oracle Cloud
- OVH Cloud
- Scaleway
- Servers.com
- Tencent Cloud
- UCloud
- Vultr

## Install

```
go install github.com/n0madic/cloudping@latest
```

## Usage

```
Usage: cloudping [--all] [--count COUNT] [--region REGION] [--location LOCATION] [--provider PROVIDER] [--timeout TIMEOUT]

Options:
  --all, -a              Scan all providers
  --count COUNT, -c COUNT
                         Number of pings to send [default: 4]
  --hide-errors, -e      Hide errors from results
  --region REGION, -r REGION
                         Filter by regions, can be specified multiple times
  --location LOCATION, -l LOCATION
                         Filter by location, can be specified multiple times
  --provider PROVIDER, -p PROVIDER
                         Choose provider for ping [default: aws]
  --timeout TIMEOUT, -t TIMEOUT
                         Timeout before ping ends [default: 5s]
  --help, -h             display this help and exit
  ```

  ```
  $ cloudping --provider vultr --region us --region ca
+----+----------+-----------+--------------------+-----------+--------+
|    | PROVIDER | REGION    | LOCATION           |       RTT | STATUS |
+----+----------+-----------+--------------------+-----------+--------+
|  1 | Vultr    | nj-us     | New Jersey, US     | 120.396ms | OK     |
|  2 | Vultr    | tor-ca    | Toronto, CA        | 134.026ms | OK     |
|  3 | Vultr    | ga-us     | Atlanta, US        | 136.809ms | OK     |
|  4 | Vultr    | fl-us     | Miami, US          | 144.553ms | OK     |
|  5 | Vultr    | il-us     | Chicago, US        | 147.464ms | OK     |
|  6 | Vultr    | tx-us     | Dallas, US         | 156.319ms | OK     |
|  7 | Vultr    | lax-ca-us | Los Angeles, US    | 181.803ms | OK     |
|  8 | Vultr    | wa-us     | Seattle, US        | 182.867ms | OK     |
|  9 | Vultr    | sjo-ca-us | Silicon Valley, US | 188.371ms | OK     |
| 10 | Vultr    | hon-hi-us | Honolulu, US       | 238.351ms | OK     |
+----+----------+-----------+--------------------+-----------+--------+
```
