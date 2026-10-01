package recon

import (
 "crypto/sha256"
 "encoding/hex"
 "encoding/json"
 "os"
 "path/filepath"
 "sync"
 "time"
)

type cacheEntry struct { Target string `json:"target"`; CreatedAt int64 `json:"created_at"`; Data json.RawMessage `json:"data"` }
type Cache struct { path string; mu sync.Mutex; entries map[string]cacheEntry }

func OpenCache(path string)(*Cache,error){
 c:=&Cache{path:path,entries:map[string]cacheEntry{}}
 if b,err:=os.ReadFile(path);err==nil {if err:=json.Unmarshal(b,&c.entries);err!=nil{return nil,err}} else if !os.IsNotExist(err){return nil,err}
 return c,nil
}
func(c *Cache)Close()error{return c.flush()}
func(c *Cache)flush()error{c.mu.Lock();defer c.mu.Unlock();if dir:=filepath.Dir(c.path);dir!="."{if err:=os.MkdirAll(dir,0755);err!=nil{return err}};b,err:=json.MarshalIndent(c.entries,"","  ");if err!=nil{return err};tmp:=c.path+".tmp";if err:=os.WriteFile(tmp,b,0600);err!=nil{return err};return os.Rename(tmp,c.path)}
func cacheKey(target string,options Options)string{b,_:=json.Marshal(options);h:=sha256.Sum256(append([]byte(target),b...));return hex.EncodeToString(h[:])}
func(c *Cache)Get(target string,options Options,maxAge time.Duration)(*Report,bool,error){c.mu.Lock();defer c.mu.Unlock();e,ok:=c.entries[cacheKey(target,options)];if !ok{return nil,false,nil};if maxAge>0&&time.Since(time.Unix(e.CreatedAt,0))>maxAge{return nil,false,nil};var r Report;if err:=json.Unmarshal(e.Data,&r);err!=nil{return nil,false,err};return &r,true,nil}
func(c *Cache)Put(target string,options Options,r *Report)error{b,err:=json.Marshal(r);if err!=nil{return err};c.mu.Lock();c.entries[cacheKey(target,options)]=cacheEntry{Target:target,CreatedAt:time.Now().Unix(),Data:b};c.mu.Unlock();return c.flush()}
func(c *Cache)Clear()error{c.mu.Lock();c.entries=map[string]cacheEntry{};c.mu.Unlock();return c.flush()}
