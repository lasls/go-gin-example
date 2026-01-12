package setting

import (
	"log"
	"strings"
	"time"

	"github.com/go-ini/ini"
)

var (
	Cfg *ini.File

	RunMode string

	HTTPPort     int
	ReadTimeout  time.Duration
	WriteTimeout time.Duration

	PageSize  int
	JwtSecret string

	// DNSPod配置
	DnsPodToken string
	DomainName  string

	// 阿里云DNS配置
	AliyunAccessKeyId     string
	AliyunAccessKeySecret string
	AliyunRegionId        string

	// 火山引擎DNS配置
	VolcengineAccessKeyId     string
	VolcengineAccessKeySecret string
	VolcengineRegionId        string

	// DNS服务器配置
	DnsServers []string
)

func init() {
	var err error
	Cfg, err = ini.Load("conf/app.ini")
	if err != nil {
		log.Fatalf("Fail to parse 'conf/app.ini': %v", err)
	}

	LoadBase()
	LoadServer()
	LoadApp()
	LoadDns()
	LoadAliyunDns()
	LoadVolcengineDns()
	LoadDnsServers()
}

func LoadBase() {
	RunMode = Cfg.Section("").Key("RUN_MODE").MustString("debug")
}

func LoadServer() {
	sec, err := Cfg.GetSection("server")
	if err != nil {
		log.Fatalf("Fail to get section 'server': %v", err)
	}

	HTTPPort = sec.Key("HTTP_PORT").MustInt(8000)
	ReadTimeout = time.Duration(sec.Key("READ_TIMEOUT").MustInt(60)) * time.Second
	WriteTimeout = time.Duration(sec.Key("WRITE_TIMEOUT").MustInt(60)) * time.Second
}

func LoadApp() {
	sec, err := Cfg.GetSection("app")
	if err != nil {
		log.Fatalf("Fail to get section 'app': %v", err)
	}

	JwtSecret = sec.Key("JWT_SECRET").MustString("!@)*#)!@U#@*!@!)")
	PageSize = sec.Key("PAGE_SIZE").MustInt(10)
}

func LoadDns() {
	sec, err := Cfg.GetSection("dns")
	if err != nil {
		log.Fatalf("Fail to get section 'dns': %v", err)
	}

	DnsPodToken = sec.Key("DNSPOD_TOKEN").MustString("")
	DomainName = sec.Key("DOMAIN_NAME").MustString("")
}

func LoadAliyunDns() {
	sec, err := Cfg.GetSection("aliyun_dns")
	if err != nil {
		log.Fatalf("Fail to get section 'aliyun_dns': %v", err)
	}

	AliyunAccessKeyId = sec.Key("ALIYUN_ACCESS_KEY_ID").MustString("")
	AliyunAccessKeySecret = sec.Key("ALIYUN_ACCESS_KEY_SECRET").MustString("")
	AliyunRegionId = sec.Key("ALIYUN_REGION_ID").MustString("cn-hangzhou")
}

func LoadVolcengineDns() {
	sec, err := Cfg.GetSection("volcengine_dns")
	if err != nil {
		log.Fatalf("Fail to get section 'volcengine_dns': %v", err)
	}

	VolcengineAccessKeyId = sec.Key("VOLCENGINE_ACCESS_KEY_ID").MustString("")
	VolcengineAccessKeySecret = sec.Key("VOLCENGINE_ACCESS_KEY_SECRET").MustString("")
	VolcengineRegionId = sec.Key("VOLCENGINE_REGION_ID").MustString("cn-north-1")
}

func LoadDnsServers() {
	sec, err := Cfg.GetSection("dns_servers")
	if err != nil {
		log.Printf("Fail to get section 'dns_servers': %v", err)
		// 设置默认值
		DnsServers = []string{"8.8.8.8:53", "8.8.4.4:53", "223.5.5.5:53", "1.1.1.1:53", "114.114.114.114:53"}
		return
	}

	// 从配置中读取DNS服务器列表
	dnsServersStr := sec.Key("DNS_SERVERS").MustString("8.8.8.8:53,8.8.4.4:53,223.5.5.5:53,1.1.1.1:53,114.114.114.114:53")
	DnsServers = strings.Split(dnsServersStr, ",")
	// 清理每个服务器地址的空白字符
	for i, server := range DnsServers {
		DnsServers[i] = strings.TrimSpace(server)
	}
}
