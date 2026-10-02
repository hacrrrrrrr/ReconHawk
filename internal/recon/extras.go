package recon

import (
 "context"
 "encoding/json"
 "fmt"
 "io"
 "net/http"
 "net/url"
 "regexp"
 "sort"
 "strings"
 "time"
)

type Subdomain struct { Host string `json:"host"`; Source string `json:"source"` }

type HistoricalURL struct { URL string `json:"url"`; Source string `json:"source"` }

type Finding struct {
 RuleID string `json:"rule_id"`
 Severity string `json:"severity"`
 URL string `json:"url"`
 Evidence string `json:"evidence"`
 Remediation string `json:"remediation"`
}
type CheckReport struct { Target string `json:"target"`; StatusCode int `json:"status_code"`; Findings []Finding `json:"findings"` }

func FetchCommonCrawl(ctx context.Context, target string)([]HistoricalURL,error){
 d:=strings.TrimSpace(strings.TrimSuffix(strings.TrimPrefix(target,"https://"),"/"))
 if i:=strings.IndexByte(d,'/');i>=0{d=d[:i]}
 endpoint:="https://index.commoncrawl.org/CC-MAIN-2026-30-index?url=*."+url.QueryEscape(d)+"/*&output=json&filter=status:200"
 req,e:=http.NewRequestWithContext(ctx,http.MethodGet,endpoint,nil);if e!=nil{return nil,e}
 resp,e:=http.DefaultClient.Do(req);if e!=nil{return nil,e};defer resp.Body.Close()
 if resp.StatusCode/100!=2{return nil,fmt.Errorf("Common Crawl returned %s",resp.Status)}
 b,e:=io.ReadAll(io.LimitReader(resp.Body,8<<20));if e!=nil{return nil,e}
 seen:=map[string]struct{}{};out:=[]HistoricalURL{}
 for _,line:=range strings.Split(string(b),"\n"){var row struct{URL string `json:"url"`};if json.Unmarshal([]byte(line),&row)==nil&&row.URL!=""{if _,ok:=seen[row.URL];!ok{seen[row.URL]=struct{}{};out=append(out,HistoricalURL{URL:row.URL,Source:"commoncrawl"})}}}
 sort.Slice(out,func(i,j int)bool{return out[i].URL<out[j].URL});if len(out)>5000{out=out[:5000]};return out,nil
}

func FetchBody(ctx context.Context, raw string, timeout time.Duration)(string,string,error){
 if timeout<=0{timeout=10*time.Second};u,e:=url.Parse(raw);if e!=nil||u.Hostname()==""{return "","",fmt.Errorf("invalid URL")}
 c:=&http.Client{Timeout:timeout};req,e:=http.NewRequestWithContext(ctx,http.MethodGet,u.String(),nil);if e!=nil{return "","",e};req.Header.Set("User-Agent","ReconHawk/0.3")
 resp,e:=c.Do(req);if e!=nil{return "","",e};defer resp.Body.Close();b,e:=io.ReadAll(io.LimitReader(resp.Body,4<<20));if e!=nil{return "","",e};return string(b),resp.Request.URL.String(),nil
}

func ExtractJavaScriptEndpoints(body,base string)[]string{
 re:=regexp.MustCompile("(?i)[\"']((?:/|https?://)[A-Za-z0-9_./?=&:%#@+~{}-]{2,})[\"']")
 u,_:=url.Parse(base);seen:=map[string]struct{}{}
 for _,m:=range re.FindAllStringSubmatch(body,-1){x:=m[1];v,e:=url.Parse(x);if e!=nil{continue};if v.IsAbs()&&(v.Scheme!="http"&&v.Scheme!="https"){continue};if !v.IsAbs(){v=u.ResolveReference(v)};v.Fragment="";seen[v.String()]=struct{}{}}
 out:=make([]string,0,len(seen));for x:=range seen{out=append(out,x)};sort.Strings(out);if len(out)>500{out=out[:500]};return out
}

func CheckTarget(ctx context.Context, raw string, timeout time.Duration)(CheckReport,error){
 if timeout<=0{timeout=8*time.Second};u,e:=normalizeTarget(raw);if e!=nil{return CheckReport{},e}
 req,e:=http.NewRequestWithContext(ctx,http.MethodGet,u,nil);if e!=nil{return CheckReport{},e};req.Header.Set("User-Agent","ReconHawk/0.3")
 resp,e:=(&http.Client{Timeout:timeout}).Do(req);if e!=nil{return CheckReport{},e};defer resp.Body.Close()
 r:=CheckReport{Target:resp.Request.URL.String(),StatusCode:resp.StatusCode,Findings:[]Finding{}}
 if resp.Header.Get("Content-Security-Policy")==""{r.Findings=append(r.Findings,Finding{"RH-SEC-001","medium",r.Target,"Content-Security-Policy header is absent","Deploy a restrictive Content-Security-Policy appropriate to the application."})}
 if resp.Request.URL.Scheme=="https"&&resp.Header.Get("Strict-Transport-Security")==""{r.Findings=append(r.Findings,Finding{"RH-SEC-002","medium",r.Target,"Strict-Transport-Security header is absent","Enable HSTS after validating HTTPS coverage."})}
 if resp.Header.Get("X-Content-Type-Options")==""{r.Findings=append(r.Findings,Finding{"RH-SEC-003","low",r.Target,"X-Content-Type-Options header is absent","Set X-Content-Type-Options: nosniff."})}
 if v:=resp.Header.Get("Access-Control-Allow-Origin");v=="*" {r.Findings=append(r.Findings,Finding{"RH-SEC-004","low",r.Target,"Wildcard Access-Control-Allow-Origin","Restrict CORS origins where cross-origin access is not intentionally public."})}
 for _,c:=range resp.Cookies(){if resp.Request.URL.Scheme=="https"&&!c.Secure{r.Findings=append(r.Findings,Finding{"RH-SEC-005","low",r.Target,"Cookie "+c.Name+" lacks Secure attribute","Set Secure on cookies sent over HTTPS."})};if !c.HttpOnly{r.Findings=append(r.Findings,Finding{"RH-SEC-006","low",r.Target,"Cookie "+c.Name+" lacks HttpOnly attribute","Use HttpOnly for cookies that do not need JavaScript access."})}}
 return r,nil
}

type Result struct {
 URL string `json:"url"`
 StatusCode int `json:"status_code,omitempty"`
 FinalURL string `json:"final_url,omitempty"`
 ContentType string `json:"content_type,omitempty"`
 Server string `json:"server,omitempty"`
 Technologies []string `json:"technologies,omitempty"`
 SecurityHeaders map[string]string `json:"security_headers,omitempty"`
 TLS *TLSInfo `json:"tls,omitempty"`
 DNS []DNSObservation `json:"dns,omitempty"`
 Findings []Finding `json:"findings,omitempty"`
}

func Fingerprint(raw string) (Result,error) {
 r,e:=Scan(Options{Target:raw,Workers:1,Timeout:10*time.Second,Depth:0})
 if e!=nil{return Result{},e}
 if len(r.Observations)==0{return Result{URL:r.Target,DNS:r.DNS},nil}
 o:=r.Observations[0]
 return Result{URL:o.URL,StatusCode:o.StatusCode,FinalURL:o.FinalURL,ContentType:o.ContentType,Server:o.Server,Technologies:o.Technologies,SecurityHeaders:o.SecurityHeaders,TLS:o.TLS,DNS:r.DNS},nil
}

func Check(raw string) (Result,error) {
 cr,e:=CheckTarget(context.Background(),raw,8*time.Second)
 if e!=nil{return Result{},e}
 return Result{URL:cr.Target,StatusCode:cr.StatusCode,Findings:cr.Findings},nil
}

func Subdomains(domain string) ([]string,error) {
 rows,e:=CRTShProvider{}.Run(context.Background(),domain)
 return rows,e
}

func Historical(domain string) ([]string,error) {
 rows,e:=FetchCommonCrawl(context.Background(),domain)
 if e!=nil{return nil,e}
 out:=make([]string,0,len(rows));for _,r:=range rows{out=append(out,r.URL)}
 return out,nil
}

func JavaScript(raw string) ([]string,error) {
 body,final,e:=FetchBody(context.Background(),raw,10*time.Second)
 if e!=nil{return nil,e}
 return ExtractJavaScriptEndpoints(body,final),nil
}
