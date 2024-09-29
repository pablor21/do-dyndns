package dyndns

import "github.com/digitalocean/godo"

type Domain struct {
	godo.Domain
	Records []godo.DomainRecord
	Exists  bool
}

func NewDomain(name string) *Domain {
	return &Domain{
		godo.Domain{
			Name: name,
		},
		[]godo.DomainRecord{},
		false,
	}
}

func NewDomainFromGodo(domain godo.Domain) *Domain {
	return &Domain{
		Domain:  domain,
		Records: []godo.DomainRecord{},
		Exists:  true,
	}
}

func (d *Domain) SetExists(exists bool) {
	d.Exists = exists
}

func (d *Domain) GetExists() bool {
	return d.Exists
}

func (d *Domain) Compare(b *Domain) int {
	if d.Name < b.Name {
		return -1
	}
	if d.Name > b.Name {
		return 1
	}
	return 0
}

type PublicIP struct {
	IP *string `json:"ip"`
}
