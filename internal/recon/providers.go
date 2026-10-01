package recon

import (
 "context"
 "encoding/json"
 "fmt"
 "net/http"
 "net/url"
 "strings"
)

type Provider interface { Name() string; Run(context.Context,string) ([]string,error) }
type CRTShProvider struct{}
func (CRTShProvider) Name() string { return "crt.sh" }
type crtEntry struct { NameValue string `json:"name_value"` }

func (CRTShProvider) Run(ctx context.Context, domain string) ([]string,error) {
 domain=strings.TrimSuffix(strings.ToLower(strings.TrimSpace(domain)),".")
 if domain=="" { return nil,fmt.Errorf("empty domain") }
 endpoint:="https://crt.sh/?q="+url.QueryEscape("%."+domain)+"&output=json"
 req,err:=http.NewRequestWithContext(ctx,http.MethodGet,endpoint,nil);if err!=nil{return nil,err}
 resp,err:=http.DefaultClient.Do(req);if err!=nil{return nil,err};defer resp.Body.Close()
 if resp.StatusCode/100!=2{return nil,fmt.Errorf("crt.sh returned %s",resp.Status)}
 var rows []crtEntry;if err:=json.NewDecoder(resp.Body).Decode(&rows);err!=nil{return nil,err}
 seen:=map[string]struct{}{};for _,r:=range rows{for _,h:=range strings.Split(r.NameValue,"\n"){h=strings.TrimSuffix(strings.TrimPrefix(strings.ToLower(strings.TrimSpace(h)),"*."),".");if h!=""&&(h==domain||strings.HasSuffix(h,"."+domain)){seen[h]=struct{}{}}}}
 out:=make([]string,0,len(seen));for h:=range seen{out=append(out,h)};sortStrings(out);return out,nil
}
func CollectSubdomains(ctx context.Context,domain string,providers []Provider)([]Subdomain,error){seen:=map[string]Subdomain{};for _,p:=range providers{hs,err:=p.Run(ctx,domain);if err!=nil{continue};for _,h:=range hs{seen[h]=Subdomain{Host:h,Source:p.Name()}}};out:=make([]Subdomain,0,len(seen));for _,s:=range seen{out=append(out,s)};for i:=1;i<len(out);i++{for j:=i;j>0&&out[j].Host<out[j-1].Host;j--{out[j],out[j-1]=out[j-1],out[j]}};return out,nil}
func sortStrings(a []string){for i:=1;i<len(a);i++{for j:=i;j>0&&a[j]<a[j-1];j--{a[j],a[j-1]=a[j-1],a[j]}}}
