# RaspberryPI 4 K3s Notes and Documentation

## Initial Setup Guide
- Review K3s documentation for any notes or updates. https://docs.k3s.io/advanced?_highlight=raspb#raspberry-pi
    - Extra install for Ubuntu Server 22 `sudo apt install linux-modules-extra-raspi`
- K3s base setup steps: https://gist.github.com/syncom/7c6e90708bc28cc9ede2c3245c203e32#steps-for-setting-up-k3s-on-ubuntu-20042-on-raspberry-pi-4-cluster
- MetalLB setup steps: https://rpi4cluster.com/k3s/k3s-nw-setting/

## RPI4 Cluster Guides
- https://rpi4cluster.com/

## Node Information
### rpi4-ctrl.local
- IP: 192.168.0.254
- Username: chris.brewin
- SSH: ssh chris.brewin@rpi4-ctrl.local

### rpi4-node01.local
- IP: 192.168.0.253
- Username: chris.brewin
- SSH: ssh chris.brewin@rpi4-node01.local

### rpi4-node2.local
- IP: 192.168.0.252
- Username: chris.brewin
- SSH: ssh chris.brewin@rpi4-node2.local

### rpi4-node3.local
- IP: 192.168.0.251
- Username: chris.brewin
- SSH: ssh chris.brewin@rpi4-node3.local