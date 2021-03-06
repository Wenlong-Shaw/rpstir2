


# Rpstir2
(Relying Party Security Technology for Internet Routing 2)

## 1. 简介
在RPKI体系中，RP软件扮演者重要的角色，是CA与边界路由器之间的桥梁。RPSTIR做为现有RP软件（rpki.net，RPSTIR，RIPE NCC，Routinator）之一，提供包括同步、验证、指导BGP路由在内的所有RP端的标准功能。

RPSTIR2基于原有RPSTIR设计思想重新开发，采用Go语言，对RPKI相关协议支持更加完善，性能得到优化。 RPSTIR2核心模块包括：同步模块，验证模块，RTR模块和数据库模块。

RPSTIR2支持的功能包括：
* 同步协议支持Rsync和RRDP(rfc8182)
* 支持Slurm(rfc8416)
* 支持RTR
* 支持Validation Reconsidered(rfc8360)


## 2. 安装
有两种方式安装RPSTIR2，包括源代码安装和Docker安装


### 2.1 源代码安装
操作系统需要 CentOS8，或Ubuntu18及以上，均需64位系统

#### 2.1.1 安装OpenSSL
OpenSSL版本需要 1.1.1b及以上，并且在编译OpenSSL时，需要设置"enable-rfc3779"，安装时需要用root安装

```shell
$ wget --no-verbose --inet4-only https://www.openssl.org/source/openssl-1.1.1f.tar.gz 
$ tar xzvf openssl-1.1.1f.tar.gz 
$ cd openssl-1.1.1f
$ ./config shared enable-rfc3779
$ make
$ make install
$ echo "export PATH=/usr/local/ssl/bin:$PATH" >> /root/.bashrc
$ source /root/.bashrc
```
 
    
#### 2.1.2 创建MySQL用户
MySQL版本需要8及以上，需要支持json。如果没有安装MySQL，需要根据操作平台，从https://dev.mysql.com/downloads/下载安装MySQL，
安装好MySQL后，请登陆MySQL，按如下脚本，创建PRSTIR2的用户和数据表。 如果有修改，请在2.1.6对应修改RPSTIR2的配置。

```mysql
CREATE USER 'rpstir2'@'localhost' IDENTIFIED WITH mysql_native_password BY 'Rpstir-123';
CREATE USER 'rpstir2'@'%' IDENTIFIED WITH mysql_native_password BY 'Rpstir-123';
flush privileges;

CREATE DATABASE rpstir2;
GRANT ALL PRIVILEGES ON rpstir2.* TO 'rpstir2'@'localhost'  with grant option;
GRANT ALL PRIVILEGES ON rpstir2.* TO 'rpstir2'@'%'  with grant option;
flush privileges;
```

#### 2.1.3 安装GoLang
GoLang版本需要1.13或以上，建议用root安装。

```shell
$ wget --no-verbose --inet4-only https://dl.google.com/go/go1.14.1.linux-amd64.tar.gz
$ tar -C /usr/local -xzf go1.14.1.linux-amd64.tar.gz
$ echo "export GOROOT=/usr/local/go" >> /root/.bashrc 
$ echo "export GOPATH=/usr/local/goext" >> /root/.bashrc 
$ echo "export PATH=$PATH:/usr/local/go/bin:/usr/local/goext/bin" >> /root/.bashrc 
$ source  /root/.bashrc
```
#### 2.1.4 创建RPSTIR2目录
请创建如下2个目录，分别是安装目录和数据目录。如果有修改，则在2.1.6修改RPSTIR2的配置

```shell
$ mkdir -p /root/rpki/ /root/rpki/rpstir2  /root/rpki/data 
```

| Directory  | Path                      |
| :--------: | ------------------------- |
| programdir | /root/rpki/rpstir2        |
| datadir    | /root/rpki/data           |


#### 2.1.5 下载RPSTIR2
进入源代码目录，通过git下载RPSTIR2。

```shell
$ cd /root/rpki/source/
$ git clone https://github.com/bgpsecurity/rpstir2.git 
```

#### 2.1.6 配置RPSTIR2
RPSTIR2的配置文件在/root/rpki/rpstir2/conf/project.conf，并按MySQL实际配置，修改MySQL参数。
其中RTR端口为8082，路由器将链接到此端口。如果端口有冲突可以修改端口。

```shell
$ cd /root/rpki/rpstir2/conf
$ vim project.conf

[rpstir2]
programdir=/root/rpki/rpstir2
datadir=/root/rpki/data
httpport=8080
httpsport=8081

[mysql]
server=127.0.0.1:3306
user=rpstir2
password=Rpstir-123
database=rpstir2

[rtr]
httpserver=127.0.0.1
httpport=8083
tcpserver=127.0.0.1
tcpport=8082
```



##### 2.1.7 部署RPSTIR2
进入RPSTIR2的bin目录，然后执行如下脚本，并部署成功后将自动启动RPSTIR2

```shell
$ cd /root/rpki/rpstir2/bin
$ chmod +x *.sh 
$ ./rpstir2-service.sh deploy
$ ./rpstir2-service.sh update
```

#### 2.1.8 配置RPSTIR2定时运行（可选）
可以通过crontab设置RPSTIR2定时同步. 

```shell
$ crontab -e
1 1 * * *  /root/rpki/rpstir2/bin/rpstir2-command.sh crontab
```
   

#### 2.2 Docker安装
##### 2.2.1 拉取MySQL镜像并配置
拉取MySQL镜像，并初始化MySQL的root密码。请注意，由于在Docker中使用了桥连接，所以MySQL数据库对外暴露的IP是172.17.0.1，端口是13306。

```shell
$ docker pull mysql
$ docker run -itd --name rpstir2_mysql -p 13306:3306 -e MYSQL_ROOT_PASSWORD=Rpstir-123 mysql
$ docker exec -it rpstir2_mysql /bin/bash
$ mysql -uroot -p
Rpstir-123
```

登录进入MySQL后，创建RPSTIR2的用户配置，类似2.1.2. 

```SQL
CREATE USER 'rpstir2'@'localhost' IDENTIFIED WITH mysql_native_password BY 'Rpstir-123';
CREATE USER 'rpstir2'@'%' IDENTIFIED WITH mysql_native_password BY 'Rpstir-123';
flush privileges;

CREATE DATABASE rpstir2;
GRANT ALL PRIVILEGES ON rpstir2.* TO 'rpstir2'@'localhost'  with grant option;
GRANT ALL PRIVILEGES ON rpstir2.* TO 'rpstir2'@'%'  with grant option;
flush privileges;
quit;
```

##### 2.2.2 拉取RPSTIR2镜像
先在主机配置好映射目录"/root/rpki/rpstir2data/data/"和"/root/rpki/rpstir2data/log"，分别对应到Docker中的RPSTIR2的数据目录和日志目录；并且配置Docker对外暴露的RPSTIR2端口。其中RTR的对应端口为18082，路由器将链接到此端口。如果端口有冲突可以修改端口。

```shell
$ cd /root/rpki/
$ mkdir -p /root/rpki/rpstir2data /root/rpki/rpstir2data/data  /root/rpki/rpstir2data/log
$ docker pull cpusoft/rpstir2_centos8
$ docker run -itd --privileged -p 18080-18090:8080-8090   -v /root/rpki/rpstir2data/data:/root/rpki/data  -v /root/rpki/rpstir2data/log:/root/rpki/rpstir2/log --name rpstir2_centos8 cpusoft/rpstir2_centos8  /usr/sbin/init
```

##### 2.2.3 更新RPSTIR2
可以登录进入Docker，执行部署和更新命令。 

```shell
$ docker exec -it rpstir2_centos8 /bin/bash
$ cd /root/rpki/rpstir2/bin 
$ chmod +x *.sh
$ ./rpstir2-service.sh deploy
$ ./rpstir2-service.sh update
```

备注1：如果需要配置crontabe，同样参考2.1.8
备注2：请注意，由于在Docker中使用了桥连接，所以MySQL数据库对外暴露的IP是172.17.0.1，所以RPSTIR2的配置文件中MySQL为"172.17.0.1:13306"


### 2.3 运行RPSTIR2
所有运行命令通过rpstir2/bin目录下的两个脚本执行：
rpstir2-serverice.sh 执行系统起停、升级等系统命令
rpstir2-command.sh 执行具体的同步、查看状态和结果、上传Slurm等程序命令，需要在系统启动后，才能执行。

#### 2.3.1 起停RPSTIR2
可以通过查看./log目录下的日志，查看有无报错

```shell
$ cd /root/rpki/rpstir2/bin
$./rpstir2-serverice.sh start 
$./rpstir2-serverice.sh stop 
```

#### 2.3.2 同步和验证RPKI
分别按照rsync或rrdp方式同步和验证RPKI数据，或者参数为sync时由系统自动执行混合同步。
注：当前建议仍然通过rsync进行验证

##### (1) rsync

```shell
$ cd /root/rpki/rpstir2/bin
$ ./rpstir2-command.sh rsync 
```

##### (2) rrdp

```shell
$ cd /root/rpki/rpstir2/bin
$ ./rpstir2-command.sh rrdp  
```

##### (3) sync

```shell
$ cd /root/rpki/rpstir2/bin
$ ./rpstir2-command.sh sync 
```


#### 2.3.3 获取同步和验证的状态
由于rsync和rrdp时间较长，并且是后台执行，需要通过以下命令获取同步和验证的结果，查看其过程是否结束.

```shell
$ cd /root/rpki/rpstir2/bin
$ ./rpstir2-command.sh states  
```

获得的json结果中，如果"state"的显示为running，则表示还在运行中，如果显示end，则表示运行完毕，此时路由器可以通过RTR端口获取RPKI数据了。


```JSON
{ "result": "ok",
  "msg": "",
  "state":
  	{ "endTime": "2019-12-19 14:07:11", 
     	"startTime": "2019-12-18 16:29:06",
     	"state": "end" 
    } 
 }
```
#### 2.3.4 结果统计
通过如下命令获取同步和验证结果。

```shell
$ cd /root/rpki/rpstir2/bin
$./rpstir2-command.sh results  
```
```JSON
{
    "cerResult": {
        "fileType": "cer",
        "validCount": 16920,
        "warningCount": 0,
        "invalidCount": 6
    },
    "crlResult": {
        "fileType": "crl",
        "validCount": 16916,
        "warningCount": 0,
        "invalidCount": 51
    },
    "mftResult": {
        "fileType": "mft",
        "validCount": 16914,
        "warningCount": 0,
        "invalidCount": 71
    },
    "roaResult": {
        "fileType": "roa",
        "validCount": 31779,
        "warningCount": 0,
        "invalidCount": 288
    }
}
```



#### 2.3.5 Slurm

```shell
$ cd /root/rpki/rpstir2/bin
$./rpstir2-command.sh slurm xxx.json  
```

#### 2.3.6 Help
两个命令的更多参数，可以分别输入help得到

```shell
$ cd /root/rpki/rpstir2/bin
$./rpstir2-service.sh help
$./rpstir2-command.sh help
```

## 3 技术咨询和Bug报告

 如果发现任何bug或者有任何问题，欢迎提出issue，或者发邮箱到 shaoqing@zdns.cn


