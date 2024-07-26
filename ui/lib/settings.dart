import 'package:flutter/material.dart';
import 'package:flutter_form_builder/flutter_form_builder.dart';
import 'package:flutter_riverpod/flutter_riverpod.dart';
import 'package:quiver/strings.dart';
import 'package:ui/providers/login.dart';
import 'package:ui/providers/settings.dart';
import 'package:ui/utils.dart';
import 'package:ui/widgets/progress_indicator.dart';
import 'package:ui/widgets/widgets.dart';
import 'package:form_builder_validators/form_builder_validators.dart';

class SystemSettingsPage extends ConsumerStatefulWidget {
  static const route = "/settings";

  const SystemSettingsPage({super.key});
  @override
  ConsumerState<ConsumerStatefulWidget> createState() {
    return _SystemSettingsPageState();
  }
}

class _SystemSettingsPageState extends ConsumerState<SystemSettingsPage> {
  final _formKey = GlobalKey<FormBuilderState>();
  final _formKey2 = GlobalKey<FormBuilderState>();
  bool? _enableAuth;

  @override
  Widget build(BuildContext context) {
    var settings = ref.watch(settingProvider);

    var tmdbSetting = settings.when(
        data: (v) {
          return Container(
              padding: const EdgeInsets.fromLTRB(40, 10, 40, 0),
              child: FormBuilder(
                key: _formKey, //设置globalKey，用于后面获取FormState
                autovalidateMode: AutovalidateMode.onUserInteraction,
                initialValue: {
                  "tmdb_api": v.tmdbApiKey,
                  "download_dir": v.downloadDIr,
                  "log_level": v.logLevel
                },
                child: Column(
                  crossAxisAlignment: CrossAxisAlignment.start,
                  children: [
                    FormBuilderTextField(
                      name: "tmdb_api",
                      autofocus: true,
                      decoration: Commons.requiredTextFieldStyle(
                          text: "TMDB Api Key", icon: const Icon(Icons.key)),
                      //
                      validator: FormBuilderValidators.required(),
                    ),
                    FormBuilderTextField(
                      name: "download_dir",
                      autofocus: true,
                      decoration: Commons.requiredTextFieldStyle(
                          text: "下载路径", icon: const Icon(Icons.folder)),
                      //
                      validator: FormBuilderValidators.required(),
                    ),
                    SizedBox(
                      width: 300,
                      child: FormBuilderDropdown(
                        name: "log_level",
                        decoration: const InputDecoration(
                          labelText: "日志级别",
                          icon: Icon(Icons.file_present_rounded),
                        ),
                        items: const [
                          DropdownMenuItem(
                              value: "debug", child: Text("DEBUG")),
                          DropdownMenuItem(value: "info", child: Text("INFO")),
                          DropdownMenuItem(
                              value: "warn", child: Text("WARNNING")),
                          DropdownMenuItem(
                              value: "error", child: Text("ERROR")),
                        ],
                        validator: FormBuilderValidators.required(),
                      ),
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
                              if (_formKey.currentState!.saveAndValidate()) {
                                var values = _formKey.currentState!.value;
                                var f = ref
                                    .read(settingProvider.notifier)
                                    .updateSettings(GeneralSetting(
                                        tmdbApiKey: values["tmdb_api"],
                                        downloadDIr: values["download_dir"],
                                        logLevel: values["log_level"]));
                                f.then((v) {
                                  Utils.showSnakeBar("更新成功");
                                }).onError((e, s) {
                                  Utils.showSnakeBar("更新失败：$e");
                                });
                              }
                            }),
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
                      child: Text(indexer.name ?? ""));
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
                    child: Text(client.name ?? ""));
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
                      child: Text(storage.name ?? ""));
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
          return FormBuilder(
              key: _formKey2,
              initialValue: {
                "user": data.user,
                "password": data.password,
                "enable": data.enable
              },
              child: Column(
                children: [
                  FormBuilderSwitch(
                      name: "enable",
                      title: const Text("开启认证"),
                      onChanged: (v) {
                        setState(() {
                          _enableAuth = v;
                        });
                      }),
                  _enableAuth!
                      ? Column(
                          children: [
                            FormBuilderTextField(
                                name: "user",
                                autovalidateMode:
                                    AutovalidateMode.onUserInteraction,
                                validator: FormBuilderValidators.required(),
                                decoration: Commons.requiredTextFieldStyle(
                                  text: "用户名",
                                  icon: const Icon(Icons.account_box),
                                )),
                            FormBuilderTextField(
                                name: "password",
                                obscureText: true,
                                enableSuggestions: false,
                                autocorrect: false,
                                autovalidateMode:
                                    AutovalidateMode.onUserInteraction,
                                validator: FormBuilderValidators.required(),
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
                            if (_formKey2.currentState!.saveAndValidate()) {
                              var values = _formKey2.currentState!.value;
                              var f = ref
                                  .read(authSettingProvider.notifier)
                                  .updateAuthSetting(_enableAuth!,
                                      values["user"], values["password"]);
                              f.then((v) {
                                Utils.showSnakeBar("更新成功");
                              }).onError((e, s) {
                                Utils.showSnakeBar("更新失败：$e");
                              });
                            }
                          }))
                ],
              ));
        },
        error: (err, trace) => Text("$err"),
        loading: () => const MyProgressIndicator());

    return ListView(
      children: [
        ExpansionTile(
          expandedAlignment: Alignment.centerLeft,
          childrenPadding: const EdgeInsets.fromLTRB(20, 0, 20, 0),
          initiallyExpanded: true,
          title: const Text("常规设置"),
          children: [tmdbSetting],
        ),
        ExpansionTile(
          expandedAlignment: Alignment.centerLeft,
          childrenPadding: const EdgeInsets.fromLTRB(20, 0, 20, 0),
          initiallyExpanded: false,
          title: const Text("索引器设置"),
          children: [indexerSetting],
        ),
        ExpansionTile(
          expandedAlignment: Alignment.centerLeft,
          childrenPadding: const EdgeInsets.fromLTRB(20, 0, 20, 0),
          initiallyExpanded: false,
          title: const Text("下载器设置"),
          children: [downloadSetting],
        ),
        ExpansionTile(
          expandedAlignment: Alignment.centerLeft,
          childrenPadding: const EdgeInsets.fromLTRB(20, 0, 50, 0),
          initiallyExpanded: false,
          title: const Text("存储设置"),
          children: [storageSetting],
        ),
        ExpansionTile(
          childrenPadding: const EdgeInsets.fromLTRB(20, 0, 20, 0),
          initiallyExpanded: false,
          title: const Text("认证设置"),
          children: [authSetting],
        ),
      ],
    );
  }

  Future<void> showIndexerDetails(Indexer indexer) {
    final _formKey = GlobalKey<FormBuilderState>();
    var selectImpl = "torznab";
    var body = FormBuilder(
      key: _formKey,
      initialValue: {
        "name": indexer.name,
        "url": indexer.url,
        "api_key": indexer.apiKey,
        "impl": "torznab"
      },
      child: Column(
        children: [
          FormBuilderDropdown(
            name: "impl",
            decoration: const InputDecoration(labelText: "类型"),
            items: const [
              DropdownMenuItem(value: "torznab", child: Text("Torznab")),
            ],
          ),
          FormBuilderTextField(
            name: "name",
            decoration: Commons.requiredTextFieldStyle(text: "名称"),
            autovalidateMode: AutovalidateMode.onUserInteraction,
            validator: FormBuilderValidators.required(),
          ),
          FormBuilderTextField(
            name: "url",
            decoration: Commons.requiredTextFieldStyle(text: "地址"),
            autovalidateMode: AutovalidateMode.onUserInteraction,
            validator: FormBuilderValidators.required(),
          ),
          FormBuilderTextField(
            name: "api_key",
            decoration: Commons.requiredTextFieldStyle(text: "API Key"),
            autovalidateMode: AutovalidateMode.onUserInteraction,
            validator: FormBuilderValidators.required(),
          ),
        ],
      ),
    );
    onDelete() async {
      return ref.read(indexersProvider.notifier).deleteIndexer(indexer.id!);
    }

    onSubmit() async {
      if (_formKey.currentState!.saveAndValidate()) {
        var values = _formKey.currentState!.value;
        return ref.read(indexersProvider.notifier).addIndexer(Indexer(
            name: values["name"],
            url: values["url"],
            apiKey: values["api_key"]));
      } else {
        throw "validation_error";
      }
    }

    return showSettingDialog(
        "索引器", indexer.id != null, body, onSubmit, onDelete);
  }

  Future<void> showDownloadClientDetails(DownloadClient client) {
    final _formKey = GlobalKey<FormBuilderState>();
    var _enableAuth = isNotBlank(client.user);
    String selectImpl = "transmission";

    final body =
        StatefulBuilder(builder: (BuildContext context, StateSetter setState) {
      return FormBuilder(
          key: _formKey,
          initialValue: {
            "name": client.name,
            "url": client.url,
            "user": client.user,
            "password": client.password,
            "impl": "transmission"
          },
          child: Column(
            children: [
              FormBuilderDropdown<String>(
                name: "impl",
                decoration: const InputDecoration(labelText: "类型"),
                onChanged: (value) {
                  setState(() {
                    selectImpl = value!;
                  });
                },
                items: const [
                  DropdownMenuItem(
                      value: "transmission", child: Text("Transmission")),
                ],
              ),
              FormBuilderTextField(
                  name: "name",
                  decoration: const InputDecoration(labelText: "名称"),
                  validator: FormBuilderValidators.required(),
                  autovalidateMode: AutovalidateMode.onUserInteraction),
              FormBuilderTextField(
                name: "url",
                decoration: const InputDecoration(
                    labelText: "地址", hintText: "http://127.0.0.1:9091"),
                autovalidateMode: AutovalidateMode.onUserInteraction,
                validator: FormBuilderValidators.required(),
              ),
              StatefulBuilder(
                  builder: (BuildContext context, StateSetter setState) {
                return Column(
                  children: [
                    FormBuilderSwitch(
                        name: "auth",
                        title: const Text("需要认证"),
                        initialValue: _enableAuth,
                        onChanged: (v) {
                          setState(() {
                            _enableAuth = v!;
                          });
                        }),
                    _enableAuth
                        ? Column(
                            children: [
                              FormBuilderTextField(
                                  name: "user",
                                  decoration: Commons.requiredTextFieldStyle(
                                      text: "用户"),
                                  validator: FormBuilderValidators.required(),
                                  autovalidateMode:
                                      AutovalidateMode.onUserInteraction),
                              FormBuilderTextField(
                                  name: "password",
                                  decoration: Commons.requiredTextFieldStyle(
                                      text: "密码"),
                                  validator: FormBuilderValidators.required(),
                                  obscureText: true,
                                  autovalidateMode:
                                      AutovalidateMode.onUserInteraction),
                            ],
                          )
                        : Container()
                  ],
                );
              })
            ],
          ));
    });
    onDelete() async {
      return ref
          .read(dwonloadClientsProvider.notifier)
          .deleteDownloadClients(client.id!);
    }

    onSubmit() async {
      if (_formKey.currentState!.saveAndValidate()) {
        var values = _formKey.currentState!.value;
        return ref.read(dwonloadClientsProvider.notifier).addDownloadClients(
            DownloadClient(
                name: values["name"],
                implementation: values["impl"],
                url: values["url"],
                user: _enableAuth ? values["user"] : null,
                password: _enableAuth ? values["password"] : null));
      } else {
        throw "validation_error";
      }
    }

    return showSettingDialog(
        "下载器", client.id != null, body, onSubmit, onDelete);
  }

  Future<void> showStorageDetails(Storage s) {
    final _formKey = GlobalKey<FormBuilderState>();

    String selectImpl = s.implementation == null ? "local" : s.implementation!;
    final widgets =
        StatefulBuilder(builder: (BuildContext context, StateSetter setState) {
      return FormBuilder(
          key: _formKey,
          autovalidateMode: AutovalidateMode.disabled,
          initialValue: {
            "name": s.name,
            "impl": s.implementation == null ? "local" : s.implementation!,
            "user": s.settings != null ? s.settings!["user"] ?? "" : "",
            "password": s.settings != null ? s.settings!["password"] ?? "" : "",
            "tv_path": s.settings != null ? s.settings!["tv_path"] ?? "" : "",
            "url": s.settings != null ? s.settings!["url"] ?? "" : "",
            "movie_path":
                s.settings != null ? s.settings!["movie_path"] ?? "" : "",
            "change_file_hash": s.settings != null
                ? s.settings!["change_file_hash"] == "true"
                    ? true
                    : false
                : false,
          },
          child: Column(
            crossAxisAlignment: CrossAxisAlignment.start,
            children: <Widget>[
              FormBuilderDropdown<String>(
                name: "impl",
                autovalidateMode: AutovalidateMode.onUserInteraction,
                decoration: const InputDecoration(labelText: "类型"),
                onChanged: (value) {
                  setState(() {
                    selectImpl = value!;
                  });
                },
                items: const [
                  DropdownMenuItem(
                    value: "local",
                    child: Text("本地存储"),
                  ),
                  DropdownMenuItem(
                    value: "webdav",
                    child: Text("webdav"),
                  )
                ],
                validator: FormBuilderValidators.required(),
              ),
              FormBuilderTextField(
                name: "name",
                autovalidateMode: AutovalidateMode.onUserInteraction,
                initialValue: s.name,
                decoration: const InputDecoration(labelText: "名称"),
                validator: FormBuilderValidators.required(),
              ),
              selectImpl != "local"
                  ? Column(
                      crossAxisAlignment: CrossAxisAlignment.start,
                      children: [
                        FormBuilderTextField(
                          name: "url",
                          autovalidateMode: AutovalidateMode.onUserInteraction,
                          decoration:
                              const InputDecoration(labelText: "Webdav地址"),
                          validator: FormBuilderValidators.required(),
                        ),
                        FormBuilderTextField(
                          name: "user",
                          autovalidateMode: AutovalidateMode.onUserInteraction,
                          decoration: const InputDecoration(labelText: "用户"),
                        ),
                        FormBuilderTextField(
                          name: "password",
                          autovalidateMode: AutovalidateMode.onUserInteraction,
                          decoration: const InputDecoration(labelText: "密码"),
                          obscureText: true,
                        ),
                        FormBuilderCheckbox(
                          name: "change_file_hash",
                          title: const Text(
                            "上传时更改文件哈希",
                            style: TextStyle(fontSize: 14),
                          ),
                        ),
                      ],
                    )
                  : Container(),
              FormBuilderTextField(
                name: "tv_path",
                autovalidateMode: AutovalidateMode.onUserInteraction,
                decoration: const InputDecoration(labelText: "电视剧路径"),
                validator: FormBuilderValidators.required(),
              ),
              FormBuilderTextField(
                name: "movie_path",
                autovalidateMode: AutovalidateMode.onUserInteraction,
                decoration: const InputDecoration(labelText: "电影路径"),
                validator: FormBuilderValidators.required(),
              )
            ],
          ));
    });
    onSubmit() async {
      if (_formKey.currentState!.saveAndValidate()) {
        final values = _formKey.currentState!.value;
        return ref.read(storageSettingProvider.notifier).addStorage(Storage(
              name: values["name"],
              implementation: selectImpl,
              settings: {
                "tv_path": values["tv_path"],
                "movie_path": values["movie_path"],
                "url": values["url"],
                "user": values["user"],
                "password": values["password"],
                "change_file_hash":
                    values["change_file_hash"] as bool ? "true" : "false"
              },
            ));
      } else {
        throw "validation_error";
      }
    }

    onDelete() async {
      return ref.read(storageSettingProvider.notifier).deleteStorage(s.id!);
    }

    return showSettingDialog('存储', s.id != null, widgets, onSubmit, onDelete);
  }

  Future<void> showSettingDialog(String title, bool showDelete, Widget body,
      Future Function() onSubmit, Future Function() onDelete) {
    return showDialog<void>(
        context: context,
        barrierDismissible: true,
        builder: (BuildContext context) {
          return AlertDialog(
            title: Text(title),
            content: SingleChildScrollView(
              child: Container(
                constraints: const BoxConstraints(maxWidth: 200),
                child: body,
              ),
            ),
            actions: <Widget>[
              showDelete
                  ? TextButton(
                      onPressed: () {
                        final f = onDelete();
                        f.then((v) {
                          Utils.showSnakeBar("删除成功");
                          Navigator.of(context).pop();
                        }).onError((e, s) {
                          Utils.showSnakeBar("删除失败：$e");
                        });
                      },
                      child: const Text(
                        '删除',
                        style: TextStyle(color: Colors.red),
                      ))
                  : const Text(""),
              TextButton(
                  onPressed: () => Navigator.of(context).pop(),
                  child: const Text('取消')),
              TextButton(
                child: const Text('确定'),
                onPressed: () {
                  final f = onSubmit();
                  f.then((v) {
                    Utils.showSnakeBar("操作成功");
                    Navigator.of(context).pop();
                  }).onError((e, s) {
                    if (e.toString() != "validation_error") {
                      Utils.showSnakeBar("操作失败：$e");
                    }
                  });
                },
              ),
            ],
          );
        });
  }
}
