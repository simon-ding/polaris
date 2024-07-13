import 'package:flutter/material.dart';
import 'package:flutter_riverpod/flutter_riverpod.dart';
import 'package:go_router/go_router.dart';
import 'package:ui/providers/APIs.dart';
import 'package:ui/providers/series_details.dart';
import 'package:ui/utils.dart';
import 'package:ui/weclome.dart';

class TvDetailsPage extends ConsumerStatefulWidget {
  static const route = "/series/:id";

  static String toRoute(int id) {
    return "/series/$id";
  }

  final String seriesId;

  const TvDetailsPage({super.key, required this.seriesId});

  @override
  ConsumerState<ConsumerStatefulWidget> createState() {
    return _TvDetailsPageState(seriesId: seriesId);
  }
}

class _TvDetailsPageState extends ConsumerState<TvDetailsPage> {
  final String seriesId;

  _TvDetailsPageState({required this.seriesId});
  Future<String>? _pendingFuture;

  @override
  void initState() {
    super.initState();
  }

  @override
  Widget build(BuildContext context) {
    var seriesDetails = ref.watch(seriesDetailsProvider(seriesId));
    return FutureBuilder(
        // We listen to the pending operation, to update the UI accordingly.
        future: _pendingFuture,
        builder: (context, snapshot) {
          return seriesDetails.when(
              data: (details) {
                Map<int, List<Widget>> m = Map();
                for (final ep in details.episodes!) {
                  var w = Container(
                    alignment: Alignment.topLeft,
                    child: Row(
                      children: [
                        SizedBox(
                          width: 70,
                          child: Text("第 ${ep.episodeNumber} 集"),
                        ),
                        SizedBox(
                          width: 100,
                          child: Opacity(
                            opacity: 0.5,
                            child: Text("${ep.airDate}"),
                          ),
                        ),
                        Text("${ep.title}", textAlign: TextAlign.left),
                        const Expanded(child: Text("")),
                        IconButton(
                            onPressed: () async {
                              var f = ref
                                  .read(
                                      seriesDetailsProvider(seriesId).notifier)
                                  .searchAndDownload(seriesId, ep.seasonNumber!,
                                      ep.episodeNumber!);
                              setState(() {
                                _pendingFuture = f;
                              });
                              if (!Utils.showError(context, snapshot)) {
                                var name = await f;
                                Utils.showSnakeBar(context, "开始下载: $name");
                              }
                            },
                            icon: const Icon(Icons.search))
                      ],
                    ),
                  );
                  if (m[ep.seasonNumber] == null) {
                    m[ep.seasonNumber!] = List.empty(growable: true);
                  }
                  m[ep.seasonNumber!]!.add(w);
                }
                List<ExpansionTile> list = List.empty(growable: true);
                for (final k in m.keys.toList().reversed) {
                  var seasonList = ExpansionTile(
                    tilePadding: const EdgeInsets.fromLTRB(10, 0, 10, 0),
                    childrenPadding: const EdgeInsets.fromLTRB(50, 0, 50, 0),
                    initiallyExpanded: k == 0 ? false : true,
                    title: Text("第 $k 季"),
                    children: m[k]!,
                  );
                  list.add(seasonList);
                }
                return ListView(
                  children: [
                    Card(
                      margin: const EdgeInsets.all(4),
                      clipBehavior: Clip.hardEdge,
                      child: Row(
                        children: <Widget>[
                          Flexible(
                            child: SizedBox(
                              width: 150,
                              height: 200,
                              child: Image.network(
                                APIs.tmdbImgBaseUrl + details!.posterPath!,
                                fit: BoxFit.contain,
                              ),
                            ),
                          ),
                          Expanded(
                            child: Row(
                              children: [
                                Expanded(
                                    child: Column(
                                  crossAxisAlignment: CrossAxisAlignment.start,
                                  children: [
                                    Text(
                                      "${details!.name}",
                                      style: const TextStyle(
                                          fontSize: 14,
                                          fontWeight: FontWeight.bold),
                                    ),
                                    const Text(""),
                                    Text(details!.overview!)
                                  ],
                                )),
                                Column(
                                  children: [
                                    IconButton(
                                        onPressed: () {
                                          ref
                                              .read(seriesDetailsProvider(
                                                      seriesId)
                                                  .notifier)
                                              .delete();
                                          context.pop();
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
                    Column(
                      children: list,
                    ),
                  ],
                );
              },
              error: (err, trace) {
                return Text("$err");
              },
              loading: () => const Center(
                  child: SizedBox(
                      width: 30,
                      height: 30,
                      child: CircularProgressIndicator())));
        });
  }
}
