import 'package:flutter/material.dart';
import 'package:flutter_riverpod/flutter_riverpod.dart';
import 'package:ui/providers/login.dart';
import 'package:ui/providers/settings.dart';
import 'package:ui/utils.dart';
import 'package:ui/widgets/progress_indicator.dart';
import 'package:ui/widgets/widgets.dart';

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

  final _tmdbApiController = TextEditingController();
  final _downloadDirController = TextEditingController();
  bool? _enableAuth;

  @override
  Widget build(BuildContext context) {
    var settings = ref.watch(settingProvider);

    var tmdbSetting = settings.when(
        data: (v) {
          _tmdbApiController.text = v.tmdbApiKey!;
          _downloadDirController.text = v.downloadDIr!;
          return Container(
              padding: const EdgeInsets.fromLTRB(40, 10, 40, 0),
              child: Form(
                key: _formKey, //设置globalKey，用于后面获取FormState
                autovalidateMode: AutovalidateMode.onUserInteraction,
                child: Column(
                  children: [
                    TextFormField(
                      autofocus: true,
                      controller: _tmdbApiController,
                      decoration: Commons.requiredTextFieldStyle(
                          text: "TMDB Api Key", icon: const Icon(Icons.key)),
                      //
                      validator: (v) {
                        return v!.trim().isNotEmpty ? null : "ApiKey 不能为空";
                      },
                      onSaved: (newValue) {},
                    ),
                    TextFormField(
                      autofocus: true,
                      controller: _downloadDirController,
                      decoration: Commons.requiredTextFieldStyle(
                          text: "下载路径", icon: const Icon(Icons.folder)),
                      //
                      validator: (v) {
                        return v!.trim().isNotEmpty ? null : "下载路径不能为空";
                      },
                      onSaved: (newValue) {},
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
                            if ((_formKey.currentState as FormState)
                                .validate()) {
                              var f = ref
                                  .read(settingProvider.notifier)
                                  .updateSettings(GeneralSetting(
                                      tmdbApiKey: _tmdbApiController.text,
                                      downloadDIr:
                                          _downloadDirController.text));
                              f.whenComplete(() {
                                Utils.showSnakeBar("更新成功");
                              }).onError((e, s) {
                                Utils.showSnakeBar("更新失败：$e");
                              });
                            }
                          },
                        ),
                      ),
                    )
                  ],
                ),
              ));
        },
        error: (err, trace) => Text("$err"),
        loading: () => const MyProgressIndicator());

    var indexers = ref.watch(indexersProvider);
    var indexerSetting = indexers.when(
        data: (value) => Wrap(
              children: List.generate(value.length + 1, (i) {
                if (i < value.length) {
                  var indexer = value[i];
                  return SettingsCard(
                      onTap: () => showIndexerDetails(indexer),
                      child: Text(indexer.name!));
                }
                return SettingsCard(
                    onTap: () => showIndexerDetails(Indexer()),
                    child: const Icon(Icons.add));
              }),
            ),
        error: (err, trace) => Text("$err"),
        loading: () => const MyProgressIndicator());

    var downloadClients = ref.watch(dwonloadClientsProvider);
    var downloadSetting = downloadClients.when(
        data: (value) => Wrap(
                children: List.generate(value.length + 1, (i) {
              if (i < value.length) {
                var client = value[i];
                return SettingsCard(
                    onTap: () => showDownloadClientDetails(client),
                    child: Text(client.name!));
              }
              return SettingsCard(
                  onTap: () => showDownloadClientDetails(DownloadClient()),
                  child: const Icon(Icons.add));
            })),
        error: (err, trace) => Text("$err"),
        loading: () => const MyProgressIndicator());

    var storageSettingData = ref.watch(storageSettingProvider);
    var storageSetting = storageSettingData.when(
        data: (value) => Wrap(
              children: List.generate(value.length + 1, (i) {
                if (i < value.length) {
                  var storage = value[i];
                  return SettingsCard(
                      onTap: () => showStorageDetails(storage),
                      child: Text(storage.name!));
                }
                return SettingsCard(
                    onTap: () => showStorageDetails(Storage()),
                    child: const Icon(Icons.add));
              }),
            ),
        error: (err, trace) => Text("$err"),
        loading: () => const MyProgressIndicator());

    var authData = ref.watch(authSettingProvider);
    TextEditingController userController = TextEditingController();
    TextEditingController passController = TextEditingController();
    var authSetting = authData.when(
        data: (data) {
          if (_enableAuth == null) {
            setState(() {
              _enableAuth = data.enable;
            });
          }
          userController.text = data.user;
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
                            controller: userController,
                            decoration: Commons.requiredTextFieldStyle(
                              text: "用户名",
                              icon: const Icon(Icons.account_box),
                            )),
                        TextFormField(
                            obscureText: true,
                            enableSuggestions: false,
                            autocorrect: false,
                            controller: passController,
                            decoration: Commons.requiredTextFieldStyle(
                              text: "密码",
                              icon: const Icon(Icons.password),
                            ))
                      ],
                    )
                  : const Column(),
              Center(
                  child: ElevatedButton(
                      child: const Text("保存"),
                      onPressed: () {
                        var f = ref
                            .read(authSettingProvider.notifier)
                            .updateAuthSetting(_enableAuth!,
                                userController.text, passController.text);
                        f.whenComplete(() {
                          Utils.showSnakeBar("更新成功");
                        }).onError((e, s) {
                          Utils.showSnakeBar("更新失败：$e");
                        });
                      }))
            ],
          );
        },
        error: (err, trace) => Text("$err"),
        loading: () => const MyProgressIndicator());

    return ListView(
      children: [
        ExpansionTile(
          expandedAlignment: Alignment.centerLeft,
          tilePadding: const EdgeInsets.fromLTRB(10, 0, 10, 0),
          childrenPadding: const EdgeInsets.fromLTRB(50, 0, 50, 0),
          initiallyExpanded: true,
          title: const Text("常规设置"),
          children: [tmdbSetting],
        ),
        ExpansionTile(
          expandedAlignment: Alignment.centerLeft,
          tilePadding: const EdgeInsets.fromLTRB(10, 0, 10, 0),
          childrenPadding: const EdgeInsets.fromLTRB(50, 0, 50, 0),
          initiallyExpanded: false,
          title: const Text("索引器设置"),
          children: [indexerSetting],
        ),
        ExpansionTile(
          expandedAlignment: Alignment.centerLeft,
          tilePadding: const EdgeInsets.fromLTRB(10, 0, 10, 0),
          childrenPadding: const EdgeInsets.fromLTRB(50, 0, 50, 0),
          initiallyExpanded: false,
          title: const Text("下载客户端设置"),
          children: [downloadSetting],
        ),
        ExpansionTile(
          expandedAlignment: Alignment.centerLeft,
          tilePadding: const EdgeInsets.fromLTRB(10, 0, 10, 0),
          childrenPadding: const EdgeInsets.fromLTRB(50, 0, 50, 0),
          initiallyExpanded: false,
          title: const Text("存储设置"),
          children: [storageSetting],
        ),
        ExpansionTile(
          tilePadding: const EdgeInsets.fromLTRB(10, 0, 10, 0),
          childrenPadding: const EdgeInsets.fromLTRB(50, 0, 50, 0),
          initiallyExpanded: false,
          title: const Text("认证设置"),
          children: [authSetting],
        ),
      ],
    );
  }

  Future<void> showIndexerDetails(Indexer indexer) {
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
                    decoration: Commons.requiredTextFieldStyle(text: "名称"),
                    controller: nameController,
                  ),
                  TextField(
                    decoration: Commons.requiredTextFieldStyle(text: "地址"),
                    controller: urlController,
                  ),
                  TextField(
                    decoration: Commons.requiredTextFieldStyle(text: "API Key"),
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
                        f.whenComplete(() {
                          Utils.showSnakeBar("操作成功");
                          Navigator.of(context).pop();
                        }).onError((e, s) {
                          Utils.showSnakeBar("操作失败：$e");
                        });
                      },
                      child: const Text(
                        '删除',
                        style: TextStyle(color: Colors.red),
                      )),
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
                  f.whenComplete(() {
                    Utils.showSnakeBar("操作成功");
                    Navigator.of(context).pop();
                  }).onError((e, s) {
                    Utils.showSnakeBar("操作失败：$e");
                  });
                },
              ),
            ],
          );
        });
  }

  Future<void> showDownloadClientDetails(DownloadClient client) {
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
                    decoration: Commons.requiredTextFieldStyle(text: "名称"),
                    controller: nameController,
                  ),
                  TextField(
                    decoration: Commons.requiredTextFieldStyle(text: "地址"),
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
                        f.whenComplete(() {
                          Utils.showSnakeBar("操作成功");
                          Navigator.of(context).pop();
                        }).onError((e, s) {
                          Utils.showSnakeBar("操作失败：$e");
                        });
                      },
                      child: const Text(
                        '删除',
                        style: TextStyle(color: Colors.red),
                      )),
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
                  f.whenComplete(() {
                    Utils.showSnakeBar("操作成功");
                    Navigator.of(context).pop();
                  }).onError((e, s) {
                    Utils.showSnakeBar("操作失败：$e");
                  });
                },
              ),
            ],
          );
        });
  }

  Future<void> showStorageDetails(Storage s) {
    return showDialog<void>(
        context: context,
        barrierDismissible: true, // user must tap button!
        builder: (BuildContext context) {
          var nameController = TextEditingController(text: s.name);
          var tvPathController = TextEditingController();
          var moviePathController = TextEditingController();
          var urlController = TextEditingController();
          var userController = TextEditingController();
          var passController = TextEditingController();
          if (s.settings != null) {
            tvPathController.text = s.settings!["tv_path"] ?? "";
            moviePathController.text = s.settings!["movie_path"] ?? "";
            urlController.text = s.settings!["url"] ?? "";
            userController.text = s.settings!["user"] ?? "";
            passController.text = s.settings!["password"] ?? "";
          }

          String selectImpl =
              s.implementation == null ? "local" : s.implementation!;
          return StatefulBuilder(builder: (context, setState) {
            return AlertDialog(
              title: const Text('存储'),
              content: SingleChildScrollView(
                child: ListBody(
                  children: <Widget>[
                    DropdownMenu(
                      label: const Text("实现"),
                      onSelected: (value) {
                        setState(() {
                          selectImpl = value!;
                        });
                      },
                      initialSelection: selectImpl,
                      dropdownMenuEntries: const [
                        DropdownMenuEntry(value: "local", label: "本地存储"),
                        DropdownMenuEntry(value: "webdav", label: "webdav")
                      ],
                    ),
                    TextField(
                      decoration: const InputDecoration(labelText: "名称"),
                      controller: nameController,
                    ),
                    selectImpl != "local"
                        ? Column(
                            children: [
                              TextField(
                                decoration: const InputDecoration(
                                    labelText: "Webdav地址"),
                                controller: urlController,
                              ),
                              TextField(
                                decoration:
                                    const InputDecoration(labelText: "用户"),
                                controller: userController,
                              ),
                              TextField(
                                decoration:
                                    const InputDecoration(labelText: "密码"),
                                controller: passController,
                              ),
                            ],
                          )
                        : Container(),
                    TextField(
                      decoration: const InputDecoration(labelText: "电视剧路径"),
                      controller: tvPathController,
                    ),
                    TextField(
                      decoration: const InputDecoration(labelText: "电影路径"),
                      controller: moviePathController,
                    )
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
                          f.whenComplete(() {
                            Utils.showSnakeBar("操作成功");
                            Navigator.of(context).pop();
                          }).onError((e, s) {
                            Utils.showSnakeBar("操作失败：$e");
                          });
                        },
                        child: const Text('删除')),
                TextButton(
                    onPressed: () => Navigator.of(context).pop(),
                    child: const Text('取消')),
                TextButton(
                  child: const Text('确定'),
                  onPressed: () async {
                    final f = ref
                        .read(storageSettingProvider.notifier)
                        .addStorage(Storage(
                          name: nameController.text,
                          implementation: selectImpl,
                          settings: {
                            "tv_path": tvPathController.text,
                            "movie_path": moviePathController.text,
                            "url": urlController.text,
                            "user": userController.text,
                            "password": passController.text
                          },
                        ));
                    f.whenComplete(() {
                      Utils.showSnakeBar("操作成功");
                      Navigator.of(context).pop();
                    }).onError((e, s) {
                      Utils.showSnakeBar("操作失败：$e");
                    });
                  },
                ),
              ],
            );
          });
        });
  }

  bool showError(AsyncSnapshot<void> snapshot) {
    final isErrored = snapshot.hasError &&
        snapshot.connectionState != ConnectionState.waiting;
    if (isErrored) {
      Utils.showSnakeBar("当前操作出错: ${snapshot.error}");
      return true;
    }
    return false;
  }
}
