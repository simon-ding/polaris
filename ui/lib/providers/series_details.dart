import 'dart:async';

import 'package:dio/dio.dart';
import 'package:flutter_riverpod/flutter_riverpod.dart';
import 'package:ui/providers/APIs.dart';
import 'package:ui/providers/server_response.dart';

var seriesDetailsProvider = AsyncNotifierProvider.autoDispose
    .family<SeriesDetailData, SeriesDetails, String>(SeriesDetailData.new);

class SeriesDetailData
    extends AutoDisposeFamilyAsyncNotifier<SeriesDetails, String> {
  @override
  FutureOr<SeriesDetails> build(String arg) async {
    var resp = await Dio().get("${APIs.seriesDetailUrl}$arg");
    var rsp = ServerResponse.fromJson(resp.data);
    if (rsp.code != 0) {
      throw rsp.message;
    }
    return SeriesDetails.fromJson(rsp.data);
  }

  Future<String> searchAndDownload(
      String seriesId, int seasonNum, int episodeNum) async {
    var resp = await Dio().post(APIs.searchAndDownloadUrl, data: {
      "id": int.parse(seriesId),
      "season": seasonNum,
      "episode": episodeNum,
    });
    var sp = ServerResponse.fromJson(resp.data);
    if (sp.code != 0) {
      throw sp.message;
    }
    var name = (sp.data as Map<String, dynamic>)["name"];
    return name;
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
