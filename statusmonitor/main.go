package main

import (
	"crypto/tls"
	"crypto/x509"
	"encoding/xml"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"time"

	yaml "gopkg.in/yaml.v2"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/apimachinery/pkg/runtime/serializer"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/kubernetes/scheme"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
	"k8s.io/client-go/util/retry"

	contrailOperatorTypes "github.com/Juniper/contrail-operator/pkg/apis/contrail/v1alpha1"
)

//NodeType definition
type NodeType string

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

type ControlConnectionInfo struct {
	Type   string   `xml:"type"`
	Name   string   `xml:"name"`
	Status string   `xml:"status"`
	Nodes  []string `xml:"server_addrs>list>element"`
}

type ProcessControlStatus struct {
	XMLName           xml.Name                `xml:"__NodeStatusUVE_list"`
	State             string                  `xml:"NodeStatusUVE>data>NodeStatus>process_status>list>ProcessStatus>state"`
	ConConnectionInfo []ControlConnectionInfo `xml:"NodeStatusUVE>data>NodeStatus>process_status>list>ProcessStatus>connection_infos>list>ConnectionInfo"`
}

type ServiceStatusControl struct {
	XMLName                  xml.Name `xml:"__BGPRouterInfo_list"`
	NumberOfXMPPPeers        int      `xml:"BGPRouterInfo>data>BgpRouterState>num_xmpp_peer"`
	NumberOfRoutingInstances int      `xml:"BGPRouterInfo>data>BgpRouterState>num_routing_instance"`
	NumStaticRoutes          int      `xml:"BGPRouterInfo>data>BgpRouterState>num_static_routes"`
	NumDownStaticRoutes      int      `xml:"BGPRouterInfo>data>BgpRouterState>num_down_static_routes"`
	NumBgpPeer               int      `xml:"BGPRouterInfo>data>BgpRouterState>num_bgp_peer"`
	NumUpBgpPeer             int      `xml:"BGPRouterInfo>data>BgpRouterState>num_up_bgp_peer"`
}

func main() {
	log.Println("Starting status monitor")
	configPtr := flag.String("config", "/config.yaml", "path to config yaml file")
	intervalPtr := flag.Int64("interval", 1, "interval for getting status")
	flag.Parse()

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
		log.Printf("error: rest client creation failed: %v", err)
		panic(err)
	}
	clientset, restClient, err := kubeClient(config)
	if err != nil {
		log.Printf("kubernates client creation failed: %v", err)
		panic(err)
	}

	ticker := time.NewTicker(time.Duration(*intervalPtr) * time.Second)
	done := make(chan bool)
	go func() {
		for {
			select {
			case t := <-ticker.C:
				log.Println("Tick at", t)
				var updatedConfig Config
				updatedConfigYaml, err := ioutil.ReadFile(*configPtr)
				if err != nil {
					panic(err)
				}
				err = yaml.Unmarshal(updatedConfigYaml, &updatedConfig)
				if err != nil {
					panic(err)
				}
				switch config.NodeType {
				case "control":
					err := getControlStatus(client, clientset, restClient, updatedConfig)
					if err != nil {
						log.Printf("warning: Error in  getControlStatus func: %v", err)
						continue
					}
				case "config":
					err := getConfigStatus(client, clientset, restClient, updatedConfig)
					if err != nil {
						log.Printf("warning: Error in  getConfigStatus func: %v", err)
						continue
					}
				}
			}
		}
	}()
	done <- true
	log.Println("Ticker stopped")
}

func getControlStatus(client http.Client, clientset *kubernetes.Clientset, restClient *rest.RESTClient, config Config) error {
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

	for _, hostname := range hostnameList {
		var url string
		var process_url string
		url = "https://" + apiServer + "/Snh_SandeshUVECacheReq?x=BgpRouterState"
		process_url = "https://" + apiServer + "/Snh_SandeshUVECacheReq?x=NodeStatus"
		resp, err := client.Get(url)
		if resp != nil {
			defer closeResp(resp)
		}
		process_resp, err_p := client.Get(process_url)
		if process_resp != nil {
			defer closeResp(process_resp)
		}
		if err != nil {
			log.Println(err)
			return err
		}
		if err_p != nil {
			log.Println(err_p)
			return err_p
		}
		log.Print(url)
		log.Print(process_url)
		if resp != nil && process_resp != nil {
			log.Printf("resp not nil %d ", resp.StatusCode)
			log.Printf("Process resp not nil %d ", process_resp.StatusCode)
			if resp.StatusCode == http.StatusOK && process_resp.StatusCode == http.StatusOK {
				bodyBytes, err := ioutil.ReadAll(resp.Body)
				bodyBytes_p, err_p := ioutil.ReadAll(process_resp.Body)
				if err != nil {
					log.Printf("Error while reading BgpRouterState response: %v", err)
					return err
				}
				if err_p != nil {
					log.Printf("Error while reading NodeStatus response: %v", err_p)
					return err_p
				}
				controlStatus, err := getControlStatusFromResponse(bodyBytes, bodyBytes_p)
				if err != nil {
					log.Printf("Error while reading ControlStatus response: %v", err)
					return err
				}
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
		log.Println(podHostnames)
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
		if err != nil {
			log.Printf("error: updateControlStatus: Failed to get status: %s", err)
			return err
		}
		controlObject.Status.ServiceStatus = controlStatusMap
		_, err = controlCient.UpdateStatus(config.NodeName, controlObject)
		if err != nil {
			log.Println(err)
		}
		return err
	})
	if retryErr != nil {
		log.Printf("Update failed: %v", retryErr)
	}
	return retryErr
}

func getControlStatusFromResponse(statusBody []byte, statusBody_p []byte) (*contrailOperatorTypes.ControlServiceStatus, error) {
	controlst := ServiceStatusControl{}
	processst := ProcessControlStatus{}

	err := xml.Unmarshal(statusBody, &controlst)
	if err != nil {
		log.Printf("Error while unmarshalling ServiceStatusControl: %v", err)
		return nil, err
	}

	err = xml.Unmarshal(statusBody_p, &processst)
	if err != nil {
		log.Printf("Error while unmarshalling ProcessControlStatus: %v", err)
		return nil, err
	}

	staticRoutes := contrailOperatorTypes.StaticRoutes{
		Down:   strconv.Itoa(controlst.NumDownStaticRoutes),
		Number: strconv.Itoa(controlst.NumStaticRoutes),
	}
	bgpPeer := contrailOperatorTypes.BGPPeer{
		Up:     strconv.Itoa(controlst.NumUpBgpPeer),
		Number: strconv.Itoa(controlst.NumBgpPeer),
	}
	numUpXMPPPeer := strconv.Itoa(controlst.NumberOfXMPPPeers)
	numRoutingInstance := strconv.Itoa(controlst.NumberOfRoutingInstances)

	connectionList := []contrailOperatorTypes.Connection{}

	for _, ConnectionInfo := range processst.ConConnectionInfo {
		connection_ps := contrailOperatorTypes.Connection{
			Type:   ConnectionInfo.Type,
			Name:   ConnectionInfo.Name,
			Status: ConnectionInfo.Status,
			Nodes:  ConnectionInfo.Nodes,
		}
		connectionList = append(connectionList, connection_ps)
	}

	controlStatus := contrailOperatorTypes.ControlServiceStatus{
		Connections:              connectionList,
		NumberOfXMPPPeers:        numUpXMPPPeer,
		NumberOfRoutingInstances: numRoutingInstance,
		StaticRoutes:             staticRoutes,
		BGPPeer:                  bgpPeer,
		State:                    processst.State,
	}

	return &controlStatus, err
}
