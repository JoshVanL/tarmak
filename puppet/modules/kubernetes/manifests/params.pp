# == Class kubernetes::params
class kubernetes::params {
  $version = '1.10.6'
  $bin_dir = '/opt/bin'
  $dest_dir = '/opt'
  $config_dir = '/etc/kubernetes'
  $run_dir = '/var/run/kubernetes'
  $apply_dir = '/etc/kubernetes/apply'
  $delete_dir = '/etc/kubernetes/delete'
  $download_dir = '/tmp'
  $systemd_dir = '/etc/systemd/system'
  $download_url = 'https://storage.googleapis.com/kubernetes-release/release/v#VERSION#/bin/linux/amd64/hyperkube'
  $log_level = '1'
  $uid = 873
  $gid = 873
  $user = 'kubernetes'
  $group = 'kubernetes'
  $curl_path = $::osfamily ? {
    'RedHat' => '/bin/curl',
    'Debian' => '/usr/bin/curl',
    default  => '/usr/bin/curl',
  }
}
