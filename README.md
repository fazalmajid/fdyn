# fdyn

## Name

*fdyn* - simple Dynamic DNS plugin for CoreDNS, using a mix of forward and
 Redis.

## Description

The *fdyn* plugin is a forked version of the standard CoreDNS [forward
plugin](https://coredns.io/plugins/forward/), with a difference: if the
returned IP from the upstream server is 0.0.0.0 (IPv4) or ::/0 (IPv6), fdyn
will look up the name in the local Redis (using `GET` of the FQDN, including
trailing period) and substitutes that instead. Thus you can set up all the
paraphernalia of SOA, NS and authority records in a proper DNS server, and
just override the dynamic IP address that is usually all you want to do.
