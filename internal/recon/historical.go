package recon

import (
 "context"
 "encoding/json"
 "fmt"
 "net/http"
 "net/url"
 "strings"
)

type HistoricalURL struct { URL string `json:"url"`; Source string `json:"source"` }

type CDXRecord struct { Original string `json:"original"` }

func FetchCommonCrawl(ctx context.Context, target string) ([]HistoricalURL,error) {
 u,err:=url.Parse(target);if err!=nil{return nil,err}
 host:=u.Hostname();if host==""{return nil,fmt.Errorf("invalid host")}
 endpoint:="https://index.commoncrawl.org/CC-MAIN-2026-30-index?url="+url.QueryEscape(host+"/*")+"&output=json&filter=status:200"
 req,err:=http.NewRequestWithContext(ctx,http.MethodGet,endpoint,nil);if err!=nil{return nil,err}
 resp,err:=http.DefaultClient.Do(req);if err!=nil{return nil,err};defer resp.Body.Close()
 var out []HistoricalURL
 dec:=json.NewDecoder(resp.Body)
 for dec.More(){var r CDXRecord;if err:=dec.Decode(&r);err!=nil{break};if strings.HasPrefix(r.Original,"http"){out=append(out,HistoricalURL{URL:r.Original,Source:"commoncrawl"})}}
 return out,nil
}
