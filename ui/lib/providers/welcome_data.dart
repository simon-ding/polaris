import 'package:dio/dio.dart';
import 'package:flutter_riverpod/flutter_riverpod.dart';
import 'package:ui/providers/APIs.dart';
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




