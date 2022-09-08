# _FRIEZA_

### 简介

`Frieza`是一个jmeter分布式工具的加强版,更有效的控制jmeter的分布式执行效果

### 相较于jmeter分布式启动方式的优势

* `Frieza`知道各个`slave`所处的运行状态
* `Frieza`会对`java`环境进行检测，并尝试安装`jdk1.8`,当安装失败后会停止运行
* `Frieza`无需修改`master`或`slave`的`properties`中的`server.rmi.ssl.disable`,`remote_hosts`
  等参数（运行jmeter命令时也无需指定-R参数）
* `Frieza`具有更强的控制`slave`负载机的运行与停止的能力，避免`master`停止，但是`slave`仍在运行的情况，避免部分`slave`异常导致无法启动分布式命令

### 配置与启动方式

1.将Frieza-master工具放置到`jmeter`的`bin`目录下      
2.赋予工具可执行权限

```text
  chmod +x Frieza-master
```

3.启动./Frieza-master

### Frieza-master交互式命令

* **slave**    
  查看所有连接成功的`slave`的运行状态

```text
    Starting        // 启动中 
    Idle            // 空闲状态
    Running         // 运行中
    Stopped         // 停止状态
    Failed          // 启动失败
```

```text
    Frieze>>> slave
    nums    Ip              status
    1       192.168.31.1    Starting
    2       192.168.31.2    Idle
    3       192.168.31.3    Idle
    4       192.168.31.4    Idle
```

* **jmeter**   
  与jmeter命令行执行方式一样，命令会自动添加`slave`中`status`处于`Idle`状态的设备

```text
    Frieze>>> jmeter -n -t demo.jmx
```

* **log**   
  读取`jmeter.log`文件的内容（查看jmeter压测数据，jmeter运行日志）

```text
    Frieze>>> log
    ...
    2022-09-02 17:49:33,504 INFO o.a.j.r.Summariser: summary +      1 in 00:00:05 =    0.2/s Avg:    46 Min:    46 Max:    46 Err:     0 (0.00%) Active: 1 Started: 1 Finished: 0
    2022-09-02 17:50:18,739 INFO o.a.j.r.Summariser: summary +      9 in 00:00:45 =    0.2/s Avg:    22 Min:    21 Max:    24 Err:     0 (0.00%) Active: 0 Started: 1 Finished: 1
    2022-09-02 17:50:18,739 INFO o.a.j.r.Summariser: summary =     10 in 00:00:50 =    0.2/s Avg:    24 Min:    21 Max:    46 Err:     0 (0.00%)
    ...
```

* **stop**    
  停止所有 `status = Running`的`slave`,使其停止运行并恢复到`Idle`状态

```text
    Frieze>>> stop
```

### 测试系统支持情况

- [x] CentOS 7.6   
- [x] CentOS 8.0
- [x] CentOS 8.2
- [x] TencentOS Server 3.1 (TK4)
- [x] TencentOS Server 2.4 (TK4)
- [ ] Ubuntu 18.04.1 LTS
- [ ] Ubuntu 20.04 LTS
- [ ] Windows Server
- [ ] Debian

### TODO
1. [ ] 更多系统支持
2. [ ] 自动根据现有的设备情况自动配置`slave`的`JVM`
3. [ ] master自动检查jmeter脚本中的results模块是否已经被禁用
4. [ ] 根据`csv`或`txt`文件以及目前`slave`的数量自动分发测试数据到各个`slave`的指定目录中
5. [ ] 检测并提示`master`和`slave`中的`jmeter`版本