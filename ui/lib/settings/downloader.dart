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
                  onTap: () => showDownloadClientDetails(DownloadClient()),
                  child: const Icon(Icons.add));
            })),
        error: (err, trace) => Text("$err"),
        loading: () => const MyProgressIndicator());
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
            "impl": client.implementation,
            "remove_completed_downloads": client.removeCompletedDownloads,
            "remove_failed_downloads": client.removeFailedDownloads,
            "priority": client.priority.toString(),
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
                  DropdownMenuItem(
                      value: "qbittorrent", child: Text("qBittorrent")),
                ],
                validator: FormBuilderValidators.required(),
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
              FormBuilderTextField(
                  name: "priority",
                  decoration: const InputDecoration(labelText: "优先级", helperText: "1-50, 1最高优先级，50最低优先级"),
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
                implementation: values["impl"],
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

    return showSettingDialog(
        context, "下载器", client.id != null, body, onSubmit, onDelete);
  }
}
