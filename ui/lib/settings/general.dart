import 'package:flutter/material.dart';
import 'package:flutter_form_builder/flutter_form_builder.dart';
import 'package:flutter_riverpod/flutter_riverpod.dart';
import 'package:form_builder_validators/form_builder_validators.dart';
import 'package:ui/providers/settings.dart';
import 'package:ui/widgets/progress_indicator.dart';
import 'package:ui/widgets/utils.dart';
import 'package:ui/widgets/widgets.dart';

class GeneralSettings extends ConsumerStatefulWidget {
  static const route = "/settings";

  const GeneralSettings({super.key});
  @override
  ConsumerState<ConsumerStatefulWidget> createState() {
    return _GeneralState();
  }
}

class _GeneralState extends ConsumerState<GeneralSettings> {
  final _formKey = GlobalKey<FormBuilderState>();

  @override
  Widget build(BuildContext context) {
    var settings = ref.watch(settingProvider);

    return settings.when(
        data: (v) {
          return FormBuilder(
            key: _formKey, //设置globalKey，用于后面获取FormState
            autovalidateMode: AutovalidateMode.onUserInteraction,
            initialValue: {
              "tmdb_api": v.tmdbApiKey,
              "download_dir": v.downloadDIr,
              "log_level": v.logLevel,
              "proxy": v.proxy,
              "enable_plexmatch": v.enablePlexmatch
            },
            child: Column(
              crossAxisAlignment: CrossAxisAlignment.start,
              children: [
                FormBuilderTextField(
                  name: "tmdb_api",
                  decoration: Commons.requiredTextFieldStyle(
                      text: "TMDB Api Key", icon: const Icon(Icons.key)),
                  //
                  validator: FormBuilderValidators.required(),
                ),
                FormBuilderTextField(
                  name: "download_dir",
                  decoration: Commons.requiredTextFieldStyle(
                      text: "下载路径",
                      icon: const Icon(Icons.folder),
                      helperText: "媒体文件临时下载路径，非最终存储路径"),
                  //
                  validator: FormBuilderValidators.required(),
                ),
                FormBuilderTextField(
                  name: "proxy",
                  decoration: const InputDecoration(
                      labelText: "代理地址",
                      icon: Icon(Icons.folder),
                      helperText: "后台联网代理地址，留空表示不启用代理"),
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
                      DropdownMenuItem(value: "debug", child: Text("DEBUG")),
                      DropdownMenuItem(value: "info", child: Text("INFO")),
                      DropdownMenuItem(value: "warn", child: Text("WARN")),
                      DropdownMenuItem(value: "error", child: Text("ERROR")),
                    ],
                    validator: FormBuilderValidators.required(),
                  ),
                ),
                SizedBox(
                  width: 300,
                  child: FormBuilderSwitch(
                      name: "enable_plexmatch", title: const Text("Plex 刮削支持")),
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
                                    logLevel: values["log_level"],
                                    proxy: values["proxy"],
                                    enablePlexmatch:
                                        values["enable_plexmatch"]))
                                .then((v) => showSnakeBar("更新成功"));
                            showLoadingWithFuture(f);
                          }
                        }),
                  ),
                )
              ],
            ),
          );
        },
        error: (err, trace) => Text("$err"),
        loading: () => const MyProgressIndicator());
  }
}
