package fdyn

import (
	"net"
	"fmt"

	"github.com/miekg/dns"
	"github.com/gomodule/redigo/redis"
)

func rewrite(msg *dns.Msg, f *Fdyn, query string) error {
	var (
		err error
				reply interface{}
		text string
	)
	
	for _, rr := range(msg.Answer) {
		rrtype := rr.Header().Rrtype
		switch rrtype {
		case dns.TypeA:
			if rr.(*dns.A).A.IsUnspecified() {
				if f.Pool == nil {
					return fmt.Errorf("no open Redis pool")
				}
				conn := f.Pool.Get()
				if conn == nil {
					return fmt.Errorf("could not get connection from Redis pool")
				}
				defer conn.Close()
				reply, err = conn.Do("GET", query)
				if err != nil {
					return err
				}
				text, err = redis.String(reply, nil)
				if err != nil {
					return err
				}
				dyn := net.ParseIP(text)
				if dyn == nil {
					return fmt.Errorf("could not parse IP %v", text)
				}
				rr.(*dns.A).A = dyn
			}
			// pass
		case dns.TypeAAAA:
			// pass
		}
	}
	return nil
}
