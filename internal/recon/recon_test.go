package recon

import "testing"

func TestNormalizeTarget(t *testing.T) {
 got,err:=normalizeTarget("example.com"); if err!=nil{t.Fatal(err)}
 if got!="https://example.com"{t.Fatalf("got %q",got)}
}
func TestRejectUnsupportedScheme(t *testing.T){if _,err:=normalizeTarget("ftp://example.com");err==nil{t.Fatal("expected error")}}
