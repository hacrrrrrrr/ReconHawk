package recon

import ("bufio";"context";"encoding/json";"fmt";"net/http";"net/url";"strings")
type HistoricalURL struct { URL string `json:"url"`; Source string `json:"source"` }

func FetchCommonCrawl(ctx context.Context,target string)([]HistoricalURL,error){
 u,err:=url.Parse(target);if err!=nil{return nil,err};host:=u.Hostname();if host==""{return nil,fmt.Errorf("invalid host")}
 req,err:=http.NewRequestWithContext(ctx,http.MethodGet,"https://index.commoncrawl.org/CC-MAIN-2026-30-index?url="+url.QueryEscape(host+"/*")+"&output=json&filter=status:200",nil);if err!=nil{return nil,err}
 resp,err:=http.DefaultClient.Do(req);if err!=nil{return nil,err};defer resp.Body.Close();if resp.StatusCode/100!=2{return nil,fmt.Errorf("Common Crawl returned %s",resp.Status)}
 seen:=map[string]struct{}{};out:=[]HistoricalURL{};s:=bufio.NewScanner(resp.Body);s.Buffer(make([]byte,65536),2097152)
 for s.Scan(){var r struct{Original string `json:"url"`};if json.Unmarshal(s.Bytes(),&r)==nil&&strings.HasPrefix(r.Original,"http"){if _,ok:=seen[r.Original];!ok{seen[r.Original]=struct{}{};out=append(out,HistoricalURL{r.Original,"commoncrawl"})}}};return out,s.Err()
}
