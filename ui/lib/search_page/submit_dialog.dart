import 'package:flutter/material.dart';
import 'package:flutter_form_builder/flutter_form_builder.dart';
import 'package:flutter_riverpod/flutter_riverpod.dart';
import 'package:ui/providers/settings.dart';
import 'package:ui/providers/welcome_data.dart';
import 'package:ui/widgets/progress_indicator.dart';
import 'package:ui/widgets/utils.dart';
import 'package:ui/widgets/widgets.dart';

class SubmitSearchResult extends ConsumerStatefulWidget {
  final SearchResult item;
  final String query;
  const SubmitSearchResult(
      {super.key, required this.item, required this.query});

  @override
  ConsumerState<ConsumerStatefulWidget> createState() {
    return _SubmitSearchResultState();
  }
}

class _SubmitSearchResultState extends ConsumerState<SubmitSearchResult> {
  final _formKey = GlobalKey<FormBuilderState>();
  bool enabledSizedLimiter = false;
  double sizeMax = 5000;

  @override
  Widget build(BuildContext context) {
    int storageSelected = 0;
    var storage = ref.watch(storageSettingProvider);
    var name = ref.watch(suggestNameDataProvider(
        (id: widget.item.id!, mediaType: widget.item.mediaType!)));

    var pathController = TextEditingController();
    return AlertDialog(
      title: Text('添加: ${widget.item.name}'),
      content: SizedBox(
        width: 500,
        height: 400,
        child: FormBuilder(
          key: _formKey,
          initialValue: const {
            "resolution": "1080p",
            "storage": null,
            "folder": "",
            "history_episodes": false,
            "eanble_size_limier": false,
            "size_limiter": RangeValues(400, 4000),
          },
          child: Column(
            crossAxisAlignment: CrossAxisAlignment.start,
            children: [
              FormBuilderDropdown(
                name: "resolution",
                decoration: const InputDecoration(labelText: "清晰度"),
                items: const [
                  DropdownMenuItem(value: "any", child: Text("不限")),
                  DropdownMenuItem(value: "720p", child: Text("720p")),
                  DropdownMenuItem(value: "1080p", child: Text("1080p")),
                  DropdownMenuItem(value: "2160p", child: Text("2160p")),
                ],
              ),
              storage.when(
                  data: (v) {
                    return StatefulBuilder(builder: (context, setState) {
                      final id1 = v.isEmpty ? 0 : v[0].id!;
                      return Column(
                        crossAxisAlignment: CrossAxisAlignment.start,
                        children: [
                          FormBuilderDropdown(
                            initialValue: id1,
                            onChanged: (v) {
                              setState(
                                () {
                                  storageSelected = v!;
                                },
                              );
                            },
                            name: "storage",
                            decoration:
                                const InputDecoration(labelText: "存储位置"),
                            items: v
                                .map((s) => DropdownMenuItem(
                                    value: s.id, child: Text(s.name!)))
                                .toList(),
                          ),
                          name.when(
                            data: (s) {
                              if (storageSelected == 0) {
                                storageSelected = id1;
                              }
                              return storageSelected == 0
                                  ? const Text("")
                                  : () {
                                      final storage = v
                                          .where((e) => e.id == storageSelected)
                                          .first;
                                      final path = widget.item.mediaType == "tv"
                                          ? storage.tvPath
                                          : storage.moviePath;

                                      pathController.text = s;
                                      return SizedBox(
                                        //width: 300,
                                        child: FormBuilderTextField(
                                          name: "folder",
                                          controller: pathController,
                                          decoration: InputDecoration(
                                              labelText: "存储路径",
                                              prefix: Text(path ?? "unknown")),
                                        ),
                                      );
                                    }();
                            },
                            error: (error, stackTrace) => Text("$error"),
                            loading: () => const MyProgressIndicator(
                              size: 20,
                            ),
                          ),
                          FormBuilderSwitch(
                            name: "eanble_size_limier",
                            title: Text(widget.item.mediaType == "tv"
                                ? "是否限制每集文件大小"
                                : "是否限制电影文件大小"),
                            onChanged: (value) {
                              setState(
                                () {
                                  enabledSizedLimiter = value!;
                                },
                              );
                            },
                          ),
                          enabledSizedLimiter
                              ? const MyRangeSlider(name: "size_limiter")
                              : const SizedBox(),
                          widget.item.mediaType == "tv"
                              ? SizedBox(
                                  width: 250,
                                  child: FormBuilderCheckbox(
                                    name: "history_episodes",
                                    title: const Text("是否下载往期剧集"),
                                  ),
                                )
                              : const SizedBox(),
                        ],
                      );
                    });
                  },
                  error: (err, trace) => Text("$err"),
                  loading: () => const MyProgressIndicator()),
            ],
          ),
        ),
      ),
      actions: <Widget>[
        TextButton(
          style: TextButton.styleFrom(
            textStyle: Theme.of(context).textTheme.labelLarge,
          ),
          child: const Text('取消'),
          onPressed: () {
            Navigator.of(context).pop();
          },
        ),
        LoadingTextButton(
          // style: TextButton.styleFrom(
          //   textStyle: Theme.of(context).textTheme.labelLarge,
          // ),
          label: const Text('确定'),
          onPressed: () async {
            if (_formKey.currentState!.saveAndValidate()) {
              final values = _formKey.currentState!.value;
              await ref
                  .read(searchPageDataProvider(widget.query).notifier)
                  .submit2Watchlist(
                      widget.item.id!,
                      values["storage"],
                      values["resolution"],
                      widget.item.mediaType!,
                      values["folder"],
                      values["history_episodes"] ?? false,
                      enabledSizedLimiter
                          ? values["size_limiter"]
                          : const RangeValues(-1, -1))
                  .then((v) {
                Navigator.of(context).pop();
                showSnakeBar("添加成功：${widget.item.name}");
              });
            }
          },
        ),
      ],
    );
  }

  String readableSize(String v) {
    if (v.endsWith("K")) {
      return v.replaceAll("K", " GB");
    }
    return "$v MB";
  }
}
