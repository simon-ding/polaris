
import 'package:flutter/material.dart';
import 'package:dio/dio.dart';
import 'package:ui/APIs.dart';

class SearchPage extends StatefulWidget {
  const SearchPage({super.key});

  @override
  State<StatefulWidget> createState() {
    return _SearchPageState();
  }
}

class _SearchPageState extends State<SearchPage> {
  List<dynamic> list = List.empty();
  final tmdbImgBaseUrl = "https://image.tmdb.org/t/p/w500/";

  void _queryResults(String q) async {
    final dio = Dio();
    var resp = await dio.get(APIs.searchUrl, queryParameters: {"query": q});
    //var dy = jsonDecode(resp.data.toString());

    print("search page results: ${resp.data}");

    var rsp = resp.data as Map<String, dynamic>;
    var data = rsp["data"] as Map<String, dynamic>;
    var results = data["results"] as List<dynamic>;

    setState(() {
      list = results;
    });
  }

  @override
  Widget build(BuildContext context) {
    var cards = List<Widget>.empty(growable: true);
    for (final item in list) {
      var m = item as Map<String, dynamic>;
      cards.add(Card(
        margin: const EdgeInsets.all(10),
        child: Row(
          children: <Widget>[
            Flexible(
              child: SizedBox(
                width: 150,
                height: 200,
                child: Image.network(
                  tmdbImgBaseUrl + m["poster_path"],
                  fit: BoxFit.contain,
                ),
              ),
            ),
            Flexible(
              child: ListTile(title: Text(m["name"])),
            ),
            Flexible(child: Text(m["overview"]))
          ],
        ),
      ));
    }

    return Expanded(
      child: Column(
        children: [
          TextField(
            autofocus: true,
            onSubmitted: (value) => _queryResults(value),
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
      ),
    );
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
