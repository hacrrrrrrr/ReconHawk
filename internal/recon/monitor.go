package recon

import (
 "context"
 "time"
)

type MonitorOptions struct { Scan Options; Interval time.Duration; Iterations int }

func Monitor(ctx context.Context,o MonitorOptions,emit func(*Report) error) error {
 if o.Interval<=0{o.Interval=5*time.Minute}
 if o.Iterations<0{o.Iterations=0}
 n:=0
 for {
  r,err:=Scan(o.Scan);if err!=nil{return err};if err:=emit(r);err!=nil{return err}
  n++;if o.Iterations>0&&n>=o.Iterations{return nil}
  timer:=time.NewTimer(o.Interval);select{case <-ctx.Done():timer.Stop();return ctx.Err();case <-timer.C:}
 }
}
