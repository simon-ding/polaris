import 'package:flutter/material.dart';
import 'package:flutter_form_builder/flutter_form_builder.dart';
import 'package:flutter_riverpod/flutter_riverpod.dart';
import 'package:go_router/go_router.dart';
import 'package:ui/providers/APIs.dart';
import 'package:ui/providers/settings.dart';
import 'package:ui/providers/welcome_data.dart';
import 'package:ui/widgets/progress_indicator.dart';
import 'package:ui/widgets/utils.dart';
import 'package:ui/widgets/widgets.dart';

class SearchPage extends ConsumerStatefulWidget {
  const SearchPage({super.key, this.query});

  static const route = "/search";
  final String? query;

  @override
  ConsumerState<ConsumerStatefulWidget> createState() {
    return _SearchPageState();
  }
}

class _SearchPageState extends ConsumerState<SearchPage> {
  List<dynamic> list = List.empty();

  @override
  Widget build(BuildContext context) {
    final q = widget.query ?? "";
    var searchList = ref.watch(searchPageDataProvider(q));

    List<Widget> res = searchList.when(
        data: (data) {
          if (data.isEmpty) {
            return [
              Container(
                  height: MediaQuery.of(context).size.height * 0.6,
                  alignment: Alignment.center,
                  child: const Text(
                    "啥都没有...",
                    style: TextStyle(fontSize: 16),
                  ))
            ];
          }
          var cards = List<Widget>.empty(growable: true);
          for (final item in data) {
            cards.add(Card(
                margin: const EdgeInsets.all(4),
                clipBehavior: Clip.hardEdge,
                child: InkWell(
                  //splashColor: Colors.blue.withAlpha(30),
                  onTap: () {
                    if (item.inWatchlist != true) {
                      _showSubmitDialog(context, item);
                    }
                  },
                  child: Row(
                    children: <Widget>[
                      Flexible(
                        child: SizedBox(
                          width: 150,
                          height: 200,
                          child: Image.network(
                            "${APIs.tmdbImgBaseUrl}${item.posterPath}",
                            fit: BoxFit.contain,
                          ),
                        ),
                      ),
                      Flexible(
                        child: Column(
                          crossAxisAlignment: CrossAxisAlignment.start,
                          children: [
                            Row(
                              children: [
                                Text(
                                  "${item.name} ${item.name != item.originalName ? item.originalName : ''} (${item.firstAirDate?.year})",
                                  style: const TextStyle(
                                      fontSize: 14,
                                      fontWeight: FontWeight.bold),
                                ),
                                const SizedBox(
                                  width: 10,
                                ),
                                item.mediaType == "tv"
                                    ? const Chip(
                                        avatar: Icon(Icons.live_tv),
                                        label: Text(
                                          "剧集",
                                        ))
                                    : const Chip(
                                        avatar: Icon(Icons.movie),
                                        label: Text("电影")),
                                item.inWatchlist == true
                                    ? const Chip(
                                        label: Icon(
                                        Icons.done,
                                        color: Colors.green,
                                      ))
                                    : const Text("")
                              ],
                            ),
                            const Text(""),
                            item.originCountry.isNotEmpty
                                ? Text("国家：${item.originCountry[0]}")
                                : Container(),
                            Text("${item.overview}")
                          ],
                        ),
                      )
                    ],
                  ),
                )));
          }
          return cards;
        },
        error: (err, trace) => [Text("$err")],
        loading: () => [const MyProgressIndicator()]);

    var f = NotificationListener(
        onNotification: (ScrollNotification scrollInfo) {
          if (scrollInfo is ScrollEndNotification &&
              scrollInfo.metrics.axisDirection == AxisDirection.down &&
              scrollInfo.metrics.pixels >= scrollInfo.metrics.maxScrollExtent) {
            ref.read(searchPageDataProvider(q).notifier).queryNextPage();
          }
          return true;
        },
        child: ListView(
          children: res,
        ));
    return Column(
      children: [
        TextField(
          autofocus: true,
          controller: TextEditingController(text: q),
          onSubmitted: (value) async {
            context.go(
                Uri(path: SearchPage.route, queryParameters: {'query': value})
                    .toString());
          },
          decoration: const InputDecoration(
              labelText: "搜索",
              hintText: "搜索剧集名称",
              prefixIcon: Icon(Icons.search)),
        ),
        Expanded(child: f)
      ],
    );
  }

  Future<void> _showSubmitDialog(BuildContext context, SearchResult item) {
    final _formKey = GlobalKey<FormBuilderState>();

    return showDialog<void>(
        context: context,
        builder: (BuildContext context) {
          return Consumer(
            builder: (context, ref, _) {
              int storageSelected = 0;
              var storage = ref.watch(storageSettingProvider);
              var name = ref.watch(suggestNameDataProvider(
                  (id: item.id!, mediaType: item.mediaType!)));

              var pathController = TextEditingController();
              return AlertDialog(
                title: Text('添加: ${item.name}'),
                content: SizedBox(
                  width: 500,
                  height: 200,
                  child: FormBuilder(
                    key: _formKey,
                    initialValue: const {
                      "resolution": "1080p",
                      "storage": null,
                      "folder": "",
                      "history_episodes": false,
                    },
                    child: Column(
                      crossAxisAlignment: CrossAxisAlignment.start,
                      children: [
                        FormBuilderDropdown(
                          name: "resolution",
                          decoration: const InputDecoration(labelText: "清晰度"),
                          items: const [
                            DropdownMenuItem(
                                value: "720p", child: Text("720p")),
                            DropdownMenuItem(
                                value: "1080p", child: Text("1080p")),
                            DropdownMenuItem(
                                value: "2160p", child: Text("2160p")),
                          ],
                        ),
                        storage.when(
                            data: (v) {
                              return StatefulBuilder(
                                  builder: (context, setState) {
                                return Column(
                                  crossAxisAlignment: CrossAxisAlignment.start,
                                  children: [
                                    FormBuilderDropdown(
                                      onChanged: (v) {
                                        setState(
                                          () {
                                            storageSelected = v!;
                                          },
                                        );
                                      },
                                      name: "storage",
                                      decoration: const InputDecoration(
                                          labelText: "存储位置"),
                                      items: v
                                          .map((s) => DropdownMenuItem(
                                              value: s.id,
                                              child: Text(s.name!)))
                                          .toList(),
                                    ),
                                    name.when(
                                      data: (s) {
                                        return storageSelected == 0
                                            ? const Text("")
                                            : () {
                                                final storage = v
                                                    .where((e) =>
                                                        e.id == storageSelected)
                                                    .first;
                                                final path =
                                                    item.mediaType == "tv"
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
                                                        prefix: Text(
                                                            path ?? "unknown")),
                                                  ),
                                                );
                                              }();
                                      },
                                      error: (error, stackTrace) =>
                                          Text("$error"),
                                      loading: () => const MyProgressIndicator(
                                        size: 20,
                                      ),
                                    ),
                                    item.mediaType == "tv"
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
                  TextButton(
                    style: TextButton.styleFrom(
                      textStyle: Theme.of(context).textTheme.labelLarge,
                    ),
                    child: const Text('确定'),
                    onPressed: () async {
                      if (_formKey.currentState!.saveAndValidate()) {
                        final values = _formKey.currentState!.value;
                        //print(values);
                        var f = ref
                            .read(searchPageDataProvider(widget.query ?? "")
                                .notifier)
                            .submit2Watchlist(
                                item.id!,
                                values["storage"],
                                values["resolution"],
                                item.mediaType!,
                                values["folder"],
                                values["history_episodes"] ?? false)
                            .then((v) {
                          Navigator.of(context).pop();
                          showSnakeBar("添加成功：${item.name}");
                        });
                        showLoadingWithFuture(f);
                      }
                    },
                  ),
                ],
              );
            },
          );
        });
  }
}

class SearchBarApp extends StatefulWidget {
  const SearchBarApp({
    super.key,
    required this.onChanged,
  });

  final ValueChanged<String> onChanged;
  @override
  State<SearchBarApp> createState() => _SearchBarAppState();
}

class _SearchBarAppState extends State<SearchBarApp> {
  @override
  Widget build(BuildContext context) {
    return SearchAnchor(
        builder: (BuildContext context, SearchController controller) {
      return SearchBar(
        controller: controller,
        padding: const WidgetStatePropertyAll<EdgeInsets>(
            EdgeInsets.symmetric(horizontal: 16.0)),
        onSubmitted: (value) => {widget.onChanged(controller.text)},
        leading: const Icon(Icons.search),
      );
    }, suggestionsBuilder: (BuildContext context, SearchController controller) {
      return List<ListTile>.generate(0, (int index) {
        final String item = 'item $index';
        return ListTile(
          title: Text(item),
          onTap: () {
            setState(() {
              controller.closeView(item);
            });
          },
        );
      });
    });
  }
}
