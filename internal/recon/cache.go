package recon

import (
 "crypto/sha256"
 "database/sql"
 _ "modernc.org/sqlite"
 "encoding/hex"
 "encoding/json"
 "fmt"
 "time"
)

type Cache struct { db *sql.DB }

func OpenCache(path string)(*Cache,error){
 db,err:=sql.Open("sqlite",path);if err!=nil{return nil,err}
 c:=&Cache{db:db};if err:=c.init();err!=nil{db.Close();return nil,err};return c,nil
}
func(c *Cache)init()error{_,err:=c.db.Exec(`CREATE TABLE IF NOT EXISTS results (
 key TEXT PRIMARY KEY, target TEXT NOT NULL, created_at INTEGER NOT NULL, data BLOB NOT NULL
); CREATE INDEX IF NOT EXISTS idx_results_target ON results(target);`);return err}
func(c *Cache)Close()error{return c.db.Close()}
func cacheKey(target string,options Options)string{b,_:=json.Marshal(options);h:=sha256.Sum256(append([]byte(target),b...));return hex.EncodeToString(h[:])}
func(c *Cache)Get(target string,options Options,maxAge time.Duration)(*Report,bool,error){
 k:=cacheKey(target,options);var created int64;var data []byte
 err:=c.db.QueryRow("SELECT created_at,data FROM results WHERE key=?",k).Scan(&created,&data)
 if err==sql.ErrNoRows{return nil,false,nil};if err!=nil{return nil,false,err}
 if maxAge>0&&time.Since(time.Unix(created,0))>maxAge{return nil,false,nil}
 var r Report;if err:=json.Unmarshal(data,&r);err!=nil{return nil,false,err};return &r,true,nil
}
func(c *Cache)Put(target string,options Options,r *Report)error{data,err:=json.Marshal(r);if err!=nil{return err};k:=cacheKey(target,options);_,err=c.db.Exec("INSERT OR REPLACE INTO results(key,target,created_at,data) VALUES(?,?,?,?)",k,target,time.Now().Unix(),data);return err}
func(c *Cache)Clear()error{if c.db==nil{return fmt.Errorf("cache is closed")};_,err:=c.db.Exec("DELETE FROM results");return err}
