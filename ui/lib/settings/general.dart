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
              "enable_plexmatch": v.enablePlexmatch,
              "allow_qiangban": v.allowQiangban,
              "enable_nfo": v.enableNfo,
              "enable_adult": v.enableAdult,
              "tv_naming_format": v.tvNamingFormat,
              "movie_naming_format": v.movieNamingFormat,
            },
            child: Column(
              crossAxisAlignment: CrossAxisAlignment.start,
              children: [
                FormBuilderTextField(
                  name: "tmdb_api",
                  decoration: const InputDecoration(
                      labelText: "TMDB Api Key", icon: Icon(Icons.key), helperText: "为防止被限流，可以提供自定义的 TMDB Api Key"),
                  //
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
                      icon: Icon(Icons.web),
                      hintText: "http://10.0.0.1:1080",
                      helperText: "后台联网代理地址，留空表示不启用代理"),
                ),
                FormBuilderTextField(
                  decoration: const InputDecoration(
                      icon: Icon(Icons.folder),
                      labelText: "电视剧路径命名规则",
                      helperText:
                          "go template语法，可用的变量为：.NameCN, .NameEN, .Year, .TmdbID"),
                  name: "tv_naming_format",
                ),
                FormBuilderTextField(
                  decoration: const InputDecoration(
                      icon: Icon(Icons.folder),
                      labelText: "电影路径命名规则",
                      helperText:
                          "go template语法，可用的变量为：.NameCN, .NameEN, .Year, .TmdbID"),
                  name: "movie_naming_format",
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
                      decoration:
                          const InputDecoration(icon: Icon(Icons.back_hand)),
                      name: "enable_adult",
                      title: const Text("是否显示成人内容")),
                ),
                SizedBox(
                  width: 300,
                  child: FormBuilderSwitch(
                      decoration:
                          const InputDecoration(icon: Icon(Icons.token)),
                      name: "enable_plexmatch",
                      title: const Text("Plex 刮削支持")),
                ),
                SizedBox(
                  width: 300,
                  child: FormBuilderSwitch(
                      decoration: const InputDecoration(
                          icon: Icon(Icons.library_books),
                          helperText: "emby/kodi等软件刮削需要"),
                      name: "enable_nfo",
                      title: const Text("nfo 文件支持")),
                ),
                SizedBox(
                  width: 300,
                  child: FormBuilderSwitch(
                      decoration: const InputDecoration(
                          icon: Icon(Icons.remove_circle)),
                      name: "allow_qiangban",
                      title: const Text("是否下载枪版资源")),
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
                                    allowQiangban: values["allow_qiangban"],
                                    enableAdult: values["enable_adult"],
                                    enableNfo: values["enable_nfo"],
                                    tvNamingFormat: values["tv_naming_format"],
                                    movieNamingFormat:
                                        values["movie_naming_format"],
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
        error: (err, trace) => PoNetworkError(err: err),
        loading: () => const MyProgressIndicator());
  }
}
