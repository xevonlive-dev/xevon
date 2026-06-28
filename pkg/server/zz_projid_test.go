package server
import "testing"
func TestProjIDAccept(t *testing.T){
  ok := []string{
    "00000000-0000-0000-defa-c01001000001",
    "proj-0002-aaaa-bbbb-cccc-ddddeeee0002",
    "proj-0003-aaaa-bbbb-cccc-ddddeeee0003",
    "8f3e5751-4f06-473e-8073-b6fdcc81a264",
  }
  bad := []string{"", "  ", "a b", "x';DROP", "\t"}
  for _,s:=range ok { if !isAcceptableProjectID(s){ t.Errorf("should ACCEPT %q",s)} }
  for _,s:=range bad { if isAcceptableProjectID(s){ t.Errorf("should REJECT %q",s)} }
}
