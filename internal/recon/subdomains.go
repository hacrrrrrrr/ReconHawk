package recon

import (
 "context"
 "net"
 "net/http"
 "net/url"
 "strings"
)

type Subdomain struct { Host string `json:"host"`; Source string `json:"source"`; Reachable bool `json:"reachable"` }

func PassiveSubdomains(ctx context.Context, domain string) []Subdomain {
 domain=strings.TrimSuffix(strings.ToLower(strings.TrimSpace(domain)),".")
 if domain=="" { return nil }
 // Provider adapters are intentionally conservative: only public, non-authenticated
 // sources configured by the caller should be added here.
 return []Subdomain{{Host:domain,Source:"base-domain",Reachable:hostReachable(ctx,domain)}}
}

func hostReachable(ctx context.Context, host string) bool {
 addrs,err:=net.DefaultResolver.LookupHost(ctx,host); if err!=nil || len(addrs)==0{return false}
 u:=url.URL{Scheme:"https",Host:host}
 req,err:=http.NewRequestWithContext(ctx,http.MethodHead,u.String(),nil);if err!=nil{return false}
 resp,err:=http.DefaultClient.Do(req);if err!=nil{return false};resp.Body.Close();return true
}
