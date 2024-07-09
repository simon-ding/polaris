import 'package:dio/dio.dart';
import 'package:flutter/material.dart';
import 'package:flutter_riverpod/flutter_riverpod.dart';
import 'package:ui/APIs.dart';
import 'package:ui/providers/welcome_data.dart';
import 'package:ui/server_response.dart';
import 'package:ui/utils.dart';

class SystemSettingsPage extends ConsumerStatefulWidget {
  static const route = "/systemsettings";

  const SystemSettingsPage({super.key});
  @override
  ConsumerState<ConsumerStatefulWidget> createState() {
    return _SystemSettingsPageState();
  }
}

class _SystemSettingsPageState extends ConsumerState<SystemSettingsPage> {
  final GlobalKey _formKey = GlobalKey<FormState>();

  List<dynamic> indexers = List.empty();

  @override
  void initState() {
    super.initState();
  }

  @override
  Widget build(BuildContext context) {
    var key = ref.watch(tmdbApiSettingProvider);

    var tmdbSetting = key.when(
        data: (data) => Container(
              padding: const EdgeInsets.fromLTRB(40, 10, 40, 0),
              child: Form(
                key: _formKey, //设置globalKey，用于后面获取FormState
                autovalidateMode: AutovalidateMode.onUserInteraction,
                child: Column(
                  children: [
                    TextFormField(
                      autofocus: true,
                      initialValue: data,
                      decoration: const InputDecoration(
                        labelText: "TMDB Api Key",
                        icon: Icon(Icons.key),
                      ),
                      //
                      validator: (v) {
                        return v!.trim().isNotEmpty ? null : "ApiKey 不能为空";
                      },
                      onSaved: (newValue) {
                        _submitSettings(context, newValue!);
                      },
                    ),
                    Center(
                      child: Padding(
                        padding: const EdgeInsets.only(top: 28.0),
                        child: ElevatedButton(
                          child: const Padding(
                            padding: EdgeInsets.all(16.0),
                            child: Text("保存"),
                          ),
                          onPressed: () {
                            // 通过_formKey.currentState 获取FormState后，
                            // 调用validate()方法校验用户名密码是否合法，校验
                            // 通过后再提交数据。
                            if ((_formKey.currentState as FormState)
                                .validate()) {
                              (_formKey.currentState as FormState).save();
                            }
                          },
                        ),
                      ),
                    )
                  ],
                ),
              ),
            ),
        error: (err, trace) => Text("$err"),
        loading: () => const CircularProgressIndicator());

    var indexers = ref.watch(indexersProvider);
    var indexerSetting = indexers.when(
        data: (value) => GridView.builder(
            itemCount: value.length + 1,
            scrollDirection: Axis.vertical,
            shrinkWrap: true,
            gridDelegate: const SliverGridDelegateWithFixedCrossAxisCount(
                crossAxisCount: 6),
            itemBuilder: (context, i) {
              if (i < value.length) {
                var indexer = value[i];
                return Card(
                    margin: const EdgeInsets.all(4),
                    clipBehavior: Clip.hardEdge,
                    child: InkWell(
                        //splashColor: Colors.blue.withAlpha(30),
                        onTap: () {
                          showIndexerDetails(context, indexer);
                        },
                        child: Center(child: Text(indexer.name!))));
              }
              return Card(
                  margin: const EdgeInsets.all(4),
                  clipBehavior: Clip.hardEdge,
                  child: InkWell(
                      //splashColor: Colors.blue.withAlpha(30),
                      onTap: () {
                        showIndexerDetails(context, Indexer());
                      },
                      child: const Center(
                        child: Icon(Icons.add),
                      )));
            }),
        error: (err, trace) => Text("$err"),
        loading: () => const CircularProgressIndicator());

    var downloadClients = ref.watch(dwonloadClientsProvider);
    var downloadSetting = downloadClients.when(
        data: (value) => GridView.builder(
            itemCount: value.length + 1,
            scrollDirection: Axis.vertical,
            shrinkWrap: true,
            gridDelegate: const SliverGridDelegateWithFixedCrossAxisCount(
                crossAxisCount: 6),
            itemBuilder: (context, i) {
              if (i < value.length) {
                var client = value[i];
                return Card(
                    margin: const EdgeInsets.all(4),
                    clipBehavior: Clip.hardEdge,
                    child: InkWell(
                        //splashColor: Colors.blue.withAlpha(30),
                        onTap: () {
                          showDownloadClientDetails(context, client);
                        },
                        child: Center(child: Text(client.name!))));
              }
              return Card(
                  margin: const EdgeInsets.all(4),
                  clipBehavior: Clip.hardEdge,
                  child: InkWell(
                      //splashColor: Colors.blue.withAlpha(30),
                      onTap: () {
                        showDownloadClientDetails(context, DownloadClient());
                      },
                      child: const Center(
                        child: Icon(Icons.add),
                      )));
            }),
        error: (err, trace) => Text("$err"),
        loading: () => const CircularProgressIndicator());

    return ListView(
      children: [
        ExpansionTile(
          tilePadding: const EdgeInsets.fromLTRB(10, 0, 10, 0),
          childrenPadding: const EdgeInsets.fromLTRB(50, 0, 50, 0),
          initiallyExpanded: true,
          title: const Text("TMDB 设置"),
          children: [tmdbSetting],
        ),
        ExpansionTile(
          tilePadding: const EdgeInsets.fromLTRB(10, 0, 10, 0),
          childrenPadding: const EdgeInsets.fromLTRB(50, 0, 50, 0),
          initiallyExpanded: true,
          title: const Text("索引器设置"),
          children: [indexerSetting],
        ),
        ExpansionTile(
          tilePadding: const EdgeInsets.fromLTRB(10, 0, 10, 0),
          childrenPadding: const EdgeInsets.fromLTRB(50, 0, 50, 0),
          initiallyExpanded: true,
          title: const Text("下载客户端设置"),
          children: [downloadSetting],
        ),
      ],
    );
  }

  void _submitSettings(BuildContext context, String v) async {
    var resp = await Dio().post(APIs.settingsUrl, data: {APIs.tmdbApiKey: v});
    var sp = ServerResponse.fromJson(resp.data as Map<String, dynamic>);
    if (sp.code != 0) {
      if (context.mounted) {
        Utils.showAlertDialog(context, sp.message);
      }
    }
  }

  Future<void> showIndexerDetails(BuildContext context, Indexer indexer) {
    var nameController = TextEditingController(text: indexer.name);
    var urlController = TextEditingController(text: indexer.url);
    var apiKeyController = TextEditingController(text: indexer.apiKey);
    return showDialog<void>(
        context: context,
        barrierDismissible: true, // user must tap button!
        builder: (BuildContext context) {
          return AlertDialog(
            title: const Text('索引器'),
            content: SingleChildScrollView(
              child: ListBody(
                children: <Widget>[
                  TextField(
                    decoration: const InputDecoration(labelText: "名称"),
                    controller: nameController,
                  ),
                  TextField(
                    decoration: const InputDecoration(labelText: "网址"),
                    controller: urlController,
                  ),
                  TextField(
                    decoration: const InputDecoration(labelText: "API Key"),
                    controller: apiKeyController,
                  ),
                ],
              ),
            ),
            actions: <Widget>[
              indexer.id == null
                  ? Text("")
                  : TextButton(
                      onPressed: () => {deleteIndexer(context, indexer.id!)},
                      child: const Text('删除')),
              TextButton(
                  onPressed: () => Navigator.of(context).pop(),
                  child: const Text('取消')),
              TextButton(
                child: const Text('确定'),
                onPressed: () {
                  addIndexer(context, nameController.text, urlController.text,
                      apiKeyController.text);
                },
              ),
            ],
          );
        });
  }

  void addIndexer(
      BuildContext context, String name, String url, String apiKey) async {
    if (name.isEmpty || url.isEmpty || apiKey.isEmpty) {
      return;
    }
    var dio = Dio();
    var resp = await dio.post(APIs.addIndexerUrl,
        data: Indexer(name: name, url: url, apiKey: apiKey).toJson());
    var sp = ServerResponse.fromJson(resp.data);
    if (sp.code != 0 && context.mounted) {
      Utils.showAlertDialog(context, sp.message);
      return;
    }
    Navigator.of(context).pop();
    ref.refresh(indexersProvider);
  }

  void deleteIndexer(BuildContext context, int id) async {
    var dio = Dio();
    var resp = await dio.delete("${APIs.delIndexerUrl}$id");
    var sp = ServerResponse.fromJson(resp.data);
    if (sp.code != 0 && context.mounted) {
      Utils.showAlertDialog(context, sp.message);
      return;
    }
    Navigator.of(context).pop();
    ref.refresh(indexersProvider);
  }

  Future<void> showDownloadClientDetails(
      BuildContext context, DownloadClient client) {
    var nameController = TextEditingController(text: client.name);
    var urlController = TextEditingController(text: client.url);

    return showDialog<void>(
        context: context,
        barrierDismissible: true, // user must tap button!
        builder: (BuildContext context) {
          return AlertDialog(
            title: const Text('下载客户端'),
            content: SingleChildScrollView(
              child: ListBody(
                children: <Widget>[
                  TextField(
                    decoration: const InputDecoration(labelText: "名称"),
                    controller: nameController,
                  ),
                  TextField(
                    decoration: const InputDecoration(labelText: "网址"),
                    controller: urlController,
                  ),
                ],
              ),
            ),
            actions: <Widget>[
              client.id == null
                  ? Text("")
                  : TextButton(
                      onPressed: () =>
                          {deleteDownloadClients(context, client.id!)},
                      child: const Text('删除')),
              TextButton(
                  onPressed: () => Navigator.of(context).pop(),
                  child: const Text('取消')),
              TextButton(
                child: const Text('确定'),
                onPressed: () {
                  addDownloadClients(
                      context, nameController.text, urlController.text);
                },
              ),
            ],
          );
        });
  }

  void addDownloadClients(BuildContext context, String name, String url) async {
    if (name.isEmpty || url.isEmpty) {
      return;
    }
    var dio = Dio();
    var resp = await dio.post(APIs.addDownloadClientUrl, data: {
      "name": name,
      "url": url,
    });
    var sp = ServerResponse.fromJson(resp.data);
    if (sp.code != 0 && context.mounted) {
      Utils.showAlertDialog(context, sp.message);
      return;
    }
    Navigator.of(context).pop();
    ref.refresh(dwonloadClientsProvider);
  }

  void deleteDownloadClients(BuildContext context, int id) async {
    var dio = Dio();
    var resp = await dio.delete("${APIs.delDownloadClientUrl}$id");
    var sp = ServerResponse.fromJson(resp.data);
    if (sp.code != 0 && context.mounted) {
      Utils.showAlertDialog(context, sp.message);
      return;
    }
    Navigator.of(context).pop();
    ref.refresh(dwonloadClientsProvider);
  }
}
