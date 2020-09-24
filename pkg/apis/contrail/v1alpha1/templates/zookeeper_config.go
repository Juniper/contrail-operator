package templates

import (
	"fmt"
	"strconv"

	core "k8s.io/api/core/v1"
)

// ZookeeperStaticConfig is the template of the Zookeeper service configuration.
var ZookeeperStaticConfig = `
dataDir=/var/lib/zookeeper
tickTime=2000
initLimit=5
syncLimit=2
maxClientCnxns=60
maxSessionTimeout=120000
admin.enableServer=false
standaloneEnabled=false
4lw.commands.whitelist=stat,ruok,conf,isro
reconfigEnabled=true
skipACL=yes
dynamicConfigFile=/var/lib/zookeeper/zoo.cfg.dynamic
`

func DynamicZookeeperConfig(pods []core.Pod, electionPort, serverPort, clientPort string) (map[string]string, error) {
	dynamicConf := make(map[string]string, 0)
	var firstServerDef string
	for _, pod := range pods {
		myidString := pod.Name[len(pod.Name)-1:]
		myidInt, err := strconv.Atoi(myidString)
		if err != nil {
			return nil, err
		}
		serverDef := firstServerDef + fmt.Sprintf("server.%d=%s:%s:participant;%s:%s\n",
			myidInt+1, pod.Status.PodIP,
			electionPort+":"+serverPort, pod.Status.PodIP, clientPort)
		if myidInt == 0 {
			firstServerDef = serverDef
		}
		dynamicConf["myid."+pod.Status.PodIP] = strconv.Itoa(myidInt + 1)
		dynamicConf["zoo.cfg.dynamic."+pod.Status.PodIP] = serverDef
	}
	return dynamicConf, nil
}

// ZookeeperLogConfig is the template of the Zookeeper Log configuration.
var ZookeeperLogConfig = `zookeeper.root.logger=INFO, CONSOLE
zookeeper.console.threshold=INFO
zookeeper.log.dir=.
zookeeper.log.file=zookeeper.log
zookeeper.log.threshold=INFO
zookeeper.log.maxfilesize=256MB
zookeeper.log.maxbackupindex=20
zookeeper.tracelog.dir=${zookeeper.log.dir}
zookeeper.tracelog.file=zookeeper_trace.log
log4j.rootLogger=${zookeeper.root.logger}
log4j.appender.CONSOLE=org.apache.log4j.ConsoleAppender
log4j.appender.CONSOLE.Threshold=${zookeeper.console.threshold}
log4j.appender.CONSOLE.layout=org.apache.log4j.PatternLayout
log4j.appender.CONSOLE.layout.ConversionPattern=%d{ISO8601} [myid:%X{myid}] - %-5p [%t:%C{1}@%L] - %m%n
log4j.appender.ROLLINGFILE=org.apache.log4j.RollingFileAppender
log4j.appender.ROLLINGFILE.Threshold=${zookeeper.log.threshold}
log4j.appender.ROLLINGFILE.File=${zookeeper.log.dir}/${zookeeper.log.file}
log4j.appender.ROLLINGFILE.MaxFileSize=${zookeeper.log.maxfilesize}
log4j.appender.ROLLINGFILE.MaxBackupIndex=${zookeeper.log.maxbackupindex}
log4j.appender.ROLLINGFILE.layout=org.apache.log4j.PatternLayout
log4j.appender.ROLLINGFILE.layout.ConversionPattern=%d{ISO8601} [myid:%X{myid}] - %-5p [%t:%C{1}@%L] - %m%n
log4j.appender.TRACEFILE=org.apache.log4j.FileAppender
log4j.appender.TRACEFILE.Threshold=TRACE
log4j.appender.TRACEFILE.File=${zookeeper.tracelog.dir}/${zookeeper.tracelog.file}
log4j.appender.TRACEFILE.layout=org.apache.log4j.PatternLayout
log4j.appender.TRACEFILE.layout.ConversionPattern=%d{ISO8601} [myid:%X{myid}] - %-5p [%t:%C{1}@%L][%x] - %m%n
`

// ZookeeperXslConfig is the template of the Zookeeper XSL configuration.
var ZookeeperXslConfig = `<?xml version="1.0"?>
<xsl:stylesheet xmlns:xsl="http://www.w3.org/1999/XSL/Transform" version="1.0">
<xsl:output method="html"/>
<xsl:template match="configuration">
<html>
<body>
<table border="1">
<tr>
<td>name</td>
<td>value</td>
<td>description</td>
</tr>
<xsl:for-each select="property">
<tr>
<td><a name="{name}"><xsl:value-of select="name"/></a></td>
<td><xsl:value-of select="value"/></td>
<td><xsl:value-of select="description"/></td>
</tr>
</xsl:for-each>
</table>
</body>
</html>
</xsl:template>
</xsl:stylesheet>
`
