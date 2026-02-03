# DevStack – OpenStack Single Node Setup

Minimal setup for running a single-node OpenStack environment using DevStack.

---

## Overview

This repository contains:
- A `local.conf` file for DevStack
- Instructions to deploy OpenStack on Ubuntu 24.04
- A small graph visualizing the [![VM Creation Flow](SetupDiagram.svg)](SetupDiagram.svg)

The setup is intended for **learning, testing, and development** purposes.

---

## System Requirements

- Virtual machine or bare-metal server
- Ubuntu 24.04 (recommended)
- Internet access
- User with `sudo` privileges

---

## Installation

### 1. Install Devstack

```bash
sudo useradd -s /bin/bash -d /opt/stack -m stack
```
```bash
sudo chmod +x /opt/stack
```
```bash
echo "stack ALL=(ALL) NOPASSWD: ALL" | sudo tee /etc/sudoers.d/stack
```
```bash
sudo -u stack -i
```

```bash
git clone https://opendev.org/openstack/devstack
```


### 2. Clone this repo and copy the config file

```bash
git clone https://github.com/dreule28/level3-cloud.git && cp level3-cloud/local.conf ~/devstack/
```


### 3. Run DevStack
```bash
cd ~/devstack
```
#### Modify the HOTS_IP with your own
```bash
nano local.conf
```
```bash
./stack.sh
```
Result after successful execution:
OpenStack services are running locally.
Horizon dashboard is available.
OpenStack CLI can be used from the system.


### 4. Install Terraform
```bash
sudo apt-get update && sudo apt-get install -y gnupg software-properties-common
```
```bash
wget -O- https://apt.releases.hashicorp.com/gpg | \
gpg --dearmor | \
sudo tee /usr/share/keyrings/hashicorp-archive-keyring.gpg > /dev/null
```
```bash
gpg --no-default-keyring \
--keyring /usr/share/keyrings/hashicorp-archive-keyring.gpg \
--fingerprint
```
```bash
echo "deb [arch=$(dpkg --print-architecture) signed-by=/usr/share/keyrings/hashicorp-archive-keyring.gpg] https://apt.releases.hashicorp.com $(grep -oP '(?<=UBUNTU_CODENAME=).*' /etc/os-release || lsb_release -cs) main" | sudo tee /etc/apt/sources.list.d/hashicorp.list
```
```bash
sudo apt update
```
```bash
sudo apt-get install terraform
```
#### Verify installation
```bash
terraform -help
```
#### Recommended alias for better usage
```bash
alias tfi='terraform init'
alias tfa='terraform apply -auto-approve'
alias tfd='terraform destroy -auto-approve'
alias tfo='terraform output'

source ~/.bashrc <-- after modifying
```
After that run
```bash
tfi
tfa
```
