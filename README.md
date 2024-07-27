# polaris

Polaris 是一个电视剧和电影的追踪软件。配置好了之后，当剧集或者电影播出后，会第一时间下载对应的资源。支持本地存储或者webdav。

![main_page](./doc/assets/main_page.png)
![detail_page](./doc/assets/detail_page.png)

交流群： https://t.me/+8R2nzrlSs2JhMDgx

## 功能

- [x] 电视剧自动追踪下载
- [x] 电影自动追踪下载
- [x] webdav 存储支持，配合 [alist](https://github.com/alist-org/alist) 或阿里云等实现更多功能

## 使用

使用此程序参考 [【快速开始】](./doc/quick_start.md)


## 对比 sonarr/radarr
* 更好的中文支持
* 对于动漫、日剧的良好支持，配合国内站点基本能匹配上对应资源
* 支持 webdav 后端存储，可以配合 alist 或者阿里云来实现下载后实时传到云上的功能。这样外出就可以不依靠家里的宽带来看电影了，或者实现个轻 NAS 功能，下载功能放在本地，数据放在云盘
* golang 实现后端，相比于 .NET 更节省资源
* 一个程序同时实现了电影、电视剧功能，不需要装两个程序
* 当然 sonarr/radarr 也是非常优秀的开源项目，目前 Polaris 功能还没有 sonarr/radarr 丰富

-------------

## 请我喝杯咖啡

<img src="./doc/assets/wechat.JPG" width=40% height=40%>
