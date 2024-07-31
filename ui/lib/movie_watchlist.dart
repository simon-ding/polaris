import 'package:flutter/material.dart';
import 'package:flutter_riverpod/flutter_riverpod.dart';
import 'package:ui/providers/activity.dart';
import 'package:ui/providers/series_details.dart';
import 'package:ui/widgets/detail_card.dart';
import 'package:ui/widgets/utils.dart';
import 'package:ui/widgets/progress_indicator.dart';
import 'package:ui/widgets/widgets.dart';

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

    return seriesDetails.when(
        data: (details) {
          return ListView(
            children: [
              DetailCard(details: details),
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
    _nestedTabController = TabController(length: 2, vsync: this);
  }

  @override
  void dispose() {
    super.dispose();
    _nestedTabController.dispose();
  }

  int selectedTab = 0;

  @override
  Widget build(BuildContext context) {
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
            return Consumer(
              builder: (context, ref, child) {
                var torrents = ref.watch(mediaTorrentsDataProvider(
                    (mediaId: widget.id, seasonNumber: 0, episodeNumber: 0)));
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
                            DataCell(
                                Text("${torrent.size?.readableFileSize()}")),
                            DataCell(Text("${torrent.seeders}")),
                            DataCell(Text("${torrent.peers}")),
                            DataCell(IconButton(
                              icon: const Icon(Icons.download),
                              onPressed: () {
                                final f = ref
                                    .read(mediaTorrentsDataProvider((
                                      mediaId: widget.id,
                                      seasonNumber: 0,
                                      episodeNumber: 0
                                    )).notifier)
                                    .download(torrent)
                                .then((v) => showSnakeBar(
                                    "开始下载：${torrent.name}"));
                                // .onError((error, trace) =>
                                //     Utils.showSnakeBar("操作失败: $error"));
                                showLoadingWithFuture(f);
                              },
                            ))
                          ]);
                        }),
                      );
                    },
                    error: (error, trace) => Text("$error"),
                    loading: () => const MyProgressIndicator());
              },
            );
          }
        })
      ],
    );
  }
}
