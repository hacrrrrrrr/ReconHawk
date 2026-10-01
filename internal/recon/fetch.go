package recon

import (
 "context"
 "fmt"
 "io"
 "net/http"
 "net/url"
)

func FetchBody(ctx context.Context, target string, timeoutSeconds int) (string,*url.URL,error) {
 if timeoutSeconds<=0 {timeoutSeconds=10}
 client:=&http.Client{Timeout:time.Duration(timeoutSeconds)*time.Second}
 req,err:=http.NewRequestWithContext(ctx,http.MethodGet,target,nil);if err!=nil{return "",nil,err}
 req.Header.Set("User-Agent","ReconHawk/0.2 (+authorized-security-research)")
 resp,err:=client.Do(req);if err!=nil{return "",nil,err};defer resp.Body.Close()
 if resp.StatusCode/100!=2{return "",resp.Request.URL,fmt.Errorf("HTTP %s",resp.Status)}
 b,err:=io.ReadAll(io.LimitReader(resp.Body,5<<20));if err!=nil{return "",resp.Request.URL,err}
 return string(b),resp.Request.URL,nil
}
