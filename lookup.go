// Package forward implements a forwarding proxy. It caches an upstream net.Conn for some time, so if the same
// client returns the upstream's Conn will be precached. Depending on how you benchmark this looks to be
// 50% faster than just opening a new connection for every client. It works with UDP and TCP and uses
// inband healthchecking.
package fdyn

import (
	"context"

	"github.com/coredns/coredns/plugin/pkg/transport"
	"github.com/coredns/coredns/request"

	"github.com/miekg/dns"
)

// Fdyn forward the request in state as-is. Unlike Lookup that adds EDNS0 suffix to the message.
// Fdyn may be called with a nil f, an error is returned in that case.
func (f *Fdyn) Fdyn(state request.Request) (*dns.Msg, error) {
	if f == nil {
		return nil, ErrNoFdyn
	}

	fails := 0
	var upstreamErr error
	for _, proxy := range f.List() {
		if proxy.Down(f.maxfails) {
			fails++
			if fails < len(f.proxies) {
				continue
			}
			// All upstream proxies are dead, assume healtcheck is complete broken and randomly
			// select an upstream to connect to.
			proxy = f.List()[0]
		}

		ret, err := proxy.Connect(context.Background(), state, f.opts)

		upstreamErr = err

		if err != nil {
			if fails < len(f.proxies) {
				continue
			}
			break
		}

		// Check if the reply is correct; if not return FormErr.
		if !state.Match(ret) {
			return state.ErrorMessage(dns.RcodeFormatError), nil
		}

		ret = state.Scrub(ret)
		return ret, err
	}

	if upstreamErr != nil {
		return nil, upstreamErr
	}

	return nil, ErrNoHealthy
}

// Lookup will use name and type to forge a new message and will send that upstream. It will
// set any EDNS0 options correctly so that downstream will be able to process the reply.
// Lookup may be called with a nil f, an error is returned in that case.
func (f *Fdyn) Lookup(state request.Request, name string, typ uint16) (*dns.Msg, error) {
	if f == nil {
		return nil, ErrNoFdyn
	}

	req := new(dns.Msg)
	req.SetQuestion(name, typ)
	state.SizeAndDo(req)

	state2 := request.Request{W: state.W, Req: req}

	return f.Fdyn(state2)
}

// NewLookup returns a Fdyn that can be used for plugin that need an upstream to resolve external names.
// Note that the caller MUST run Close on the forward to stop the health checking goroutines.
func NewLookup(addr []string) *Fdyn {
	f := New()
	for i := range addr {
		p := NewProxy(addr[i], transport.DNS)
		f.SetProxy(p)
	}
	return f
}
