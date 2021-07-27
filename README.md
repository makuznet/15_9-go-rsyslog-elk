# Golang app with log lib sends logs to ELK
> This repo creates Server1 VPS equipped with Go app logging its input and RSyslog sending this log to Server2 VPS that receives and visualises the log with Elasticsearch and Kibana.

## Usage
## Big picture
Server1: Golang app called numdub using phuslu log lib sends logs in JSON format to the local Syslog.  
```shell
Jul 27 14:54:18 srv1 numdub[1032]: {"time":"2021-07-27T14:54:18.572Z","level":"info","message":"number 1"}
Jul 27 14:54:21 srv1 numdub[1032]: {"time":"2021-07-27T14:54:21.363Z","level":"info","message":"number 2"}
Jul 27 14:54:27 srv1 numdub[1032]: {"time":"2021-07-27T14:54:27.690Z","level":"info","message":"number 3"}
```
Server1: Local RSyslog sends numdub logs to RSyslog on Server2.
See `playbooks/roles/rsyslog/files/10-send-to-server.conf` for configuration details.
```shell
$ActionQueueType Direct # send immediately
$ActionResumeRetryCount -1 #try sending endlessly
$ActionQueueSaveOnShutdown on # Write to the disk in case of shutdown
*.* @@192.168.8.20:514 # @@ — use tcp when sending logs to the RSyslog server
```
Server2: Local RSyslog receives logs from Server1. Local RSyslog is equipped with omelasticsearch module that sends logs to localy installed Elasticsearch.
See `playbooks/roles/rsyslog-server/files/10-remote-logger.conf` for configuration details.  

Server2: Locally installed Kibana shows logs after configuring index template. 

## Logging in Golang
I tried Logrus, Zerolog, and phuslu log lib.
Finally, I imported phuslu as this is much easier to work with than with former libs.
```go
log.Info().Msgf("number %d", num)
```
Info() means time and severity.  
Msgf() allows include vars.  
%d means print int var mentioned after comma that is `num` in my case.  

Syntax changes when using a log lib.  
Standard log:
```go
log.Fatal(http.ListenAndServe(":8080", nil))
```
With phuslu log lib:
```go
log.Fatal().Err(http.ListenAndServe(":8080", nil))
```


## RSyslog
I configured RSyslog on both servers using these two articles in Russian:
- [RSYSLOG + ELASTICSEARCH НАСТРОЙКА RSYSLOG](https://www.casp.ru/2016/10/14/Настройка-rsyslog-storage/)  
- [ОТПРАВКА JSON ЧЕРЕЗ RSYSLOG В ELASTICSEARCH](https://www.casp.ru/2016/10/14/json-over-rsyslog-to-elasticsearch/)  

Reading `playbooks/roles/rsyslog-server/files/10-remote-logger.conf` file I can't say I understand templates lines in full.  
More time and experiments are needed to get them.

### Clone this repo
git clone https://github...

### Terraform
To roll out and configure two VPSes got to a project folderrun
```shell
terraform init
terraform apply --auto-approve
```
### Ansible
#### Vault password
> You'll be asked for a Vault password when Terraform launches Ansible as id_rsa and id_rsa.pub are encrypted.
```shell
ansible-playbook -i terraform-ya/inventory.yml playbooks/main.yml --ask-vault-pass
```
To decrypt files:
```shell
ansible-vault decrypt foo.yml bar.yml baz.yml
```
#### Roles
See playbooks/main.yml to know what roles belong to what VPS. See roles/`role_name`/tasks/main.yml  to get through installation process.  
Roles filebeat and logstash are absent in the final version of playbook.  
They are kept for future use.  

#### Add a line to a file task
```shell
- name: Add a line to a file
  ansible.builtin.lineinfile:
    path: /etc/elasticsearch/elasticsearch.yml
    line: "discovery.type: single-node"
```
### Build Golang binary for Linux
- Download and install Golang according to the [doc](https://golang.org/doc/install);  
- Edit the .go file with your editor;
- Compile the .go file for Linux amd64;
```shell
GOOS=linux GOARCH=amd64 go build numdub.go
```
- And move it to corresponding Ansible dir, eg. `roles/numdub/files`;

### Kibana 
#### Creating an index
Kibana needs an index to show logs.  
Open http://`srv2_ip_addr`:5601 in your web browser.    

Type `index pattern` to the `Search Elastic` search form in the middle of the top of your Kibana window.  
Choose `Index Pattern` then `Create index pattern`.  
Then tick `syslogjson-2021.07.27` and push `Next step` button and follow instructions. 

#### See logs
To see logs open http://`srv2_ip_addr`:5601 in your web browser.  
Menu (top-left) > Analytics > Discover.


## Installation
### Yandex OAuth token
[Yandex.OAuth](https://oauth.yandex.com)

### Yandex CLI (MacOS)
```bash
curl https://storage.yandexcloud.net/yandexcloud-yc/install.sh | bash
brew install bash-completion
echo 'source /Users/makuznet/yandex-cloud/completion.zsh.inc' >>  ~/.zshrc
source "/Users/makuznet/.bash_profile"
yc init # provide your yandex token
yc config profile get <your_profile_name> 
```
### Terraform (MacOS)
```bash
/bin/bash -c "$(curl -fsSL https://raw.githubusercontent.com/Homebrew/install/HEAD/install.sh)"
brew install terraform
```
### Ansible (MacOS)
```bash
https://www.python.org/ftp/python/3.9.5/python-3.9.5-macosx10.9.pkg
python get-pip.py
pip install ansible
```

## Acknowledgments
This repo was inspired by [skillfactory.ru](https://skillfactory.ru/devops#syllabus) team

## See also 
- [Logrus: Golang log lib ](https://github.com/sirupsen/logrus#readme)  
- [Golang: How to write log to a file](https://stackoverflow.com/questions/19965795/how-to-write-log-to-file)  
- [Golang: Can't find main module, but found .git/config](https://stackoverflow.com/questions/67306638/go-test-results-in-go-cannot-find-main-module-but-found-git-config-in-users)  
- [Ansible: add gpg-key](https://docs.ansible.com/ansible/latest/collections/ansible/builtin/apt_key_module.html)  
- [Ansible: Add apt repositories](https://docs.ansible.com/ansible/latest/collections/ansible/builtin/apt_repository_module.html)  
- [Elasticsearch: Installation Debian](https://www.elastic.co/guide/en/elasticsearch/reference/current/deb.html)  
- [Kibana: Installation Debian](https://www.elastic.co/guide/en/kibana/current/deb.html)  
- [RSYSLOG + ELASTICSEARCH НАСТРОЙКА RSYSLOG](https://www.casp.ru/2016/10/14/Настройка-rsyslog-storage/)  
- [ОТПРАВКА JSON ЧЕРЕЗ RSYSLOG В ELASTICSEARCH](https://www.casp.ru/2016/10/14/json-over-rsyslog-to-elasticsearch/)  
- [Parsing JSON (CEE) Logs and Sending them to Elasticsearch](https://www.rsyslog.com/json-elasticsearch/)  
- [Rsyslog to Elasticsearch](https://serverascode.com/2016/11/11/ryslog-to-elasticsearch.html)  
- [How to ship JSON logs via Rsyslog](https://techpunch.co.uk/development/how-to-shop-json-logs-via-rsyslog)  
- [Recipe: rsyslog + Elasticsearch + Kibana](https://sematext.com/blog/recipe-rsyslog-elasticsearch-kibana/)  
- [Rsyslog configuration](https://selivan.github.io/2017/02/07/rsyslog-log-forward-save-filename-handle-multi-line-failover.html#configuration-examples)  
- [Rsyslog на Debian](https://www.k-max.name/linux/rsyslog-na-debian-nastrojka-servera/)  
- [Syslog на Debian](https://www.k-max.name/linux/syslogd-and-logrotate/)  
- [RSyslog Documentation: omelasticsearch: Elasticsearch Output Module](https://www.rsyslog.com/doc/v8-stable/configuration/modules/omelasticsearch.html#server)  
- [Filebeat quick start: installation and configuration](https://www.elastic.co/guide/en/beats/filebeat/7.13/filebeat-installation-configuration.html#filebeat-installation-configuration)  
- [Elastic Stack and Product Documentation](https://www.elastic.co/guide/index.html)  
- [Getting started with Elastic stack](https://www.elastic.co/guide/en/elastic-stack-get-started/7.13/get-started-elastic-stack.html)  
- [Set up minimal security for Elasticsearch](https://www.elastic.co/guide/en/elasticsearch/reference/7.13/security-minimal-setup.html)  
- [Configure Filebeat](https://www.elastic.co/guide/en/beats/filebeat/7.13/configuring-howto-filebeat.html)  
- [Getting Started with Logstash](https://www.elastic.co/guide/en/logstash/current/getting-started-with-logstash.html)  
- [For Logstash: How To Install Java with Apt on Debian 10](https://www.digitalocean.com/community/tutorials/how-to-install-java-with-apt-on-debian-10)

## License
Follow all involved parties licenses terms and conditions.