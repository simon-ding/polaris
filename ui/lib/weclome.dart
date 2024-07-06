import 'package:dio/dio.dart';
import 'package:flutter/material.dart';
import 'package:go_router/go_router.dart';
import 'package:ui/APIs.dart';
import 'package:ui/server_response.dart';
import 'package:ui/tv_details.dart';

class WelcomePage extends StatefulWidget {
  const WelcomePage({super.key});
  static const route = "/welcome";

  @override
  State<StatefulWidget> createState() {
    return _WeclomePageState();
  }
}

class _WeclomePageState extends State<WelcomePage> {
  var favList = List.empty(growable: true);

  @override
  Widget build(BuildContext context) {
    _onRefresh();
    return GridView.builder(
        itemCount: favList.length,
        gridDelegate:
            const SliverGridDelegateWithFixedCrossAxisCount(crossAxisCount: 4),
        itemBuilder: (context, i) {
          var item = TvSeries.fromJson(favList[i]);
          return Container(
            child: Card(
                margin: const EdgeInsets.all(4),
                clipBehavior: Clip.hardEdge,
                child: InkWell(
                  //splashColor: Colors.blue.withAlpha(30),
                  onTap: () {
                    context.go(TvDetailsPage.toRoute(item.id!));
                    //showDialog(context: context, builder: builder)
                  },
                  child: Column(
                    children: <Widget>[
                      Flexible(
                        child: SizedBox(
                          width: 300,
                          height: 600,
                          child: Image.network(
                            APIs.tmdbImgBaseUrl + item.posterPath!,
                            fit: BoxFit.contain,
                          ),
                        ),
                      ),
                      Flexible(
                        child: Text(
                          item.name!,
                          style: const TextStyle(
                              fontSize: 14, fontWeight: FontWeight.bold),
                        ),
                      )
                    ],
                  ),
                )),
          );
        });
  }

  Future<void> _onRefresh() async {
    if (favList.isNotEmpty) {
      return;
    }
    var resp = await Dio().get(APIs.watchlistUrl);
    var sp = ServerResponse.fromJson(resp.data);
    setState(() {
      favList = sp.data as List;
    });
  }
}

class TvSeries {
  int? id;
  int? tmdbId;
  String? name;
  String? originalName;
  String? overview;
  String? path;
  String? posterPath;

  TvSeries(
      {this.id,
      this.tmdbId,
      this.name,
      this.originalName,
      this.overview,
      this.path,
      this.posterPath});

  TvSeries.fromJson(Map<String, dynamic> json) {
    id = json['id'];
    tmdbId = json['tmdb_id'];
    name = json['name'];
    originalName = json['original_name'];
    overview = json['overview'];
    path = json['path'];
    posterPath = json["poster_path"];
  }
}
