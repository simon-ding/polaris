import 'package:flutter/material.dart';
import 'package:flutter_riverpod/flutter_riverpod.dart';
import 'package:go_router/go_router.dart';
import 'package:ui/providers/APIs.dart';
import 'package:ui/providers/welcome_data.dart';
import 'package:ui/search_page/submit_dialog.dart';
import 'package:ui/widgets/progress_indicator.dart';
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
        error: (err, trace) => [PoError(msg: "网络错误，请确认TMDB Key正确配置，并且能够正常连接到TMDB网站", err: err)],
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
    return showDialog<void>(
        context: context,
        builder: (BuildContext context) {
          return SubmitSearchResult(
            item: item,
            query: widget.query!,
          );
        });
  }
}
