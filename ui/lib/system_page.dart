import 'package:flutter/material.dart';
import 'package:flutter_riverpod/flutter_riverpod.dart';
import 'package:ui/providers/APIs.dart';

import 'package:ui/providers/settings.dart';
import 'package:ui/utils.dart';
import 'package:ui/widgets/progress_indicator.dart';
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
    return SingleChildScrollView(
      child: Column(
        children: [
          ExpansionTile(
            expandedCrossAxisAlignment: CrossAxisAlignment.stretch,
            initiallyExpanded: true,
            childrenPadding: EdgeInsets.all(20),
            title: Text("日志"),
            children: [
              logs.when(
                  data: (list) {
                    return DataTable(
                        columns: const [
                          DataColumn(label: Text("日志")),
                          DataColumn(label: Text("大小")),
                          DataColumn(label: Text("*"))
                        ],
                        rows: List.generate(list.length, (i) {
                          final item = list[i];
                          final uri =
                              Uri.parse("${APIs.logsBaseUrl}${item.name}");

                          return DataRow(cells: [
                            DataCell(Text(item.name ?? "")),
                            DataCell(Text((item.size??0).readableFileSize())),
                            DataCell(InkWell(
                              child: Icon(Icons.download),
                              onTap: () => launchUrl(uri,
                                  webViewConfiguration: WebViewConfiguration(
                                      headers: APIs.authHeaders)),
                            ))
                          ]);
                        }));
                  },
                  error: (err, trace) => Text("$err"),
                  loading: () => const MyProgressIndicator())
            ],
          ),
          ExpansionTile(
            title: Text("关于"),
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
                                InkWell(
                                  child: Text(v.homepage ?? "",
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
    );
  }
}
