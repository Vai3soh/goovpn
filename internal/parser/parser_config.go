package parser

import (
	"bufio"
	"os"
	"regexp"
	"strings"

	"github.com/Vai3soh/goovpn/entity"
)

type Config struct {
	entity.Profile
}

type Option func(*Config)

func NewConfig(opts ...Option) *Config {
	c := &Config{}
	for _, opt := range opts {
		opt(c)
	}
	return c
}

func WithBody(body string) Option {
	return func(c *Config) {
		c.Body = body
	}
}

func WithPath(path string) Option {
	return func(c *Config) {
		c.Path = path
	}
}

func (c *Config) SetPath(path string) {
	c.Path = path
}

func (c *Config) SetBody(body string) {
	c.Body = body
}

func (c *Config) GetBody() string {
	return c.Body
}

func (c *Config) RemoveSpaceLines() {

	c.Body = regexp.MustCompile(`(?m)^\s+`).ReplaceAllString(c.Body, "")
}

func (c *Config) RemoveCommentLines() {
	c.Body = regexp.MustCompile(`(?m)^#.+|^#`).ReplaceAllString(c.Body, "")
}

func (c *Config) RemoveEmptyString() {
	c.Body = regexp.MustCompile(`(?m)^\s*$`).ReplaceAllString(c.Body, "")
}

func (c *Config) RemoveNotCertsAndKeys() {
	c.Body = regexp.MustCompile(`(?m)[^,;]+deadline`).ReplaceAllString(c.Body, "")
}

func (c *Config) RemoveCertsAndKeys() {
	r := `(?m)deadline[^,;]+`
	c.Body = regexp.MustCompile(r).ReplaceAllString(c.Body, "")
}

func contains(s []string, str string) bool {
	for _, v := range s {
		if v == str {
			return true
		}
	}
	return false
}

func (c *Config) CheckConfigUseFiles() bool {
	r := `(?m)^cert\s+|^ca\s+|^key\s+|^tls-auth\s+|^tls-crypt\s+` //
	return regexp.MustCompile(r).MatchString(c.Body)
}

func (c *Config) AddStringToConfig(inFile *os.File) string {

	defer inFile.Close()
	scanner := bufio.NewScanner(inFile)
	scanner.Split(bufio.ScanLines)

	var strs []string
	substr := `deadline`
	for scanner.Scan() {
		if !contains(strs, substr) {
			r := `(?m)^ca\s+|^key\s+|^cert\s+|^tls-auth\s+|^tls-crypt\s+`
			if regexp.MustCompile(r).MatchString(scanner.Text()) {
				strs = append(strs, substr)
			}
		}
		strs = append(strs, scanner.Text())
	}
	body := strings.Join(strs, "\n")
	return body
}

func (c *Config) SearchFilesPaths() map[string]string {

	CertKeysMap := make(map[string]string)
	reg := regexp.MustCompile(`(?m)(^\w+|\w+-\w+)\s+(.+\S\S)`)
	matches := reg.FindAllStringSubmatch(c.Body, -1)
	for _, match := range matches {
		CertKeysMap[match[1]] = match[2]
	}
	return CertKeysMap
}

func (c *Config) MergeCertsAndKeys(key string) string {
	return "\n<" + key + ">\n" + c.Body + "</" + key + ">\n"
}

func (c *Config) GetAuthpathFileName() string {
	authpathFileName := regexp.MustCompile(`auth-user-pass(.*)`).FindStringSubmatch(c.GetBody())
	s := strings.Replace(authpathFileName[1], "'", "", 2)
	s = strings.Replace(s, "\"", "", 2)
	s = strings.Trim(s, " ")
	return s
}

func (c *Config) GetUserAndPass() (string, string) {
	userAndPass := strings.Split(c.Body, "\n")
	return userAndPass[0], userAndPass[1]
}

func (c *Config) RemoveComments() string {
	s := regexp.MustCompile(`(#.+)`).ReplaceAllString(c.Body, "")
	return s
}

func (c *Config) CheckStringAuthUserPass() bool {
	s := c.RemoveComments()
	return strings.Contains(s, `auth-user-pass`)
}
