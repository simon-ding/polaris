import 'package:flutter/material.dart';
import 'package:flutter_form_builder/flutter_form_builder.dart';
import 'package:flutter_riverpod/flutter_riverpod.dart';
import 'package:form_builder_validators/form_builder_validators.dart';
import 'package:quiver/strings.dart';
import 'package:ui/providers/settings.dart';
import 'package:ui/providers/size_limiter.dart';
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
    return Column(
      crossAxisAlignment: CrossAxisAlignment.start,
      children: [
        downloadClients.when(
            data: (value) => Wrap(
                    children: List.generate(value.length + 1, (i) {
                  if (i < value.length) {
                    var client = value[i];
                    return SettingsCard(
                        onTap: () => showDownloadClientDetails(client),
                        child: Text(client.name ?? ""));
                  }
                  return SettingsCard(
                      onTap: () => showSelections(),
                      child: const Icon(Icons.add));
                })),
            error: (err, trace) => PoNetworkError(err: err),
            loading: () => const MyProgressIndicator()),
        Divider(),
        getSizeLimiterWidget()
      ],
    );
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
            "use_nat_traversal": client.useNatTraversal,
          },
          child: Column(
            children: [
              FormBuilderTextField(
                  name: "name",
                  enabled: client.implementation != "buildin",
                  decoration: const InputDecoration(labelText: "名称"),
                  validator: FormBuilderValidators.required(),
                  autovalidateMode: AutovalidateMode.onUserInteraction),
              FormBuilderTextField(
                name: "url",
                enabled: client.implementation != "buildin",
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
                  name: "use_nat_traversal",
                  enabled: client.implementation == "qbittorrent",
                  title: const Text("使用内置STUN NAT穿透"),
                  decoration: InputDecoration(helperText: "内建的NAT穿透功能帮助BT客户端上传(会自动更改下载器的监听地址)"),
                  ),
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
                        enabled: client.implementation != "buildin",
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
                useNatTraversal: values["use_nat_traversal"],
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
        context,
        title,
        client.idExists() && client.implementation != "buildin",
        body,
        onSubmit,
        onDelete);
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

  Widget getSizeLimiterWidget() {
    var data = ref.watch(mediaSizeLimiterDataProvider);
    final _formKey = GlobalKey<FormBuilderState>();

    return Container(
      padding: EdgeInsets.only(left: 20, right: 20, top: 20),
      child: data.when(
          data: (value) {
            return FormBuilder(
              key: _formKey,
              initialValue: {
                "tv_720p_min": toMbString(value.tvLimiter!.p720p!.minSize!),
                "tv_720p_max": toMbString(value.tvLimiter!.p720p!.maxSize!),
                "tv_1080p_min": toMbString(value.tvLimiter!.p1080p!.minSize!),
                "tv_1080p_max": toMbString(value.tvLimiter!.p1080p!.maxSize!),
                "tv_2160p_min": toMbString(value.tvLimiter!.p2160p!.minSize!),
                "tv_2160p_max": toMbString(value.tvLimiter!.p2160p!.maxSize!),
                "movie_720p_min":
                    toMbString(value.movieLimiter!.p720p!.minSize!),
                "movie_720p_max":
                    toMbString(value.movieLimiter!.p720p!.maxSize!),
                "movie_1080p_min":
                    toMbString(value.movieLimiter!.p1080p!.minSize!),
                "movie_1080p_max":
                    toMbString(value.movieLimiter!.p1080p!.maxSize!),
                "movie_2160p_min":
                    toMbString(value.movieLimiter!.p2160p!.minSize!),
                "movie_2160p_max":
                    toMbString(value.movieLimiter!.p2160p!.maxSize!),
              },
              child: Column(
                crossAxisAlignment: CrossAxisAlignment.start,
                children: [
                  Text(
                    "剧集大小限制",
                    style: TextStyle(fontSize: 18),
                  ),
                  Divider(),
                  minMaxRow("  720p", "tv_720p_min", "tv_720p_max"),
                  minMaxRow("1080p", "tv_1080p_min", "tv_1080p_max"),
                  minMaxRow("2160p", "tv_2160p_min", "tv_2160p_max"),
                  Text(
                    "电影大小限制",
                    style: TextStyle(fontSize: 18),
                  ),
                  Divider(),
                  minMaxRow("  720p", "movie_720p_min", "movie_720p_max"),
                  minMaxRow("1080p", "movie_1080p_min", "movie_1080p_max"),
                  minMaxRow("2160p", "movie_2160p_min", "movie_2160p_max"),
                  Center(
                    child: Padding(
                      padding: EdgeInsets.all(20),
                      child: LoadingElevatedButton(
                        onPressed: () async {
                          if (_formKey.currentState!.saveAndValidate()) {
                            var values = _formKey.currentState!.value;

                            return ref
                                .read(mediaSizeLimiterDataProvider.notifier)
                                .submit(MediaSizeLimiter(
                                    tvLimiter: SizeLimiter(
                                      p720p: ResLimiter(
                                          minSize:
                                              toByteInt(values["tv_720p_min"]),
                                          maxSize:
                                              toByteInt(values["tv_720p_max"])),
                                      p1080p: ResLimiter(
                                          minSize:
                                              toByteInt(values["tv_1080p_min"]),
                                          maxSize: toByteInt(
                                              values["tv_1080p_max"])),
                                      p2160p: ResLimiter(
                                          minSize:
                                              toByteInt(values["tv_2160p_min"]),
                                          maxSize: toByteInt(
                                              values["tv_2160p_max"])),
                                    ),
                                    movieLimiter: SizeLimiter(
                                      p720p: ResLimiter(
                                          minSize: toByteInt(
                                              values["movie_720p_min"]),
                                          maxSize: toByteInt(
                                              values["movie_720p_max"])),
                                      p1080p: ResLimiter(
                                          minSize: toByteInt(
                                              values["movie_1080p_min"]),
                                          maxSize: toByteInt(
                                              values["movie_1080p_max"])),
                                      p2160p: ResLimiter(
                                          minSize: toByteInt(
                                              values["movie_2160p_min"]),
                                          maxSize: toByteInt(
                                              values["movie_2160p_max"])),
                                    )));
                          } else {
                            throw "validation_error";
                          }
                        },
                        label: Text("保存"),
                      ),
                    ),
                  )
                ],
              ),
            );
          },
          error: (err, trace) => Container(),
          loading: () => const MyProgressIndicator()),
    );
  }

  Widget minMaxRow(String title, String nameMin, String nameMax) {
    return Row(
      children: [
        Flexible(flex: 2, child: Container()),
        Flexible(
            flex: 2,
            child: Text(
              title,
              style: TextStyle(fontSize: 16),
            )),
        Flexible(flex: 1, child: Container()),
        Flexible(
            flex: 6,
            child: FormBuilderTextField(
              name: nameMin,
              decoration: InputDecoration(suffixText: "MB", labelText: "最小"),
              validator: FormBuilderValidators.compose([
                FormBuilderValidators.required(),
                FormBuilderValidators.numeric()
              ]),
            )),
        Flexible(flex: 1, child: Text("    -    ")),
        Flexible(
            flex: 6,
            child: FormBuilderTextField(
              name: nameMax,
              decoration: InputDecoration(suffixText: "MB", labelText: "最大"),
              validator: FormBuilderValidators.compose([
                FormBuilderValidators.required(),
                FormBuilderValidators.numeric()
              ]),
            )),
        Flexible(flex: 2, child: Container()),
      ],
    );
  }
}

String toMbString(int size) {
  return (size / 1000 / 1000).toString();
}

int toByteInt(String s) {
  return int.parse(s) * 1000 * 1000;
}
