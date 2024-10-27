# 配置

要正确使用此程序，需要配置好以下设置：

### 1. TMDB设置

1. 因为此程序需要使用到 TMDB 的数据，使用此程序首先要申请一个 TMDB 的 Api Key. 申请教程请 google [tmdb api key申请](https://www.google.com/search?q=tmdb+api+key%E7%94%B3%E8%AF%B7)

2. 拿到 TMDB Api Key之后，请填到 *设置 -> 常规设置 -> TMDB Api Key里*

**注意：** TMDB可能需要翻墙才能使用，参考 [TMDB 访问问题](./tmdb.md)

### 2. 索引器

使用配置页面索引器配置或者prowlarr设置，其中一个即可。

#### jackett配置参考 [jackett](./jackett.md)

#### prowlarr设置

1) 取得prowlarr的url和api key， api key在 *Prowlarr 设置 -> 通用 -> API 密钥* 处取得
2) 对应参数填到 polaris程序，*设置->prowlarr设置*当中

### 下载器

资源下载器，目前可支持 tansmission/qbittorrent，请配置好对应配置

![transmission](./assets/downloader.png)

### 存储设置

默认配置了名为 local 的本地存储，如果你不知道怎么配置。请使用默认配置

![local_storage](./assets/local_storage.png)

类型里选择 webdav 可以使用 webdav 存储，配合 alist/clouddrive 等，可以实现存储到云盘里的功能。

![webdav](./assets/webdav_storage.png)