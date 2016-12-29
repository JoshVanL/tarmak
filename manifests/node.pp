define calico::node (
  String  $calico_node_version,
  Integer $etcd_count,
  Integer $calico_etcd_port,
)
{
  include ::systemd
  include k8s

  wget::fetch { "calicoctl-v${calico::calico_node_version}":
    source      => "https://github.com/projectcalico/calico-containers/releases/download/v${calico_node_version}/calicoctl",
    destination => '/opt/cni/bin/',
    require     => Class['calico'],
    before      => File['/opt/cni/bin/calicoctl'],
  }

  file { '/opt/cni/bin/calicoctl':
    ensure => file,
    mode   => '0755',
  }

  file { '/etc/calico/calico.env':
    ensure  => file,
    content => template('calico/calico.env.erb'),
    require => Class['calico'],
  }

  file { '/usr/lib/systemd/system/calico-node.service':
    ensure  => file,
    content => template('calico/calico-node.service.erb'),
  } ~>
  Exec['systemctl-daemon-reload']

  service { 'calico-node':
    ensure  => running,
    enable  => true,
    require => [ Class['k8s'], File['/etc/calico/calico.env'], File['/usr/lib/systemd/system/calico-node.service'] ],
  }

  file { '/usr/local/sbin/calico_filter_hack.sh':
    ensure  => file,
    content => template('calico/calico_filter_hack.sh.erb'),
    mode    => '0750',
  }

  exec { 'Modify calico filter':
    command => '/usr/local/sbin/calico_filter_hack.sh set',
    unless  => '/usr/local/sbin/calico_filter_hack.sh test',
    require => Service['calico-node'],
  }
}
