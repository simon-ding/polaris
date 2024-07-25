import 'package:flutter/material.dart';
import 'package:flutter_riverpod/flutter_riverpod.dart';
import 'package:go_router/go_router.dart';
import 'package:ui/providers/APIs.dart';
import 'package:ui/providers/activity.dart';
import 'package:ui/providers/series_details.dart';
import 'package:ui/providers/settings.dart';
import 'package:ui/providers/welcome_data.dart';
import 'package:ui/utils.dart';
import 'package:ui/welcome_page.dart';
import 'package:ui/widgets/progress_indicator.dart';

class MovieDetailsPage extends ConsumerStatefulWidget {
  static const route = "/movie/:id";

  static String toRoute(int id) {
    return "/movie/$id";
  }

  final String id;

  const MovieDetailsPage({super.key, required this.id});

  @override
  ConsumerState<ConsumerStatefulWidget> createState() {
    return _MovieDetailsPageState();
  }
}

class _MovieDetailsPageState extends ConsumerState<MovieDetailsPage> {
  @override
  Widget build(BuildContext context) {
    var seriesDetails = ref.watch(mediaDetailsProvider(widget.id));
    var storage = ref.watch(storageSettingProvider);

    return seriesDetails.when(
        data: (details) {
          return ListView(
            children: [
              Card(
                margin: const EdgeInsets.all(4),
                clipBehavior: Clip.hardEdge,
                child: Container(
                    decoration: BoxDecoration(
                        image: DecorationImage(
                            fit: BoxFit.fitWidth,
                            opacity: 0.5,
                            image: NetworkImage(
                                "${APIs.imagesUrl}/${details.id}/backdrop.jpg",
                                headers: APIs.authHeaders))),
                    child: Padding(
                      padding: const EdgeInsets.all(10),
                      child: Row(
                        children: <Widget>[
                          Flexible(
                            flex: 1,
                            child: Padding(
                              padding: const EdgeInsets.all(10),
                              child: Image.network(
                                "${APIs.imagesUrl}/${details.id}/poster.jpg",
                                fit: BoxFit.contain,
                                headers: APIs.authHeaders,
                              ),
                            ),
                          ),
                          Expanded(
                            flex: 6,
                            child: Row(
                              children: [
                                Expanded(
                                    child: Column(
                                  crossAxisAlignment: CrossAxisAlignment.start,
                                  children: [
                                    Row(
                                      children: [
                                        Text("${details.resolution}"),
                                        const SizedBox(
                                          width: 30,
                                        ),
                                        storage.when(
                                            data: (value) {
                                              for (final s in value) {
                                                if (s.id == details.storageId) {
                                                  return Text(
                                                      "${s.name}(${s.implementation})");
                                                }
                                              }
                                              return const Text("未知存储");
                                            },
                                            error: (error, stackTrace) =>
                                                Text("$error"),
                                            loading: () =>
                                                const MyProgressIndicator()),
                                      ],
                                    ),
                                    const Divider(thickness: 1, height: 1),
                                    Text(
                                      "${details.name} (${details.airDate!.split("-")[0]})",
                                      style: const TextStyle(
                                          fontSize: 20,
                                          fontWeight: FontWeight.bold),
                                    ),
                                    const Text(""),
                                    Text(
                                      details.overview!,
                                    ),
                                  ],
                                )),
                                Column(
                                  children: [
                                    IconButton(
                                        onPressed: () {
                                          ref
                                              .read(mediaDetailsProvider(
                                                      widget.id)
                                                  .notifier)
                                              .delete()
                                              .then((v) => context
                                                  .go(WelcomePage.routeMoivie))
                                              .onError((error, trace) =>
                                                  Utils.showSnakeBar(
                                                      "删除失败：$error"));
                                        },
                                        icon: const Icon(Icons.delete))
                                  ],
                                )
                              ],
                            ),
                          ),
                        ],
                      ),
                    )),
              ),
              NestedTabBar(
                id: widget.id,
              )
            ],
          );
        },
        error: (err, trace) {
          return Text("$err");
        },
        loading: () => const MyProgressIndicator());
  }
}

class NestedTabBar extends ConsumerStatefulWidget {
  final String id;

  const NestedTabBar({super.key, required this.id});

  @override
  _NestedTabBarState createState() => _NestedTabBarState();
}

class _NestedTabBarState extends ConsumerState<NestedTabBar>
    with TickerProviderStateMixin {
  late TabController _nestedTabController;
  @override
  void initState() {
    super.initState();
    _nestedTabController = new TabController(length: 2, vsync: this);
  }

  @override
  void dispose() {
    super.dispose();
    _nestedTabController.dispose();
  }

  int selectedTab = 0;

  @override
  Widget build(BuildContext context) {
    var torrents = ref.watch(movieTorrentsDataProvider(widget.id));
    var histories = ref.watch(mediaHistoryDataProvider(widget.id));

    return Column(
      mainAxisAlignment: MainAxisAlignment.spaceAround,
      crossAxisAlignment: CrossAxisAlignment.stretch,
      children: <Widget>[
        TabBar(
          controller: _nestedTabController,
          isScrollable: true,
          onTap: (value) {
            setState(() {
              selectedTab = value;
            });
          },
          tabs: const <Widget>[
            Tab(
              text: "下载记录",
            ),
            Tab(
              text: "资源",
            ),
          ],
        ),
        Builder(builder: (context) {
          if (selectedTab == 0) {
            return histories.when(
                data: (v) {
                  if (v.isEmpty) {
                    return const Center(
                      child: Text("无下载记录"),
                    );
                  }
                  return DataTable(
                      columns: const [
                        DataColumn(label: Text("#"), numeric: true),
                        DataColumn(label: Text("名称")),
                        DataColumn(label: Text("下载时间")),
                      ],
                      rows: List.generate(v.length, (i) {
                        final activity = v[i];
                        return DataRow(cells: [
                          DataCell(Text("${activity.id}")),
                          DataCell(Text("${activity.sourceTitle}")),
                          DataCell(Text("${activity.date!.toLocal()}")),
                        ]);
                      }));
                },
                error: (error, trace) => Text("$error"),
                loading: () => const MyProgressIndicator());
          } else {
            return torrents.when(
                data: (v) {
                  if (v.isEmpty) {
                    return const Center(
                      child: Text("无可用资源"),
                    );
                  }

                  return DataTable(
                    columns: const [
                      DataColumn(label: Text("名称")),
                      DataColumn(label: Text("大小")),
                      DataColumn(label: Text("seeders")),
                      DataColumn(label: Text("peers")),
                      DataColumn(label: Text("操作"))
                    ],
                    rows: List.generate(v.length, (i) {
                      final torrent = v[i];
                      return DataRow(cells: [
                        DataCell(Text("${torrent.name}")),
                        DataCell(Text("${torrent.size?.readableFileSize()}")),
                        DataCell(Text("${torrent.seeders}")),
                        DataCell(Text("${torrent.peers}")),
                        DataCell(IconButton(
                          icon: const Icon(Icons.download),
                          onPressed: () {
                            ref
                                .read(movieTorrentsDataProvider(widget.id)
                                    .notifier)
                                .download(torrent)
                                .then((v) =>
                                    Utils.showSnakeBar("开始下载：${torrent.name}"))
                                .onError((error, trace) =>
                                    Utils.showSnakeBar("操作失败: $error"));
                          },
                        ))
                      ]);
                    }),
                  );
                },
                error: (error, trace) => Text("$error"),
                loading: () => const MyProgressIndicator());
          }
        })
      ],
    );
  }
}
