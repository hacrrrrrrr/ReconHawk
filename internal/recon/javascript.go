package recon

import (
 "net/url"
 "regexp"
 "sort"
 "strings"
)

var jsURLRE=regexp.MustCompile(`(?i)(?:src|href)s*=s*["']([^"']+.js(?:?[^"']*)?)["']`)

func ExtractJavaScriptEndpoints(body string, base *url.URL) []string {
 seen:=map[string]struct{}{}
 for _,m:=range jsURLRE.FindAllStringSubmatch(body,1000) {
  u,err:=url.Parse(strings.TrimSpace(m[1]));if err!=nil{continue}
  u=base.ResolveReference(u);if u.Scheme!="http"&&u.Scheme!="https"{continue}
  u.Fragment="";seen[u.String()]=struct{}{}
 }
 out:=make([]string,0,len(seen));for x:=range seen{out=append(out,x)};sort.Strings(out);return out
}
