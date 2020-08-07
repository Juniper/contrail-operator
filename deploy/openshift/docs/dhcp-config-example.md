## Configure DHCP server for Openshift 4 deployment

**Install dhcp package**

```
sudo yum install -y dhcp
```

**Modify dhcpd.conf**

```
$ sudo cat /etc/dhcp/dhcpd.conf
authoritative;
ddns-update-style interim;
allow booting;
allow bootp;
allow unknown-clients;
ignore client-updates;
default-lease-time 14400;
max-lease-time 14400;

subnet 192.168.100.0 netmask 255.255.255.0 {
  range 192.168.100.50 192.168.100.80;
  option routers                  192.168.100.1;
  option broadcast-address        192.168.100.255;
  option subnet-mask              255.255.255.0;
  option domain-name-servers      192.168.100.2;
  option domain-search            "ocp4.testbed.io";
  filename "pxelinux.0";
  next-server 192.168.100.2;
}

# Static entries
host bootstrap { hardware ethernet 52:54:00:60:72:67; fixed-address 192.168.100.50; option host-name "bootstrap.ocp4.testbed.io"; }
host master0 { hardware ethernet 52:54:00:e7:9d:67; fixed-address 192.168.100.51; option host-name "master0.ocp4.testbed.io"; }
host master1 { hardware ethernet 52:54:00:80:16:23; fixed-address 192.168.100.52; option host-name "master1.ocp4.testbed.io"; }
host master2 { hardware ethernet 52:54:00:d5:1c:39; fixed-address 192.168.100.53; option host-name "master2.ocp4.testbed.io"; }
host worker0 { hardware ethernet 52:54:00:f4:26:a1; fixed-address 192.168.100.54; option host-name "worker0.ocp4.testbed.io"; }
host worker1 { hardware ethernet 52:54:00:82:90:00; fixed-address 192.168.100.55; option host-name "worker1.ocp4.testbed.io"; }
host worker2 { hardware ethernet 52:54:00:8e:10:34; fixed-address 192.168.100.56; option host-name "worker2.ocp4.testbed.io"; }
```

**Start and enable dhcpd service**

```
sudo systemctl enable --now dhcpd
```
