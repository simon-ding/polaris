import 'package:flutter/material.dart';
import 'package:flutter_riverpod/flutter_riverpod.dart';
import 'package:go_router/go_router.dart';
import 'package:ui/providers/APIs.dart';
import 'package:ui/providers/series_details.dart';
import 'package:ui/providers/settings.dart';
import 'package:ui/providers/welcome_data.dart';
import 'package:ui/utils.dart';
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
    var torrents = ref.watch(movieTorrentsDataProvider(widget.id));
    var storage = ref.watch(storageSettingProvider);

    return seriesDetails.when(
        data: (details) {
          return ListView(
            children: [
              Card(
                margin: const EdgeInsets.all(4),
                clipBehavior: Clip.hardEdge,
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
                                          .read(
                                              mediaDetailsProvider(widget.id).notifier)
                                          .delete();
                                      context.go(MovieDetailsPage.route);
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
              torrents.when(
                  data: (v) {
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
                            onPressed: () async {
                              await ref
                                  .read(movieTorrentsDataProvider(widget.id)
                                      .notifier)
                                  .download(torrent.link!);
                            },
                          ))
                        ]);
                      }),
                    );
                  },
                  error: (error, trace) => Text("$error"),
                  loading: () => const MyProgressIndicator()),
            ],
          );
        },
        error: (err, trace) {
          return Text("$err");
        },
        loading: () => const MyProgressIndicator());
  }
}
