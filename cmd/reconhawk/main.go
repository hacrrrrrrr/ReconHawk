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
 default: fmt.Fprintf(os.Stderr,"unknown command %q\n",os.Args[1]); help(); os.Exit(2)
 }
}

func scan(args []string) {
 fs:=flag.NewFlagSet("scan",flag.ExitOnError)
 fs.Usage=func(){fmt.Println("Usage: reconhawk scan <target> [flags]");fs.PrintDefaults()}
 workers:=fs.Int("workers",8,"maximum concurrent requests")
 timeout:=fs.Duration("timeout",8*time.Second,"request timeout")
 depth:=fs.Int("depth",1,"same-origin crawl depth")
 output:=fs.String("output","","write JSON report to file")
 jsonl:=fs.String("jsonl","","write observations as JSONL")
 fs.Parse(args)
 if fs.NArg()!=1 {fs.Usage();os.Exit(2)}
 r,e:=recon.Scan(recon.Options{Target:fs.Arg(0),Workers:*workers,Timeout:*timeout,Depth:*depth});if e!=nil{fatal(e)}
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
 fmt.Println("  fingerprint   technology identification")
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
