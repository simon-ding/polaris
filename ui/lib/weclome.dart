import 'package:flutter/material.dart';
import 'package:ui/APIs.dart';

class WelcomePage extends StatefulWidget {
  const WelcomePage({super.key});

  @override
  State<StatefulWidget> createState() {
    return _WeclomePageState();
  }
}

class _WeclomePageState extends State<WelcomePage> {
  var favList = List.empty(growable: true);

  @override
  Widget build(BuildContext context) {
    var cards = List<Widget>.empty(growable: true);
    for (final item in favList) {
      var m = item as Map<String, dynamic>;
      cards.add(Card(
          margin: const EdgeInsets.all(4),
          clipBehavior: Clip.hardEdge,
          child: InkWell(
            //splashColor: Colors.blue.withAlpha(30),
            onTap: () {
              //showDialog(context: context, builder: builder)
              debugPrint('Card tapped.');
            },
            child: Row(
              children: <Widget>[
                Flexible(
                  child: SizedBox(
                    width: 150,
                    height: 200,
                    child: Image.network(
                      APIs.tmdbImgBaseUrl + m["poster_path"],
                      fit: BoxFit.contain,
                    ),
                  ),
                ),
                Flexible(
                  child: Column(
                    crossAxisAlignment: CrossAxisAlignment.start,
                    children: [
                      Text(
                        m["name"],
                        style: const TextStyle(
                            fontSize: 14, fontWeight: FontWeight.bold),
                      ),
                      const Text(""),
                      Text(m["overview"])
                    ],
                  ),
                )
              ],
            ),
          )));
    }

    return Expanded(
        child: RefreshIndicator(
            onRefresh: _onRefresh,
            child: Expanded(
                child: ListView(
              children: cards,
            ))));
  }

  Future<void> _onRefresh() async {}
}
