package main

var providers = provider{
	"aws": {
		name:         "Amazon Web Services",
		hostTemplate: "s3.%s.amazonaws.com",
		regions: []*region{
			{name: "af-south-1", location: "Cape Town, ZA"},
			{name: "ap-east-1", location: "Hong Kong, CN"},
			{name: "ap-northeast-1", location: "Tokyo, JP"},
			{name: "ap-northeast-2", location: "Seoul, KR"},
			{name: "ap-northeast-3", location: "Osaka, JP"},
			{name: "ap-south-1", location: "Mumbai, IN"},
			{name: "ap-south-2", location: "Hyderabad, IN"},
			{name: "ap-southeast-1", location: "Singapore, SG"},
			{name: "ap-southeast-2", location: "Sydney, AU"},
			{name: "ap-southeast-3", location: "Jakarta, ID"},
			{name: "ap-southeast-4", location: "Melbourne, AU"},
			{name: "ap-southeast-5", location: "Cyberjaya, MY"},
			{name: "ca-central-1", location: "Montreal, CA"},
			{name: "ca-west-1", location: "Calgary, CA"},
			{name: "eu-central-1", location: "Frankfurt, DE"},
			{name: "eu-central-2", location: "Zurich, CH"},
			{name: "eu-north-1", location: "Stockholm, SE"},
			{name: "eu-south-1", location: "Milan, IT"},
			{name: "eu-south-2", location: "Zaragoza, ES"},
			{name: "eu-west-1", location: "Dublin, IE"},
			{name: "eu-west-2", location: "London, UK"},
			{name: "eu-west-3", location: "Paris, FR"},
			{name: "il-central-1", location: "Tel Aviv, IL"},
			{name: "me-central-1", location: "Middle East, UAE"},
			{name: "me-south-1", location: "Bahrain, BH"},
			{name: "sa-east-1", location: "Sao Paulo, BR"},
			{name: "us-east-1", location: "North Virginia, US"},
			{name: "us-east-2", location: "Ohio, US"},
			{name: "us-west-1", location: "North California, US"},
			{name: "us-west-2", location: "Oregon, US"},
		},
	},
	"alibaba": {
		name:         "Alibaba Cloud",
		hostTemplate: "oss-%s.aliyuncs.com",
		regions: []*region{
			{name: "ap-southeast-1", location: "Singapore, SG"},
			{name: "ap-southeast-2", location: "Sydney, AU"},
			{name: "ap-southeast-3", location: "Kuala Lumpur, MY"},
			{name: "ap-southeast-5", location: "Jakarta, ID"},
			{name: "ap-southeast-6", location: "Manila, PH"},
			{name: "ap-southeast-7", location: "Bangkok, TH"},
			{name: "ap-northeast-1", location: "Tokyo, JP"},
			{name: "ap-northeast-2", location: "Seoul, KR"},
			{name: "ap-south-1", location: "Mumbai, IN"},
			{name: "cn-beijing", location: "Beijing, CN"},
			{name: "cn-chengdu", location: "Chengdu, CN"},
			{name: "cn-guangzhou", location: "Guangzhou, CN"},
			{name: "cn-fuzhou", location: "Fuzhou, CN"},
			{name: "cn-hangzhou", location: "Hangzhou, CN"},
			{name: "cn-heyuan", location: "Heyuan, CN"},
			{name: "cn-hongkong", location: "Hong Kong, CN"},
			{name: "cn-huhehaote", location: "Huhehaote, CN"},
			{name: "cn-nanjing", location: "Nanjing, CN"},
			{name: "cn-qingdao", location: "Qingdao, CN"},
			{name: "cn-shanghai", location: "Shanghai, CN"},
			{name: "cn-shenzhen", location: "Shenzhen, CN"},
			{name: "cn-wuhan-lr", location: "Wuhan, CN"},
			{name: "cn-wulanchabu", location: "Ulanqab, CN"},
			{name: "cn-zhangjiakou", location: "Zhangjiakou, CN"},
			{name: "eu-central-1", location: "Frankfurt, DE"},
			{name: "eu-west-1", location: "London, UK"},
			{name: "me-central-1", location: "Riyadh, SA"},
			{name: "me-east-1", location: "Dubai, AE"},
			{name: "us-east-1", location: "Virginia, US"},
			{name: "us-west-1", location: "California, US"},
		},
	},
	"azure": {
		name:         "Microsoft Azure",
		hostTemplate: "https://s3%s.blob.core.windows.net/public/latency-test.json",
		regions: []*region{
			{name: "australiacentral", location: "Canberra, AU"},
			{name: "australiaeast", location: "New South Wales, AU"},
			{name: "australiasoutheast", location: "Victoria, AU"},
			{name: "brazilsouth", location: "Sao Paulo, BR"},
			{name: "canadacentral", location: "Toronto, CA", endpoint: "https://speedtestcac.blob.core.windows.net/cb.js"},
			{name: "canadaeast", location: "Quebec City, CA"},
			{name: "centralindia", location: "Pune, IN"},
			{name: "centralus", location: "Iowa, US"},
			{name: "centraluseuap", location: "Des Moines, US", endpoint: "https://centraluseuap.blob.core.windows.net/cb.js"},
			{name: "chinaeast", location: "Shanghai, CN", endpoint: "https://astchinaeast.blob.core.chinacloudapi.cn/public/latency-test.json"},
			{name: "chinanorth", location: "Beijing, CN", endpoint: "https://astchinanorth.blob.core.chinacloudapi.cn/public/latency-test.json"},
			{name: "eastasia", location: "Hong Kong, CN"},
			{name: "eastus", location: "Virginia, US"},
			{name: "eastus2", location: "Virginia, US"},
			{name: "francecentral", location: "Paris, FR"},
			{name: "germanywestcentral", location: "Frankfurt, DE"},
			{name: "japaneast", location: "Tokyo, JP"},
			{name: "japanwest", location: "Osaka, JP"},
			{name: "koreacentral", location: "Seoul, KR"},
			{name: "koreasouth", location: "Busan, KR"},
			{name: "northcentralus", location: "Illinois, US"},
			{name: "northeurope", location: "Dublin, IE"},
			{name: "norwayeast", location: "Oslo, NO"},
			{name: "polandcentral", location: "Warsaw, PL", endpoint: "https://speedtestplc.blob.core.windows.net/cb.js"},
			{name: "qatarcentral", location: "Qatar, QA"},
			{name: "southafricanorth", location: "Johannesburg, ZA"},
			{name: "southcentralus", location: "Texas, US"},
			{name: "southeastasia", location: "Singapore, SG"},
			{name: "southindia", location: "Chennai, IN"},
			{name: "swedencentral", location: "Gaevle, SE"},
			{name: "swedensouth", location: "Malmoe, SE", endpoint: "https://swedensouth.blob.core.windows.net/cb.js"},
			{name: "switzerlandnorth", location: "Zurich, CH"},
			{name: "uaenorth", location: "Dubai, AE"},
			{name: "uksouth", location: "London, UK"},
			{name: "ukwest", location: "Cardiff, UK"},
			{name: "westcentralus", location: "Wyoming, US"},
			{name: "westeurope", location: "Amsterdam, NL"},
			{name: "westindia", location: "Mumbai, IN"},
			{name: "westus", location: "California, US"},
			{name: "westus2", location: "Washington, US"},
		},
	},
	"baidu": {
		name:         "Baidu Cloud",
		hostTemplate: "s3.%s.bcebos.com",
		regions: []*region{
			{name: "bd", location: "Baoding, CN"},
			{name: "bj", location: "Beijing, CN"},
			{name: "fsh", location: "Shanghai, CN"},
			{name: "fwh", location: "Wuhan, CN"},
			{name: "gz", location: "Guangzhou, CN"},
			{name: "hkg", location: "Hong Kong, CN"},
			{name: "sin", location: "Singapore, SG"},
			{name: "su", location: "Suzhou, CN"},
		},
	},
	"coreweave": {
		name:         "CoreWeave",
		hostTemplate: "https://http.speedtest.%s.coreweave.com/ping",
		regions: []*region{
			{name: "las1", location: "Las Vegas, US"},
			{name: "lga1", location: "New York, US"},
			{name: "ord1", location: "Chicago, US"},
		},
	},
	"digitalocean": {
		name:         "Digital Ocean",
		hostTemplate: "%s.digitaloceanspaces.com",
		regions: []*region{
			// {name: "ams2", location: "Amsterdam, NL"},
			{name: "ams3", location: "Amsterdam, NL"},
			{name: "blr1", location: "Bangalore, IN"},
			{name: "fra1", location: "Frankfurt, DE"},
			// {name: "lon1", location: "London, UK"},
			// {name: "nyc1", location: "New York, US"},
			// {name: "nyc2", location: "New York, US"},
			{name: "nyc3", location: "New York, US"},
			// {name: "sfo1", location: "San Francisco, US"},
			{name: "sfo2", location: "San Francisco, US"},
			{name: "sfo3", location: "San Francisco, US"},
			{name: "sgp1", location: "Singapore, SG"},
			// {name: "tor1", location: "Toronto, CA"},
		},
	},
	"dns": {
		name: "Public DNS servers",
		regions: []*region{
			{name: "AdGuard DNS", endpoint: "94.140.14.14", location: "primary"},
			{name: "AdGuard DNS", endpoint: "94.140.15.15", location: "secondary"},
			{name: "Alternate DNS", endpoint: "76.76.19.19", location: "primary"},
			{name: "Alternate DNS", endpoint: "76.223.122.150", location: "secondary"},
			{name: "Cloudflare", endpoint: "1.1.1.1", location: "primary"},
			{name: "Cloudflare", endpoint: "1.0.0.1", location: "secondary"},
			{name: "Comodo", endpoint: "8.26.56.26", location: "primary"},
			{name: "Comodo", endpoint: "8.20.247.20", location: "secondary"},
			{name: "ControlD", endpoint: "76.76.2.0", location: "primary"},
			{name: "ControlD", endpoint: "76.76.10.0", location: "secondary"},
			{name: "Google", endpoint: "8.8.8.8", location: "primary"},
			{name: "Google", endpoint: "8.8.4.4", location: "secondary"},
			{name: "NextDNS", endpoint: "45.90.28.190", location: "primary"},
			{name: "NextDNS", endpoint: "45.90.30.190", location: "secondary"},
			{name: "OpenDNS", endpoint: "208.67.222.222", location: "primary"},
			{name: "OpenDNS", endpoint: "208.67.220.220", location: "secondary"},
			{name: "Quad9", endpoint: "9.9.9.9", location: "primary"},
			{name: "Quad9", endpoint: "149.112.112.112", location: "secondary"},
		},
	},
	"gcore": {
		name:         "Gcore Cloud",
		hostTemplate: "https://%s-latency.tools.gcore.com",
		regions: []*region{
			{name: "kal-2", location: "Almaty, KZ"},
			{name: "dra4-2", location: "Amsterdam, NL"},
			{name: "drc-2", location: "Chicago, US"},
			{name: "darz-2", location: "Darmstadt, DE"},
			{name: "dx1-2", location: "Dubai, AE"},
			{name: "drf-2", location: "Frankfurt, DE"},
			{name: "gnc-2", location: "Hong Kong, CN"},
			{name: "tii-2", location: "Istanbul, TR"},
			{name: "jb1-2", location: "Johannesburg, ZA"},
			{name: "ld3-2", location: "London, UK"},
			{name: "ed-16", location: "Luxembourg, LU"},
			{name: "anx-2", location: "Manassas, US"},
			{name: "ww", location: "Mumbai, IN", endpoint: "ww-speedtest.tools.gcore.com"},
			{name: "cwl1-2", location: "Newport, UK"},
			{name: "pa5-2", location: "Paris, FR"},
			{name: "scl2-2", location: "Santa Clara, US"},
			{name: "sp3-2", location: "Sao Paolo, BR"},
			{name: "sgc-2", location: "Singapore, SG"},
			{name: "sy4-2", location: "Sydney, AU"},
			{name: "cc1-2", location: "Tokyo, JP"},
			{name: "wa2-2", location: "Warsaw, PL"},
		},
	},
	"gcp": {
		name: "Google Cloud Platform",
		regions: []*region{
			{name: "asia-east1", location: "Changhua County, TW", endpoint: "https://asia-east1-5tkroniexa-de.a.run.app/api/ping"},
			{name: "asia-east2", location: "Hong Kong, CN", endpoint: "https://asia-east2-5tkroniexa-df.a.run.app/api/ping"},
			{name: "asia-northeast1", location: "Tokyo, JP", endpoint: "https://asia-northeast1-5tkroniexa-an.a.run.app/api/ping"},
			{name: "asia-northeast2", location: "Osaka, JP", endpoint: "https://asia-northeast2-5tkroniexa-dt.a.run.app/api/ping"},
			{name: "asia-northeast3", location: "Seoul, KR", endpoint: "https://asia-northeast3-5tkroniexa-du.a.run.app/api/ping"},
			{name: "asia-south1", location: "Mumbai, IN", endpoint: "https://asia-south1-5tkroniexa-el.a.run.app/api/ping"},
			{name: "asia-south2", location: "Delhi, IN", endpoint: "https://asia-south2-5tkroniexa-em.a.run.app/api/ping"},
			{name: "asia-southeast1", location: "Singapore, SG", endpoint: "https://asia-southeast1-5tkroniexa-as.a.run.app/api/ping"},
			{name: "asia-southeast2", location: "Jakarta, ID", endpoint: "https://asia-southeast2-5tkroniexa-et.a.run.app/api/ping"},
			{name: "australia-southeast1", location: "Sydney, AU", endpoint: "https://australia-southeast1-5tkroniexa-ts.a.run.app/api/ping"},
			{name: "australia-southeast2", location: "Melbourne, AU", endpoint: "https://australia-southeast2-5tkroniexa-km.a.run.app/api/ping"},
			{name: "europe-central2", location: "Warsaw, PL", endpoint: "https://europe-central2-5tkroniexa-lm.a.run.app/api/ping"},
			{name: "europe-north1", location: "Hamina, FI", endpoint: "https://europe-north1-5tkroniexa-lz.a.run.app/api/ping"},
			{name: "europe-southwest1", location: "Madrid, ES", endpoint: "https://europe-southwest1-5tkroniexa-no.a.run.app/api/ping"},
			{name: "europe-west1", location: "St. Ghislain, BE", endpoint: "https://europe-west1-5tkroniexa-ew.a.run.app/api/ping"},
			{name: "europe-west2", location: "London, UK", endpoint: "https://europe-west2-5tkroniexa-nw.a.run.app/api/ping"},
			{name: "europe-west3", location: "Frankfurt, DE", endpoint: "https://europe-west3-5tkroniexa-ey.a.run.app/api/ping"},
			{name: "europe-west4", location: "Eemshaven, NL", endpoint: "https://europe-west4-5tkroniexa-ez.a.run.app/api/ping"},
			{name: "europe-west6", location: "Zurich, CH", endpoint: "https://europe-west6-5tkroniexa-oa.a.run.app/api/ping"},
			{name: "europe-west8", location: "Milan, IT", endpoint: "https://europe-west8-5tkroniexa-oc.a.run.app/api/ping"},
			{name: "europe-west9", location: "Paris, FR", endpoint: "https://europe-west9-5tkroniexa-od.a.run.app/api/ping"},
			{name: "me-west1", location: "Tel Aviv, IL", endpoint: "https://me-west1-5tkroniexa-zf.a.run.app/api/ping"},
			{name: "northamerica-northeast1", location: "Montreal, CA", endpoint: "https://northamerica-northeast1-5tkroniexa-nn.a.run.app/api/ping"},
			{name: "northamerica-northeast2", location: "Toronto, CA", endpoint: "https://northamerica-northeast2-5tkroniexa-pd.a.run.app/api/ping"},
			{name: "southamerica-east1", location: "Osasco, BR", endpoint: "https://southamerica-east1-5tkroniexa-rj.a.run.app/api/ping"},
			{name: "southamerica-west1", location: "Santiago, CL", endpoint: "https://southamerica-west1-5tkroniexa-tl.a.run.app/api/ping"},
			{name: "us-central1", location: "Council Bluffs, US", endpoint: "https://us-central1-5tkroniexa-uc.a.run.app/api/ping"},
			{name: "us-east1", location: "Moncks Corner, US", endpoint: "https://us-east1-5tkroniexa-ue.a.run.app/api/ping"},
			{name: "us-east4", location: "Ashburn, US", endpoint: "https://us-east4-5tkroniexa-uk.a.run.app/api/ping"},
			{name: "us-east5", location: "Columbus, US", endpoint: "https://us-east5-5tkroniexa-ul.a.run.app/api/ping"},
			{name: "us-south1", location: "Dallas, US", endpoint: "https://us-south1-5tkroniexa-vp.a.run.app/api/ping"},
			{name: "us-west1", location: "The Dalles, US", endpoint: "https://us-west1-5tkroniexa-uw.a.run.app/api/ping"},
			{name: "us-west2", location: "Los Angeles, US", endpoint: "https://us-west2-5tkroniexa-wl.a.run.app/api/ping"},
			{name: "us-west3", location: "Salt Lake City, US", endpoint: "https://us-west3-5tkroniexa-wm.a.run.app/api/ping"},
			{name: "us-west4", location: "Las Vegas, US", endpoint: "https://us-west4-5tkroniexa-wn.a.run.app/api/ping"},
		},
	},
	"hetzner": {
		name:         "Hetzner Cloud",
		hostTemplate: "%s.icmp.hetzner.com",
		regions: []*region{
			{name: "ash", location: "Ashburn, US"},
			{name: "fsn", location: "Frankfurt, DE"},
			{name: "hel", location: "Helsinki, FI"},
			{name: "hil", location: "Hillsboro, US"},
			{name: "nbg", location: "Nurnberg, DE"},
		},
	},
	"huawei": {
		name:         "Huawei Cloud",
		hostTemplate: "https://dns.%s.myhuaweicloud.com",
		regions: []*region{
			{name: "ae-ad-1", location: "Abu Dhabi, AE"},
			{name: "af-south-1", location: "Johannesburg, ZA"},
			{name: "ap-southeast-1", location: "Hong Kong, CN"},
			{name: "ap-southeast-2", location: "Bangkok, TH"},
			{name: "ap-southeast-3", location: "Singapore, SG"},
			{name: "ap-southeast-4", location: "Jakarta, ID"},
			{name: "cn-east-2", location: "Shanghai, CN"},
			{name: "cn-east-3", location: "Shanghai, CN"},
			{name: "cn-north-1", location: "Beijing, CN"},
			{name: "cn-north-4", location: "Beijing, CN"},
			{name: "cn-south-1", location: "Guangzhou, CN"},
			{name: "eu-west-0", location: "Paris, FR"},
			{name: "eu-west-101", location: "Dublin, IE"},
			{name: "la-north-2", location: "Mexico City, MX"},
			{name: "la-south-2", location: "Santiago, CL"},
			{name: "me-east-1", location: "Riyadh, SA"},
			{name: "my-kualalumpur-1", location: "Kuala Lumpur, MY"},
			{name: "na-mexico-1", location: "Mexico City, MX"},
			{name: "ru-northwest-2", location: "Moscow, RU"},
			{name: "sa-argentina-1", location: "Buenos Aires, AR"},
			{name: "sa-brazil-1", location: "Sao Paulo, BR"},
			{name: "sa-peru-1", location: "Lima, PE"},
			{name: "tr-west-1", location: "Istanbul, TR"},
		},
	},
	"ibm": {
		name:         "IBM Cloud",
		hostTemplate: "s3.%s.cloud-object-storage.appdomain.cloud",
		regions: []*region{
			{name: "ams03", location: "Amsterdam, NL"},
			{name: "au-syd", location: "Sydney, AU"},
			{name: "br-sao", location: "Sao Paulo, BR"},
			{name: "che01", location: "Chennai, IN"},
			{name: "eu-de", location: "Frankfurt, DE"},
			{name: "eu-es", location: "Madrid, ES"},
			{name: "eu-gb", location: "London, UK"},
			// {name: "hkg02", location: "Hong Kong, CN"},
			{name: "jp-osa", location: "Osaka, JP"},
			{name: "jp-tok", location: "Tokyo, JP"},
			// {name: "mex01", location: "Mexico City, MX"},
			{name: "mil01", location: "Milan, IT"},
			{name: "mon01", location: "Montreal, CA"},
			{name: "par01", location: "Paris, FR"},
			{name: "sao01", location: "Sao Paulo, BR"},
			// {name: "seo01", location: "Seoul, KR"},
			{name: "sjc04", location: "San Jose, US"},
			{name: "sng01", location: "Singapore, SG"},
			{name: "tor01", location: "Toronto, CA"},
			{name: "us-east", location: "Washington, US"},
			{name: "us-south", location: "Dallas, US"},
		},
	},
	"kingsoft": {
		name:         "Kingsoft Cloud",
		hostTemplate: "ks3-%s.ksyuncs.com",
		regions: []*region{
			{name: "cn-beijing", location: "Beijing, CN"},
			{name: "cn-hk", location: "Hong Kong, CN", code: "cn-hk-1"},
			{name: "cn-shanghai", location: "Shanghai, CN"},
			{name: "cn-guangzhou", location: "Guangzhou, CN"},
			{name: "jr-beijing", location: "Beijing, CN"},
			{name: "jr-shanghai", location: "Shanghai, CN"},
			{name: "gov-beijing", location: "Beijing, CN"},
			{name: "rus", location: "Moscow, RU"},
			{name: "sgp", location: "Singapore, SG"},
		},
	},
	"linode": {
		name:         "Linode",
		hostTemplate: "speedtest.%s.linode.com",
		regions: []*region{
			{name: "ap-northeast", code: "tokyo2", location: "Tokyo, JP"},
			{name: "ap-south", code: "singapore", location: "Singapore, SG"},
			{name: "ap-southeast", code: "syd1", location: "Sydney, AU"},
			{name: "ap-west", code: "mumbai1", location: "Mumbai, IN"},
			{name: "ca-central", code: "toronto1", location: "Toronto, CA"},
			{name: "eu-central", code: "frankfurt", location: "Frankfurt, DE"},
			{name: "eu-west", code: "london", location: "London, UK"},
			{name: "fr-par", code: "paris", location: "Paris, FR"},
			{name: "it-mil", code: "milan", location: "Milan, IT"},
			{name: "se-sto", code: "stockholm", location: "Stockholm, SE"},
			{name: "us-central", code: "dallas", location: "Dallas, US"},
			{name: "us-east", code: "newark", location: "Newark, US"},
			{name: "us-iad", code: "washington", location: "Washington, US"},
			{name: "us-ord", code: "chicago", location: "Chicago, US"},
			{name: "us-sea", code: "seattle", location: "Seattle, US"},
			{name: "us-southeast", code: "atlanta", location: "Atlanta, US"},
			{name: "us-west", code: "fremont", location: "Fremont, US"},
		},
	},
	"oracle": {
		name:         "Oracle Cloud",
		hostTemplate: "objectstorage.%s.oraclecloud.com",
		regions: []*region{
			{name: "af-johannesburg-1", location: "Johannesburg, ZA"},
			{name: "ap-chuncheon-1", location: "Chuncheon, KR"},
			{name: "ap-hyderabad-1", location: "Hyderabad, IN"},
			{name: "ap-melbourne-1", location: "Melbourne, AU"},
			{name: "ap-mumbai-1", location: "Mumbai, IN"},
			{name: "ap-osaka-1", location: "Osaka, JP"},
			{name: "ap-seoul-1", location: "Seoul, KR"},
			{name: "ap-singapore-1", location: "Singapore, SG"},
			{name: "ap-sydney-1", location: "Sydney, AU"},
			{name: "ap-tokyo-1", location: "Tokyo, JP"},
			{name: "ca-montreal-1", location: "Montreal, CA"},
			{name: "ca-toronto-1", location: "Toronto, CA"},
			{name: "eu-amsterdam-1", location: "Amsterdam, NL"},
			{name: "eu-frankfurt-1", location: "Frankfurt, DE"},
			{name: "eu-madrid-1", location: "Madrid, ES"},
			{name: "eu-marseille-1", location: "Marseille, FR"},
			{name: "eu-milan-1", location: "Milan, IT"},
			{name: "eu-paris-1", location: "Paris, FR"},
			{name: "eu-stockholm-1", location: "Stockholm, SE"},
			{name: "eu-zurich-1", location: "Zurich, CH"},
			{name: "il-jerusalem-1", location: "Jerusalem, IL"},
			{name: "me-abudhabi-1", location: "Abu Dhabi, AE"},
			{name: "me-dubai-1", location: "Dubai, AE"},
			{name: "me-jeddah-1", location: "Jeddah, SA"},
			{name: "mx-monterrey-1", location: "Monterrey, MX"},
			{name: "mx-queretaro-1", location: "Queretaro, MX"},
			{name: "sa-santiago-1", location: "Santiago, CL"},
			{name: "sa-saopaulo-1", location: "Sao Paulo, BR"},
			{name: "sa-vinhedo-1", location: "Vinhedo, BR"},
			{name: "uk-cardiff-1", location: "Newport, UK"},
			{name: "uk-london-1", location: "London, UK"},
			{name: "us-ashburn-1", location: "Ashburn, US"},
			{name: "us-chicago-1", location: "Chicago, US"},
			{name: "us-phoenix-1", location: "Phoenix, US"},
			{name: "us-sanjose-1", location: "San Jose, US"},
		},
	},
	"ovh": {
		name:         "OVH Cloud",
		hostTemplate: "s3.%s.cloud.ovh.net",
		regions: []*region{
			{name: "bhs", location: "Beauharnois, CA"},
			{name: "de", location: "Frankfurt, DE"},
			{name: "gra", location: "Gravelines, FR"},
			{name: "hil", location: "Hillsboro, US", endpoint: "hil.proof.ovh.us"},
			{name: "lim", location: "Limburg, DE", endpoint: "s3.lim.io.cloud.ovh.net"},
			{name: "rbx", location: "Roubaix, FR", endpoint: "s3.rbx.io.cloud.ovh.net"},
			{name: "sbg", location: "Strasbourg, FR"},
			{name: "sgp", location: "Singapore, SG"},
			{name: "syd", location: "Sydney, AU"},
			{name: "uk", location: "Erith, UK"},
			{name: "vin", location: "Vint Hill, US", endpoint: "vin.proof.ovh.us"},
			{name: "waw", location: "Warsaw, PL"},
		},
	},
	"qing": {
		name:         "QingCloud",
		hostTemplate: "%s.qingstor.com",
		regions: []*region{
			{name: "gd2", location: "Guangdong, CN"},
			{name: "pek3a", location: "Beijing, CN"},
			{name: "pek3b", location: "Beijing, CN"},
			{name: "sh1a", location: "Shanghai, CN"},
		},
	},
	"qiniu": {
		name:         "Qiniu Cloud",
		hostTemplate: "up-%s.qiniup.com",
		regions: []*region{
			{name: "z0", location: "Huadong, CN"},
			{name: "z1", location: "Huabei, CN"},
			{name: "z2", location: "Huanan, CN"},
			{name: "na0", location: "USA"},
			{name: "as0", location: "Southeast Asia"},
			{name: "cn-east-2", location: "Zhejiang, CN"},
		},
	},
	"scaleway": {
		name:         "Scaleway",
		hostTemplate: "s3.%s.scw.cloud",
		regions: []*region{
			{name: "fr-par", location: "Paris, FR"},
			{name: "nl-ams", location: "Amsterdam, NL"},
			{name: "pl-waw", location: "Warsaw, PL"},
		},
	},
	"servers": {
		name:         "Servers.com",
		hostTemplate: "test.%s.servers.com",
		regions: []*region{
			{name: "ams1", location: "Amsterdam, NL"},
			{name: "ams2", location: "Haarlem, NL"},
			{name: "ams3", location: "Amsterdam, NL"},
			{name: "ams4", location: "Haarlem, NL"},
			{name: "ams5", location: "Haarlem, NL"},
			{name: "dfw1", location: "Dallas, US"},
			{name: "dfw2", location: "Dallas, US"},
			{name: "dfw3", location: "Dallas, US"},
			{name: "hkg1", location: "Hong Kong, CN"},
			{name: "lon1", location: "Slough, UK"},
			{name: "lux1", location: "Roost, LU"},
			{name: "lux2", location: "Roost, LU"},
			{name: "sin1", location: "Singapore, SG"},
			{name: "sjc1", location: "San Jose, US"},
			{name: "was1", location: "Washington, US"},
			{name: "was2", location: "Washington, US"},
		},
	},
	"tencent": {
		name:         "Tencent Cloud",
		hostTemplate: "cos.%s.myqcloud.com",
		regions: []*region{
			{name: "ap-bangkok", location: "Bangkok, TH"},
			{name: "ap-beijing", location: "Beijing, CN"},
			{name: "ap-chengdu", location: "Chengdu, CN"},
			{name: "ap-chongqing", location: "Chongqing, CN"},
			{name: "ap-guangzhou", location: "Guangzhou, CN"},
			{name: "ap-hongkong", location: "Hong Kong, CN"},
			{name: "ap-jakarta", location: "Jakarta, ID"},
			{name: "ap-mumbai", location: "Mumbai, IN"},
			{name: "ap-nanjing", location: "Nanjing, CN"},
			{name: "ap-seoul", location: "Seoul, KR"},
			{name: "ap-shanghai-fsi", location: "Shanghai, CN"},
			{name: "ap-shanghai", location: "Shanghai, CN"},
			{name: "ap-shenzhen-fsi", location: "Shenzhen, CN"},
			{name: "ap-singapore", location: "Singapore, SG"},
			{name: "ap-tokyo", location: "Tokyo, JP"},
			{name: "eu-frankfurt", location: "Frankfurt, DE"},
			{name: "na-ashburn", location: "Ashburn, US"},
			{name: "na-siliconvalley", location: "Silicon Valley, US"},
			{name: "na-toronto", location: "Toronto, US"},
			{name: "sa-saopaulo", location: "Sao Paulo, BR"},
		},
	},
	"ucloud": {
		name: "UCloud",
		regions: []*region{
			{name: "afr-nigeria", location: "Lagos, NG", endpoint: "feitsui-los.afr-nigeria.ufileos.com"},
			{name: "bra-saopaulo", location: "Sao Paulo, BR", endpoint: "feitsui-gru.bra-saopaulo.ufileos.com"},
			{name: "cn-bj", location: "Beijing, CN", endpoint: "feitsui-bjs.cn-bj.ufileos.com"},
			{name: "cn-gd", location: "Guangzhou, CN", endpoint: "feitsui-can.cn-gd.ufileos.com"},
			{name: "cn-sh2", location: "Shanghai, CN", endpoint: "feitsui-sha2.cn-sh2.ufileos.com"},
			{name: "ge-fra", location: "Frankfurt, DE", endpoint: "feitsui-fra.ge-fra.ufileos.com"},
			{name: "hk", location: "Hong Kong, CN", endpoint: "feitsui-hk.hk.ufileos.com"},
			{name: "idn-jakarta", location: "Jakarta, IN", endpoint: "feitsui-cgk.idn-jakarta.ufileos.com"},
			{name: "ind-mumbai", location: "Mumbai, IN", endpoint: "feitsui-bom.ind-mumbai.ufileos.com"},
			{name: "jpn-tky", location: "Tokyo, JP", endpoint: "s3-jpn-tky.ufileos.com"},
			{name: "kr-seoul", location: "Seoul, KR", endpoint: "feitsui-icn.kr-seoul.ufileos.com"},
			{name: "sg", location: "Singapore, SG", endpoint: "feitsui-sg.sg.ufileos.com"},
			{name: "th-bkk", location: "Bangkok, TH", endpoint: "s3-th-bkk.ufileos.com"},
			{name: "tw-tp", location: "Taipei, TW", endpoint: "feitsui-tpe.tw-tp.ufileos.com"},
			{name: "uae-dubai", location: "Dubai, UA", endpoint: "feitsui-dxb.uae-dubai.ufileos.com"},
			{name: "uk-london", location: "London, UK", endpoint: "s3-uk-london.ufileos.com"},
			{name: "us-ca", location: "Los Angeles, US", endpoint: "feitsui-lax.us-ca.ufileos.com"},
			{name: "us-ws", location: "Washington, US", endpoint: "feitsui-us-ws.us-ws.ufileos.com"},
			{name: "vn-sng", location: "Ho Chi Minh City, VN", endpoint: "feitsui-sgn.vn-sng.ufileos.com"},
		},
	},
	"upcloud": {
		name:         "UpCloud",
		hostTemplate: "%s.upcloudobjects.com",
		regions: []*region{
			{name: "au-syd1", location: "Sydney, AU"},
			{name: "de-fra1", location: "Frankfurt, DE"},
			{name: "es-mad1", location: "Madrid, ES"},
			{name: "fi-hel1", location: "Helsinki, FI", endpoint: "authdns1.fi-hel1.upcloud.com"},
			{name: "fi-hel2", location: "Helsinki, FI"},
			{name: "nl-ams1", location: "Amsterdam, NL"},
			{name: "pl-waw1", location: "Warsaw, PL"},
			{name: "sg-sin1", location: "Singapore, SG"},
			{name: "uk-lon1", location: "London, UK"},
			{name: "us-chi1", location: "Chicago, US"},
			{name: "us-nyc1", location: "New York, US"},
			{name: "us-sjo1", location: "San Jose, US"},
		},
	},
	"vultr": {
		name:         "Vultr",
		hostTemplate: "%s-ping.vultr.com",
		regions: []*region{
			{name: "ams-nl", location: "Amsterdam, NL"},
			{name: "blr-in", location: "Bangalore, IN"},
			{name: "bom-in", location: "Mumbai, IN"},
			{name: "del-in", location: "Delhi, IN"},
			{name: "fl-us", location: "Miami, US"},
			{name: "fra-de", location: "Frankfurt, DE"},
			{name: "ga-us", location: "Atlanta, US"},
			{name: "hnd-jp", location: "Tokyo, JP"},
			{name: "hon-hi-us", location: "Honolulu, US"},
			{name: "il-us", location: "Chicago, US"},
			{name: "jnb-za", location: "Johannesburg, ZA"},
			{name: "lax-ca-us", location: "Los Angeles, US"},
			{name: "lon-gb", location: "London, UK"},
			{name: "mad-es", location: "Madrid, ES"},
			{name: "man-uk", location: "Manchester, UK"},
			{name: "mel-au", location: "Melbourne, AU"},
			{name: "mex-mx", location: "Mexico City, MX"},
			{name: "nj-us", location: "New Jersey, US"},
			{name: "osk-jp", location: "Osaka, JP"},
			{name: "par-fr", location: "Paris, FR"},
			{name: "sao-br", location: "Sao Paulo, BR"},
			{name: "scl-cl", location: "Santiago, CL"},
			{name: "sel-kor", location: "Seoul, KR"},
			{name: "sgp", location: "Singapore, SG"},
			{name: "sjo-ca-us", location: "Silicon Valley, US"},
			{name: "sto-se", location: "Stockholm, SE"},
			{name: "syd-au", location: "Sydney, AU"},
			{name: "tor-ca", location: "Toronto, CA"},
			{name: "tlv-il", location: "Tel Aviv, IL"},
			{name: "tx-us", location: "Dallas, US"},
			{name: "wa-us", location: "Seattle, US"},
			{name: "waw-pl", location: "Warsaw, PL"},
		},
	},
}
