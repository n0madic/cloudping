package main

import (
	"bufio"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net"
	"net/http"
	"os"
	"path/filepath"
	"regexp"
	"runtime"
	"sort"
	"strings"
	"sync"
	"time"
)

// CLI flags
type config struct {
	Provider   string        `json:"-"`
	Timeout    time.Duration `json:"-"`
	Concurrent int           `json:"-"`
	JSON       bool          `json:"-"`
}

// discoveredRegion represents a region found by a discovery strategy.
type discoveredRegion struct {
	Name     string `json:"name"`
	Location string `json:"location,omitempty"`
	Status   string `json:"status"` // "ok", "failed", "skipped"
	Detail   string `json:"detail,omitempty"`
}

// providerResult holds the discovery output for one provider.
type providerResult struct {
	Provider    string             `json:"provider"`
	DisplayName string            `json:"display_name"`
	Source      string             `json:"source"`
	Candidates  int               `json:"candidates"`
	Known       int               `json:"known"`
	New         []discoveredRegion `json:"new,omitempty"`
	Missing     []discoveredRegion `json:"missing,omitempty"`
	Skipped     []discoveredRegion `json:"skipped,omitempty"`
}

// knownRegionInfo holds parsed region data from providers.go.
type knownRegionInfo struct {
	Names     map[string]bool
	Locations map[string]string // name -> location
}

// discoveryFunc fetches candidate regions and tests them.
type discoveryFunc func(known knownRegionInfo, cfg config) (*providerResult, error)

var discoveryRegistry = map[string]discoveryFunc{
	// API-based
	"aws":    discoverAWS,
	"vultr":  discoverVultr,
	"linode": discoverLinode,
	"oracle": discoverOracle,
	// DNS-probe
	"azure":        discoverAzure,
	"digitalocean": discoverDigitalOcean,
	"hetzner":      discoverHetzner,
	"tencent":      discoverTencent,
	"alibaba":      discoverAlibaba,
	"huawei":       discoverHuawei,
	"servers":      discoverServers,
}

func main() {
	cfg := config{
		Timeout:    3 * time.Second,
		Concurrent: 20,
	}
	parseFlags(&cfg)

	// Find providers.go relative to this tool
	providersPath := findProvidersGo()
	if providersPath == "" {
		fmt.Fprintf(os.Stderr, "Error: cannot find providers.go\n")
		os.Exit(1)
	}

	knownRegions, err := parseProvidersGo(providersPath)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error parsing providers.go: %v\n", err)
		os.Exit(1)
	}

	// Determine which providers to scan
	providerKeys := make([]string, 0, len(discoveryRegistry))
	if cfg.Provider != "" {
		if _, ok := discoveryRegistry[cfg.Provider]; !ok {
			fmt.Fprintf(os.Stderr, "Error: no discovery strategy for provider %q\n", cfg.Provider)
			fmt.Fprintf(os.Stderr, "Available providers: %s\n", strings.Join(registryKeys(), ", "))
			os.Exit(1)
		}
		providerKeys = append(providerKeys, cfg.Provider)
	} else {
		for k := range discoveryRegistry {
			providerKeys = append(providerKeys, k)
		}
		sort.Strings(providerKeys)
	}

	var results []*providerResult
	for _, key := range providerKeys {
		known := knownRegions[key]
		if known.Names == nil {
			known.Names = make(map[string]bool)
		}
		if known.Locations == nil {
			known.Locations = make(map[string]string)
		}
		if !cfg.JSON {
			fmt.Fprintf(os.Stderr, "Discovering %s...\n", key)
		}
		result, err := discoveryRegistry[key](known, cfg)
		if err != nil {
			fmt.Fprintf(os.Stderr, "  Error: %v\n", err)
			continue
		}
		results = append(results, result)
	}

	if cfg.JSON {
		enc := json.NewEncoder(os.Stdout)
		enc.SetIndent("", "  ")
		enc.Encode(results)
		return
	}

	for _, r := range results {
		printResult(r)
	}
}

func parseFlags(cfg *config) {
	args := os.Args[1:]
	for i := 0; i < len(args); i++ {
		switch args[i] {
		case "--provider":
			i++
			if i < len(args) {
				cfg.Provider = strings.ToLower(args[i])
			}
		case "--timeout":
			i++
			if i < len(args) {
				d, err := time.ParseDuration(args[i])
				if err == nil {
					cfg.Timeout = d
				}
			}
		case "--concurrent":
			i++
			if i < len(args) {
				fmt.Sscanf(args[i], "%d", &cfg.Concurrent)
			}
		case "--json":
			cfg.JSON = true
		case "--help", "-h":
			fmt.Println("Usage: discover [flags]")
			fmt.Println()
			fmt.Println("Flags:")
			fmt.Println("  --provider <name>    Only discover for a specific provider (default: all)")
			fmt.Println("  --timeout <duration> Endpoint test timeout (default: 3s)")
			fmt.Println("  --concurrent <int>   Max concurrent probes (default: 20)")
			fmt.Println("  --json               Output as JSON")
			fmt.Println()
			fmt.Println("Available providers:", strings.Join(registryKeys(), ", "))
			os.Exit(0)
		}
	}
}

func registryKeys() []string {
	keys := make([]string, 0, len(discoveryRegistry))
	for k := range discoveryRegistry {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	return keys
}

// findProvidersGo walks up from the binary/working directory to find providers.go.
func findProvidersGo() string {
	candidates := []string{
		"providers.go",
		"../../providers.go",
	}

	if exe, err := os.Executable(); err == nil {
		candidates = append(candidates, filepath.Join(filepath.Dir(exe), "providers.go"))
		candidates = append(candidates, filepath.Join(filepath.Dir(exe), "..", "..", "providers.go"))
	}

	_, thisFile, _, ok := runtime.Caller(0)
	if ok {
		candidates = append(candidates, filepath.Join(filepath.Dir(thisFile), "..", "..", "providers.go"))
	}

	for _, c := range candidates {
		abs, err := filepath.Abs(c)
		if err != nil {
			continue
		}
		if _, err := os.Stat(abs); err == nil {
			return abs
		}
	}
	return ""
}

// parseProvidersGo extracts known region names and locations per provider from providers.go.
func parseProvidersGo(path string) (map[string]knownRegionInfo, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	result := make(map[string]knownRegionInfo)

	providerRe := regexp.MustCompile(`"(\w[\w-]*)"\s*:\s*\{`)
	regionRe := regexp.MustCompile(`\{name:\s*"([^"]+)"`)
	locationRe := regexp.MustCompile(`location:\s*"([^"]+)"`)

	lines := strings.Split(string(data), "\n")
	var currentProvider string
	braceDepth := 0

	for _, line := range lines {
		trimmed := strings.TrimSpace(line)

		// Skip commented-out regions
		if strings.HasPrefix(trimmed, "//") {
			continue
		}

		if currentProvider == "" {
			if m := providerRe.FindStringSubmatch(trimmed); m != nil {
				currentProvider = m[1]
				result[currentProvider] = knownRegionInfo{
					Names:     make(map[string]bool),
					Locations: make(map[string]string),
				}
				braceDepth = 1
			}
		} else {
			braceDepth += strings.Count(trimmed, "{") - strings.Count(trimmed, "}")
			if braceDepth <= 0 {
				currentProvider = ""
				continue
			}
			if m := regionRe.FindStringSubmatch(trimmed); m != nil {
				info := result[currentProvider]
				info.Names[m[1]] = true
				if lm := locationRe.FindStringSubmatch(trimmed); lm != nil {
					info.Locations[m[1]] = lm[1]
				}
				result[currentProvider] = info
			}
		}
	}

	return result, nil
}

// --- Concurrent endpoint testing ---

type probeResult struct {
	Name   string
	OK     bool
	Method string // "dns" or "http"
}

func testEndpointsDNS(endpoints map[string]string, cfg config) map[string]probeResult {
	return testEndpoints(endpoints, cfg, func(host string, timeout time.Duration) bool {
		ctx, cancel := context.WithTimeout(context.Background(), timeout)
		defer cancel()
		addrs, err := net.DefaultResolver.LookupHost(ctx, host)
		return err == nil && len(addrs) > 0
	}, "dns")
}

func testEndpointsHTTP(endpoints map[string]string, cfg config) map[string]probeResult {
	return testEndpoints(endpoints, cfg, func(url string, timeout time.Duration) bool {
		client := &http.Client{Timeout: timeout}
		resp, err := client.Head(url)
		if err != nil {
			return false
		}
		resp.Body.Close()
		return true
	}, "http")
}

func testEndpoints(endpoints map[string]string, cfg config, probe func(string, time.Duration) bool, method string) map[string]probeResult {
	results := make(map[string]probeResult)
	var mu sync.Mutex
	var wg sync.WaitGroup
	sem := make(chan struct{}, cfg.Concurrent)

	for name, endpoint := range endpoints {
		wg.Add(1)
		sem <- struct{}{}
		go func(n, ep string) {
			defer wg.Done()
			defer func() { <-sem }()
			ok := probe(ep, cfg.Timeout)
			mu.Lock()
			results[n] = probeResult{Name: n, OK: ok, Method: method}
			mu.Unlock()
		}(name, endpoint)
	}
	wg.Wait()
	return results
}

// --- Provider display names ---

var providerDisplayNames = map[string]string{
	"aws":          "Amazon Web Services",
	"vultr":        "Vultr",
	"linode":       "Linode",
	"oracle":       "Oracle Cloud",
	"azure":        "Microsoft Azure",
	"digitalocean": "Digital Ocean",
	"hetzner":      "Hetzner Cloud",
	"tencent":      "Tencent Cloud",
	"alibaba":      "Alibaba Cloud",
	"huawei":       "Huawei Cloud",
	"servers":      "Servers.com",
}

// --- API-based discovery ---

func discoverAWS(known knownRegionInfo, cfg config) (*providerResult, error) {
	type ipRanges struct {
		Prefixes []struct {
			Region  string `json:"region"`
			Service string `json:"service"`
		} `json:"prefixes"`
	}

	body, err := httpGet("https://ip-ranges.amazonaws.com/ip-ranges.json", cfg.Timeout)
	if err != nil {
		return nil, fmt.Errorf("fetching ip-ranges.json: %w", err)
	}

	var data ipRanges
	if err := json.Unmarshal(body, &data); err != nil {
		return nil, fmt.Errorf("parsing ip-ranges.json: %w", err)
	}

	regionSet := make(map[string]bool)
	for _, p := range data.Prefixes {
		if p.Service == "S3" && p.Region != "GLOBAL" {
			regionSet[p.Region] = true
		}
	}

	var skipped []discoveredRegion
	endpoints := make(map[string]string)
	for r := range regionSet {
		if strings.HasPrefix(r, "cn-") {
			skipped = append(skipped, discoveredRegion{Name: r, Status: "skipped", Detail: "China region"})
			continue
		}
		if strings.HasPrefix(r, "us-gov-") {
			skipped = append(skipped, discoveredRegion{Name: r, Status: "skipped", Detail: "Government region"})
			continue
		}
		if !known.Names[r] {
			endpoints[r] = fmt.Sprintf("s3.%s.amazonaws.com", r)
		}
	}

	tested := testEndpointsDNS(endpoints, cfg)

	var newRegions []discoveredRegion
	for name, result := range tested {
		status := "ok"
		if !result.OK {
			status = "failed"
		}
		newRegions = append(newRegions, discoveredRegion{
			Name:   name,
			Status: status,
			Detail: fmt.Sprintf("%s:%s", result.Method, strings.ToUpper(status)),
		})
	}

	sort.Slice(skipped, func(i, j int) bool { return skipped[i].Name < skipped[j].Name })
	sort.Slice(newRegions, func(i, j int) bool { return newRegions[i].Name < newRegions[j].Name })

	var missing []discoveredRegion
	for r := range known.Names {
		if !regionSet[r] {
			missing = append(missing, discoveredRegion{Name: r, Status: "missing", Detail: "not in ip-ranges.json"})
		}
	}
	sort.Slice(missing, func(i, j int) bool { return missing[i].Name < missing[j].Name })

	return &providerResult{
		Provider:    "aws",
		DisplayName: providerDisplayNames["aws"],
		Source:      "ip-ranges.amazonaws.com",
		Candidates:  len(regionSet),
		Known:       len(known.Names),
		New:         newRegions,
		Missing:     missing,
		Skipped:     skipped,
	}, nil
}

func discoverVultr(known knownRegionInfo, cfg config) (*providerResult, error) {
	type vultrRegion struct {
		ID      string `json:"id"`
		City    string `json:"city"`
		Country string `json:"country"`
	}
	type vultrResponse struct {
		Regions []vultrRegion `json:"regions"`
	}

	body, err := httpGet("https://api.vultr.com/v2/regions", cfg.Timeout)
	if err != nil {
		return nil, fmt.Errorf("fetching vultr regions: %w", err)
	}

	var data vultrResponse
	if err := json.Unmarshal(body, &data); err != nil {
		return nil, fmt.Errorf("parsing vultr regions: %w", err)
	}

	// Build a set of known locations for matching (normalized to lowercase)
	knownLocations := make(map[string]bool)
	for _, loc := range known.Locations {
		knownLocations[strings.ToLower(loc)] = true
	}

	// Check each API region against known locations
	var newRegions []discoveredRegion
	for _, r := range data.Regions {
		apiLocation := fmt.Sprintf("%s, %s", r.City, r.Country)
		normalizedLoc := strings.ToLower(apiLocation)

		// Check if this location is already covered by a known region
		if knownLocations[normalizedLoc] {
			continue
		}

		// Also check partial match (city name match, ASCII-normalized, bidirectional)
		cityMatch := false
		normalizedCity := normalizeASCII(strings.ToLower(r.City))
		for _, knownLoc := range known.Locations {
			normalizedKnown := normalizeASCII(strings.ToLower(knownLoc))
			if strings.Contains(normalizedKnown, normalizedCity) || strings.Contains(normalizedCity, normalizeASCII(strings.ToLower(strings.SplitN(knownLoc, ",", 2)[0]))) {
				cityMatch = true
				break
			}
		}
		if cityMatch {
			continue
		}

		// This is a genuinely new location - test the endpoint
		hostname := fmt.Sprintf("%s-ping.vultr.com", r.ID)
		ctx, cancel := context.WithTimeout(context.Background(), cfg.Timeout)
		addrs, lookupErr := net.DefaultResolver.LookupHost(ctx, hostname)
		cancel()
		dnsOK := lookupErr == nil && len(addrs) > 0

		// Also try the {city}-{country} format used in providers.go
		altID := fmt.Sprintf("%s-%s", strings.ToLower(r.City[:3]), strings.ToLower(r.Country))
		altHostname := fmt.Sprintf("%s-ping.vultr.com", altID)
		ctx2, cancel2 := context.WithTimeout(context.Background(), cfg.Timeout)
		altAddrs, altErr := net.DefaultResolver.LookupHost(ctx2, altHostname)
		cancel2()
		altOK := altErr == nil && len(altAddrs) > 0

		status := "failed"
		detail := fmt.Sprintf("dns:FAILED (tried: %s, %s)", hostname, altHostname)
		if dnsOK {
			status = "ok"
			detail = fmt.Sprintf("dns:OK (hostname: %s)", hostname)
		} else if altOK {
			status = "ok"
			detail = fmt.Sprintf("dns:OK (hostname: %s)", altHostname)
		}

		newRegions = append(newRegions, discoveredRegion{
			Name:     r.ID,
			Location: apiLocation,
			Status:   status,
			Detail:   detail,
		})
	}
	sort.Slice(newRegions, func(i, j int) bool { return newRegions[i].Name < newRegions[j].Name })

	return &providerResult{
		Provider:    "vultr",
		DisplayName: providerDisplayNames["vultr"],
		Source:      "api.vultr.com",
		Candidates:  len(data.Regions),
		Known:       len(known.Names),
		New:         newRegions,
	}, nil
}

func discoverLinode(known knownRegionInfo, cfg config) (*providerResult, error) {
	type linodeRegion struct {
		ID      string `json:"id"`
		Label   string `json:"label"`
		Country string `json:"country"`
	}
	type linodeResponse struct {
		Data []linodeRegion `json:"data"`
	}

	body, err := httpGet("https://api.linode.com/v4/regions", cfg.Timeout)
	if err != nil {
		return nil, fmt.Errorf("fetching linode regions: %w", err)
	}

	var data linodeResponse
	if err := json.Unmarshal(body, &data); err != nil {
		return nil, fmt.Errorf("parsing linode regions: %w", err)
	}

	endpoints := make(map[string]string)
	regionInfo := make(map[string]linodeRegion)
	for _, r := range data.Data {
		regionInfo[r.ID] = r
		if !known.Names[r.ID] {
			endpoints[r.ID] = fmt.Sprintf("speedtest.%s.linode.com", r.ID)
		}
	}

	tested := testEndpointsDNS(endpoints, cfg)

	var newRegions []discoveredRegion
	for id, result := range tested {
		info := regionInfo[id]
		status := "ok"
		if !result.OK {
			status = "failed"
		}
		newRegions = append(newRegions, discoveredRegion{
			Name:     id,
			Location: fmt.Sprintf("%s, %s", info.Label, strings.ToUpper(info.Country)),
			Status:   status,
			Detail:   fmt.Sprintf("speedtest:%s", strings.ToUpper(status)),
		})
	}
	sort.Slice(newRegions, func(i, j int) bool { return newRegions[i].Name < newRegions[j].Name })

	apiRegionSet := make(map[string]bool)
	for _, r := range data.Data {
		apiRegionSet[r.ID] = true
	}
	var missing []discoveredRegion
	for r := range known.Names {
		if !apiRegionSet[r] {
			missing = append(missing, discoveredRegion{Name: r, Status: "missing", Detail: "not in Linode API"})
		}
	}
	sort.Slice(missing, func(i, j int) bool { return missing[i].Name < missing[j].Name })

	return &providerResult{
		Provider:    "linode",
		DisplayName: providerDisplayNames["linode"],
		Source:      "api.linode.com",
		Candidates:  len(data.Data),
		Known:       len(known.Names),
		New:         newRegions,
		Missing:     missing,
	}, nil
}

func discoverOracle(known knownRegionInfo, cfg config) (*providerResult, error) {
	body, err := httpGet("https://raw.githubusercontent.com/oracle/oci-go-sdk/master/common/regions.go", cfg.Timeout*3)
	if err != nil {
		return nil, fmt.Errorf("fetching OCI regions.go: %w", err)
	}

	// Match patterns like: Region = "us-ashburn-1" or "us-gov-ashburn-1"
	regionRe := regexp.MustCompile(`"([a-z]{2,3}(?:-[a-z]+)+-\d+)"`)
	matches := regionRe.FindAllStringSubmatch(string(body), -1)

	regionSet := make(map[string]bool)
	for _, m := range matches {
		regionSet[m[1]] = true
	}

	var skipped []discoveredRegion
	endpoints := make(map[string]string)
	for r := range regionSet {
		if strings.Contains(r, "-gov-") || strings.Contains(r, "-dcc-") ||
			strings.Contains(r, "-langley-") || strings.Contains(r, "-luke-") {
			skipped = append(skipped, discoveredRegion{Name: r, Status: "skipped", Detail: "Government/classified region"})
			continue
		}
		if !known.Names[r] {
			endpoints[r] = fmt.Sprintf("objectstorage.%s.oraclecloud.com", r)
		}
	}

	tested := testEndpointsDNS(endpoints, cfg)

	var newRegions []discoveredRegion
	for name, result := range tested {
		status := "ok"
		if !result.OK {
			status = "failed"
		}
		newRegions = append(newRegions, discoveredRegion{
			Name:   name,
			Status: status,
			Detail: fmt.Sprintf("%s:%s", result.Method, strings.ToUpper(status)),
		})
	}

	sort.Slice(skipped, func(i, j int) bool { return skipped[i].Name < skipped[j].Name })
	sort.Slice(newRegions, func(i, j int) bool { return newRegions[i].Name < newRegions[j].Name })

	return &providerResult{
		Provider:    "oracle",
		DisplayName: providerDisplayNames["oracle"],
		Source:      "OCI Go SDK (GitHub)",
		Candidates:  len(regionSet),
		Known:       len(known.Names),
		New:         newRegions,
		Skipped:     skipped,
	}, nil
}

// --- DNS-probe based discovery ---

func discoverAzure(known knownRegionInfo, cfg config) (*providerResult, error) {
	candidates := []string{
		"australiacentral", "australiacentral2", "australiaeast", "australiasoutheast",
		"brazilsouth", "brazilsoutheast", "brazilus",
		"canadacentral", "canadaeast",
		"centralindia", "centralus", "centraluseuap",
		"chilecentral",
		"eastasia", "eastus", "eastus2", "eastus2euap", "eastusstg",
		"francecentral", "francesouth",
		"germanywestcentral", "germanynorth",
		"indonesiacentral",
		"israelcentral",
		"italynorth",
		"japaneast", "japanwest",
		"jioindiacentral", "jioindiawest",
		"koreacentral", "koreasouth",
		"malaysiasouth", "malaysiawest",
		"mexicocentral",
		"newzealandnorth",
		"northcentralus", "northeurope",
		"norwayeast", "norwaywest",
		"polandcentral",
		"qatarcentral",
		"southafricanorth", "southafricawest",
		"southcentralus", "southeastasia", "southindia",
		"spaincentral",
		"swedencentral", "swedensouth",
		"switzerlandnorth", "switzerlandwest",
		"taiwannorth", "taiwannorthwest",
		"uaecentral", "uaenorth",
		"uksouth", "ukwest",
		"westcentralus", "westeurope", "westindia",
		"westus", "westus2", "westus3",
	}

	endpoints := make(map[string]string)
	for _, c := range candidates {
		if !known.Names[c] {
			endpoints[c] = fmt.Sprintf("s3%s.blob.core.windows.net", c)
		}
	}

	tested := testEndpointsDNS(endpoints, cfg)

	var newRegions []discoveredRegion
	for name, result := range tested {
		if !result.OK {
			continue
		}
		newRegions = append(newRegions, discoveredRegion{
			Name:   name,
			Status: "ok",
			Detail: "dns:OK",
		})
	}
	sort.Slice(newRegions, func(i, j int) bool { return newRegions[i].Name < newRegions[j].Name })

	return &providerResult{
		Provider:    "azure",
		DisplayName: providerDisplayNames["azure"],
		Source:      "DNS probing (Azure blob patterns)",
		Candidates:  len(candidates),
		Known:       len(known.Names),
		New:         newRegions,
	}, nil
}

func discoverDigitalOcean(known knownRegionInfo, cfg config) (*providerResult, error) {
	baseCodes := []string{
		"ams", "blr", "fra", "lon", "nyc", "sfo", "sgp", "syd", "tor", "atl",
		"bom", "del", "dfw", "iad", "lax", "mia", "ord", "sea", "sin",
		"hkg", "icn", "jnb", "osl", "par", "waw", "mad",
	}

	var candidates []string
	for _, base := range baseCodes {
		for i := 1; i <= 5; i++ {
			candidates = append(candidates, fmt.Sprintf("%s%d", base, i))
		}
	}

	endpoints := make(map[string]string)
	for _, c := range candidates {
		if !known.Names[c] {
			endpoints[c] = fmt.Sprintf("%s.digitaloceanspaces.com", c)
		}
	}

	tested := testEndpointsDNS(endpoints, cfg)

	var newRegions []discoveredRegion
	for name, result := range tested {
		if !result.OK {
			continue
		}
		newRegions = append(newRegions, discoveredRegion{
			Name:   name,
			Status: "ok",
			Detail: "dns:OK",
		})
	}
	sort.Slice(newRegions, func(i, j int) bool { return newRegions[i].Name < newRegions[j].Name })

	return &providerResult{
		Provider:    "digitalocean",
		DisplayName: providerDisplayNames["digitalocean"],
		Source:      "DNS probing (Spaces endpoints)",
		Candidates:  len(candidates),
		Known:       len(known.Names),
		New:         newRegions,
	}, nil
}

func discoverHetzner(known knownRegionInfo, cfg config) (*providerResult, error) {
	candidates := []string{
		"ash", "fsn", "hel", "hil", "nbg", "sin",
		// Additional potential codes
		"lax", "ams", "lon", "par", "fra", "sgp", "syd",
		"tok", "osk", "sel", "bom", "del", "jnb",
		"sao", "mia", "dfw", "ord", "sea", "sjc",
		"tor", "waw", "mad", "mil", "vie",
	}

	endpoints := make(map[string]string)
	for _, c := range candidates {
		if !known.Names[c] {
			endpoints[c] = fmt.Sprintf("%s.icmp.hetzner.com", c)
		}
	}

	tested := testEndpointsDNS(endpoints, cfg)

	var newRegions []discoveredRegion
	for name, result := range tested {
		if !result.OK {
			continue
		}
		newRegions = append(newRegions, discoveredRegion{
			Name:   name,
			Status: "ok",
			Detail: "dns:OK",
		})
	}
	sort.Slice(newRegions, func(i, j int) bool { return newRegions[i].Name < newRegions[j].Name })

	return &providerResult{
		Provider:    "hetzner",
		DisplayName: providerDisplayNames["hetzner"],
		Source:      "DNS probing (ICMP endpoints)",
		Candidates:  len(candidates),
		Known:       len(known.Names),
		New:         newRegions,
	}, nil
}

func discoverTencent(known knownRegionInfo, cfg config) (*providerResult, error) {
	candidates := []string{
		"ap-bangkok", "ap-beijing", "ap-chengdu", "ap-chongqing",
		"ap-guangzhou", "ap-hongkong", "ap-jakarta", "ap-mumbai",
		"ap-nanjing", "ap-seoul", "ap-shanghai", "ap-shanghai-fsi",
		"ap-shenzhen-fsi", "ap-singapore", "ap-tokyo",
		"eu-frankfurt", "eu-moscow",
		"me-saudi-arabia",
		"na-ashburn", "na-siliconvalley", "na-toronto",
		"sa-saopaulo",
		// Additional potential regions
		"ap-taipei", "ap-manila", "ap-osaka",
		"eu-london", "eu-paris", "eu-amsterdam", "eu-stockholm",
		"me-dubai", "me-bahrain",
		"sa-santiago", "sa-buenosaires",
		"af-johannesburg",
	}

	endpoints := make(map[string]string)
	for _, c := range candidates {
		if !known.Names[c] {
			endpoints[c] = fmt.Sprintf("cos.%s.myqcloud.com", c)
		}
	}

	tested := testEndpointsDNS(endpoints, cfg)

	var newRegions []discoveredRegion
	for name, result := range tested {
		if !result.OK {
			continue
		}
		newRegions = append(newRegions, discoveredRegion{
			Name:   name,
			Status: "ok",
			Detail: "dns:OK",
		})
	}
	sort.Slice(newRegions, func(i, j int) bool { return newRegions[i].Name < newRegions[j].Name })

	return &providerResult{
		Provider:    "tencent",
		DisplayName: providerDisplayNames["tencent"],
		Source:      "DNS probing (COS endpoints)",
		Candidates:  len(candidates),
		Known:       len(known.Names),
		New:         newRegions,
	}, nil
}

func discoverAlibaba(known knownRegionInfo, cfg config) (*providerResult, error) {
	candidates := []string{
		"ap-southeast-1", "ap-southeast-2", "ap-southeast-3",
		"ap-southeast-5", "ap-southeast-6", "ap-southeast-7",
		"ap-northeast-1", "ap-northeast-2",
		"ap-south-1",
		"cn-beijing", "cn-chengdu", "cn-guangzhou", "cn-fuzhou",
		"cn-hangzhou", "cn-heyuan", "cn-hongkong", "cn-huhehaote",
		"cn-nanjing", "cn-qingdao", "cn-shanghai", "cn-shenzhen",
		"cn-wuhan-lr", "cn-wulanchabu", "cn-zhangjiakou",
		"eu-central-1", "eu-west-1",
		"me-central-1", "me-east-1",
		"na-south-1",
		"us-east-1", "us-west-1",
		// Additional potential regions
		"ap-southeast-4", "ap-south-2",
		"cn-zhengzhou", "cn-changsha", "cn-kunming",
		"eu-north-1", "eu-south-1",
		"me-west-1",
		"sa-east-1",
	}

	endpoints := make(map[string]string)
	for _, c := range candidates {
		if !known.Names[c] {
			endpoints[c] = fmt.Sprintf("oss-%s.aliyuncs.com", c)
		}
	}

	tested := testEndpointsDNS(endpoints, cfg)

	var newRegions []discoveredRegion
	for name, result := range tested {
		if !result.OK {
			continue
		}
		newRegions = append(newRegions, discoveredRegion{
			Name:   name,
			Status: "ok",
			Detail: "dns:OK",
		})
	}
	sort.Slice(newRegions, func(i, j int) bool { return newRegions[i].Name < newRegions[j].Name })

	return &providerResult{
		Provider:    "alibaba",
		DisplayName: providerDisplayNames["alibaba"],
		Source:      "DNS probing (OSS endpoints)",
		Candidates:  len(candidates),
		Known:       len(known.Names),
		New:         newRegions,
	}, nil
}

func discoverHuawei(known knownRegionInfo, cfg config) (*providerResult, error) {
	candidates := []string{
		"ae-ad-1", "af-south-1",
		"ap-southeast-1", "ap-southeast-2", "ap-southeast-3",
		"ap-southeast-4", "ap-southeast-5",
		"cn-east-2", "cn-east-3", "cn-east-4", "cn-east-5",
		"cn-north-1", "cn-north-4", "cn-north-9",
		"cn-south-1", "cn-south-2", "cn-south-4",
		"cn-southwest-2",
		"eu-west-0", "eu-west-101",
		"la-north-2", "la-south-2",
		"me-east-1",
		"my-kualalumpur-1",
		"na-mexico-1",
		"ru-northwest-2",
		"sa-argentina-1", "sa-brazil-1", "sa-peru-1",
		"tr-west-1",
		// Additional potential regions
		"ap-southeast-6", "cn-north-2", "cn-south-3",
		"eu-east-1", "eu-north-1",
	}

	endpoints := make(map[string]string)
	for _, c := range candidates {
		if !known.Names[c] {
			endpoints[c] = fmt.Sprintf("dns.%s.myhuaweicloud.com", c)
		}
	}

	// Huawei uses HTTPS endpoints, try DNS first then HTTP fallback
	tested := testEndpointsDNS(endpoints, cfg)

	// HTTP fallback for DNS failures
	httpEndpoints := make(map[string]string)
	for name, result := range tested {
		if !result.OK {
			httpEndpoints[name] = fmt.Sprintf("https://dns.%s.myhuaweicloud.com", name)
		}
	}
	if len(httpEndpoints) > 0 {
		httpTested := testEndpointsHTTP(httpEndpoints, cfg)
		for name, result := range httpTested {
			if result.OK {
				tested[name] = probeResult{Name: name, OK: true, Method: "http"}
			}
		}
	}

	var newRegions []discoveredRegion
	for name, result := range tested {
		if !result.OK {
			continue
		}
		newRegions = append(newRegions, discoveredRegion{
			Name:   name,
			Status: "ok",
			Detail: fmt.Sprintf("%s:OK", result.Method),
		})
	}
	sort.Slice(newRegions, func(i, j int) bool { return newRegions[i].Name < newRegions[j].Name })

	return &providerResult{
		Provider:    "huawei",
		DisplayName: providerDisplayNames["huawei"],
		Source:      "DNS probing (Huawei Cloud endpoints)",
		Candidates:  len(candidates),
		Known:       len(known.Names),
		New:         newRegions,
	}, nil
}

func discoverServers(known knownRegionInfo, cfg config) (*providerResult, error) {
	bases := []string{
		"ams", "dfw", "hkg", "lon", "lux", "mia", "sao", "sin", "sjc", "was",
		// Additional potential bases
		"fra", "par", "tok", "syd", "sel", "bom", "jnb",
		"ord", "lax", "sea", "tor", "waw", "mad",
	}

	var candidates []string
	for _, base := range bases {
		for i := 1; i <= 10; i++ {
			candidates = append(candidates, fmt.Sprintf("%s%d", base, i))
		}
	}

	endpoints := make(map[string]string)
	for _, c := range candidates {
		if !known.Names[c] {
			endpoints[c] = fmt.Sprintf("test.%s.servers.com", c)
		}
	}

	tested := testEndpointsDNS(endpoints, cfg)

	var newRegions []discoveredRegion
	for name, result := range tested {
		if !result.OK {
			continue
		}
		newRegions = append(newRegions, discoveredRegion{
			Name:   name,
			Status: "ok",
			Detail: "dns:OK",
		})
	}
	sort.Slice(newRegions, func(i, j int) bool { return newRegions[i].Name < newRegions[j].Name })

	return &providerResult{
		Provider:    "servers",
		DisplayName: providerDisplayNames["servers"],
		Source:      "DNS probing (test.*.servers.com)",
		Candidates:  len(candidates),
		Known:       len(known.Names),
		New:         newRegions,
	}, nil
}

// --- Utilities ---

// normalizeASCII strips common diacritics for lenient matching.
func normalizeASCII(s string) string {
	replacements := map[rune]rune{
		'á': 'a', 'à': 'a', 'ã': 'a', 'â': 'a', 'ä': 'a',
		'é': 'e', 'è': 'e', 'ê': 'e', 'ë': 'e',
		'í': 'i', 'ì': 'i', 'î': 'i', 'ï': 'i',
		'ó': 'o', 'ò': 'o', 'õ': 'o', 'ô': 'o', 'ö': 'o',
		'ú': 'u', 'ù': 'u', 'û': 'u', 'ü': 'u',
		'ñ': 'n', 'ç': 'c',
	}
	var b strings.Builder
	for _, r := range s {
		if rep, ok := replacements[r]; ok {
			b.WriteRune(rep)
		} else {
			b.WriteRune(r)
		}
	}
	return b.String()
}

func httpGet(url string, timeout time.Duration) ([]byte, error) {
	client := &http.Client{Timeout: timeout * 3}
	resp, err := client.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("HTTP %d", resp.StatusCode)
	}

	return io.ReadAll(bufio.NewReader(io.LimitReader(resp.Body, 50*1024*1024)))
}

// printResult formats and prints a provider discovery result.
func printResult(r *providerResult) {
	fmt.Printf("\n=== %s (%s) ===\n", r.Provider, r.DisplayName)
	fmt.Printf("  Source: %s (%d candidates)\n", r.Source, r.Candidates)
	fmt.Printf("  Known: %d regions in providers.go\n", r.Known)

	if len(r.New) > 0 {
		hasOK := false
		for _, n := range r.New {
			if n.Status == "ok" {
				hasOK = true
				break
			}
		}
		if hasOK {
			fmt.Printf("\n  NEW regions found:\n")
			for _, n := range r.New {
				if n.Status == "ok" {
					loc := ""
					if n.Location != "" {
						loc = "  " + n.Location
					}
					fmt.Printf("    + %-25s %s%s\n", n.Name, n.Detail, loc)
				}
			}
		}

		hasFailed := false
		for _, n := range r.New {
			if n.Status == "failed" {
				hasFailed = true
				break
			}
		}
		if hasFailed {
			fmt.Printf("\n  NEW but endpoint unreachable:\n")
			for _, n := range r.New {
				if n.Status == "failed" {
					loc := ""
					if n.Location != "" {
						loc = "  " + n.Location
					}
					fmt.Printf("    ? %-25s %s%s\n", n.Name, n.Detail, loc)
				}
			}
		}
	} else {
		fmt.Printf("\n  No new regions found.\n")
	}

	if len(r.Missing) > 0 {
		fmt.Printf("\n  MISSING from source (in providers.go but not in API):\n")
		for _, m := range r.Missing {
			fmt.Printf("    - %-25s %s\n", m.Name, m.Detail)
		}
	}

	if len(r.Skipped) > 0 {
		fmt.Printf("\n  Skipped (filtered):\n")
		for _, s := range r.Skipped {
			fmt.Printf("    ~ %-25s (%s)\n", s.Name, s.Detail)
		}
	}

	fmt.Println()
}
