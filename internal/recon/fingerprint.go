package recon

import (
 "net/http"
 "strings"
)

type Technology struct { Name string `json:"name"`; Evidence string `json:"evidence"` }

func FingerprintTechnology(resp *http.Response, body []byte) []Technology {
 var out []Technology
 add:=func(n,e string){out=append(out,Technology{Name:n,Evidence:e})}
 if strings.Contains(strings.ToLower(resp.Header.Get("server")),"nginx"){add("nginx","Server header")}
 if strings.Contains(strings.ToLower(resp.Header.Get("server")),"apache"){add("Apache","Server header")}
 if v:=resp.Header.Get("x-powered-by");v!=""{add(v,"X-Powered-By header")}
 b:=strings.ToLower(string(body))
 if strings.Contains(b,"wp-content"){add("WordPress","wp-content marker")}
 if strings.Contains(b,"__next_data__"){add("Next.js","__NEXT_DATA__ marker")}
 if strings.Contains(b,"data-reactroot"){add("React","data-reactroot marker")}
 return out
}
