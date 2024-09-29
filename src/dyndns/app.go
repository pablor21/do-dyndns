package dyndns

import (
	"context"
	"encoding/json"
	"flag"
	"io"
	"log"
	"net/http"
	"time"

	"github.com/digitalocean/godo"
	"github.com/pablor21/do-dyndns/dyndns/config"
)

func Run() {
	// get flag -c for config file
	cfgFile := flag.String("c", "./config.yml", "config file")
	flag.Parse()
	config := config.LoadConfig(*cfgFile)
	app := NewApp(config)
	app.Run()
}

type App struct {
	config  *config.Config
	client  *godo.Client
	domains []*Domain
}

func NewApp(config *config.Config) *App {
	return &App{
		config: config,
		client: godo.NewFromToken(config.DoConfig.Token),
	}
}

func (a *App) Run() {
	ctx := context.Background()

	originalDomains, _, err := a.client.Domains.List(ctx, nil)
	if err != nil {
		panic(err)
	}

	godoDomains := make([]*Domain, 0)
	for _, domain := range originalDomains {
		godoDomains = append(godoDomains, NewDomainFromGodo(domain))
	}

	domains := make([]*Domain, 0)

	for _, domainName := range a.config.Domains {
		found := false
		for _, godoDomain := range godoDomains {
			if godoDomain.Name == domainName {
				domains = append(domains, godoDomain)
				found = true
				break
			}
		}

		if !found {
			g, _, err := a.client.Domains.Create(ctx, &godo.DomainCreateRequest{
				Name: domainName,
			})
			if err != nil {
				log.Fatal(err)
			} else {
				domains = append(domains, NewDomainFromGodo(*g))
			}
		}

	}

	a.domains = domains

	a.SyncDomains(ctx)

	ticker := time.NewTicker(time.Duration(a.config.Interval) * time.Second)
	for range ticker.C {
		a.SyncDomains(ctx)
	}

}

func (a *App) SyncDomains(ctx context.Context) {

	publicIp, err := a.GetPublicIp(ctx)
	if err != nil {
		log.Default().Printf("Error getting public ip: %s", err)
		return
	}
	for _, domain := range a.domains {
		records, _, err := a.client.Domains.RecordsByType(ctx, domain.Name, "A", nil)
		if err != nil {
			log.Default().Printf("Error getting records for domain %s: %s", domain.Name, err)
			continue
		}
		domain.Records = records

		a.SyncDomainIp(ctx, domain, *publicIp.IP)
	}
}

func (a *App) SyncDomainIp(ctx context.Context, domain *Domain, ip string) {
	log.Default().Printf("Checking domain %s", domain.Name)
	if len(domain.Records) == 0 {
		_, _, err := a.client.Domains.CreateRecord(ctx, domain.Name, &godo.DomainRecordEditRequest{
			Type: "A",
			Name: "@",
			Data: ip,
			TTL:  1800,
		})
		if err != nil {
			log.Default().Printf("Error creating record for domain %s: %s", domain.Name, err)
		} else {
			log.Default().Printf("Created record for domain %s with ip %s", domain.Name, ip)
		}
		return
	}

	// check the current ip
	currentIp := domain.Records[0].Data
	if currentIp == ip {
		return
	}

	// update the ip
	record := &godo.DomainRecordEditRequest{
		Data: ip,
		Type: "A",
		TTL:  1800,
	}

	_, _, err := a.client.Domains.EditRecord(ctx, domain.Name, domain.Records[0].ID, record)
	if err != nil {
		log.Default().Printf("Error updating record for domain %s: %s", domain.Name,
			err)
	} else {
		log.Default().Printf("Updated record for domain %s to %s", domain.Name, ip)
	}
}

func (a *App) GetPublicIp(ctx context.Context) (*PublicIP, error) {
	var p *PublicIP
	url := a.config.IpClient.Uri
	client := &http.Client{}
	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return nil, err
	}

	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}

	if resp.Body != nil {
		defer resp.Body.Close()
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	err = json.Unmarshal(body, &p)
	if err != nil {
		return nil, err
	}

	return p, nil
}
