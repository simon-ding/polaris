import 'package:flutter/material.dart';
import 'package:flutter_riverpod/flutter_riverpod.dart';
import 'package:ui/providers/APIs.dart';

import 'package:ui/providers/settings.dart';
import 'package:ui/widgets/utils.dart';
import 'package:ui/widgets/progress_indicator.dart';
import 'package:ui/widgets/widgets.dart';
import 'package:url_launcher/url_launcher.dart';

class SystemPage extends ConsumerStatefulWidget {
  static const route = "/system";

  const SystemPage({super.key});
  @override
  ConsumerState<ConsumerStatefulWidget> createState() {
    return _SystemPageState();
  }
}

class _SystemPageState extends ConsumerState<SystemPage> {
  @override
  Widget build(BuildContext context) {
    final logs = ref.watch(logFileDataProvider);
    final about = ref.watch(aboutDataProvider);
    return SelectionArea(
        child: SingleChildScrollView(
      child: Column(
        children: [
          ExpansionTile(
            expandedCrossAxisAlignment: CrossAxisAlignment.stretch,
            initiallyExpanded: true,
            childrenPadding: const EdgeInsets.all(20),
            title: const Text("日志"),
            children: [
              logs.when(
                  data: (list) {
                    return DataTable(
                        columns: const [
                          DataColumn(label: Text("日志")),
                          DataColumn(label: Text("大小")),
                          DataColumn(label: Text("下载"))
                        ],
                        rows: List.generate(list.length, (i) {
                          final item = list[i];
                          final uri =
                              Uri.parse("${APIs.logsBaseUrl}${item.name}");

                          return DataRow(cells: [
                            DataCell(Text(item.name ?? "")),
                            DataCell(Text((item.size ?? 0).readableFileSize())),
                            DataCell(InkWell(
                              child: const Icon(Icons.download),
                              onTap: () => launchUrl(uri),
                            ))
                          ]);
                        }));
                  },
                  error: (err, trace) => Text("$err"),
                  loading: () => const MyProgressIndicator())
            ],
          ),
          ExpansionTile(
            expandedCrossAxisAlignment: CrossAxisAlignment.stretch,
            initiallyExpanded: true,
            childrenPadding: const EdgeInsets.all(20),
            title: const Text("定时任务"),
            children: [
              logs.when(
                  data: (list) {
                    return DataTable(columns: const [
                      DataColumn(label: Text("任务")),
                      DataColumn(label: Text("间隔")),
                      DataColumn(label: Text("手动触发"))
                    ], rows: [
                      DataRow(cells: [
                        const DataCell(Text("检查并处理已完成任务")),
                        const DataCell(Text("每分钟")),
                        DataCell(LoadingIconButton(
                          icon: Icons.not_started,
                          onPressed: () =>
                              APIs.triggerCronJob("check_running_tasks"),
                        ))
                      ]),
                      DataRow(cells: [
                        const DataCell(Text("下载监控中的剧集和电影")),
                        const DataCell(Text("每小时")),
                        DataCell(LoadingIconButton(
                          icon: Icons.not_started,
                          onPressed: () => APIs.triggerCronJob(
                              "check_available_medias_to_download"),
                        ))
                      ]),
                      DataRow(cells: [
                        const DataCell(Text("更新监控列表")),
                        const DataCell(Text("每小时")),
                        DataCell(IconButton(
                          icon: const Icon(Icons.not_started),
                          onPressed: () =>
                              APIs.triggerCronJob("update_import_lists"),
                        ))
                      ]),
                      DataRow(cells: [
                        const DataCell(Text("更新最新剧集信息")),
                        const DataCell(Text("每天2次")),
                        DataCell(LoadingIconButton(
                          icon: Icons.not_started,
                          onPressed: () =>
                              APIs.triggerCronJob("check_series_new_release"),
                        ))
                      ]),
                    ]);
                  },
                  error: (err, trace) => Text("$err"),
                  loading: () => const MyProgressIndicator())
            ],
          ),
          ExpansionTile(
            title: const Text("关于"),
            expandedCrossAxisAlignment: CrossAxisAlignment.center,
            initiallyExpanded: true,
            children: [
              about.when(
                  data: (v) {
                    final uri = Uri.parse(v.chatGroup ?? "");
                    final homepage = Uri.parse(v.homepage ?? "");
                    return Row(
                      children: [
                        const Expanded(
                            child: Column(
                          crossAxisAlignment: CrossAxisAlignment.end,
                          children: [
                            SizedBox(
                              height: 20,
                            ),
                            Text(
                              "#",
                              style: TextStyle(height: 2.5),
                            ),
                            Text("版本", style: TextStyle(height: 2.5)),
                            Text("主页", style: TextStyle(height: 2.5)),
                            Text("讨论组", style: TextStyle(height: 2.5)),
                            Text("go version", style: TextStyle(height: 2.5)),
                            Text("uptime", style: TextStyle(height: 2.5)),
                            SizedBox(
                              height: 20,
                            ),
                          ],
                        )),
                        const SizedBox(
                          width: 20,
                        ),
                        Expanded(
                            flex: 2,
                            child: Column(
                              crossAxisAlignment: CrossAxisAlignment.start,
                              children: [
                                const SizedBox(
                                  height: 20,
                                ),
                                Text(v.intro ?? "",
                                    style: const TextStyle(height: 2.5)),
                                Text(v.version ?? "",
                                    style: const TextStyle(height: 2.5)),
                                InkWell(
                                  child: Text(v.homepage ?? "",
                                      softWrap: false,
                                      style: const TextStyle(height: 2.5)),
                                  onTap: () => launchUrl(homepage),
                                ),
                                InkWell(
                                  child: const Text("Telegram",
                                      style: TextStyle(height: 2.5)),
                                  onTap: () => launchUrl(uri),
                                ),
                                Text("${v.goVersion}",
                                    style: const TextStyle(height: 2.5)),
                                Text("${v.uptime}",
                                    style: const TextStyle(height: 2.5)),
                                const SizedBox(
                                  height: 20,
                                ),
                              ],
                            )),
                      ],
                    );
                  },
                  error: (err, trace) => Text("$err"),
                  loading: () => const MyProgressIndicator())
            ],
          )
        ],
      ),
    ));
  }
}
