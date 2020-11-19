package main

import (
	"encoding/json"
	"errors"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"net"
	"os"
	"strconv"
	"strings"
	"syscall"
	"time"

	"gopkg.in/yaml.v2"

	contrail "github.com/Juniper/contrail-go-api"

	contrailTypes "github.com/Juniper/contrail-operator/contrail-provisioner/contrail-go-types"
	"github.com/Juniper/contrail-operator/contrail-provisioner/types"
	"github.com/Juniper/contrail-operator/contrail-provisioner/vrouternodes"
)

// APIServer struct contains API Server configuration
type APIServer struct {
	APIPort       string     `yaml:"apiPort,omitempty"`
	APIServerList []string   `yaml:"apiServerList,omitempty"`
	Encryption    encryption `yaml:"encryption,omitempty"`
}

type encryption struct {
	CA       string `yaml:"ca,omitempty"`
	Cert     string `yaml:"cert,omitempty"`
	Key      string `yaml:"key,omitempty"`
	Insecure bool   `yaml:"insecure,omitempty"`
}

type KeystoneAuthParameters struct {
	AdminUsername string     `yaml:"admin_user,omitempty"`
	AdminPassword string     `yaml:"admin_password,omitempty"`
	AuthUrl       string     `yaml:"auth_url,omitempty"`
	TenantName    string     `yaml:"tenant_name,omitempty"`
	Encryption    encryption `yaml:"encryption,omitempty"`
}

type EcmpHashingIncludeFields struct {
	HashingConfigured bool `json:"hashingConfigured,omitempty"`
	SourceIp          bool `json:"sourceIp,omitempty"`
	DestinationIp     bool `json:"destinationIp,omitempty"`
	IpProtocol        bool `json:"ipProtocol,omitempty"`
	SourcePort        bool `json:"sourcePort,omitempty"`
	DestinationPort   bool `json:"destinationPort,omitempty"`
}

type GlobalVrouterConfiguration struct {
	EcmpHashingIncludeFields   EcmpHashingIncludeFields `json:"ecmpHashingIncludeFields,omitempty"`
	EncapsulationPriorities    string                   `json:"encapPriority,omitempty"`
	VxlanNetworkIdentifierMode string                   `json:"vxlanNetworkIdentifierMode,omitempty"`
}

func nodeManager(nodesPtr *string, nodeType string, contrailClient *contrail.Client) {
	fmt.Printf("%s %s updated\n", nodeType, *nodesPtr)
	nodeYaml, err := ioutil.ReadFile(*nodesPtr)
	if err != nil {
		panic(err)
	}
	switch nodeType {
	case "control":
		var nodeList []*types.ControlNode
		err = yaml.Unmarshal(nodeYaml, &nodeList)
		if err != nil {
			panic(err)
		}
		if err = controlNodes(contrailClient, nodeList); err != nil {
			panic(err)
		}
	case "analytics":
		var nodeList []*types.AnalyticsNode
		err = yaml.Unmarshal(nodeYaml, &nodeList)
		if err != nil {
			panic(err)
		}
		if err = analyticsNodes(contrailClient, nodeList); err != nil {
			panic(err)
		}
	case "config":
		var nodeList []*types.ConfigNode
		err = yaml.Unmarshal(nodeYaml, &nodeList)
		if err != nil {
			panic(err)
		}
		if err = configNodes(contrailClient, nodeList); err != nil {
			panic(err)
		}
	case "vrouter":
		var nodeList []*types.VrouterNode
		err = yaml.Unmarshal(nodeYaml, &nodeList)
		if err != nil {
			panic(err)
		}
		if err = vrouternodes.ReconcileVrouterNodes(contrailClient, nodeList); err != nil {
			panic(err)
		}
	case "database":
		var nodeList []*types.DatabaseNode
		err = yaml.Unmarshal(nodeYaml, &nodeList)
		if err != nil {
			panic(err)
		}
		if err = databaseNodes(contrailClient, nodeList); err != nil {
			panic(err)
		}
	}
}

func check(err error) {
	if err != nil {
		log.Fatalf("error: %v", err)
	}
}

func main() {

	controlNodesPtr := flag.String("controlNodes", "/provision.yaml", "path to control nodes yaml file")
	configNodesPtr := flag.String("configNodes", "/provision.yaml", "path to config nodes yaml file")
	analyticsNodesPtr := flag.String("analyticsNodes", "/provision.yaml", "path to analytics nodes yaml file")
	vrouterNodesPtr := flag.String("vrouterNodes", "/provision.yaml", "path to vrouter nodes yaml file")
	databaseNodesPtr := flag.String("databaseNodes", "/provision.yaml", "path to database nodes yaml file")
	apiserverPtr := flag.String("apiserver", "/provision.yaml", "path to apiserver yaml file")
	keystoneAuthConfPtr := flag.String("keystoneAuthConf", "/provision.yaml", "path to keystone authentication configuration file")
	globalVrouterConfPtr := flag.String("globalVrouterConf", "/provision.yaml", "path to global vrouter configuration file")
	modePtr := flag.String("mode", "watch", "watch/run")
	flag.Parse()

	if *modePtr == "watch" {

		var apiServer APIServer
		apiServerYaml, err := ioutil.ReadFile(*apiserverPtr)
		if err != nil {
			panic(err)
		}
		err = yaml.Unmarshal(apiServerYaml, &apiServer)
		if err != nil {
			panic(err)
		}

		var keystoneAuthParameters *KeystoneAuthParameters = &KeystoneAuthParameters{}
		if _, err := os.Stat(*keystoneAuthConfPtr); err == nil {
			keystoneAuthParameters = getKeystoneAuthParametersFromFile(*keystoneAuthConfPtr)
		}

		var contrailClient *contrail.Client
		err = retry(5, 10*time.Second, func() (err error) {
			contrailClient, err = getAPIClient(&apiServer, keystoneAuthParameters)
			return

		})
		if err != nil {
			if !connectionError(err) {
				panic(err)
			}
		}

		globalVrouterConfiguration := &GlobalVrouterConfiguration{}
		if _, err := os.Stat(*globalVrouterConfPtr); err == nil {
			globalVrouterConfiguration = getGlobalVrouterConfigFromFile(*globalVrouterConfPtr)
		}
		globalVrouterConfFQName := []string{"default-global-system-config", "default-global-vrouter-config"}
		encapPriority := strings.Split(globalVrouterConfiguration.EncapsulationPriorities, ",")
		encapPriorityObj := &contrailTypes.EncapsulationPrioritiesType{Encapsulation: encapPriority}
		ecmpObj := globalVrouterConfiguration.EcmpHashingIncludeFields
		ecmpHashingIncludeFieldsObj := &contrailTypes.EcmpHashingIncludeFields{ecmpObj.HashingConfigured, ecmpObj.SourceIp, ecmpObj.DestinationIp, ecmpObj.IpProtocol, ecmpObj.SourcePort, ecmpObj.DestinationPort}
		GlobalVrouterConfig := &contrailTypes.GlobalVrouterConfig{}
		GlobalVrouterConfig.SetFQName("", globalVrouterConfFQName)
		GlobalVrouterConfig.SetEncapsulationPriorities(encapPriorityObj)
		GlobalVrouterConfig.SetEcmpHashingIncludeFields(ecmpHashingIncludeFieldsObj)
		GlobalVrouterConfig.SetVxlanNetworkIdentifierMode(globalVrouterConfiguration.VxlanNetworkIdentifierMode)
		if err = contrailClient.Create(GlobalVrouterConfig); err != nil {
			if !strings.Contains(err.Error(), "409 Conflict") {
				panic(err)
			}
			obj, err := contrailClient.FindByName("global-vrouter-config", strings.Join(globalVrouterConfFQName, ":"))
			if err != nil {
				panic(err)
			}
			obj.(*contrailTypes.GlobalVrouterConfig).SetEncapsulationPriorities(encapPriorityObj)
			obj.(*contrailTypes.GlobalVrouterConfig).SetEcmpHashingIncludeFields(ecmpHashingIncludeFieldsObj)
			obj.(*contrailTypes.GlobalVrouterConfig).SetVxlanNetworkIdentifierMode(globalVrouterConfiguration.VxlanNetworkIdentifierMode)
			if err = contrailClient.Update(obj); err != nil {
				panic(err)
			}
		}

		fmt.Println("start watcher")
		done := make(chan bool)

		if controlNodesPtr != nil {
			fmt.Println("initial control node run")
			_, err := os.Stat(*controlNodesPtr)
			if !os.IsNotExist(err) {
				nodeManager(controlNodesPtr, "control", contrailClient)
			} else if os.IsNotExist(err) {
				controlNodes(contrailClient, []*types.ControlNode{})
			}
			fmt.Println("setting up control node watcher")
			watchFile := strings.Split(*controlNodesPtr, "/")
			watchPath := strings.TrimSuffix(*controlNodesPtr, watchFile[len(watchFile)-1])
			nodeWatcher, err := WatchFile(watchPath, time.Second, func() {
				fmt.Println("control node event")
				_, err := os.Stat(*controlNodesPtr)
				if !os.IsNotExist(err) {
					nodeManager(controlNodesPtr, "control", contrailClient)
				} else if os.IsNotExist(err) {
					controlNodes(contrailClient, []*types.ControlNode{})
				}
			})
			check(err)

			defer func() {
				nodeWatcher.Close()
			}()
		}

		if vrouterNodesPtr != nil {
			fmt.Println("initial vrouter node run")
			_, err := os.Stat(*vrouterNodesPtr)
			if !os.IsNotExist(err) {
				nodeManager(vrouterNodesPtr, "vrouter", contrailClient)
			} else if os.IsNotExist(err) {
				vrouternodes.ReconcileVrouterNodes(contrailClient, []*types.VrouterNode{})
			}
			fmt.Println("setting up vrouter node watcher")
			watchFile := strings.Split(*vrouterNodesPtr, "/")
			watchPath := strings.TrimSuffix(*vrouterNodesPtr, watchFile[len(watchFile)-1])
			nodeWatcher, err := WatchFile(watchPath, time.Second, func() {
				fmt.Println("vrouter node event")
				_, err := os.Stat(*vrouterNodesPtr)
				if !os.IsNotExist(err) {
					nodeManager(vrouterNodesPtr, "vrouter", contrailClient)
				} else if os.IsNotExist(err) {
					vrouternodes.ReconcileVrouterNodes(contrailClient, []*types.VrouterNode{})
				}
			})
			check(err)

			defer func() {
				nodeWatcher.Close()
			}()
		}

		if analyticsNodesPtr != nil {
			fmt.Println("initial analytics node run")
			_, err := os.Stat(*analyticsNodesPtr)
			if !os.IsNotExist(err) {
				nodeManager(analyticsNodesPtr, "analytics", contrailClient)
			} else if os.IsNotExist(err) {
				analyticsNodes(contrailClient, []*types.AnalyticsNode{})
			}
			fmt.Println("setting up analytics node watcher")
			watchFile := strings.Split(*analyticsNodesPtr, "/")
			watchPath := strings.TrimSuffix(*analyticsNodesPtr, watchFile[len(watchFile)-1])
			nodeWatcher, err := WatchFile(watchPath, time.Second, func() {
				fmt.Println("analytics node event")
				_, err := os.Stat(*analyticsNodesPtr)
				if !os.IsNotExist(err) {
					nodeManager(analyticsNodesPtr, "analytics", contrailClient)
				} else if os.IsNotExist(err) {
					analyticsNodes(contrailClient, []*types.AnalyticsNode{})
				}
			})
			check(err)

			defer func() {
				nodeWatcher.Close()
			}()
		}

		if configNodesPtr != nil {
			fmt.Println("initial config node run")
			_, err := os.Stat(*configNodesPtr)
			if !os.IsNotExist(err) {
				nodeManager(configNodesPtr, "config", contrailClient)
			} else if os.IsNotExist(err) {
				configNodes(contrailClient, []*types.ConfigNode{})
			}
			fmt.Println("setting up config node watcher")
			watchFile := strings.Split(*configNodesPtr, "/")
			watchPath := strings.TrimSuffix(*configNodesPtr, watchFile[len(watchFile)-1])
			nodeWatcher, err := WatchFile(watchPath, time.Second, func() {
				fmt.Println("config node event")
				_, err := os.Stat(*configNodesPtr)
				if !os.IsNotExist(err) {
					nodeManager(configNodesPtr, "config", contrailClient)
				} else if os.IsNotExist(err) {
					configNodes(contrailClient, []*types.ConfigNode{})
				}
			})
			check(err)

			defer func() {
				nodeWatcher.Close()
			}()
		}

		if databaseNodesPtr != nil {
			fmt.Println("initial database node run")
			_, err := os.Stat(*databaseNodesPtr)
			if !os.IsNotExist(err) {
				nodeManager(databaseNodesPtr, "database", contrailClient)
			} else if os.IsNotExist(err) {
				databaseNodes(contrailClient, []*types.DatabaseNode{})
			}
			fmt.Println("setting up database node watcher")
			watchFile := strings.Split(*databaseNodesPtr, "/")
			watchPath := strings.TrimSuffix(*databaseNodesPtr, watchFile[len(watchFile)-1])
			nodeWatcher, err := WatchFile(watchPath, time.Second, func() {
				fmt.Println("database node event")
				_, err := os.Stat(*databaseNodesPtr)
				if !os.IsNotExist(err) {
					nodeManager(databaseNodesPtr, "database", contrailClient)
				} else if os.IsNotExist(err) {
					databaseNodes(contrailClient, []*types.DatabaseNode{})
				}
			})
			check(err)

			defer func() {
				nodeWatcher.Close()
			}()
		}

		<-done
	}

	if *modePtr == "run" {

		var apiServer APIServer

		apiServerYaml, err := ioutil.ReadFile(*apiserverPtr)
		if err != nil {
			panic(err)
		}
		err = yaml.Unmarshal(apiServerYaml, &apiServer)
		if err != nil {
			panic(err)
		}

		var keystoneAuthParameters *KeystoneAuthParameters = &KeystoneAuthParameters{}
		if _, err := os.Stat(*keystoneAuthConfPtr); err == nil {
			keystoneAuthParameters = getKeystoneAuthParametersFromFile(*keystoneAuthConfPtr)
		}

		contrailClient, err := getAPIClient(&apiServer, keystoneAuthParameters)
		if err != nil {
			panic(err.Error())
		}

		if controlNodesPtr != nil {
			var controlNodeList []*types.ControlNode
			controlNodeYaml, err := ioutil.ReadFile(*controlNodesPtr)
			if err != nil {
				panic(err)
			}
			err = yaml.Unmarshal(controlNodeYaml, &controlNodeList)
			if err != nil {
				panic(err)
			}
			err = retry(5, 10*time.Second, func() (err error) {
				err = controlNodes(contrailClient, controlNodeList)
				return
			})
			if err != nil {
				panic(err)
			}
		}

		if configNodesPtr != nil {
			var configNodeList []*types.ConfigNode
			configNodeYaml, err := ioutil.ReadFile(*configNodesPtr)
			if err != nil {
				panic(err)
			}
			err = yaml.Unmarshal(configNodeYaml, &configNodeList)
			if err != nil {
				panic(err)
			}
			if err = configNodes(contrailClient, configNodeList); err != nil {
				panic(err)
			}
		}

		if analyticsNodesPtr != nil {
			var analyticsNodeList []*types.AnalyticsNode
			analyticsNodeYaml, err := ioutil.ReadFile(*analyticsNodesPtr)
			if err != nil {
				panic(err)
			}
			err = yaml.Unmarshal(analyticsNodeYaml, &analyticsNodeList)
			if err != nil {
				panic(err)
			}
			if err = analyticsNodes(contrailClient, analyticsNodeList); err != nil {
				panic(err)
			}
		}

		if vrouterNodesPtr != nil {
			var vrouterNodeList []*types.VrouterNode
			vrouterNodeYaml, err := ioutil.ReadFile(*vrouterNodesPtr)
			if err != nil {
				panic(err)
			}
			err = yaml.Unmarshal(vrouterNodeYaml, &vrouterNodeList)
			if err != nil {
				panic(err)
			}
			if err = vrouternodes.ReconcileVrouterNodes(contrailClient, vrouterNodeList); err != nil {
				panic(err)
			}
		}

		if databaseNodesPtr != nil {
			var databaseNodeList []*types.DatabaseNode
			databaseNodeYaml, err := ioutil.ReadFile(*databaseNodesPtr)
			if err != nil {
				panic(err)
			}
			err = yaml.Unmarshal(databaseNodeYaml, &databaseNodeList)
			if err != nil {
				panic(err)
			}
			if err = databaseNodes(contrailClient, databaseNodeList); err != nil {
				panic(err)
			}
		}

	}
}
func retry(attempts int, sleep time.Duration, f func() error) (err error) {
	for i := 0; ; i++ {
		err = f()
		if err == nil {
			return
		}
		if attempts != 0 {
			if i >= (attempts - 1) {
				break
			}
		}

		time.Sleep(sleep)

		fmt.Println("retrying after error:", err)
	}
	return err
}

func connectionError(err error) bool {
	if err == nil {
		fmt.Println("Ok")
		return false

	} else if netError, ok := err.(net.Error); ok && netError.Timeout() {
		fmt.Println("Timeout")
		return true
	}
	unwrappedError := errors.Unwrap(err)
	switch t := unwrappedError.(type) {
	case *net.OpError:
		if t.Op == "dial" {
			fmt.Println("Unknown host")
			return true
		} else if t.Op == "read" {
			fmt.Println("Connection refused")
			return true
		}

	case syscall.Errno:
		if t == syscall.ECONNREFUSED {
			fmt.Println("Connection refused")
			return true
		}

	default:
		fmt.Println(t)
	}
	return false
}

func getAPIClient(apiServerObj *APIServer, keystoneAuthParameters *KeystoneAuthParameters) (*contrail.Client, error) {
	var contrailClient *contrail.Client
	for _, apiServer := range apiServerObj.APIServerList {
		apiServerSlice := strings.Split(apiServer, ":")
		apiPortInt, err := strconv.Atoi(apiServerSlice[1])
		if err != nil {
			return contrailClient, err
		}
		fmt.Printf("api server %s:%d\n", apiServerSlice[0], apiPortInt)
		contrailClient := contrail.NewClient(apiServerSlice[0], apiPortInt)
		err = contrailClient.AddEncryption(apiServerObj.Encryption.CA, apiServerObj.Encryption.Key, apiServerObj.Encryption.Cert, true)
		if err != nil {
			return nil, err
		}
		if keystoneAuthParameters.AuthUrl != "" {
			setupAuthKeystone(contrailClient, keystoneAuthParameters)
		}
		//contrailClient.AddHTTPParameter(1)
		_, err = contrailClient.List("global-system-config")
		if err == nil {
			return contrailClient, nil
		}
	}
	return contrailClient, fmt.Errorf("%s", "cannot get api server")

}

func controlNodes(contrailClient *contrail.Client, nodeList []*types.ControlNode) error {
	var actionMap = make(map[string]string)
	nodeType := "bgp-router"
	vncNodes := []*types.ControlNode{}
	vncNodeList, err := contrailClient.List(nodeType)
	if err != nil {
		return err
	}
	for _, vncNode := range vncNodeList {
		obj, err := contrailClient.ReadListResult(nodeType, &vncNode)
		if err != nil {
			return err
		}
		typedNode := obj.(*contrailTypes.BgpRouter)
		bgpRouterParamters := typedNode.GetBgpRouterParameters()
		if bgpRouterParamters.RouterType == "control-node" {
			node := &types.ControlNode{
				IPAddress: bgpRouterParamters.Address,
				Hostname:  typedNode.GetName(),
				ASN:       bgpRouterParamters.AutonomousSystem,
			}
			vncNodes = append(vncNodes, node)
		}
	}
	for _, node := range nodeList {
		actionMap[node.Hostname] = "create"
	}
	for _, vncNode := range vncNodes {
		if _, ok := actionMap[vncNode.Hostname]; ok {
			for _, node := range nodeList {
				if node.Hostname == vncNode.Hostname {
					actionMap[node.Hostname] = "noop"
					if node.IPAddress != vncNode.IPAddress {
						actionMap[node.Hostname] = "update"
					}
					if node.ASN != vncNode.ASN {
						actionMap[node.Hostname] = "update"
					}
				}
			}
		} else {

			actionMap[vncNode.Hostname] = "delete"
		}
	}
	for k, v := range actionMap {
		switch v {
		case "update":
			for _, node := range nodeList {
				if node.Hostname == k {
					fmt.Println("updating node ", node.Hostname)
					err = node.Update(nodeList, k, contrailClient)
					if err != nil {
						return err
					}
				}
			}
		case "create":
			for _, node := range nodeList {
				if node.Hostname == k {
					fmt.Println("creating node ", node.Hostname)
					err = node.Create(nodeList, node.Hostname, contrailClient)
					if err != nil {
						return err
					}
				}
			}
		case "delete":
			node := &types.ControlNode{}
			err = node.Delete(k, contrailClient)
			if err != nil {
				return err
			}
			fmt.Println("deleting node ", k)
		}
	}
	return nil
}

func configNodes(contrailClient *contrail.Client, nodeList []*types.ConfigNode) error {
	var actionMap = make(map[string]string)
	nodeType := "config-node"
	vncNodes := []*types.ConfigNode{}
	vncNodeList, err := contrailClient.List(nodeType)
	if err != nil {
		return err
	}
	for _, vncNode := range vncNodeList {
		obj, err := contrailClient.ReadListResult(nodeType, &vncNode)
		if err != nil {
			return err
		}
		typedNode := obj.(*contrailTypes.ConfigNode)

		node := &types.ConfigNode{
			IPAddress: typedNode.GetConfigNodeIpAddress(),
			Hostname:  typedNode.GetName(),
		}
		vncNodes = append(vncNodes, node)
	}
	for _, node := range nodeList {
		actionMap[node.Hostname] = "create"
	}
	for _, vncNode := range vncNodes {
		if _, ok := actionMap[vncNode.Hostname]; ok {
			for _, node := range nodeList {
				if node.Hostname == vncNode.Hostname {
					actionMap[node.Hostname] = "noop"
					if node.IPAddress != vncNode.IPAddress {
						actionMap[node.Hostname] = "update"
					}
				}
			}
		} else {
			actionMap[vncNode.Hostname] = "delete"
		}
	}
	for k, v := range actionMap {
		switch v {
		case "update":
			for _, node := range nodeList {
				if node.Hostname == k {
					fmt.Println("updating node ", node.Hostname)
					err = node.Update(nodeList, k, contrailClient)
					if err != nil {
						return err
					}
				}
			}
		case "create":
			for _, node := range nodeList {
				if node.Hostname == k {
					fmt.Println("creating node ", node.Hostname)
					err = node.Create(nodeList, node.Hostname, contrailClient)
					if err != nil {
						return err
					}
				}
			}
		case "delete":
			node := &types.ConfigNode{}
			err = node.Delete(k, contrailClient)
			if err != nil {
				return err
			}
			fmt.Println("deleting node ", k)
		}
	}
	return nil
}

func analyticsNodes(contrailClient *contrail.Client, nodeList []*types.AnalyticsNode) error {
	var actionMap = make(map[string]string)
	nodeType := "analytics-node"
	vncNodes := []*types.AnalyticsNode{}
	vncNodeList, err := contrailClient.List(nodeType)
	if err != nil {
		return err
	}
	for _, vncNode := range vncNodeList {
		obj, err := contrailClient.ReadListResult(nodeType, &vncNode)
		if err != nil {
			return err
		}
		typedNode := obj.(*contrailTypes.AnalyticsNode)

		node := &types.AnalyticsNode{
			IPAddress: typedNode.GetAnalyticsNodeIpAddress(),
			Hostname:  typedNode.GetName(),
		}
		vncNodes = append(vncNodes, node)
	}
	for _, node := range nodeList {
		actionMap[node.Hostname] = "create"
	}
	for _, vncNode := range vncNodes {
		if _, ok := actionMap[vncNode.Hostname]; ok {
			for _, node := range nodeList {
				if node.Hostname == vncNode.Hostname {
					actionMap[node.Hostname] = "noop"
					if node.IPAddress != vncNode.IPAddress {
						actionMap[node.Hostname] = "update"
					}
				}
			}
		} else {
			actionMap[vncNode.Hostname] = "delete"
		}
	}
	for k, v := range actionMap {
		switch v {
		case "update":
			for _, node := range nodeList {
				if node.Hostname == k {
					fmt.Println("updating node ", node.Hostname)
					err = node.Update(nodeList, k, contrailClient)
					if err != nil {
						return err
					}
				}
			}
		case "create":
			for _, node := range nodeList {
				if node.Hostname == k {
					fmt.Println("creating node ", node.Hostname)
					err = node.Create(nodeList, node.Hostname, contrailClient)
					if err != nil {
						return err
					}
				}
			}
		case "delete":
			node := &types.ConfigNode{}
			err = node.Delete(k, contrailClient)
			if err != nil {
				return err
			}
			fmt.Println("deleting node ", k)
		}
	}
	return nil
}

func databaseNodes(contrailClient *contrail.Client, nodeList []*types.DatabaseNode) error {
	var actionMap = make(map[string]string)
	nodeType := "database-node"
	vncNodes := []*types.DatabaseNode{}
	vncNodeList, err := contrailClient.List(nodeType)
	if err != nil {
		return err
	}
	log.Printf("VncNodeList: %v\n", vncNodeList)
	for _, vncNode := range vncNodeList {
		obj, err := contrailClient.ReadListResult(nodeType, &vncNode)
		if err != nil {
			return err
		}
		typedNode := obj.(*contrailTypes.DatabaseNode)

		node := &types.DatabaseNode{
			IPAddress: typedNode.GetDatabaseNodeIpAddress(),
			Hostname:  typedNode.GetName(),
		}
		vncNodes = append(vncNodes, node)
	}
	for _, node := range nodeList {
		actionMap[node.Hostname] = "create"
	}
	log.Printf("VncNodes: %v\n", vncNodes)

	for _, vncNode := range vncNodes {
		if _, ok := actionMap[vncNode.Hostname]; ok {
			for _, node := range nodeList {
				if node.Hostname == vncNode.Hostname {
					actionMap[node.Hostname] = "noop"
					if node.IPAddress != vncNode.IPAddress {
						actionMap[node.Hostname] = "update"
					}
				}
			}
		} else {
			actionMap[vncNode.Hostname] = "delete"
		}
	}
	for k, v := range actionMap {
		log.Printf("actionMapValue: %v\n", v)
		switch v {
		case "update":
			for _, node := range nodeList {
				if node.Hostname == k {
					log.Println("updating node ", node.Hostname)
					err = node.Update(nodeList, k, contrailClient)
					if err != nil {
						return err
					}
				}
			}
		case "create":
			for _, node := range nodeList {
				if node.Hostname == k {
					log.Println("creating node ", node.Hostname)
					err = node.Create(nodeList, node.Hostname, contrailClient)
					if err != nil {
						return err
					}
				}
			}
		case "delete":
			node := &types.DatabaseNode{}
			err = node.Delete(k, contrailClient)
			if err != nil {
				return err
			}
			log.Println("deleting node ", k)
		}
	}
	return nil
}

func setupAuthKeystone(client *contrail.Client, keystoneAuthParameters *KeystoneAuthParameters) {
	var keystone *contrail.KeepaliveKeystoneClient
	if strings.HasPrefix(keystoneAuthParameters.AuthUrl, "https") {
		// AddEncryption expected http url in older versions of contrail-go-api
		// https://github.com/Juniper/contrail-go-api/commit/4c876ba038a8ecec211376133375d467b6098202
		keystone = contrail.NewKeepaliveKeystoneClient(
			strings.Replace(keystoneAuthParameters.AuthUrl, "https", "http", 1),
			keystoneAuthParameters.TenantName,
			keystoneAuthParameters.AdminUsername,
			keystoneAuthParameters.AdminPassword,
			"",
		)
		err := keystone.AddEncryption(
			keystoneAuthParameters.Encryption.CA,
			keystoneAuthParameters.Encryption.Key,
			keystoneAuthParameters.Encryption.Cert,
			keystoneAuthParameters.Encryption.Insecure)

		if err != nil {
			fmt.Fprintln(os.Stderr, err)
			os.Exit(1)
		}
	} else {
		keystone = contrail.NewKeepaliveKeystoneClient(
			keystoneAuthParameters.AuthUrl,
			keystoneAuthParameters.TenantName,
			keystoneAuthParameters.AdminUsername,
			keystoneAuthParameters.AdminPassword,
			"",
		)
	}
	err := keystone.AuthenticateV3()
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
	client.SetAuthenticator(keystone)

}

func getKeystoneAuthParametersFromFile(authParamsFilePath string) *KeystoneAuthParameters {
	var keystoneAuthParameters *KeystoneAuthParameters
	keystoneAuthYaml, err := ioutil.ReadFile(authParamsFilePath)
	if err != nil {
		panic(err)
	}
	if err = yaml.Unmarshal(keystoneAuthYaml, &keystoneAuthParameters); err != nil {
		panic(err)
	}
	return keystoneAuthParameters
}

func getGlobalVrouterConfigFromFile(globalVrouterFilePath string) *GlobalVrouterConfiguration {
	var globalVrouterConfig *GlobalVrouterConfiguration
	globalVrouterJson, err := ioutil.ReadFile(globalVrouterFilePath)
	if err != nil {
		panic(err)
	}
	if err = json.Unmarshal([]byte(globalVrouterJson), &globalVrouterConfig); err != nil {
		panic(err)
	}
	return globalVrouterConfig
}
