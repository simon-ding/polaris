import 'package:flutter/material.dart';
import 'package:flutter_form_builder/flutter_form_builder.dart';
import 'package:flutter_riverpod/flutter_riverpod.dart';
import 'package:form_builder_validators/form_builder_validators.dart';
import 'package:quiver/strings.dart';
import 'package:ui/providers/settings.dart';
import 'package:ui/settings/dialog.dart';
import 'package:ui/widgets/progress_indicator.dart';
import 'package:ui/widgets/widgets.dart';

class DownloaderSettings extends ConsumerStatefulWidget {
  static const route = "/settings";

  const DownloaderSettings({super.key});
  @override
  ConsumerState<ConsumerStatefulWidget> createState() {
    return _DownloaderState();
  }
}

class _DownloaderState extends ConsumerState<DownloaderSettings> {
  @override
  Widget build(BuildContext context) {
    var downloadClients = ref.watch(dwonloadClientsProvider);
    return downloadClients.when(
        data: (value) => Wrap(
                children: List.generate(value.length + 1, (i) {
              if (i < value.length) {
                var client = value[i];
                return SettingsCard(
                    onTap: () => showDownloadClientDetails(client),
                    child: Text(client.name ?? ""));
              }
              return SettingsCard(
                  onTap: () => showSelections(), child: const Icon(Icons.add));
            })),
        error: (err, trace) => PoNetworkError(err: err),
        loading: () => const MyProgressIndicator());
  }

  Future<void> showDownloadClientDetails(DownloadClient client) {
    final _formKey = GlobalKey<FormBuilderState>();
    var _enableAuth = isNotBlank(client.user);

    final body =
        StatefulBuilder(builder: (BuildContext context, StateSetter setState) {
      return FormBuilder(
          key: _formKey,
          initialValue: {
            "name": client.name,
            "url": client.url,
            "user": client.user,
            "password": client.password,
            "remove_completed_downloads": client.removeCompletedDownloads,
            "remove_failed_downloads": client.removeFailedDownloads,
            "priority": client.priority.toString(),
          },
          child: Column(
            children: [
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
              FormBuilderTextField(
                  name: "priority",
                  decoration: const InputDecoration(
                      labelText: "优先级", helperText: "1-50, 1最高优先级，50最低优先级"),
                  validator: FormBuilderValidators.integer(),
                  autovalidateMode: AutovalidateMode.onUserInteraction),
              FormBuilderSwitch(
                  name: "remove_completed_downloads",
                  title: const Text("任务完成后删除")),
              FormBuilderSwitch(
                  name: "remove_failed_downloads",
                  title: const Text("任务失败后删除")),
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
                implementation: client.implementation,
                url: values["url"],
                user: _enableAuth ? values["user"] : null,
                password: _enableAuth ? values["password"] : null,
                priority: int.parse(values["priority"]),
                removeCompletedDownloads: values["remove_completed_downloads"],
                removeFailedDownloads: values["remove_failed_downloads"]));
      } else {
        throw "validation_error";
      }
    }

    var title = "下载器";
    if (client.implementation == "transmission") {
      title = "Transmission";
    } else if (client.implementation == "qbittorrent") {
      title = "qBittorrent";
    }

    return showSettingDialog(
        context, title, client.id != null, body, onSubmit, onDelete);
  }

  Future<void> showSelections() {
    return showDialog<void>(
        context: context,
        barrierDismissible: true,
        builder: (BuildContext context) {
          return AlertDialog(
            content: SizedBox(
              height: 500,
              width: 500,
              child: Wrap(
                children: [
                  SettingsCard(
                    child: InkWell(
                      child: const Center(
                        child: Text("Transmission"),
                      ),
                      onTap: () {
                        Navigator.of(context).pop();
                        showDownloadClientDetails(DownloadClient(
                            implementation: "transmission",
                            name: "Transmission"));
                      },
                    ),
                  ),
                  SettingsCard(
                    child: InkWell(
                      child: const Center(
                        child: Text("qBittorrent"),
                      ),
                      onTap: () {
                        Navigator.of(context).pop();
                        showDownloadClientDetails(DownloadClient(
                            implementation: "qbittorrent",
                            name: "qBittorrent"));
                      },
                    ),
                  )
                ],
              ),
            ),
          );
        });
  }
}
