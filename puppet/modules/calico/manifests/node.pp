# Calico Node
#
# Calico Node contains a Daemon Set that spinsup the overlay network on every
# workern node.
# @param metrics_port Port for felix metrics endpoint, 0 disables metrics collection
class calico::node (
  String $node_version = '3.1.4',
  String $cni_version = '3.1.4',
  Enum['always', 'cross-subnet', 'off'] $ipv4_pool_ipip_mode = 'always',
  Integer[0,65535] $metrics_port = 9091,
)
{
  include ::kubernetes
  include ::calico

  $namespace = $::calico::namespace
  $mtu = $::calico::mtu
  $ipv4_pool_cidr = $::calico::pod_network
  $backend = $::calico::backend

  $typha_enabled  = $::calico::typha_enabled
  $typha_replicas = $::calico::typha_replicas


  $authorization_mode = $::kubernetes::_authorization_mode
  if member($authorization_mode, 'RBAC'){
    $rbac_enabled = true
  } else {
    $rbac_enabled = false
  }

  if versioncmp($::kubernetes::version, '1.6.0') >= 0 {
    $version_before_1_6 = false
  } else {
    $version_before_1_6 = true
  }

  if $backend == 'etcd' {
    $etcd_cert_path = $::calico::etcd_cert_path
    $etcd_proto = $::calico::etcd_proto
    $node_image = 'quay.io/calico/node'
    $cni_image = 'quay.io/calico/cni'

    $manifests = ''

  } else {
    $node_image = 'calico/node'
    $cni_image = 'calico/cni'

    if $typha_enabled {
      $manifests = [
          template('calico/node-crd.yaml.erb'),
          template('calico/node-typha.yaml.erb'),
        ]
    } else {
      $manifests = [
          template('calico/node-crd.yaml.erb'),
        ]
    }
  }

  kubernetes::apply{'calico-node':
    ensure    => 'present',
    manifests => [
      template('calico/node-rbac.yaml.erb'),
      template('calico/node-daemonset.yaml.erb'),
      $manifests,
    ],
  }

}
