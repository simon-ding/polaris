import 'package:flutter/material.dart';
import 'package:dio/dio.dart';
import 'package:flutter_riverpod/flutter_riverpod.dart';
import 'package:ui/providers/APIs.dart';
import 'package:ui/providers/welcome_data.dart';
import 'package:ui/server_response.dart';
import 'package:ui/utils.dart';

class SearchPage extends ConsumerStatefulWidget {
  const SearchPage({super.key});

  static const route = "/search";

  @override
  ConsumerState<ConsumerStatefulWidget> createState() {
    return _SearchPageState();
  }
}

class _SearchPageState extends ConsumerState<SearchPage> {
  List<dynamic> list = List.empty();

  void _queryResults(BuildContext context, String q) async {
    final dio = Dio();
    var resp = await dio.get(APIs.searchUrl, queryParameters: {"query": q});
    //var dy = jsonDecode(resp.data.toString());

    print("search page results: ${resp.data}");
    var rsp = ServerResponse.fromJson(resp.data as Map<String, dynamic>);
    if (rsp.code != 0 && context.mounted) {
      Utils.showAlertDialog(context, rsp.message);
      return;
    }

    var data = rsp.data as Map<String, dynamic>;
    var results = data["results"] as List<dynamic>;

    setState(() {
      list = results;
    });
  }

  @override
  Widget build(BuildContext context) {
    var cards = List<Widget>.empty(growable: true);
    for (final item in list) {
      var m = SearchResult.fromJson(item);
      cards.add(Card(
          margin: const EdgeInsets.all(4),
          clipBehavior: Clip.hardEdge,
          child: InkWell(
            //splashColor: Colors.blue.withAlpha(30),
            onTap: () {
              //showDialog(context: context, builder: builder)
              _showSubmitDialog(context, m);
            },
            child: Row(
              children: <Widget>[
                Flexible(
                  child: SizedBox(
                    width: 150,
                    height: 200,
                    child: Image.network(
                      APIs.tmdbImgBaseUrl + m.posterPath!,
                      fit: BoxFit.contain,
                    ),
                  ),
                ),
                Flexible(
                  child: Column(
                    crossAxisAlignment: CrossAxisAlignment.start,
                    children: [
                      Text(
                        "${m.name} (${m.firstAirDate?.split("-")[0]})",
                        style: const TextStyle(
                            fontSize: 14, fontWeight: FontWeight.bold),
                      ),
                      const Text(""),
                      Text(m.overview!)
                    ],
                  ),
                )
              ],
            ),
          )));
    }

    return Column(
      children: [
        TextField(
          autofocus: true,
          onSubmitted: (value) => _queryResults(context, value),
          decoration: const InputDecoration(
              labelText: "搜索",
              hintText: "搜索剧集名称",
              prefixIcon: Icon(Icons.search)),
        ),
        Expanded(
            child: ListView(
          children: cards,
        ))
      ],
    );
  }

  Future<void> _showSubmitDialog(BuildContext context, SearchResult item) {
    return showDialog<void>(
        context: context,
        builder: (BuildContext context) {
          return AlertDialog(
            title: const Text('添加剧集'),
            content: Text("是否添加剧集: ${item.name}"),
            actions: <Widget>[
              TextButton(
                style: TextButton.styleFrom(
                  textStyle: Theme.of(context).textTheme.labelLarge,
                ),
                child: const Text('取消'),
                onPressed: () {
                  Navigator.of(context).pop();
                },
              ),
              TextButton(
                style: TextButton.styleFrom(
                  textStyle: Theme.of(context).textTheme.labelLarge,
                ),
                child: const Text('确定'),
                onPressed: () {
                  _submit2Watchlist(context, item.id!);
                  Navigator.of(context).pop();
                },
              ),
            ],
          );
        });
  }

  void _submit2Watchlist(BuildContext context, int id) async {
    var resp = await Dio()
        .post(APIs.watchlistUrl, data: {"id": id, "folder": "/downloads"});
    var sp = ServerResponse.fromJson(resp.data);
    if (sp.code != 0 && context.mounted) {
      Utils.showAlertDialog(context, sp.message);
    }
    ref.refresh(welcomePageDataProvider);
  }
}

class SearchBarApp extends StatefulWidget {
  const SearchBarApp({
    super.key,
    required this.onChanged,
  });

  final ValueChanged<String> onChanged;
  @override
  State<SearchBarApp> createState() => _SearchBarAppState();
}

class _SearchBarAppState extends State<SearchBarApp> {
  @override
  Widget build(BuildContext context) {
    return SearchAnchor(
        builder: (BuildContext context, SearchController controller) {
      return SearchBar(
        controller: controller,
        padding: const WidgetStatePropertyAll<EdgeInsets>(
            EdgeInsets.symmetric(horizontal: 16.0)),
        onSubmitted: (value) => {widget.onChanged(controller.text)},
        leading: const Icon(Icons.search),
      );
    }, suggestionsBuilder: (BuildContext context, SearchController controller) {
      return List<ListTile>.generate(0, (int index) {
        final String item = 'item $index';
        return ListTile(
          title: Text(item),
          onTap: () {
            setState(() {
              controller.closeView(item);
            });
          },
        );
      });
    });
  }
}

class SearchResult {
  String? originalName;
  int? id;
  String? name;
  int? voteCount;
  double? voteAverage;
  String? posterPath;
  String? firstAirDate;
  double? popularity;
  List<int>? genreIds;
  String? originalLanguage;
  String? backdropPath;
  String? overview;
  List<String>? originCountry;

  SearchResult(
      {this.originalName,
      this.id,
      this.name,
      this.voteCount,
      this.voteAverage,
      this.posterPath,
      this.firstAirDate,
      this.popularity,
      this.genreIds,
      this.originalLanguage,
      this.backdropPath,
      this.overview,
      this.originCountry});

  SearchResult.fromJson(Map<String, dynamic> json) {
    originalName = json['original_name'];
    id = json['id'];
    name = json['name'];
    voteCount = json['vote_count'];
    voteAverage = json['vote_average'];
    posterPath = json['poster_path'];
    firstAirDate = json['first_air_date'];
    popularity = json['popularity'];
    genreIds = json['genre_ids'].cast<int>();
    originalLanguage = json['original_language'];
    backdropPath = json['backdrop_path'];
    overview = json['overview'];
    originCountry = json['origin_country'].cast<String>();
  }
}
