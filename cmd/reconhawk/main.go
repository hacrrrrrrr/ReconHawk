package main

import (
 "context"
 "encoding/json"
 "flag"
 "fmt"
 "os"
 "time"
 "github.com/hacrrrrrrr/ReconHawk/internal/recon"
)

const version = "0.2.0"

func main() {
 if len(os.Args) < 2 { help(); return }
 switch os.Args[1] {
 case "help","--help","-h": help()
 case "version","--version": fmt.Println("ReconHawk",version)
 case "scan": scan(os.Args[2:])
 case "subdomains": subdomains(os.Args[2:])
 case "historical": historical(os.Args[2:])
 case "js": js(os.Args[2:])
 case "fingerprint": fingerprint(os.Args[2:])
 case "monitor": monitor(os.Args[2:])
 case "cache": cacheCmd(os.Args[2:])
 default: fmt.Fprintf(os.Stderr,"unknown command %q\n",os.Args[1]); help(); os.Exit(2)
 }
}

func scan(args []string) {
 fs:=flag.NewFlagSet("scan",flag.ExitOnError)
 cacheFile:=fs.String("cache","reconhawk-cache.json","persistent result cache file")
 fs.Usage=func(){fmt.Println("Usage: reconhawk scan <target> [flags]");fs.PrintDefaults()}
 workers:=fs.Int("workers",8,"maximum concurrent requests")
 timeout:=fs.Duration("timeout",8*time.Second,"request timeout")
 depth:=fs.Int("depth",1,"same-origin crawl depth")
 output:=fs.String("output","","write JSON report to file")
 jsonl:=fs.String("jsonl","","write observations as JSONL")
 fs.Parse(args)
 if fs.NArg()!=1 {fs.Usage();os.Exit(2)}
 opts:=recon.Options{Target:fs.Arg(0),Workers:*workers,Timeout:*timeout,Depth:*depth}
 c,e:=recon.OpenCache(*cacheFile);if e!=nil{fatal(e)};defer c.Close()
 r,hit,e:=c.Get(fs.Arg(0),opts,30*time.Minute);if e!=nil{fatal(e)}
 if !hit {r,e=recon.Scan(opts);if e!=nil{fatal(e)};if e:=c.Put(fs.Arg(0),opts,r);e!=nil{fatal(e)}}
 writeJSON(r,*output)
 if *jsonl!="" {f,e:=os.Create(*jsonl);if e!=nil{fatal(e)};defer f.Close();enc:=json.NewEncoder(f);for _,o:=range r.Observations{if e:=enc.Encode(o);e!=nil{fatal(e)}}}
}

func subdomains(args []string) {
 fs:=flag.NewFlagSet("subdomains",flag.ExitOnError);fs.Usage=func(){fmt.Println("Usage: reconhawk subdomains <domain>")}
 fs.Parse(args);if fs.NArg()!=1{fs.Usage();os.Exit(2)}
 ctx,cancel:=context.WithTimeout(context.Background(),20*time.Second);defer cancel()
 r,e:=recon.CollectSubdomains(ctx,fs.Arg(0),[]recon.Provider{recon.CRTShProvider{}});if e!=nil{fatal(e)};writeJSON(r,"")
}

func historical(args []string) {
 fs:=flag.NewFlagSet("historical",flag.ExitOnError);fs.Usage=func(){fmt.Println("Usage: reconhawk historical <domain-or-url>")}
 fs.Parse(args);if fs.NArg()!=1{fs.Usage();os.Exit(2)}
 ctx,cancel:=context.WithTimeout(context.Background(),30*time.Second);defer cancel()
 r,e:=recon.FetchCommonCrawl(ctx,fs.Arg(0));if e!=nil{fatal(e)};writeJSON(r,"")
}

func js(args []string) {
 fs:=flag.NewFlagSet("js",flag.ExitOnError);fs.Usage=func(){fmt.Println("Usage: reconhawk js <url>")}
 fs.Parse(args);if fs.NArg()!=1{fs.Usage();os.Exit(2)}
 ctx,cancel:=context.WithTimeout(context.Background(),15*time.Second);defer cancel()
 body,base,e:=recon.FetchBody(ctx,fs.Arg(0),15);if e!=nil{fatal(e)}
 writeJSON(recon.ExtractJavaScriptEndpoints(body,base),"")
}

func fingerprint(args []string) {
 fs:=flag.NewFlagSet("fingerprint",flag.ExitOnError);fs.Usage=func(){fmt.Println("Usage: reconhawk fingerprint <url>")}
 fs.Parse(args);if fs.NArg()!=1{fs.Usage();os.Exit(2)}
 r,e:=recon.Scan(recon.Options{Target:fs.Arg(0),Workers:1,Timeout:10*time.Second,Depth:0});if e!=nil{fatal(e)};writeJSON(r.Observations,"")
}


func monitor(args []string) {
 fs:=flag.NewFlagSet("monitor",flag.ExitOnError)
 fs.Usage=func(){fmt.Println("Usage: reconhawk monitor <target> [flags]");fs.PrintDefaults()}
 interval:=fs.Duration("interval",5*time.Minute,"time between scans")
 iterations:=fs.Int("iterations",0,"number of scans; 0 means until interrupted")
 workers:=fs.Int("workers",4,"maximum concurrent requests")
 timeout:=fs.Duration("timeout",8*time.Second,"request timeout")
 depth:=fs.Int("depth",1,"same-origin crawl depth")
 fs.Parse(args)
 if fs.NArg()!=1{fs.Usage();os.Exit(2)}
 ctx:=context.Background()
 err:=recon.Monitor(ctx,recon.MonitorOptions{Scan:recon.Options{Target:fs.Arg(0),Workers:*workers,Timeout:*timeout,Depth:*depth},Interval:*interval,Iterations:*iterations},func(r *recon.Report)error{return writeJSONValue(r)})
 if err!=nil&&err!=context.Canceled{fatal(err)}
}

func writeJSONValue(v any) error {
 b,e:=json.Marshal(v);if e!=nil{return e}
 _,e=fmt.Println(string(b));return e
}


func cacheCmd(args []string) {
 fs:=flag.NewFlagSet("cache",flag.ExitOnError)
 fs.Usage=func(){fmt.Println("Usage: reconhawk cache <clear> [--file path]");fs.PrintDefaults()}
 file:=fs.String("file","reconhawk-cache.json","persistent cache file")
 fs.Parse(args)
 if fs.NArg()!=1 || fs.Arg(0)!="clear" {fs.Usage();os.Exit(2)}
 c,e:=recon.OpenCache(*file);if e!=nil{fatal(e)};defer c.Close()
 if e:=c.Clear();e!=nil{fatal(e)}
 fmt.Println("cache cleared:",*file)
}

func writeJSON(v any,path string){b,e:=json.MarshalIndent(v,"","  ");if e!=nil{fatal(e)};b=append(b,'\n');if path==""{fmt.Print(string(b));return};if e=os.WriteFile(path,b,0644);e!=nil{fatal(e)}}
func fatal(e error){fmt.Fprintln(os.Stderr,"error:",e);os.Exit(1)}

func help(){
 fmt.Printf("ReconHawk %s - automated authorized reconnaissance\n\n",version)
 fmt.Println("Usage: reconhawk <command> [arguments]")
 fmt.Println("\nCommands:")
 fmt.Println("  scan          HTTP/DNS discovery, crawling and metadata")
 fmt.Println("  subdomains    passive certificate-transparency discovery")
 fmt.Println("  historical    historical URLs from Common Crawl")
 fmt.Println("  js             JavaScript endpoint extraction")
 fmt.Println("  fingerprint   technology identification
  monitor       repeatedly scan an authorized target
  cache         manage persistent local results")
 fmt.Println("  version        print version")
 fmt.Println("  help           show help")
 fmt.Println("\nExamples:")
 fmt.Println("  reconhawk scan example.com --depth 2 --workers 8 --output report.json")
 fmt.Println("  reconhawk subdomains example.com")
 fmt.Println("  reconhawk historical example.com")
 fmt.Println("  reconhawk js https://example.com")
 fmt.Println("  reconhawk fingerprint https://example.com")
 fmt.Println("\nUse only against assets you own or are explicitly authorized to test.")
}
