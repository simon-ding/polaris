import 'package:dio/dio.dart';
import 'package:flutter_riverpod/flutter_riverpod.dart';
import 'package:ui/APIs.dart';
import 'package:ui/server_response.dart';

var seriesDetailsProvider = FutureProvider.family((ref, seriesId) async {
  var resp = await Dio().get("${APIs.seriesDetailUrl}$seriesId");
  var rsp = ServerResponse.fromJson(resp.data);
  if (rsp.code != 0) {
    throw rsp.message;
  }
  return SeriesDetails.fromJson(rsp.data);
});

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
