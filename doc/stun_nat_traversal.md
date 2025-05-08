# 利用STUN进行NAT内网穿透

可以在下载器选项里打开 *使用内置STUN NAT穿透* 功能，即使处在NAT网络环境下，BT/PT也可以满速上传。打开后Polaris自动更改下载客户端的监听端口，并代理BT的上传流量。

要想正常使用此功能，需要具备以下几个条件：

1. 所在的NAT网络非对称NAT（Symmetric NAT），可以使用 [NatTypeTester](https://github.com/HMBSbige/NatTypeTester/releases/) 检查自己的网络的NAT类型
2. 下载器设置选项中下载器地址为下载器docker的实际地址，而非映射地址。达到这一目标可以使用host网络创建下载器，也可以利用docker-compose自带的域名解析来实现

