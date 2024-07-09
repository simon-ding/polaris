import 'package:dio/dio.dart';
import 'package:flutter_riverpod/flutter_riverpod.dart';
import 'package:ui/APIs.dart';
import 'package:ui/server_response.dart';

final welcomePageDataProvider = FutureProvider((ref) async {
  var resp = await Dio().get(APIs.watchlistUrl);
  var sp = ServerResponse.fromJson(resp.data);
  List<TvSeries> favList = List.empty(growable: true);
  for (var item in sp.data as List) {
    var tv = TvSeries.fromJson(item);
    favList.add(tv);
  }
  return favList;
});

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

var tmdbApiSettingProvider = FutureProvider(
  (ref) async {
    final dio = Dio();
    var resp = await dio
        .get(APIs.settingsUrl, queryParameters: {"key": APIs.tmdbApiKey});
    var rrr = resp.data as Map<String, dynamic>;
    var data = rrr["data"] as Map<String, dynamic>;
    var key = data[APIs.tmdbApiKey] as String;

    return key;
  },
);

var indexersProvider = FutureProvider((ref) async {
  final dio = Dio();
  var resp = await dio.get(APIs.allIndexersUrl);
  var sp = ServerResponse.fromJson(resp.data);
  if (sp.code != 0) {
    throw sp.message;
  }
  List<Indexer> indexers = List.empty(growable: true);
  for (final item in sp.data as List) {
    indexers.add(Indexer.fromJson(item));
  }
  print("indexers: ${indexers[0].name}");
  return indexers;
});

class Indexer {
  String? name;
  String? url;
  String? apiKey;
  int? id;

  Indexer({this.name, this.url, this.apiKey});

  Indexer.fromJson(Map<String, dynamic> json) {
    name = json['name'];
    url = json['url'];
    apiKey = json['api_key'];
    id = json["id"];
  }
  Map<String, dynamic> toJson() {
    final Map<String, dynamic> data = new Map<String, dynamic>();
    data['name'] = this.name;
    data['url'] = this.url;
    data['api_key'] = this.apiKey;
    return data;
  }
}
