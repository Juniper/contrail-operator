ansible-playbook -i inventory/mycluster/inventory.yaml cluster.yml \
  -e "kube_network_plugin=cni" \
  -e "kube_network_plugin_multus=false" \
  -e "override_system_hostname=false" \
  -e "kubectl_localhost=true" \
  -e "kubeconfig_localhost=true" \
  -e "container_manager=crio" \
  -e "kubelet_deployment_type=host" \
  -e "download_container=false" \
  -e "etcd_deployment_type=host"
  
