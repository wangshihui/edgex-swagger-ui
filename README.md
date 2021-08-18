# edgex 设备能力注册中心 设备抽象定义 # 

## 基于edgex的app service sdk 实现 ## 

### why ###

有两个alternative的方案

- 基于edgex的device sdk实现，定义抽象设备, edgex的device sdk提供了设备自动发现接口, 但是在进行设备注册时,会同时将设备注册到 core-metadata

- 脱离edgex,完全自己实现, edgex 提供了诸多开箱即用的工具,如果从头实现的话,会有许多额外的工作需要处理

### Goal ###

实现设备能力的抽象,设备和设备的能力两个概念分离

主要的两个目的

- 在用户侧可以根据需要的能力去找设备,

- 可以将设备的角色进行转换, 可以在系统中自行组装有任意能力的设备, 称为抽象设备
  
- 可以将抽象设备和一个或者多个物理设备连接.物理设备的能力可以间接通过抽象设备来调用 

### API 网关 ###

接口服务注册到consul , 本模块支持横向扩展，注册到consul,通过kong 对外提供服务


## refrence ##

参考了edgex的官方demo
https://github.com/edgexfoundry/edgex-examples/tree/main/application-services/custom/send-command

参考了edgex的部分核心模块代码