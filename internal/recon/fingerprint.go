package recon

import ("net/http";"strings")
type Technology struct { Name string `json:"name"`; Evidence string `json:"evidence"` }
func FingerprintTechnology(resp *http.Response,body []byte)[]Technology{var out []Technology;seen:=map[string]bool{};add:=func(n,e string){if !seen[n]{seen[n]=true;out=append(out,Technology{n,e})}};server:=strings.ToLower(resp.Header.Get("server"));if strings.Contains(server,"nginx"){add("nginx","Server header")};if strings.Contains(server,"apache"){add("Apache","Server header")};if strings.Contains(server,"cloudflare"){add("Cloudflare","Server header")};if v:=resp.Header.Get("x-powered-by");v!=""{add(v,"X-Powered-By header")};b:=strings.ToLower(string(body));for n,m:=range map[string]string{"WordPress":"wp-content","Next.js":"__next_data__","React":"data-reactroot","Vue.js":"__vue__","Laravel":"laravel_session"}{if strings.Contains(b,m){add(n,m+" marker")}};return out}
