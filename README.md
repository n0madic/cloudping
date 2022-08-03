# cloudping

This util performs a ping to measure the network latency from your location to the various cloud data centers around the world.

### Supported providers

- Alibaba Cloud
- Amazon Web Services
- Digital Ocean
- Google Cloud Platform
- Hetzner Cloud
- IBM Cloud
- Linode
- Microsoft Azure
- Oracle Cloud
- OVH Cloud
- Scaleway
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
+----+-------+-----------+--------------------+-----------+--------+
|    | CLOUD | REGION    | LOCATION           |       RTT | STATUS |
+----+-------+-----------+--------------------+-----------+--------+
|  1 | Vultr | nj-us     | New Jersey, US     | 123.422ms | OK     |
|  2 | Vultr | tor-ca    | Toronto, CA        | 134.604ms | OK     |
|  3 | Vultr | ga-us     | Atlanta, US        | 137.333ms | OK     |
|  4 | Vultr | fl-us     | Miami, US          | 145.317ms | OK     |
|  5 | Vultr | il-us     | Chicago, US        | 149.162ms | OK     |
|  6 | Vultr | tx-us     | Dallas, US         | 161.395ms | OK     |
|  7 | Vultr | lax-ca-us | Los Angeles, US    | 185.347ms | OK     |
|  8 | Vultr | wa-us     | Seattle, US        | 185.407ms | OK     |
|  9 | Vultr | sjo-ca-us | Silicon Valley, US | 185.909ms | OK     |
| 10 | Vultr | hon-hi-us | Honolulu, US       | 238.706ms | OK     |
+----+-------+-----------+--------------------+-----------+--------+
```
