## Configure Web (HTTP) for OpenShift 4 deployment

The webserver is used to store RHCOS images and ignition files for the PXE server.

**Install required packages**

```
sudo yum -y install syslinux httpd wget
```

**Configure firewall rules**

```
sudo firewall-cmd --add-service={http,https} --permanent
sudo firewall-cmd --add-port=8080/tcp --permanent
sudo firewall-cmd --reload
```

**Configure webserver**

Update /etc/httpd/conf/httpd.conf to change port to 8080

Restart httpd service

```
sudo systemctl restart httpd
```

**Download RHCOS images**

```
export OC_VERSION="4.4"
export BUILD_VERSION="4.4.3-x86_64"
VERSION="latest"
cd /var/lib/tftpboot/rhcos
sudo wget https://mirror.openshift.com/pub/openshift-v4/dependencies/rhcos/${OC_VERSION}/${VERSION}/rhcos-${BUILD_VERSION}-installer-kernel-x86_64
sudo wget https://mirror.openshift.com/pub/openshift-v4/dependencies/rhcos/${OC_VERSION}/${VERSION}/rhcos-${BUILD_VERSION}-installer-initramfs.x86_64.img
sudo mv rhcos-${BUILD_VERSION}-installer-kernel-x86_64 kernel
sudo mv rhcos-${BUILD_VERSION}-installer-initramfs.x86_64.img initramfs.img
sudo chmod 555 kernel
sudo chmod 555 initramfs.img
sudo restorecon -vR .

sudo mkdir -p /var/www/html/install
sudo mkdir -p /var/www/html/ignition
sudo chmod 775 -R /var/www/html/install
sudo chmod 775 -R /var/www/html/ignition

cd /var/www/html/install
sudo wget https://mirror.openshift.com/pub/openshift-v4/dependencies/rhcos/${OC_VERSION}/${VERSION}/rhcos-${BUILD_VERSION}-metal.x86_64.raw.gz
sudo mv rhcos-${BUILD_VERSION}-metal.x86_64.raw.gz bios.raw.gz
sudo chmod 555 bios.raw.gz
sudo restorecon -vR /var/www/html
```
