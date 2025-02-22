package rules

import (
	"strings"

	C "github.com/yuche/clashr/constant"
)

type Domain struct {
	domain  string
	adapter string
}

func (d *Domain) RuleType() C.RuleType {
	return C.Domain
}

func (d *Domain) IsMatch(metadata *C.Metadata) bool {
	if metadata.AddrType != C.AtypDomainName {
		return false
	}
	return metadata.Host == d.domain
}

func (d *Domain) Adapter() string {
	return d.adapter
}

func (d *Domain) Payload() string {
	return d.domain
}

func NewDomain(domain string, adapter string) *Domain {
	return &Domain{
		domain:  strings.ToLower(domain),
		adapter: adapter,
	}
}
