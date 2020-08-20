package main

import (
	"crypto/tls"
	"crypto/x509"
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"reflect"
	"time"

	yaml "gopkg.in/yaml.v2"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/apimachinery/pkg/runtime/serializer"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/kubernetes/scheme"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
	"k8s.io/client-go/util/retry"

	contrailOperatorTypes "github.com/Juniper/contrail-operator/pkg/apis/contrail/v1alpha1"
	"github.com/Juniper/contrail-operator/statusmonitor/uves"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

//NodeType definition
type NodeType string

const (
	//VROUTER nodetype
	VROUTER NodeType = "vrouter"
	//CONFIG nodetype
	CONFIG = "config"
	//ANALYTICS nodetype
	ANALYTICS = "analytics"
	//CONTROL nodetype
	CONTROL = "control"
)

//Config struct
type Config struct {
	APIServerList  []string   `yaml:"apiServerList,omitempty"`
	Encryption     encryption `yaml:"encryption,omitempty"`
	NodeType       NodeType   `yaml:"nodeType,omitempty"`
	Interval       int64      `yaml:"interval,omitempty"`
	Hostname       string     `yaml:"hostname,omitempty"`
	InCluster      *bool      `yaml:"inCluster,omitempty"`
	KubeConfigPath string     `yaml:"kubeConfigPath,omitempty"`
	NodeName       string     `yaml:"nodeName,omitempty"`
	Namespace      string     `yaml:"namespace,omitempty"`
	PodName        string     `yaml:"podName,omitempty"`
}

type encryption struct {
	CA       *string `yaml:"ca,omitempty"`
	Cert     *string `yaml:"cert,omitempty"`
	Key      *string `yaml:"key,omitempty"`
	Insecure bool    `yaml:"insecure,omitempty"`
}

func check(err error) {
	if err != nil {
		log.Fatalf("error: %v", err)
	}
}

func main() {
	log.Println("Starting status monitor")
	configPtr := flag.String("config", "/config.yaml", "path to config yaml file")
	intervalPtr := flag.Int64("interval", 1, "interval for getting status")
	flag.Parse()

	ticker := time.NewTicker(time.Duration(*intervalPtr) * time.Second)
	done := make(chan bool)
	go func() {
		for {
			select {
			case t := <-ticker.C:
				fmt.Println("Tick at", t)
				var config Config
				configYaml, err := ioutil.ReadFile(*configPtr)
				if err != nil {
					panic(err)
				}
				err = yaml.Unmarshal(configYaml, &config)
				if err != nil {
					panic(err)
				}
				client, err := CreateRestClient(config)
				if err != nil {
					log.Printf("rest client creation failed: %v", err)
					continue
				}
				clientset, restClient, err := kubeClient(config)
				if err != nil {
					log.Printf("kubernates client creation failed: %v", err)
					continue
				}
				ticker = time.NewTicker(time.Duration(config.Interval) * time.Second)
				switch config.NodeType {
				case "control":
					err := getControlStatus(client, clientset, restClient, config)
					if err != nil{
						log.Printf("warning: Error in  getControlStatus func: %v", err)
						continue
					}
				case "config":
					getConfigStatus(client, clientset, restClient, config)
				}
			}
		}
	}()
	done <- true
	fmt.Println("Ticker stopped")
}

func getControlStatus(client http.Client, clientset *kubernetes.Clientset, restClient *rest.RESTClient, config Config) (error){
	var controlStatusMap = make(map[string]contrailOperatorTypes.ControlServiceStatus)
	for _, apiServer := range config.APIServerList {
		hostnameList, err := getPods(config, clientset)
		if err != nil {
			log.Printf("warning: Unable to get the hostnameList in getControlStatus func: %v", err)
			return err
		}

		err = GetControlStatusFromApiServer(apiServer, &config, &client, hostnameList, controlStatusMap)
		if err != nil {
			log.Printf("warning: Failed in GetControlStatusFromApiServer func: %v", err)
			return err
		}
		if len(controlStatusMap) == len(hostnameList) {
			break
		}
	}
	err := updateControlStatus(&config, controlStatusMap, restClient)
	if err != nil {
		log.Printf("warning: Getting error in updateControlStatus: %v", err)
		return err
	}
	return nil
}

func CreateRestClient(config Config) (http.Client, error) {
	tlsConfig := &tls.Config{
		InsecureSkipVerify: config.Encryption.Insecure,
	}
	if config.Encryption.Key != nil && config.Encryption.Cert != nil && config.Encryption.CA != nil {
		caCert, err := ioutil.ReadFile(*config.Encryption.CA)
		pool := x509.NewCertPool()
		pool.AppendCertsFromPEM(caCert)
		tlsConfig.RootCAs = pool
		cer, err := tls.LoadX509KeyPair(*config.Encryption.Cert, *config.Encryption.Key)
		if err != nil {
			log.Println(err)
			return http.Client{}, err
		}
		tlsConfig.Certificates = []tls.Certificate{cer}
	}
	transport := http.Transport{
		TLSClientConfig: tlsConfig,
	}
	client := http.Client{
		Transport: &transport,
	}
	return client, nil
}

func GetControlStatusFromApiServer(apiServer string, config *Config, client *http.Client, hostnameList []string,
	controlStatusMap map[string]contrailOperatorTypes.ControlServiceStatus) error {
	var nodeType string
	switch config.NodeType {
	case "config":
		nodeType = "config-node"
	case "control":
		nodeType = "control-node"
	case "analytics":
		nodeType = "analytics-node"
	case "vrouter":
		nodeType = "vrouter"
	}
	for _, hostname := range hostnameList {
		var url string
		url = "https://" + apiServer + "/analytics/uves/" + nodeType + "/" + hostname
		resp, err := client.Get(url)
		if err != nil {
			log.Println(err)
			return err
		}
		log.Print(url)
		defer closeResp(resp)
		if resp != nil {
			log.Printf("resp not nil %d ", resp.StatusCode)
			if resp.StatusCode == http.StatusOK {
				bodyBytes, err := ioutil.ReadAll(resp.Body)
				if err != nil {
					log.Println(err)
					return err
				}
				controlStatus := getControlStatusFromResponse(bodyBytes)
				controlStatusMap[hostname] = *controlStatus
			}
		}
	}
	return nil
}

func closeResp(resp *http.Response) {
	err := resp.Body.Close()
	if err != nil {
		log.Printf("closing http session failed: %v", err)
	}
}
func homeDir() string {
	if h := os.Getenv("HOME"); h != "" {
		return h
	}
	return os.Getenv("USERPROFILE") // windows
}

func getPods(config Config, clientSet kubernetes.Interface) ([]string, error) {
	var podHostnames []string
	podList, err := clientSet.CoreV1().Pods(config.Namespace).List(metav1.ListOptions{LabelSelector: string(config.NodeType) + "=" + config.NodeName})
	if err != nil {
		return podHostnames, fmt.Errorf("get pods failed: %v", err)
	}
	var podAnnotations map[string]string

	for _, pod := range podList.Items {
		podAnnotations = pod.GetAnnotations()
		podHostnames = append(podHostnames, podAnnotations["hostname"])
	}
	if len(podHostnames) > 0 {
		fmt.Println(podHostnames)
	}
	return podHostnames, nil
}

func kubeClient(config Config) (*kubernetes.Clientset, *rest.RESTClient, error) {

	var err error
	clientset := &kubernetes.Clientset{}
	restClient := &rest.RESTClient{}
	kubeConfig := &rest.Config{}
	if config.InCluster != nil && !*config.InCluster {
		var kubeConfigPath string
		if config.KubeConfigPath != "" {
			kubeConfigPath = config.KubeConfigPath
		} else {
			kubeConfigPath = filepath.Join(homeDir(), ".kube", "config")
		}
		kubeConfig, err = clientcmd.BuildConfigFromFlags("", kubeConfigPath)
		if err != nil {
			return clientset, restClient, err
		}

	} else {
		kubeConfig, err = rest.InClusterConfig()
		if err != nil {
			return clientset, restClient, err
		}
		kubeConfig.CAFile = ""
		kubeConfig.TLSClientConfig.Insecure = true
	}
	// create the clientset
	err = contrailOperatorTypes.SchemeBuilder.AddToScheme(scheme.Scheme)
	if err != nil {
		return nil, nil, err
	}

	crdConfig := kubeConfig
	crdConfig.ContentConfig.GroupVersion = &schema.GroupVersion{Group: contrailOperatorTypes.SchemeGroupVersion.Group, Version: contrailOperatorTypes.SchemeGroupVersion.Version}
	crdConfig.APIPath = "/apis"

	crdConfig.NegotiatedSerializer = serializer.WithoutConversionCodecFactory{CodecFactory: scheme.Codecs}
	crdConfig.UserAgent = rest.DefaultKubernetesUserAgent()

	restClient, err = rest.UnversionedRESTClientFor(crdConfig)
	if err != nil {
		return clientset, restClient, err
	}
	clientset, err = kubernetes.NewForConfig(crdConfig)
	if err != nil {
		return clientset, restClient, err
	}
	return clientset, restClient, nil
}

type controlClient struct {
	restClient rest.Interface
	ns         string
}

type vrouterClient struct {
	restClient rest.Interface
	ns         string
}

func (c *controlClient) Get(name string, opts metav1.GetOptions) (*contrailOperatorTypes.Control, error) {
	result := contrailOperatorTypes.Control{}
	err := c.restClient.
		Get().
		Namespace(c.ns).
		Resource("controls").
		Name(name).
		VersionedParams(&opts, scheme.ParameterCodec).
		Do().
		Into(&result)

	return &result, err
}

func (c *controlClient) UpdateStatus(name string, object *contrailOperatorTypes.Control) (*contrailOperatorTypes.Control, error) {
	result := contrailOperatorTypes.Control{}
	err := c.restClient.
		Put().
		Namespace(c.ns).
		Resource("controls").
		Name(name).
		SubResource("status").
		Body(object).
		Do().
		Into(&result)

	return &result, err
}

func (c *vrouterClient) Get(name string, opts metav1.GetOptions) (*contrailOperatorTypes.Vrouter, error) {
	result := contrailOperatorTypes.Vrouter{}
	err := c.restClient.
		Get().
		Namespace(c.ns).
		Resource("vrouters").
		Name(name).
		VersionedParams(&opts, scheme.ParameterCodec).
		Do().
		Into(&result)

	return &result, err
}

func (c *vrouterClient) UpdateStatus(name string, object *contrailOperatorTypes.Vrouter) (*contrailOperatorTypes.Vrouter, error) {
	result := contrailOperatorTypes.Vrouter{}
	err := c.restClient.
		Put().
		Namespace(c.ns).
		Resource("vrouters").
		Name(name).
		SubResource("status").
		Body(object).
		Do().
		Into(&result)

	return &result, err
}

func updateControlStatus(config *Config, controlStatusMap map[string]contrailOperatorTypes.ControlServiceStatus, restClient *rest.RESTClient) error {
	retryErr := retry.RetryOnConflict(retry.DefaultRetry, func() error {

		controlCient := &controlClient{
			ns:         config.Namespace,
			restClient: restClient,
		}
		controlObject, err := controlCient.Get(config.NodeName, metav1.GetOptions{})
		check(err)
		update := false
		if controlObject.Status.ServiceStatus == nil {
			controlObject.Status.ServiceStatus = controlStatusMap
			update = true
		} else {
			same := reflect.DeepEqual(controlObject.Status.ServiceStatus, controlStatusMap)
			if !same {
				controlObject.Status.ServiceStatus = controlStatusMap
				update = true

			} else {
				controlObject.Status.ServiceStatus = controlStatusMap
				update = true
			}
		}
		if update {
			_, err = controlCient.UpdateStatus(config.NodeName, controlObject)
			if err != nil {
				fmt.Println(err)
			}
			return err
		}
		return nil
	})
	if retryErr != nil {
		panic(fmt.Errorf("Update failed: %v", retryErr))
	}
	return nil
}

func getControlStatusFromResponse(statusBody []byte) *contrailOperatorTypes.ControlServiceStatus {
	controlUVEStatus := &uves.ControlUVEStatus{}
	err := json.Unmarshal(statusBody, controlUVEStatus)
	if err != nil {
		log.Fatal(err)
	}
	connectionList := []contrailOperatorTypes.Connection{}

	bgpRouterList := typeSwitch(controlUVEStatus.BgpRouterState.BgpRouterIPList)
	bgpRouterConnection := contrailOperatorTypes.Connection{
		Type:  "BGPRouter",
		Nodes: bgpRouterList,
	}
	connectionList = append(connectionList, bgpRouterConnection)
	if len(connectionList) > 0 && len(controlUVEStatus.NodeStatus.ProcessStatus.List.ProcessStatus) > 0 {
		for _, connectionInfo := range controlUVEStatus.NodeStatus.ProcessStatus.List.ProcessStatus[0].ConnectionInfos.List.ConnectionInfo {
			nodeList := typeSwitch(connectionInfo.ServerAddrs.List.Element)
			connection := contrailOperatorTypes.Connection{
				Type:   connectionInfo.Type.Text,
				Name:   connectionInfo.Name.Text,
				Status: connectionInfo.Status.Text,
				Nodes:  nodeList,
			}
			connectionList = append(connectionList, connection)
		}
	}

	numUpXMPPPeer := "0"
	numRoutingInstance := "0"
	staticRoutes := contrailOperatorTypes.StaticRoutes{}
	bgpPeer := contrailOperatorTypes.BGPPeer{}
	state := "down"

	if controlUVEStatus != nil {
		numDownStaticRoutes := controlUVEStatus.BgpRouterState.NumDownStaticRoutes.Status()
		numStaticRoutes := controlUVEStatus.BgpRouterState.NumStaticRoutes.Status()
		staticRoutes = contrailOperatorTypes.StaticRoutes{
			Down:   numDownStaticRoutes,
			Number: numStaticRoutes,
		}
		numUpBgpPeer := controlUVEStatus.BgpRouterState.NumUpBgpPeer.Status()
		numBgpPeer := controlUVEStatus.BgpRouterState.NumBgpPeer.Status()
		bgpPeer = contrailOperatorTypes.BGPPeer{
			Up:     numUpBgpPeer,
			Number: numBgpPeer,
		}
		numUpXMPPPeer = controlUVEStatus.BgpRouterState.NumUpXMPPPeer.Status()
		numRoutingInstance = controlUVEStatus.BgpRouterState.NumRoutingInstance.Status()
		if len(controlUVEStatus.NodeStatus.ProcessStatus.List.ProcessStatus) > 0 {
			state = controlUVEStatus.NodeStatus.ProcessStatus.List.ProcessStatus[0].State.Text
		}
	}

	controlStatus := contrailOperatorTypes.ControlServiceStatus{
		Connections:              connectionList,
		NumberOfXMPPPeers:        numUpXMPPPeer,
		NumberOfRoutingInstances: numRoutingInstance,
		StaticRoutes:             staticRoutes,
		BGPPeer:                  bgpPeer,
		State:                    state,
	}
	return &controlStatus
}

func typeSwitch(tst interface{}) []string {
	var nodeList []string
	switch v := tst.(type) {
	case interface{}:
		inter, ok := v.([][]interface{})
		if ok {
			for _, element := range inter {
				for _, element2 := range element {
					x, ok := element2.(map[string]interface{})
					if ok {
						y, ok := x["list"].(map[string]interface{})
						if ok {
							z, ok := y["element"].([]interface{})
							if ok {
								for _, zz := range z {
									nodeList = append(nodeList, zz.(string))
								}
							}
						}
					}
				}
			}
		} else {
			inter2, ok := v.([]interface{})
			if ok {
				for _, element := range inter2 {
					nodeList = append(nodeList, element.(string))
				}
			} else {
				inter3, ok := v.(interface{})
				if ok {
					nodeList = append(nodeList, inter3.(string))
				}
			}
		}
	case string:
		fmt.Println("String:", v)
	case [][]interface{}:
		fmt.Println("[][]interface{}:", v)
	default:
		fmt.Println("unknown")
	}
	return nodeList
}
