## jackett 索引器配置

1. 打开 jackett 主页后，点击页面上面的 Add indexer，会出现 BT/PT 站点列表，选择你需要的站点点击+号添加。如果是PT，请自行配置好相关配置

![add indexer](./assets/add_indexer.png)

![search add](./assets/search_add.png)


2. 添加后主页即会显示相应的BT/PT站点，点击 *Copy Torznab Feed* 即得到了我们需要的地址

![copy feed](./assets/copy_feed.png)

3. 回到我们的主程序 Polaris 当中，点击 *设置 -> 索引器设置* -> 点击+号增加新的索引器，输入一个名称，拷贝我们第2步得到的地址到地址栏

![polaris add indexer](./assets/polaris_add_indexer.png)

4. 选相框中我们可以看到，还需要一个 API Key，我们回到 Jackett 中，在页面右上角，复制我们需要的 API Key：
![api key](./assets/jackett_api_key.png)

5. 恭喜！你已经成功完成了索引器配置。如需要更多的站点，请重复相同的操作完成配置
