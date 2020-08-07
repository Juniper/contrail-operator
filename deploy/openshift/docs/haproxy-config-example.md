## Configure HA Proxy for OpenShift 4 deployment

**Install ha proxy package**

```
sudo yum install -y haproxy
```
**Configure firewall rules**

```
sudo firewall-cmd --add-port={80/tcp,443/tcp,6443/tcp,22623/tcp,32700/tcp} --permanent
sudo firewall-cmd --reload
```

**Backup default haproxy config**

```
sudo cp /etc/haproxy/haproxy.cfg /etc/haproxy/haproxy.cfg.bak
```

**Modify haproxy.cfg**

```
$ sudo cat /etc/haproxy/haproxy.cfg
#---------------------------------------------------------------------
# Example configuration for a possible web application.  See the
# full configuration options online.
#
#   http://haproxy.1wt.eu/download/1.4/doc/configuration.txt
#
#---------------------------------------------------------------------

#---------------------------------------------------------------------
# Global settings
#---------------------------------------------------------------------
global
    # to have these messages end up in /var/log/haproxy.log you will
    # need to:
    #
    # 1) configure syslog to accept network log events.  This is done
    #    by adding the '-r' option to the SYSLOGD_OPTIONS in
    #    /etc/sysconfig/syslog
    #
    # 2) configure local2 events to go to the /var/log/haproxy.log
    #   file. A line like the following can be added to
    #   /etc/sysconfig/syslog
    #
    #    local2.*                       /var/log/haproxy.log
    #
    log         127.0.0.1 local2

    chroot      /var/lib/haproxy
    pidfile     /var/run/haproxy.pid
    maxconn     4000
    user        haproxy
    group       haproxy
    daemon

    # turn on stats unix socket
    stats socket /var/lib/haproxy/stats

#---------------------------------------------------------------------
# common defaults that all the 'listen' and 'backend' sections will
# use if not designated in their block
#---------------------------------------------------------------------
defaults
    mode                    http
    log                     global
    option                  httplog
    option                  dontlognull
    option http-server-close
    option forwardfor       except 127.0.0.0/8
    option                  redispatch
    retries                 3
    timeout http-request    10s
    timeout queue           1m
    timeout connect         10s
    timeout client          1m
    timeout server          1m
    timeout http-keep-alive 10s
    timeout check           10s
    maxconn                 3000

#---------------------------------------------------------------------

listen stats
    bind :9000
    mode http
    stats enable
    stats uri /
    monitor-uri /healthz


frontend openshift-api-server
    bind *:6443
    default_backend openshift-api-server
    mode tcp
    option tcplog

backend openshift-api-server
    balance source
    mode tcp
    server bootstrap 192.168.100.50:6443 check
    server master0 192.168.100.51:6443 check
    server master1 192.168.100.52:6443 check
    server master2 192.168.100.53:6443 check

frontend machine-config-server
    bind *:22623
    default_backend machine-config-server
    mode tcp
    option tcplog

backend machine-config-server
    balance source
    mode tcp
    server bootstrap 192.168.100.50:22623 check
    server master0 192.168.100.51:22623 check
    server master1 192.168.100.52:22623 check
    server master2 192.168.100.53:22623 check

frontend ingress-http
    bind *:80
    default_backend ingress-http
    mode tcp
    option tcplog

backend ingress-http
    balance source
    mode tcp
    server worker0-http-router0 192.168.100.54:80 check
    server worker1-http-router1 192.168.100.55:80 check
    server worker2-http-router2 192.168.100.56:80 check

frontend ingress-https
    bind *:443
    default_backend ingress-https
    mode tcp
    option tcplog

backend ingress-https
    balance source
    mode tcp
    server worker0-https-router0 192.168.100.54:443 check
    server worker1-https-router1 192.168.100.55:443 check
    server worker2-https-router2 192.168.100.56:443 check
```

**Set semanage ports for selinux**

```
sudo semanage port  -a 22623 -t http_port_t -p tcp
sudo semanage port  -a 6443 -t http_port_t -p tcp
sudo semanage port  -a 32700 -t http_port_t -p tcp
sudo semanage port  -l  | grep -w http_port_t
```
**Start and enable HAproxy service**

```
sudo systemctl enable --now haproxy
```
