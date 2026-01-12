package dns

import (
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
	"time"

	"github.com/miekg/dns"
)

// DomainVerificationService 域名验证服务
type DomainVerificationService struct {
	dnsManager *DnsManager
}

// VerificationRequest 验证请求
type VerificationRequest struct {
	Domain string `json:"domain"`
	Type   string `json:"type"` // "dns", "file", "cname"
}

// VerificationResponse 验证响应
type VerificationResponse struct {
	Success     bool   `json:"success"`
	Message     string `json:"message"`
	Challenge   string `json:"challenge,omitempty"`    // 验证挑战内容
	RecordType  string `json:"record_type,omitempty"`  // 记录类型，如TXT、CNAME
	RecordName  string `json:"record_name,omitempty"`  // 记录名称，如 _cdnauth.example.com
	FileName    string `json:"file_name,omitempty"`    // 文件验证时的文件名
	FileContent string `json:"file_content,omitempty"` // 文件验证时的内容
}

// NewDomainVerificationService 创建域名验证服务
func NewDomainVerificationService(dnsManager *DnsManager) *DomainVerificationService {
	return &DomainVerificationService{
		dnsManager: dnsManager,
	}
}

// GenerateChallenge 生成验证挑战
func (s *DomainVerificationService) GenerateChallenge(domain, verificationType string) (*VerificationResponse, error) {
	switch strings.ToLower(verificationType) {
	case "dns":
		return s.generateDNSChallenge(domain)
	case "cname":
		return s.generateCNAMEChallenge(domain)
	case "file":
		return s.generateFileChallenge(domain)
	default:
		return nil, fmt.Errorf("不支持的验证类型: %s", verificationType)
	}
}

// generateDNSChallenge 生成DNS验证挑战
func (s *DomainVerificationService) generateDNSChallenge(domain string) (*VerificationResponse, error) {
	// 生成唯一的验证token
	token := s.generateToken()

	// 构造验证记录名，格式为 _cdnauth.域名
	recordName := "_cdnauth." + domain

	return &VerificationResponse{
		Success:    true,
		Message:    "请在您的DNS设置中添加以下TXT记录",
		Challenge:  token,
		RecordType: "TXT",
		RecordName: recordName,
	}, nil
}

// generateCNAMEChallenge 生成CNAME验证挑战
func (s *DomainVerificationService) generateCNAMEChallenge(domain string) (*VerificationResponse, error) {
	// 生成唯一的验证token
	token := s.generateToken()

	// 构造验证记录名，格式为 _cdnauth.域名
	recordName := "_cdnauth." + domain

	// 使用一个固定的验证服务器地址
	verificationServer := fmt.Sprintf("checkpoint-%s.xldns.com", token[:8])

	return &VerificationResponse{
		Success:    true,
		Message:    "请在您的DNS设置中添加以下CNAME记录",
		Challenge:  verificationServer,
		RecordType: "CNAME",
		RecordName: recordName,
	}, nil
}

// generateFileChallenge 生成文件验证挑战
func (s *DomainVerificationService) generateFileChallenge(domain string) (*VerificationResponse, error) {
	// 生成唯一的验证token
	token := s.generateToken()

	// 生成文件名和内容
	fileName := fmt.Sprintf("cdn-auth-%s.txt", token[:16])
	fileContent := token

	return &VerificationResponse{
		Success:     true,
		Message:     "请将以下内容保存为文件并上传到网站根目录",
		FileName:    fileName,
		FileContent: fileContent,
	}, nil
}

// VerifyDomain 执行域名验证
func (s *DomainVerificationService) VerifyDomain(domain, verificationType, challenge string) (*VerificationResponse, error) {
	switch strings.ToLower(verificationType) {
	case "dns":
		return s.verifyDNS(domain, challenge)
	case "cname":
		return s.verifyCNAME(domain, challenge)
	case "file":
		return s.verifyFile(domain, challenge)
	default:
		return nil, fmt.Errorf("不支持的验证类型: %s", verificationType)
	}
}

// verifyDNS 验证DNS记录
func (s *DomainVerificationService) verifyDNS(domain, expectedToken string) (*VerificationResponse, error) {
	// 构造查询的记录名
	recordName := "_cdnauth." + domain

	// 查询TXT记录
	txtRecords, err := s.queryTXTRecord(recordName)
	if err != nil {
		return &VerificationResponse{
			Success: false,
			Message: fmt.Sprintf("查询DNS记录失败: %v", err),
		}, nil
	}

	// 检查返回的TXT记录中是否包含预期的token
	for _, txt := range txtRecords {
		if strings.Contains(txt, expectedToken) {
			return &VerificationResponse{
				Success: true,
				Message: "DNS验证成功",
			}, nil
		}
	}

	return &VerificationResponse{
		Success: false,
		Message: "未找到匹配的TXT记录，请确认已正确添加DNS记录并等待生效",
	}, nil
}

// verifyCNAME 验证CNAME记录
func (s *DomainVerificationService) verifyCNAME(domain, expectedTarget string) (*VerificationResponse, error) {
	// 构造查询的记录名
	recordName := "_cdnauth." + domain

	// 查询CNAME记录
	cnameRecord, err := s.queryCNAMERecord(recordName)
	if err != nil {
		return &VerificationResponse{
			Success: false,
			Message: fmt.Sprintf("查询CNAME记录失败: %v", err),
		}, nil
	}

	// 检查返回的CNAME记录是否匹配, 截取尾部的点(点号.表示根域名，在DNS解析中代表完整域名)
	if cnameRecord == expectedTarget || strings.TrimSuffix(cnameRecord, ".") == expectedTarget {
		return &VerificationResponse{
			Success: true,
			Message: "CNAME验证成功",
		}, nil
	}

	return &VerificationResponse{
		Success: false,
		Message: fmt.Sprintf("CNAME记录不匹配，期望: %s，实际: %s", expectedTarget, cnameRecord),
	}, nil
}

// verifyFile 验证文件
func (s *DomainVerificationService) verifyFile(domain, expectedContent string) (*VerificationResponse, error) {
	// 构造文件URL
	fileURL := fmt.Sprintf("http://%s/.well-known/cdn-auth-%s.txt", domain, expectedContent[:16])

	// 尝试HTTP访问
	content, err := s.fetchFileContent(fileURL)
	if err != nil {
		// 如果HTTP失败，尝试HTTPS
		fileURL = fmt.Sprintf("https://%s/.well-known/cdn-auth-%s.txt", domain, expectedContent[:16])
		content, err = s.fetchFileContent(fileURL)
		if err != nil {
			return &VerificationResponse{
				Success: false,
				Message: fmt.Sprintf("无法访问验证文件: %v", err),
			}, nil
		}
	}

	if strings.TrimSpace(content) == expectedContent {
		return &VerificationResponse{
			Success: true,
			Message: "文件验证成功",
		}, nil
	}

	return &VerificationResponse{
		Success: false,
		Message: "文件内容不匹配",
	}, nil
}

// queryTXTRecord 查询TXT记录
func (s *DomainVerificationService) queryTXTRecord(domain string) ([]string, error) {
	c := new(dns.Client)
	m := new(dns.Msg)
	m.SetQuestion(dns.Fqdn(domain), dns.TypeTXT)

	r, _, err := c.Exchange(m, "8.8.8.8:53") // 使用Google DNS
	if err != nil {
		return nil, err
	}

	if r.Rcode != dns.RcodeSuccess {
		return nil, fmt.Errorf("DNS查询失败: %s", dns.RcodeToString[r.Rcode])
	}

	var txtRecords []string
	for _, ans := range r.Answer {
		if t, ok := ans.(*dns.TXT); ok {
			txtRecords = append(txtRecords, strings.Join(t.Txt, " "))
		}
	}

	return txtRecords, nil
}

// queryCNAMERecord 查询CNAME记录
func (s *DomainVerificationService) queryCNAMERecord(domain string) (string, error) {
	c := new(dns.Client)
	m := new(dns.Msg)
	m.SetQuestion(dns.Fqdn(domain), dns.TypeCNAME)

	r, _, err := c.Exchange(m, "8.8.8.8:53") // 使用Google DNS
	if err != nil {
		return "", err
	}

	if r.Rcode != dns.RcodeSuccess {
		return "", fmt.Errorf("DNS查询失败: %s", dns.RcodeToString[r.Rcode])
	}

	for _, ans := range r.Answer {
		if c, ok := ans.(*dns.CNAME); ok {
			return c.Target, nil
		}
	}

	return "", fmt.Errorf("未找到CNAME记录")
}

// fetchFileContent 获取文件内容
func (s *DomainVerificationService) fetchFileContent(url string) (string, error) {
	client := &http.Client{
		Timeout: 10 * time.Second,
	}

	resp, err := client.Get(url)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("HTTP状态码: %d", resp.StatusCode)
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	return string(body), nil
}

// generateToken 生成随机token
func (s *DomainVerificationService) generateToken() string {
	bytes := make([]byte, 16)
	if _, err := rand.Read(bytes); err != nil {
		// 如果随机数生成失败，使用时间戳作为备选
		return hex.EncodeToString([]byte(fmt.Sprintf("%d", time.Now().UnixNano())))
	}
	return hex.EncodeToString(bytes)
}
