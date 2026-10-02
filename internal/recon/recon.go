package recon

import ("context"; "crypto/tls"; "errors"; "fmt"; "io"; "net"; "net/http"; "net/url"; "regexp"; "sort"; "strings"; "sync"; "time")

type Options struct { Target string; Workers int; Timeout time.Duration; Depth int }
type Report struct { Target string; StartedAt time.Time; FinishedAt time.Time; Observations []Observation; DiscoveredURLs []string; DNS []DNSObservation }
type DNSObservation struct { Host string; Addresses []string; Error string }
type Observation struct { URL string; StatusCode int; ContentType string; Server string; FinalURL string; ContentLength int64; Technologies []string; SecurityHeaders map[string]string; TLS *TLSInfo; Links []string; Error string }
type TLSInfo struct { Version string; Cipher string }

func Scan(o Options) (*Report,error) {
 target,err:=normalizeTarget(o.Target); if err!=nil{return nil,err}; if o.Workers<1{o.Workers=1}; if o.Workers>32{o.Workers=32}; if o.Timeout<=0{o.Timeout=8*time.Second}; if o.Depth<0{o.Depth=0}; if o.Depth>3{o.Depth=3}
 start:=time.Now(); client:=&http.Client{Timeout:o.Timeout,Transport:&http.Transport{Proxy:http.ProxyFromEnvironment,TLSClientConfig:&tls.Config{MinVersion:tls.VersionTLS12}}}; scope,_:=url.Parse(target)
 seeds:=[]string{target,join(scope,"/robots.txt"),join(scope,"/sitemap.xml"),join(scope,"/.well-known/security.txt")}; obs:=crawl(client,seeds,o.Workers,o.Depth,scope); seen:=map[string]struct{}{}; for _,x:=range obs{for _,u:=range x.Links{seen[u]=struct{}{}}}; urls:=make([]string,0,len(seen)); for u:=range seen{urls=append(urls,u)}; sort.Strings(urls)
 return &Report{Target:target,StartedAt:start,FinishedAt:time.Now(),Observations:obs,DiscoveredURLs:urls,DNS:resolve(scope.Hostname())},nil
}
func normalizeTarget(raw string)(string,error){raw=strings.TrimSpace(raw);if raw==""{return "",errors.New("target is empty")};if !strings.Contains(raw,"://"){raw="https://"+raw};u,e:=url.Parse(raw);if e!=nil||u.Host==""||(u.Scheme!="http"&&u.Scheme!="https"){return "",fmt.Errorf("invalid target %q",raw)};u.Fragment="";return u.String(),nil}
func join(b *url.URL,p string)string{u:=*b;u.Path=p;u.RawQuery="";u.Fragment="";return u.String()}
func crawl(c *http.Client,seeds []string,w,d int,s *url.URL)[]Observation{q:=seeds;seen:=map[string]struct{}{};var all []Observation;for level:=0;level<=d;level++{var batch []string;for _,u:=range q{if _,ok:=seen[u];ok{continue};seen[u]=struct{}{};batch=append(batch,u)};if len(batch)==0{break};obs:=probeMany(c,batch,w,s);all=append(all,obs...);q=nil;for _,x:=range obs{q=append(q,x.Links...)}};return all}
func probeMany(c *http.Client,urls []string,w int,s *url.URL)[]Observation{type job struct{i int;u string};jobs:=make(chan job);res:=make(chan struct{i int;o Observation});var wg sync.WaitGroup;for i:=0;i<w;i++{wg.Add(1);go func(){defer wg.Done();for j:=range jobs{res<-struct{i int;o Observation}{j.i,probe(c,j.u,s)}}}()};go func(){for i,u:=range urls{jobs<-job{i,u}};close(jobs);wg.Wait();close(res)}();out:=make([]Observation,len(urls));for r:=range res{out[r.i]=r.o};return out}
func probe(c *http.Client,target string,s *url.URL)Observation{o:=Observation{URL:target,SecurityHeaders:map[string]string{}};req,e:=http.NewRequestWithContext(context.Background(),http.MethodGet,target,nil);if e!=nil{o.Error=e.Error();return o};req.Header.Set("User-Agent","ReconHawk/0.3 (+authorized-security-research)");r,e:=c.Do(req);if e!=nil{o.Error=e.Error();return o};defer r.Body.Close();body,e:=io.ReadAll(io.LimitReader(r.Body,2<<20));if e!=nil{o.Error=e.Error();return o};o.StatusCode=r.StatusCode;o.ContentType=r.Header.Get("Content-Type");o.Server=r.Header.Get("Server");o.FinalURL=r.Request.URL.String();o.ContentLength=int64(len(body));if r.TLS!=nil{o.TLS=&TLSInfo{Version:tlsName(r.TLS.Version),Cipher:tls.CipherSuiteName(r.TLS.CipherSuite)}};o.SecurityHeaders=headers(r.Header);if strings.Contains(strings.ToLower(o.ContentType),"text/html"){o.Links=links(string(body),s)};return o}
func resolve(h string)[]DNSObservation{a,e:=net.LookupHost(h);o:=DNSObservation{Host:h,Addresses:a};if e!=nil{o.Error=e.Error()};sort.Strings(o.Addresses);return []DNSObservation{o}}
func headers(h http.Header)map[string]string{out:=map[string]string{};for _,k:=range []string{"Content-Security-Policy","Strict-Transport-Security","X-Content-Type-Options","X-Frame-Options","Referrer-Policy","Permissions-Policy"}{if v:=h.Get(k);v!=""{out[k]=v}};return out}
func tlsName(v uint16)string{if v==tls.VersionTLS13{return "TLS 1.3"};if v==tls.VersionTLS12{return "TLS 1.2"};return fmt.Sprintf("0x%x",v)}
func links(body string,s *url.URL)[]string {
 re:=regexp.MustCompile("(?i)(?:href|src)=[\\\"']([^\\\"'#]+)")
 out:=[]string{}
 for _,m:=range re.FindAllStringSubmatch(body,-1) {
  u,e:=url.Parse(strings.TrimSpace(m[1]))
  if e!=nil {continue}
  u=s.ResolveReference(u)
  if (u.Scheme=="http"||u.Scheme=="https")&&strings.EqualFold(u.Hostname(),s.Hostname()) {u.Fragment="";out=append(out,u.String())}
 }
 return unique(out)
}
func unique(in []string)[]string{m:=map[string]struct{}{};out:=[]string{};for _,x:=range in{if _,ok:=m[x];ok{continue};m[x]=struct{}{};out=append(out,x)};sort.Strings(out);return out}
