import 'package:dio/dio.dart';
import 'package:flutter/material.dart';
import 'package:ui/APIs.dart';
import 'package:ui/server_response.dart';
import 'package:ui/utils.dart';

class TvDetailsPage extends StatefulWidget {
  static const route = "/series/:id";

  static String toRoute(int id) {
    return "/series/$id";
  }

  final String seriesId;

  const TvDetailsPage({super.key, required this.seriesId});

  @override
  State<StatefulWidget> createState() {
    return _TvDetailsPageState(seriesId: seriesId);
  }
}

class _TvDetailsPageState extends State<TvDetailsPage> {
  final String seriesId;

  _TvDetailsPageState({required this.seriesId});

  SeriesDetails? details;
  @override
  Widget build(BuildContext context) {
    _querySeriesDetails(context);

    if (details == null) {
      return const Center(
        child: Text("nothing here"),
      );
    }

    Map<int, List<Widget>> m = Map();
    for (final ep in details!.episodes!) {
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
            IconButton(onPressed: () {}, icon: const Icon(Icons.search))
          ],
        ),
      );
      if (m[ep.seasonNumber] == null) {
        m[ep.seasonNumber!] = List.empty(growable: true);
      }
      m[ep.seasonNumber!]!.add(w);
    }
    List<ExpansionTile> list = List.empty(growable: true);
    for (final k in m.keys) {
      bool _customTileExpanded = false;
      var seasonList = ExpansionTile(
        tilePadding: const EdgeInsets.fromLTRB(10, 0, 10, 0),
        childrenPadding: const EdgeInsets.fromLTRB(50, 0, 50, 0),
        initiallyExpanded: k == 0 ? false : true,
        title: Text("第 $k 季"),
        trailing: Icon(
          _customTileExpanded
              ? Icons.arrow_drop_down_circle
              : Icons.arrow_drop_down,
        ),
        children: m[k]!,
        onExpansionChanged: (bool expanded) {
          setState(() {
            _customTileExpanded = expanded;
          });
        },
      );
      list.add(seasonList);
    }

    return Column(
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
              Flexible(
                child: Column(
                  crossAxisAlignment: CrossAxisAlignment.start,
                  children: [
                    Text(
                      "${details!.name}",
                      style: const TextStyle(
                          fontSize: 14, fontWeight: FontWeight.bold),
                    ),
                    const Text(""),
                    Text(details!.overview!)
                  ],
                ),
              ),
            ],
          ),
        ),
        Expanded(
          child: ListView(
            children: list,
          ),
        ),
      ],
    );
  }

  void _querySeriesDetails(BuildContext context) async {
    if (details != null) {
      return;
    }
    var resp = await Dio().get("${APIs.seriesDetailUrl}$seriesId");
    var rsp = ServerResponse.fromJson(resp.data);
    if (rsp.code != 0 && context.mounted) {
      Utils.showAlertDialog(context, rsp.message);
    }
    setState(() {
      details = SeriesDetails.fromJson(rsp.data);
    });
  }
}

class SeriesDetails {
  int? id;
  int? tmdbId;
  String? name;
  String? originalName;
  String? overview;
  String? path;
  String? posterPath;
  String? createdAt;
  List<Episodes>? episodes;

  SeriesDetails(
      {this.id,
      this.tmdbId,
      this.name,
      this.originalName,
      this.overview,
      this.path,
      this.posterPath,
      this.createdAt,
      this.episodes});

  SeriesDetails.fromJson(Map<String, dynamic> json) {
    id = json['id'];
    tmdbId = json['tmdb_id'];
    name = json['name'];
    originalName = json['original_name'];
    overview = json['overview'];
    path = json['path'];
    posterPath = json['poster_path'];
    createdAt = json['created_at'];
    if (json['episodes'] != null) {
      episodes = <Episodes>[];
      json['episodes'].forEach((v) {
        episodes!.add(Episodes.fromJson(v));
      });
    }
  }
}

class Episodes {
  int? id;
  int? seriesId;
  int? episodeNumber;
  String? title;
  String? airDate;
  int? seasonNumber;
  String? overview;

  Episodes(
      {this.id,
      this.seriesId,
      this.episodeNumber,
      this.title,
      this.airDate,
      this.seasonNumber,
      this.overview});

  Episodes.fromJson(Map<String, dynamic> json) {
    id = json['id'];
    seriesId = json['series_id'];
    episodeNumber = json['episode_number'];
    title = json['title'];
    airDate = json['air_date'];
    seasonNumber = json['season_number'];
    overview = json['overview'];
  }
}
