import 'package:flutter/material.dart';
import 'package:flutter_riverpod/flutter_riverpod.dart';
import 'package:ui/providers/series_details.dart';
import 'package:ui/widgets/progress_indicator.dart';
import 'package:ui/widgets/utils.dart';
import 'package:ui/widgets/widgets.dart';

class ResourceList extends ConsumerWidget {
  final String mediaId;
  final int seasonNum;
  final int episodeNum;

  const ResourceList(
      {super.key,
      required this.mediaId,
      this.seasonNum = 0,
      this.episodeNum = 0});

  @override
  Widget build(BuildContext context, WidgetRef ref) {
    final torrents = ref.watch(mediaTorrentsDataProvider((
      mediaId: mediaId,
      seasonNumber: seasonNum,
      episodeNumber: episodeNum
    )));
    var widgets = torrents.when(
        data: (v) {
          bool hasPrivate = false;
          for (final item in v) {
            if (item.isPrivate == true) {
              hasPrivate = true;
            }
          }
          final columns = [
            const DataColumn(label: Text("名称")),
            const DataColumn(label: Text("大小")),
            const DataColumn(label: Text("S/P")),
            const DataColumn(label: Text("来源")),
          ];
          if (hasPrivate) {
            columns.add(const DataColumn(label: Text("消耗")));
          }
          columns.add(const DataColumn(label: Text("下载")));

          return DataTable(
              dataTextStyle: const TextStyle(fontSize: 12),
              columns: columns,
              rows: List.generate(v.length, (i) {
                final torrent = v[i];
                final rows = [
                  DataCell(Text("${torrent.name}")),
                  DataCell(Text("${torrent.size?.readableFileSize()}")),
                  DataCell(Text("${torrent.seeders}/${torrent.peers}")),
                  DataCell(Text(torrent.source ?? "-")),
                ];
                if (hasPrivate) {
                  rows.add(DataCell(Text(torrent.isPrivate == true
                      ? "${torrent.downloadFactor}dl/${torrent.uploadFactor}up"
                      : "-")));
                }

                rows.add(DataCell(LoadingIconButton(
                  icon: Icons.download,
                  onPressed: () async {
                    await ref
                        .read(mediaTorrentsDataProvider((
                          mediaId: mediaId,
                          seasonNumber: seasonNum,
                          episodeNumber: episodeNum
                        )).notifier)
                        .download(torrent)
                        .then((v) => showSnakeBar("开始下载：${torrent.name}"));
                  },
                )));
                return DataRow(cells: rows);
              }));
        },
        error: (err, trace) {
          return "$err".contains("no resource found")
              ? const Center(
                  child: Text("没有资源"),
                )
              : Text("$err");
        },
        loading: () => const MyProgressIndicator());
    return isSmallScreen(context)
        ? SingleChildScrollView(
            scrollDirection: Axis.horizontal, child: widgets)
        : widgets;
  }
}
