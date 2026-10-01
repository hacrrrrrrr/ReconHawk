package recon

import (
 "bufio"
 "context"
 "encoding/json"
 "fmt"
 "os/exec"
 "strings"
)

type PluginObservation struct { Name string `json:"name"`; Data map[string]any `json:"data"` }

func RunPlugin(ctx context.Context,path,target string)(PluginObservation,error){
 if strings.TrimSpace(path)==""{return PluginObservation{},fmt.Errorf("plugin path is empty")}
 cmd:=exec.CommandContext(ctx,path,target)
 stdout,err:=cmd.StdoutPipe();if err!=nil{return PluginObservation{},err}
 if err:=cmd.Start();err!=nil{return PluginObservation{},err}
 scanner:=bufio.NewScanner(stdout);scanner.Buffer(make([]byte,4096),2<<20)
 if !scanner.Scan(){_ = cmd.Wait();return PluginObservation{},fmt.Errorf("plugin produced no JSON")}
 line:=scanner.Text();if err:=cmd.Wait();err!=nil{return PluginObservation{},err}
 var data map[string]any;if err:=json.Unmarshal([]byte(line),&data);err!=nil{return PluginObservation{},fmt.Errorf("plugin must emit one JSON object: %w",err)}
 return PluginObservation{Name:path,Data:data},nil
}
