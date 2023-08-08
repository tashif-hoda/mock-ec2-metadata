# mock-ec2-metadata [![Build Status](https://travis-ci.org/NYTimes/mock-ec2-metadata.svg?branch=master)](https://travis-ci.org/NYTimes/mock-ec2-metadata)


A simple service (written in [go](https://golang.org/) using [gizmo](https://github.com/NYTimes/gizmo)) to mock the ec2 metadata service. This is usefully for development images (like vagrant or packer) that require Instance base IAM permission or other metadata information.

For example, [cob](https://github.com/henrysher/cob) and [s3-iam](https://github.com/seporaitis/yum-s3-iam) can both use s3 as a yum repo. Both of these systems rely on the instances the proper credentials to have authorization to the s3 repos that yum uses.


The metadata service normal listens on a special private ip address `169.254.169.254`. This is a special address that will not exist on your system. One option is to bind an alias to the loopback iterface. This can be done with the following command:

```console
/sbin/ifconfig lo:1 inet 169.254.169.254 netmask 255.255.255.255 up
```

To change it back

```console
/sbin/ifconfig lo:1 inet 127.0.0.1 netmask 255.0.0.0 up
```

Many services assume that use the metadata service uses a default port 80 and do not allow configuration or override. A simple IP talbes rule and IP forwarding can get around that, as follows:

```console
$ echo 1 > /proc/sys/net/ipv4/ip_forward
$ iptables -t nat -A OUTPUT -p tcp -d 169.254.169.254/32 --dport 80  -j DNAT --to-destination 169.254.169.254:8111
$ service iptables save
```

## Configuration
All configuration is contained in either `./mock-ec2-metadata-config.json` or `/etc/mock-ec2-metadata-config.json`, the former overriding the latter.

Currently the support URLs for the metadata service are:

  * http://169.254.169.254/latest/meta-data/
  * http://169.254.169.254/latest/meta-data/ami-id
  * http://169.254.169.254/latest/meta-data/ami-launch-index
  * http://169.254.169.254/latest/meta-data/ami-manifest-path
  * http://169.254.169.254/latest/meta-data/placement/availability-zone
  * http://169.254.169.254/latest/meta-data/hostname
  * http://169.254.169.254/latest/meta-data/instance-action
  * http://169.254.169.254/latest/meta-data/instance-id
  * http://169.254.169.254/latest/meta-data/instance-type
  * http://169.254.169.254/latest/meta-data/local-hostname
  * http://169.254.169.254/latest/meta-data/local-ipv4
  * http://169.254.169.254/latest/meta-data/mac
  * http://169.254.169.254/latest/meta-data/profile
  * http://169.254.169.254/latest/meta-data/reservation-id
  * http://169.254.169.254/latest/meta-data/security-groups
  * http://169.254.169.254/latest/meta-data/iam/security-credentials
  * http://169.254.169.254/latest/meta-data/network/interfaces/macs/{mac_address}/security-group-ids


## Getting started

This project requires Go to be installed. On OS X with Homebrew you can just run `brew install go`.

Running it then should be as simple as:

```console
$ make deps
$ make
$ ./bin/mock-ec2-metadata
```


### Testing

``make test``

## License

See `LICENSE`

## Contributing

See `CONTRIBUTING.md` for more details.
