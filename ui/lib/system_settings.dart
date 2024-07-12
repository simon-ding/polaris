import 'package:flutter/material.dart';
import 'package:flutter_riverpod/flutter_riverpod.dart';
import 'package:quiver/strings.dart';
import 'package:ui/providers/login.dart';
import 'package:ui/providers/settings.dart';
import 'package:ui/utils.dart';

import 'providers/APIs.dart';

class SystemSettingsPage extends ConsumerStatefulWidget {
  static const route = "/settings";

  const SystemSettingsPage({super.key});
  @override
  ConsumerState<ConsumerStatefulWidget> createState() {
    return _SystemSettingsPageState();
  }
}

class _SystemSettingsPageState extends ConsumerState<SystemSettingsPage> {
  final GlobalKey _formKey = GlobalKey<FormState>();
  Future<void>? _pendingTmdb;
  Future<void>? _pendingIndexer;
  Future<void>? _pendingDownloadClient;
  Future<void>? _pendingStorage;
  final _tmdbApiController = TextEditingController();
  final _downloadDirController = TextEditingController();
  bool? _enableAuth;

  @override
  Widget build(BuildContext context) {
    var tmdbKey = ref.watch(settingProvider(APIs.tmdbApiKey));
    var dirKey = ref.watch(settingProvider(APIs.downloadDirKey));

    var tmdbSetting = FutureBuilder(
        // We listen to the pending operation, to update the UI accordingly.
        future: _pendingTmdb,
        builder: (context, snapshot) {
          return Container(
            padding: const EdgeInsets.fromLTRB(40, 10, 40, 0),
            child: Form(
              key: _formKey, //设置globalKey，用于后面获取FormState
              autovalidateMode: AutovalidateMode.onUserInteraction,
              child: Column(
                children: [
                  tmdbKey.when(
                      data: (value) {
                        _tmdbApiController.text = value;
                        return TextFormField(
                          autofocus: true,
                          controller: _tmdbApiController,
                          decoration: const InputDecoration(
                            labelText: "TMDB Api Key",
                            icon: Icon(Icons.key),
                          ),
                          //
                          validator: (v) {
                            return v!.trim().isNotEmpty ? null : "ApiKey 不能为空";
                          },
                          onSaved: (newValue) {},
                        );
                      },
                      error: (err, trace) => Text("$err"),
                      loading: () => const CircularProgressIndicator()),
                  dirKey.when(
                      data: (data) {
                        _downloadDirController.text = data;
                        return TextFormField(
                          autofocus: true,
                          controller: _downloadDirController,
                          decoration: const InputDecoration(
                            labelText: "下载路径",
                            icon: Icon(Icons.folder),
                          ),
                          //
                          validator: (v) {
                            return v!.trim().isNotEmpty ? null : "ApiKey 不能为空";
                          },
                          onSaved: (newValue) {},
                        );
                      },
                      error: (err, trace) => Text("$err"),
                      loading: () => const CircularProgressIndicator()),
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
                          if ((_formKey.currentState as FormState).validate()) {
                            var furture = ref
                                .read(settingProvider(APIs.tmdbApiKey).notifier)
                                .updateSettings(_tmdbApiController.text);
                            ref
                                .read(settingProvider(APIs.downloadDirKey)
                                    .notifier)
                                .updateSettings(_downloadDirController.text);
                            setState(() {
                              _pendingTmdb = furture;
                            });
                            showError(snapshot);
                          }
                        },
                      ),
                    ),
                  )
                ],
              ),
            ),
          );
        });

    var indexers = ref.watch(indexersProvider);
    var indexerSetting = FutureBuilder(
        // We listen to the pending operation, to update the UI accordingly.
        future: _pendingIndexer,
        builder: (context, snapshot) {
          return indexers.when(
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
                                showIndexerDetails(snapshot, context, indexer);
                              },
                              child: Center(child: Text(indexer.name!))));
                    }
                    return Card(
                        margin: const EdgeInsets.all(4),
                        clipBehavior: Clip.hardEdge,
                        child: InkWell(
                            //splashColor: Colors.blue.withAlpha(30),
                            onTap: () {
                              showIndexerDetails(snapshot, context, Indexer());
                            },
                            child: const Center(
                              child: Icon(Icons.add),
                            )));
                  }),
              error: (err, trace) => Text("$err"),
              loading: () => const CircularProgressIndicator());
        });

    var downloadClients = ref.watch(dwonloadClientsProvider);
    var downloadSetting = FutureBuilder(
        // We listen to the pending operation, to update the UI accordingly.
        future: _pendingDownloadClient,
        builder: (context, snapshot) {
          return downloadClients.when(
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
                                showDownloadClientDetails(
                                    snapshot, context, client);
                              },
                              child: Center(child: Text(client.name!))));
                    }
                    return Card(
                        margin: const EdgeInsets.all(4),
                        clipBehavior: Clip.hardEdge,
                        child: InkWell(
                            //splashColor: Colors.blue.withAlpha(30),
                            onTap: () {
                              showDownloadClientDetails(
                                  snapshot, context, DownloadClient());
                            },
                            child: const Center(
                              child: Icon(Icons.add),
                            )));
                  }),
              error: (err, trace) => Text("$err"),
              loading: () => const CircularProgressIndicator());
        });

    var storageSettingData = ref.watch(storageSettingProvider);
    var storageSetting = FutureBuilder(
        // We listen to the pending operation, to update the UI accordingly.
        future: _pendingStorage,
        builder: (context, snapshot) {
          return storageSettingData.when(
              data: (value) => GridView.builder(
                  itemCount: value.length + 1,
                  scrollDirection: Axis.vertical,
                  shrinkWrap: true,
                  gridDelegate: const SliverGridDelegateWithFixedCrossAxisCount(
                      crossAxisCount: 6),
                  itemBuilder: (context, i) {
                    if (i < value.length) {
                      var storage = value[i];
                      return Card(
                          margin: const EdgeInsets.all(4),
                          clipBehavior: Clip.hardEdge,
                          child: InkWell(
                              //splashColor: Colors.blue.withAlpha(30),
                              onTap: () {
                                showStorageDetails(snapshot, context, storage);
                              },
                              child: Center(child: Text(storage.name!))));
                    }
                    return Card(
                        margin: const EdgeInsets.all(4),
                        clipBehavior: Clip.hardEdge,
                        child: InkWell(
                            //splashColor: Colors.blue.withAlpha(30),
                            onTap: () {
                              showStorageDetails(snapshot, context, Storage());
                            },
                            child: const Center(
                              child: Icon(Icons.add),
                            )));
                  }),
              error: (err, trace) => Text("$err"),
              loading: () => const CircularProgressIndicator());
        });

    var authData = ref.watch(authSettingProvider);
    TextEditingController _userController = TextEditingController();
    TextEditingController _passController = TextEditingController();
    var authSetting = authData.when(
        data: (data) {
          if (_enableAuth == null) {
            setState(() {
              _enableAuth = data.enable;
            });
          }
          _userController.text = data.user;
          return Column(
            children: [
              SwitchListTile(
                  title: const Text("开启认证"),
                  value: _enableAuth!,
                  onChanged: (v) {
                    setState(() {
                      _enableAuth = v;
                    });
                  }),
              _enableAuth!
                  ? Column(
                      children: [
                        TextFormField(
                            controller: _userController,
                            decoration: const InputDecoration(
                              labelText: "用户名",
                              icon: Icon(Icons.verified_user),
                            )),
                        TextFormField(
                            controller: _passController,
                            decoration: const InputDecoration(
                              labelText: "密码",
                              icon: Icon(Icons.verified_user),
                            ))
                      ],
                    )
                  : const Column(),
              Center(
                  child: ElevatedButton(
                      child: const Text("保存"),
                      onPressed: () {
                        ref
                            .read(authSettingProvider.notifier)
                            .updateAuthSetting(_enableAuth!,
                                _userController.text, _passController.text);
                      }))
            ],
          );
        },
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
        ExpansionTile(
          tilePadding: const EdgeInsets.fromLTRB(10, 0, 10, 0),
          childrenPadding: const EdgeInsets.fromLTRB(50, 0, 50, 0),
          initiallyExpanded: true,
          title: const Text("存储设置"),
          children: [storageSetting],
        ),
        ExpansionTile(
          tilePadding: const EdgeInsets.fromLTRB(10, 0, 10, 0),
          childrenPadding: const EdgeInsets.fromLTRB(50, 0, 50, 0),
          initiallyExpanded: true,
          title: const Text("认证设置"),
          children: [authSetting],
        ),
      ],
    );
  }

  Future<void> showIndexerDetails(
      AsyncSnapshot<void> snapshot, BuildContext context, Indexer indexer) {
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
                  ? const Text("")
                  : TextButton(
                      onPressed: () {
                        var f = ref
                            .read(indexersProvider.notifier)
                            .deleteIndexer(indexer.id!);
                        setState(() {
                          _pendingIndexer = f;
                        });
                        if (!showError(snapshot)) {
                          Navigator.of(context).pop();
                        }
                      },
                      child: const Text('删除')),
              TextButton(
                  onPressed: () => Navigator.of(context).pop(),
                  child: const Text('取消')),
              TextButton(
                child: const Text('确定'),
                onPressed: () {
                  var f = ref.read(indexersProvider.notifier).addIndexer(
                      Indexer(
                          name: nameController.text,
                          url: urlController.text,
                          apiKey: apiKeyController.text));
                  setState(() {
                    _pendingIndexer = f;
                  });

                  if (!showError(snapshot)) {
                    Navigator.of(context).pop();
                  }
                },
              ),
            ],
          );
        });
  }

  Future<void> showDownloadClientDetails(AsyncSnapshot<void> snapshot,
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
                  ? const Text("")
                  : TextButton(
                      onPressed: () {
                        var f = ref
                            .read(dwonloadClientsProvider.notifier)
                            .deleteDownloadClients(client.id!);
                        setState(() {
                          _pendingDownloadClient = f;
                        });
                        if (!showError(snapshot)) {
                          Navigator.of(context).pop();
                        }
                      },
                      child: const Text('删除')),
              TextButton(
                  onPressed: () => Navigator.of(context).pop(),
                  child: const Text('取消')),
              TextButton(
                child: const Text('确定'),
                onPressed: () {
                  var f = ref
                      .read(dwonloadClientsProvider.notifier)
                      .addDownloadClients(
                          nameController.text, urlController.text);
                  setState(() {
                    _pendingDownloadClient = f;
                  });
                  if (!showError(snapshot)) {
                    Navigator.of(context).pop();
                  }
                },
              ),
            ],
          );
        });
  }

  Future<void> showStorageDetails(
      AsyncSnapshot<void> snapshot, BuildContext context, Storage s) {
    var nameController = TextEditingController(text: s.name);
    var implController = TextEditingController(
        text: isBlank(s.implementation) ? "transmission" : s.implementation);
    var pathController = TextEditingController(text: s.path);
    var userController = TextEditingController(text: s.user);
    var passController = TextEditingController(text: s.password);

    return showDialog<void>(
        context: context,
        barrierDismissible: true, // user must tap button!
        builder: (BuildContext context) {
          return AlertDialog(
            title: const Text('存储'),
            content: SingleChildScrollView(
              child: ListBody(
                children: <Widget>[
                  TextField(
                    decoration: const InputDecoration(labelText: "名称"),
                    controller: nameController,
                  ),
                  TextField(
                    decoration: const InputDecoration(labelText: "实现"),
                    controller: implController,
                  ),
                  TextField(
                    decoration: const InputDecoration(labelText: "路径"),
                    controller: pathController,
                  ),
                  TextField(
                    decoration: const InputDecoration(labelText: "用户"),
                    controller: userController,
                  ),
                  TextField(
                    decoration: const InputDecoration(labelText: "密码"),
                    controller: passController,
                  ),
                ],
              ),
            ),
            actions: <Widget>[
              s.id == null
                  ? const Text("")
                  : TextButton(
                      onPressed: () {
                        var f = ref
                            .read(storageSettingProvider.notifier)
                            .deleteStorage(s.id!);
                        setState(() {
                          _pendingStorage = f;
                        });
                        if (!showError(snapshot)) {
                          Navigator.of(context).pop();
                        }
                      },
                      child: const Text('删除')),
              TextButton(
                  onPressed: () => Navigator.of(context).pop(),
                  child: const Text('取消')),
              TextButton(
                child: const Text('确定'),
                onPressed: () {
                  var f = ref.read(storageSettingProvider.notifier).addStorage(
                      Storage(
                          name: nameController.text,
                          implementation: implController.text,
                          path: pathController.text,
                          user: userController.text,
                          password: passController.text));
                  setState(() {
                    _pendingStorage = f;
                  });
                  if (!showError(snapshot)) {
                    Navigator.of(context).pop();
                  }
                },
              ),
            ],
          );
        });
  }

  bool showError(AsyncSnapshot<void> snapshot) {
    final isErrored = snapshot.hasError &&
        snapshot.connectionState != ConnectionState.waiting;
    if (isErrored) {
      Utils.showSnakeBar(context, "当前操作出错: ${snapshot.error}");
      return true;
    }
    return false;
  }
}
