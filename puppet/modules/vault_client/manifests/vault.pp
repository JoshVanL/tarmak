define vault_client::vault (
  String $user = 'root',
  String $group = 'root',
)
{
  $script_name = 'vault'

  file { "${::vault_client::local_bin_dir}/${script_name}.sh":
    ensure  => file,
    content => template('vault_client/vault.sh'),
    notify  => Script["${script_name}.sh"],
    owner   => $user,
    group   => $group,
    mode    => '0644',
  }
}