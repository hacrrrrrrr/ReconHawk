package recon

import ("net/url";"regexp";"sort";"strings")
var endpointRE=regexp.MustCompile(`(?i)(?:https?://|/)[^"'\s<>]+`)
func ExtractJavaScriptEndpoints(body string,base *url.URL)[]string{seen:=map[string]struct{}{};for _,raw:=range endpointRE.FindAllString(body,3000){u,e:=url.Parse(strings.TrimSpace(raw));if e!=nil{continue};if base!=nil{u=base.ResolveReference(u)};if u.Scheme!="http"&&u.Scheme!="https"{continue};if base!=nil&&!strings.EqualFold(u.Hostname(),base.Hostname()){continue};u.Fragment="";seen[u.String()]=struct{}{}};out:=make([]string,0,len(seen));for x:=range seen{out=append(out,x)};sort.Strings(out);return out}
