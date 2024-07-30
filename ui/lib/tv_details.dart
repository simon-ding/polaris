import 'package:flutter/material.dart';
import 'package:flutter_riverpod/flutter_riverpod.dart';
import 'package:go_router/go_router.dart';
import 'package:ui/providers/APIs.dart';
import 'package:ui/providers/series_details.dart';
import 'package:ui/providers/settings.dart';
import 'package:ui/welcome_page.dart';
import 'package:ui/widgets/utils.dart';
import 'package:ui/widgets/progress_indicator.dart';
import 'package:ui/widgets/widgets.dart';

class TvDetailsPage extends ConsumerStatefulWidget {
  static const route = "/series/:id";

  static String toRoute(int id) {
    return "/series/$id";
  }

  final String seriesId;

  const TvDetailsPage({super.key, required this.seriesId});

  @override
  ConsumerState<ConsumerStatefulWidget> createState() {
    return _TvDetailsPageState();
  }
}

class _TvDetailsPageState extends ConsumerState<TvDetailsPage> {
  @override
  void initState() {
    super.initState();
  }

  @override
  Widget build(BuildContext context) {
    var seriesDetails = ref.watch(mediaDetailsProvider(widget.seriesId));
    var storage = ref.watch(storageSettingProvider);
    return seriesDetails.when(
        data: (details) {
          Map<int, List<DataRow>> m = {};
          for (final ep in details.episodes!) {
            var row = DataRow(cells: [
              DataCell(Text("${ep.episodeNumber}")),
              DataCell(Text("${ep.title}")),
              DataCell(Opacity(
                opacity: 0.5,
                child: Text(ep.airDate ?? "-"),
              )),
              DataCell(
                Opacity(
                    opacity: 0.7,
                    child: ep.status == "downloading"
                        ? const Tooltip(
                            message: "下载中",
                            child: Icon(Icons.downloading),
                          )
                        : (ep.status == "downloaded"
                            ? const Tooltip(
                                message: "已下载",
                                child: Icon(Icons.download_done),
                              )
                            : const Tooltip(
                                message: "未下载",
                                child: Icon(Icons.warning_amber_rounded),
                              ))),
              ),
              DataCell(Row(
                children: [
                  Tooltip(
                    message: "搜索下载对应剧集",
                    child: IconButton(
                        onPressed: () {
                          var f = ref
                              .read(mediaDetailsProvider(widget.seriesId)
                                  .notifier)
                              .searchAndDownload(widget.seriesId,
                                  ep.seasonNumber!, ep.episodeNumber!);
                          showLoadingWithFuture(f);
                        },
                        icon: const Icon(Icons.download)),
                  ),
                  const SizedBox(
                    width: 10,
                  ),
                  IconButton(
                      onPressed: () => showAvailableTorrents(widget.seriesId,
                          ep.seasonNumber ?? 0, ep.episodeNumber ?? 0),
                      icon: const Icon(Icons.manage_search))
                ],
              ))
            ]);

            if (m[ep.seasonNumber] == null) {
              m[ep.seasonNumber!] = List.empty(growable: true);
            }
            m[ep.seasonNumber!]!.add(row);
          }
          List<ExpansionTile> list = List.empty(growable: true);
          for (final k in m.keys.toList().reversed) {
            var seasonList = ExpansionTile(
              tilePadding: const EdgeInsets.fromLTRB(10, 0, 10, 0),
              //childrenPadding: const EdgeInsets.fromLTRB(50, 0, 50, 0),
              initiallyExpanded: false,
              title: k == 0 ? const Text("特别篇") : Text("第 $k 季"),
              expandedCrossAxisAlignment: CrossAxisAlignment.stretch,
              children: [
                DataTable(columns: [
                  const DataColumn(label: Text("#")),
                  const DataColumn(
                    label: Text("标题"),
                  ),
                  const DataColumn(label: Text("播出时间")),
                  const DataColumn(label: Text("状态")),
                  DataColumn(
                      label: Row(
                    children: [
                      Tooltip(
                        message: "搜索下载全部剧集",
                        child: IconButton(
                            onPressed: () {
                              final f = ref
                                  .read(mediaDetailsProvider(widget.seriesId)
                                      .notifier)
                                  .searchAndDownload(widget.seriesId, k, 0).then((v) => showSnakeBar("开始下载: $v"));
                              showLoadingWithFuture(f);
                            },
                            icon: const Icon(Icons.download)),
                      ),
                      const SizedBox(
                        width: 10,
                      ),
                      IconButton(
                          onPressed: () =>
                              showAvailableTorrents(widget.seriesId, k, 0),
                          icon: const Icon(Icons.manage_search))
                    ],
                  ))
                ], rows: m[k]!),
              ],
            );
            list.add(seasonList);
          }
          return ListView(
            children: [
              Card(
                margin: const EdgeInsets.all(4),
                clipBehavior: Clip.hardEdge,
                child: Container(
                  decoration: BoxDecoration(
                      image: DecorationImage(
                          fit: BoxFit.cover,
                          opacity: 0.3,
                          colorFilter: ColorFilter.mode(
                              Colors.black.withOpacity(0.3), BlendMode.dstATop),
                          image: NetworkImage(
                              "${APIs.imagesUrl}/${details.id}/backdrop.jpg"))),
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
                                fit: BoxFit.contain),
                          ),
                        ),
                        Flexible(
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
                                          loading: () => const Text("")),
                                    ],
                                  ),
                                  const Divider(thickness: 1, height: 1),
                                  Text(
                                    "${details.name} ${details.name != details.originalName ? details.originalName : ''} (${details.airDate!.split("-")[0]})",
                                    style: const TextStyle(
                                        fontSize: 20,
                                        fontWeight: FontWeight.bold),
                                  ),
                                  const Text(""),
                                  Text(
                                    details.overview ?? "",
                                  ),
                                ],
                              )),
                              Column(
                                children: [
                                  IconButton(
                                      onPressed: () {
                                        var f = ref
                                            .read(mediaDetailsProvider(
                                                    widget.seriesId)
                                                .notifier)
                                            .delete().then((v) => context.go(WelcomePage.routeTv));
                                        showLoadingWithFuture(f);
                                      },
                                      icon: const Icon(Icons.delete))
                                ],
                              )
                            ],
                          ),
                        ),
                      ],
                    ),
                  ),
                ),
              ),
              Column(
                children: list,
              ),
            ],
          );
        },
        error: (err, trace) {
          return Text("$err");
        },
        loading: () => const MyProgressIndicator());
  }

  Future<void> showAvailableTorrents(String id, int season, int episode) {
    return showDialog<void>(
      context: context,
      barrierDismissible: true,
      builder: (BuildContext context) {
        return Consumer(builder: (context, ref, _) {
          final torrents = ref.watch(mediaTorrentsDataProvider(
              (mediaId: id, seasonNumber: season, episodeNumber: episode)));

          return AlertDialog(
              //title: Text("资源"),
              content: SelectionArea(
            child: SizedBox(
              width: 800,
              height: 400,
              child: torrents.when(
                  data: (v) {
                    return SingleChildScrollView(
                        child: DataTable(
                            dataTextStyle:
                                const TextStyle(fontSize: 12, height: 0),
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
                                DataCell(Text(
                                    "${torrent.size?.readableFileSize()}")),
                                DataCell(Text("${torrent.seeders}")),
                                DataCell(Text("${torrent.peers}")),
                                DataCell(IconButton(
                                  icon: const Icon(Icons.download),
                                  onPressed: () async {
                                    var f = ref
                                        .read(mediaTorrentsDataProvider((
                                          mediaId: id,
                                          seasonNumber: season,
                                          episodeNumber: episode
                                        )).notifier)
                                        .download(torrent).then((v) => showSnakeBar("开始下载：${torrent.name}"));
                                    showLoadingWithFuture(f);
                                  },
                                ))
                              ]);
                            })));
                  },
                  error: (err, trace) {
                    return Text("$err");
                  },
                  loading: () => const MyProgressIndicator()),
            ),
          ));
        });
      },
    );
  }
}
