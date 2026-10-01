package recon

import (
 "context"
 "fmt"
 "net/http"
 "strings"
 "time"
)

type Severity string
const (
 SeverityInfo Severity = "info"
 SeverityLow Severity = "low"
 SeverityMedium Severity = "medium"
)

type Finding struct {
 ID string `json:"id"`
 Title string `json:"title"`
 Severity Severity `json:"severity"`
 URL string `json:"url"`
 Evidence string `json:"evidence"`
 Remediation string `json:"remediation"`
}

type SecurityReport struct {
 Target string `json:"target"`
 StartedAt time.Time `json:"started_at"`
 FinishedAt time.Time `json:"finished_at"`
 Findings []Finding `json:"findings"`
}

func CheckTarget(ctx context.Context,target string,timeout time.Duration) (*SecurityReport,error) {
 if timeout<=0 {timeout=8*time.Second}
 start:=time.Now()
 target,err=normalizeTarget(target);if err!=nil{return nil,err}
 client:=&http.Client{Timeout:timeout}
 req,err:=http.NewRequestWithContext(ctx,http.MethodGet,target,nil);if err!=nil{return nil,err}
 req.Header.Set("User-Agent","ReconHawk/0.3 (+authorized-security-research)")
 resp,err:=client.Do(req);if err!=nil{return nil,err};defer resp.Body.Close()
 findings:=[]Finding{}
 add:=func(id,title string,sev Severity,evidence,fix string){findings=append(findings,Finding{ID:id,Title:title,Severity:sev,URL:resp.Request.URL.String(),Evidence:evidence,Remediation:fix})}
 if resp.Request.URL.Scheme=="https" && resp.Header.Get("Strict-Transport-Security")=="" {add("RH-SEC-001","HSTS header missing",Severity.Low,"Strict-Transport-Security header was not present","Enable HSTS after confirming the site is ready for HTTPS-only operation.")}
 if resp.Header.Get("Content-Security-Policy")=="" {add("RH-SEC-002","Content-Security-Policy header missing",Severity.Low,"Content-Security-Policy header was not present","Deploy a suitable CSP for the application.")}
 if resp.Header.Get("X-Content-Type-Options")=="" {add("RH-SEC-003","X-Content-Type-Options missing",Severity.Low,"X-Content-Type-Options header was not present","Set X-Content-Type-Options: nosniff.")}
 if resp.Header.Get("Access-Control-Allow-Origin")=="*" {add("RH-SEC-004","Wildcard CORS policy",Severity.Medium,"Access-Control-Allow-Origin: *","Restrict cross-origin access to explicitly trusted origins, especially for credentialed APIs.")}
 if server:=resp.Header.Get("Server");server!="" {add("RH-INFO-001","Server header disclosed",Severity.Info,fmt.Sprintf("Server: %s",server),"Remove unnecessary version/product disclosure from public response headers.")}
 for _,cookie:=range resp.Header.Values("Set-Cookie") {
  low:=strings.ToLower(cookie)
  if !strings.Contains(low,"secure") {add("RH-SEC-005","Cookie missing Secure attribute",Severity.Low,"A Set-Cookie response did not include Secure","Set Secure on cookies that should only travel over HTTPS.")}
  if !strings.Contains(low,"httponly") {add("RH-SEC-006","Cookie missing HttpOnly attribute",Severity.Low,"A Set-Cookie response did not include HttpOnly","Set HttpOnly on cookies that do not need JavaScript access.")}
 }
 return &SecurityReport{Target:target,StartedAt:start,FinishedAt:time.Now(),Findings:findings},nil
}
