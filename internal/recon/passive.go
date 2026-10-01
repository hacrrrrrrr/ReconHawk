package recon

import (
 "context"
 "encoding/json"
 "fmt"
 "net/http"
 "net/url"
 "strings"
)

type PassiveProvider interface {
 Name() string
 Run(context.Context, string) ([]string, error)
}

type HackerTargetProvider struct{}
func (HackerTargetProvider) Name() string { return "hackertarget" }
func (HackerTargetProvider) Run(ctx context.Context, domain string) ([]string,error) {
 endpoint := "https://api.hackertarget.com/hostsearch/?q="+url.QueryEscape(domain)
 req,err:=http.NewRequestWithContext(ctx,http.MethodGet,endpoint,nil);if err!=nil{return nil,err}
 resp,err:=http.DefaultClient.Do(req);if err!=nil{return nil,err};defer resp.Body.Close()
 if resp.StatusCode/100!=2{return nil,fmt.Errorf("hackertarget returned %s",resp.Status)}
 var out []string
 dec:=json.NewDecoder(resp.Body)
 _=dec
 body:=make([]byte,0)
 buf:=make([]byte,64*1024)
 for {n,e:=resp.Body.Read(buf);if n>0{body=append(body,buf[:n]...)};if e!=nil{break}}
 for _,line:=range strings.Split(string(body),"\n"){parts:=strings.SplitN(strings.TrimSpace(line),",",2);if len(parts)>0&&strings.Contains(parts[0],"."){out=append(out,parts[0])}}
 return uniqueStrings(out),nil
}

type AlienVaultOTXProvider struct{}
func (AlienVaultOTXProvider) Name() string { return "alienvault-otx" }
func (AlienVaultOTXProvider) Run(ctx context.Context,domain string)([]string,error){
 endpoint:="https://otx.alienvault.com/api/v1/indicators/domain/"+url.PathEscape(domain)+"/passive_dns"
 req,err:=http.NewRequestWithContext(ctx,http.MethodGet,endpoint,nil);if err!=nil{return nil,err}
 resp,err:=http.DefaultClient.Do(req);if err!=nil{return nil,err};defer resp.Body.Close()
 if resp.StatusCode/100!=2{return nil,fmt.Errorf("otx returned %s",resp.Status)}
 var data struct{PassiveDNS []struct{Hostname string `json:"hostname"`} `json:"passive_dns"`}
 if err:=json.NewDecoder(resp.Body).Decode(&data);err!=nil{return nil,err}
 out:=[]string{};for _,x:=range data.PassiveDNS{if x.Hostname!=""{out=append(out,x.Hostname)}};return uniqueStrings(out),nil
}

func CollectPassive(ctx context.Context,domain string,providers []PassiveProvider) []Subdomain {
 seen:=map[string]Subdomain{}
 for _,p:=range providers {
  hosts,err:=p.Run(ctx,domain);if err!=nil{continue}
  for _,h:=range hosts {h=strings.TrimSuffix(strings.ToLower(strings.TrimSpace(h)),".");if h!=""{seen[h]=Subdomain{Host:h,Source:p.Name()}}}
 }
 out:=make([]Subdomain,0,len(seen));for _,x:=range seen{out=append(out,x)};return out
}
func uniqueStrings(in []string)[]string{m:=map[string]struct{}{};out:=[]string{};for _,x:=range in{if _,ok:=m[x];!ok{m[x]=struct{}{};out=append(out,x)}};return out}
