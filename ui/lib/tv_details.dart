import 'package:flutter/material.dart';
import 'package:flutter_riverpod/flutter_riverpod.dart';
import 'package:ui/providers/series_details.dart';
import 'package:ui/widgets/detail_card.dart';
import 'package:ui/widgets/resource_list.dart';
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
    final screenWidth = MediaQuery.of(context).size.width;
    var seriesDetails = ref.watch(mediaDetailsProvider(widget.seriesId));
    return SelectionArea(
        child: seriesDetails.when(
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
                            ? const IconButton(
                                tooltip: "下载中",
                                onPressed: null,
                                icon: Icon(Icons.downloading))
                            : (ep.status == "downloaded"
                                ? const IconButton(
                                    tooltip: "已下载",
                                    onPressed: null,
                                    icon: Icon(Icons.download_done))
                                : (ep.monitored == true
                                    ? IconButton(
                                        tooltip: "监控中",
                                        onPressed: () {
                                          ref
                                              .read(mediaDetailsProvider(
                                                      widget.seriesId)
                                                  .notifier)
                                              .changeMonitoringStatus(
                                                  ep.id!, false);
                                        },
                                        icon: const Icon(Icons.alarm))
                                    : Opacity(
                                        opacity: 0.7,
                                        child: IconButton(
                                            tooltip: "未监控",
                                            onPressed: () {
                                              ref
                                                  .read(mediaDetailsProvider(
                                                          widget.seriesId)
                                                      .notifier)
                                                  .changeMonitoringStatus(
                                                      ep.id!, true);
                                            },
                                            icon: const Icon(Icons.alarm_off)),
                                      )))),
                  ),
                  DataCell(Row(
                    children: [
                      LoadingIconButton(
                          tooltip: "搜索下载对应剧集",
                          onPressed: () async {
                            await ref
                                .read(mediaDetailsProvider(widget.seriesId)
                                    .notifier)
                                .searchAndDownload(widget.seriesId,
                                    ep.seasonNumber!, ep.episodeNumber!)
                                .then((v) => showSnakeBar("开始下载: $v"));
                          },
                          icon: Icons.download),
                      const SizedBox(
                        width: 10,
                      ),
                      Tooltip(
                        message: "查看可用资源",
                        child: IconButton(
                            onPressed: () => showAvailableTorrents(
                                widget.seriesId,
                                ep.seasonNumber ?? 0,
                                ep.episodeNumber ?? 0),
                            icon: const Icon(Icons.manage_search)),
                      )
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
                final seasonEpisodes = DataTable(columns: [
                  const DataColumn(label: Text("#")),
                  const DataColumn(
                    label: Text("标题"),
                  ),
                  const DataColumn(label: Text("播出时间")),
                  const DataColumn(label: Text("状态")),
                  DataColumn(
                      label: Row(
                    children: [
                      LoadingIconButton(
                          tooltip: "搜索下载全部剧集",
                          onPressed: () async {
                            await ref
                                .read(mediaDetailsProvider(widget.seriesId)
                                    .notifier)
                                .searchAndDownload(widget.seriesId, k, 0)
                                .then((v) => showSnakeBar("开始下载: $v"));
                            //showLoadingWithFuture(f);
                          },
                          icon: Icons.download),
                      const SizedBox(
                        width: 10,
                      ),
                      Tooltip(
                        message: "查看可用资源",
                        child: IconButton(
                            onPressed: () =>
                                showAvailableTorrents(widget.seriesId, k, 0),
                            icon: const Icon(Icons.manage_search)),
                      )
                    ],
                  ))
                ], rows: m[k]!);

                var seasonList = ExpansionTile(
                  tilePadding: const EdgeInsets.fromLTRB(10, 0, 10, 0),
                  //childrenPadding: const EdgeInsets.fromLTRB(50, 0, 50, 0),
                  initiallyExpanded: false,
                  title: k == 0 ? const Text("特别篇") : Text("第 $k 季"),
                  expandedCrossAxisAlignment: CrossAxisAlignment.stretch,
                  children: [
                    screenWidth < 600
                        ? SingleChildScrollView(
                            scrollDirection: Axis.horizontal,
                            child: seasonEpisodes,
                          )
                        : seasonEpisodes
                  ],
                );
                list.add(seasonList);
              }
              return ListView(
                children: [
                  DetailCard(details: details),
                  Column(
                    children: list,
                  ),
                ],
              );
            },
            error: (err, trace) => PoNetworkError(err: err),
            loading: () => const MyProgressIndicator()));
  }

  Future<void> showAvailableTorrents(String id, int season, int episode) {
    return showDialog<void>(
      context: context,
      barrierDismissible: true,
      builder: (BuildContext context) {
        return AlertDialog(
            //title: Text("资源"),
            content: SelectionArea(
          child: SizedBox(
            width: MediaQuery.of(context).size.width * 0.7,
            height: MediaQuery.of(context).size.height * 0.6,
            child: SingleChildScrollView(
              child: ResourceList(
                mediaId: id,
                seasonNum: season,
                episodeNum: episode,
              ),
            ),
          ),
        ));
      },
    );
  }
}
