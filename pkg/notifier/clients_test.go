package notifier

import "testing"

func TestServerChan(t *testing.T) {
	s, err := NewServerChanClient(``)
	if err != nil {
		t.Error(err)
		return
	}
	err = s.SendMsg("test")
	if err != nil {
		t.Error(err)
	}
}
