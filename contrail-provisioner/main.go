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

	contrailtypes "github.com/Juniper/contrail-operator/contrail-provisioner/contrail-go-types"
	"github.com/Juniper/contrail-operator/contrail-provisioner/contrailclient"
	"github.com/Juniper/contrail-operator/contrail-provisioner/contrailnode"
	"github.com/Juniper/contrail-operator/contrail-provisioner/nodemanager"
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

// LinkLocalServiceEntryType struct defines link local service
type LinkLocalServiceEntryType struct {
	LinkLocalServiceName   string   `json:"linkLocalServiceName,omitempty"`
	LinkLocalServiceIP     string   `json:"linkLocalServiceIP,omitempty"`
	LinkLocalServicePort   int      `json:"linkLocalServicePort,omitempty"`
	IPFabricDNSServiceName string   `json:"ipFabricDNSServiceName,omitempty"`
	IPFabricServicePort    int      `json:"ipFabricServicePort,omitempty"`
	IPFabricServiceIP      []string `json:"ipFabricServiceIP,omitempty"`
}

// LinkLocalServicesTypes struct contains list of link local services definitions
type LinkLocalServicesTypes struct {
	LinkLocalServicesEntries []LinkLocalServiceEntryType `json:"linkLocalServicesEntries,omitempty"`
}

type GlobalVrouterConfiguration struct {
	EcmpHashingIncludeFields   EcmpHashingIncludeFields `json:"ecmpHashingIncludeFields,omitempty"`
	EncapsulationPriorities    string                   `json:"encapPriority,omitempty"`
	VxlanNetworkIdentifierMode string                   `json:"vxlanNetworkIdentifierMode,omitempty"`
	LinkLocalServices          LinkLocalServicesTypes   `json:"linkLocalServices,omitempty"`
}

const RequiredAnnotationsKey = "managed_by"
const MaxRetryAttempts = 5
const BackoffTimeSeconds = 10

func loadBytesFromFile(filePath string) []byte {
	var data []byte
	_, err := os.Stat(filePath)
	if os.IsNotExist(err) {
		return data
	}
	data, err = ioutil.ReadFile(filePath)
	if err != nil {
		panic(err)
	}
	return data
}

func check(err error) {
	if err != nil {
		log.Fatalf("error: %v", err)
	}
}

func runNodeManager(filePath string, nodeType contrailnode.ContrailNodeType, contrailClient contrailclient.ApiClient, requiredAnnotations map[string]string) {
	requiredNodesData := loadBytesFromFile(filePath)
	err := retry(MaxRetryAttempts, BackoffTimeSeconds*time.Second, func() (err error) {
		return nodemanager.ManageNodes(requiredNodesData, requiredAnnotations, nodeType, contrailClient)
	})
	if err != nil {
		log.Fatalf("%s node manager failed after %d attempts with error: %s\n", nodeType, MaxRetryAttempts, err)
	}
}

func setupNodeFileWatcher(filePath string, nodeType contrailnode.ContrailNodeType, contrailClient contrailclient.ApiClient, requiredAnnotations map[string]string) *FileWatcher {
	log.Printf("Initial run of node manager for %s\n", nodeType)
	runNodeManager(filePath, nodeType, contrailClient, requiredAnnotations)
	log.Printf("Setting up file watcher for %s listed in %s\n", nodeType, filePath)
	watchFile := strings.Split(filePath, "/")
	watchPath := strings.TrimSuffix(filePath, watchFile[len(watchFile)-1])
	nodeWatcher, err := WatchFile(watchPath, time.Second, func() {
		log.Printf("%s node event\n", nodeType)
		runNodeManager(filePath, nodeType, contrailClient, requiredAnnotations)
	})
	check(err)
	return nodeWatcher
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
	requiredAnnotationsPtr := flag.String("requiredAnnotations", "/etc/provision/metadata/managed_by", "path to file with required annotation value")
	modePtr := flag.String("mode", "watch", "watch/run")
	flag.Parse()

	requiredAnnotationValue := string(loadBytesFromFile(*requiredAnnotationsPtr))
	requiredAnnotations := map[string]string{RequiredAnnotationsKey: requiredAnnotationValue}

	log.Printf("Required annotations for all objects managed by contrail-provisioner: %v", requiredAnnotations)

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
		encapPriorityObj := &contrailtypes.EncapsulationPrioritiesType{Encapsulation: encapPriority}
		ecmpObj := globalVrouterConfiguration.EcmpHashingIncludeFields
		ecmpHashingIncludeFieldsObj := &contrailtypes.EcmpHashingIncludeFields{ecmpObj.HashingConfigured, ecmpObj.SourceIp, ecmpObj.DestinationIp, ecmpObj.IpProtocol, ecmpObj.SourcePort, ecmpObj.DestinationPort}
		GlobalVrouterConfig := &contrailtypes.GlobalVrouterConfig{}
		linkLocalServicesTypesObj := operatorLinkLocalToContrailType(globalVrouterConfiguration.LinkLocalServices.LinkLocalServicesEntries)
		GlobalVrouterConfig.SetFQName("", globalVrouterConfFQName)
		GlobalVrouterConfig.SetEncapsulationPriorities(encapPriorityObj)
		GlobalVrouterConfig.SetEcmpHashingIncludeFields(ecmpHashingIncludeFieldsObj)
		GlobalVrouterConfig.SetVxlanNetworkIdentifierMode(globalVrouterConfiguration.VxlanNetworkIdentifierMode)
		GlobalVrouterConfig.SetLinklocalServices(&linkLocalServicesTypesObj)
		if err = contrailClient.Create(GlobalVrouterConfig); err != nil {
			if !strings.Contains(err.Error(), "409 Conflict") {
				panic(err)
			}
			obj, err := contrailClient.FindByName("global-vrouter-config", strings.Join(globalVrouterConfFQName, ":"))
			if err != nil {
				panic(err)
			}
			obj.(*contrailtypes.GlobalVrouterConfig).SetEncapsulationPriorities(encapPriorityObj)
			obj.(*contrailtypes.GlobalVrouterConfig).SetEcmpHashingIncludeFields(ecmpHashingIncludeFieldsObj)
			obj.(*contrailtypes.GlobalVrouterConfig).SetVxlanNetworkIdentifierMode(globalVrouterConfiguration.VxlanNetworkIdentifierMode)
			obj.(*contrailtypes.GlobalVrouterConfig).SetLinklocalServices(&linkLocalServicesTypesObj)
			if err = contrailClient.Update(obj); err != nil {
				panic(err)
			}
		}

		log.Println("start watcher")
		done := make(chan bool)

		if controlNodesPtr != nil {
			nodeWatcher := setupNodeFileWatcher(*controlNodesPtr, contrailnode.ControlNode, contrailClient, requiredAnnotations)
			defer func() {
				nodeWatcher.Close()
			}()
		}

		if vrouterNodesPtr != nil {
			nodeWatcher := setupNodeFileWatcher(*vrouterNodesPtr, contrailnode.VrouterNode, contrailClient, requiredAnnotations)
			defer func() {
				nodeWatcher.Close()
			}()
		}

		if analyticsNodesPtr != nil {
			nodeWatcher := setupNodeFileWatcher(*analyticsNodesPtr, contrailnode.AnalyticsNode, contrailClient, requiredAnnotations)
			defer func() {
				nodeWatcher.Close()
			}()
		}

		if configNodesPtr != nil {
			nodeWatcher := setupNodeFileWatcher(*configNodesPtr, contrailnode.ConfigNode, contrailClient, requiredAnnotations)
			defer func() {
				nodeWatcher.Close()
			}()
		}

		if databaseNodesPtr != nil {
			nodeWatcher := setupNodeFileWatcher(*databaseNodesPtr, contrailnode.DatabaseNode, contrailClient, requiredAnnotations)
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
			runNodeManager(*controlNodesPtr, contrailnode.ControlNode, contrailClient, requiredAnnotations)
		}

		if vrouterNodesPtr != nil {
			runNodeManager(*vrouterNodesPtr, contrailnode.VrouterNode, contrailClient, requiredAnnotations)
		}

		if configNodesPtr != nil {
			runNodeManager(*configNodesPtr, contrailnode.ConfigNode, contrailClient, requiredAnnotations)
		}

		if analyticsNodesPtr != nil {
			runNodeManager(*analyticsNodesPtr, contrailnode.AnalyticsNode, contrailClient, requiredAnnotations)
		}

		if databaseNodesPtr != nil {
			runNodeManager(*databaseNodesPtr, contrailnode.DatabaseNode, contrailClient, requiredAnnotations)
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

		log.Println("retrying after error:", err)
	}
	return err
}

func connectionError(err error) bool {
	if err == nil {
		log.Println("Ok")
		return false

	} else if netError, ok := err.(net.Error); ok && netError.Timeout() {
		log.Println("Timeout")
		return true
	}
	unwrappedError := errors.Unwrap(err)
	switch t := unwrappedError.(type) {
	case *net.OpError:
		if t.Op == "dial" {
			log.Println("Unknown host")
			return true
		} else if t.Op == "read" {
			log.Println("Connection refused")
			return true
		}

	case syscall.Errno:
		if t == syscall.ECONNREFUSED {
			log.Println("Connection refused")
			return true
		}

	default:
		log.Println(t)
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
		log.Printf("api server %s:%d\n", apiServerSlice[0], apiPortInt)
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

func operatorLinkLocalToContrailType(links []LinkLocalServiceEntryType) contrailtypes.LinklocalServicesTypes {
	var entries []contrailtypes.LinklocalServiceEntryType
	for _, entry := range links {
		entries = append(entries,
			contrailtypes.LinklocalServiceEntryType{entry.LinkLocalServiceName,
				entry.LinkLocalServiceIP,
				entry.LinkLocalServicePort,
				entry.IPFabricDNSServiceName,
				entry.IPFabricServicePort,
				entry.IPFabricServiceIP})
	}
	return contrailtypes.LinklocalServicesTypes{entries}
}
