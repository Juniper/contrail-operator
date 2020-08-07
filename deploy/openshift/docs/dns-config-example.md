## Configure DNS server for Openshift 4 deployment

**Install bind server packages**

```
sudo yum -y install bind bind-utils
```

**Configure firewall rules**

```
sudo firewall-cmd --add-services=dns --permanent
sudo firewall-cmd --reload
```

**Modify named.conf**

```
$ sudo cat /etc/named.conf
//
// named.conf
//
// Provided by Red Hat bind package to configure the ISC BIND named(8) DNS
// server as a caching only nameserver (as a localhost DNS resolver only).
//
// See /usr/share/doc/bind*/sample/ for example named configuration files.
//

options {
	listen-on port 53 { any; };
	listen-on-v6 port 53 { ::1; };
	directory 	"/var/named";
	dump-file 	"/var/named/data/cache_dump.db";
	statistics-file "/var/named/data/named_stats.txt";
	memstatistics-file "/var/named/data/named_mem_stats.txt";
	allow-query     { any; };

	/*
	 - If you are building an AUTHORITATIVE DNS server, do NOT enable recursion.
	 - If you are building a RECURSIVE (caching) DNS server, you need to enable
	   recursion.
	 - If your recursive DNS server has a public IP address, you MUST enable access
	   control to limit queries to your legitimate users. Failing to do so will
	   cause your server to become part of large scale DNS amplification
	   attacks. Implementing BCP38 within your network would greatly
	   reduce such attack surface
	*/
	recursion yes;

	/* Fowarders */
	forward only;
	forwarders { 8.8.8.8; 8.8.4.4; };

	dnssec-enable yes;
	dnssec-validation no;

	managed-keys-directory "/var/named/dynamic";

	pid-file "/run/named/named.pid";
	session-keyfile "/run/named/session.key";

	/* https://fedoraproject.org/wiki/Changes/CryptoPolicy */
	/* include "/etc/crypto-policies/back-ends/bind.config"; */
};

logging {
        channel default_debug {
                file "data/named.run";
                severity dynamic;
        };
};

zone "." IN {
	type hint;
	file "named.ca";
};

########### Add what's between these comments ###########
zone "ocp4.testbed.io" IN {
	type	master;
	file	"zonefile.db";
};

zone "100.168.192.in-addr.arpa" IN {
	type	master;
	file	"reverse.db";
};
########################################################

include "/etc/named.rfc1912.zones";
include "/etc/named.root.key";
```

**Create zone file**

```
$ sudo cat /var/named/zonefile.db
$TTL 1W
@	IN	SOA	ns1.ocp4.testbed.io.	root (
			2020052600	; serial
			3H		; refresh (3 hours)
			30M		; retry (30 minutes)
			2W		; expiry (2 weeks)
			1W )		; minimum (1 week)
	IN	NS	ns1.ocp4.testbed.io.
	IN	MX 10	smtp.ocp4.testbed.io.
;
;
ns1	IN	A	192.168.100.2
smtp	IN	A	192.168.100.2
;
helper	IN	A	192.168.100.2
;
; The api points to the IP of your load balancer
api		IN	A	192.168.100.2
api-int		IN	A	192.168.100.2
;
; The wildcard also points to the load balancer
*.apps		IN	A	192.168.100.2
;
; Create entry for the bootstrap host
bootstrap	IN	A	192.168.100.50
;
; Create entries for the master hosts
master0		IN	A	192.168.100.51
master1		IN	A	192.168.100.52
master2		IN	A	192.168.100.53
;
; Create entries for the worker hosts
worker0		IN	A	192.168.100.54
worker1		IN	A	192.168.100.55
worker2		IN	A	192.168.100.56
;
; The ETCd cluster lives on the masters...so point these to the IP of the masters
etcd-0	IN	A	192.168.100.51
etcd-1	IN	A	192.168.100.52
etcd-2	IN	A	192.168.100.53
;
; The SRV records are IMPORTANT....make sure you get these right...note the trailing dot at the end...
_etcd-server-ssl._tcp	IN	SRV	0 10 2380 etcd-0.ocp4.testbed.io.
_etcd-server-ssl._tcp	IN	SRV	0 10 2380 etcd-1.ocp4.testbed.io.
_etcd-server-ssl._tcp	IN	SRV	0 10 2380 etcd-2.ocp4.testbed.io.
;
;EOF
```

**Create reverse zone file**

```
$ sudo cat /var/named/reverse.db
$TTL 1W
@	IN	SOA	ns1.ocp4.testbed.io.	root (
			2020052600	; serial
			3H		; refresh (3 hours)
			30M		; retry (30 minutes)
			2W		; expiry (2 weeks)
			1W )		; minimum (1 week)
	IN	NS	ns1.ocp4.testbed.io.
;
; syntax is "last octet" and the host must have fqdn with trailing dot
1       IN      PTR     helper.ocp4.testbed.io.

51	IN	PTR	master0.ocp4.testbed.io.
52	IN	PTR	master1.ocp4.testbed.io.
53	IN	PTR	master2.ocp4.testbed.io.
;
50	IN	PTR	bootstrap.ocp4.testbed.io.
;
2	IN	PTR	api.ocp4.testbed.io.
2	IN	PTR	api-int.ocp4.testbed.io.
;
54	IN	PTR	worker0.ocp4.testbed.io.
55	IN	PTR	worker1.ocp4.testbed.io.
56	IN	PTR	worker2.ocp4.testbed.io.
;
;EOF
```

**Test the bind server configuration**

```
named-checkconf /etc/named.conf
```

**Start and enable the bind service**

```
sudo systemctl enable --now named
```

**Test dns resolution**

```
dig @localhost etcd-0.ocp4.testbed.io
```

**Test reverse pointer**

```
dig @localhost -t srv _etcd-server-ssl._tcp.ocp4.testbed.io
```
