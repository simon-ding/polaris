import 'dart:async';

import 'package:flutter_riverpod/flutter_riverpod.dart';
import 'package:ui/providers/APIs.dart';
import 'package:ui/providers/server_response.dart';
import 'package:ui/providers/settings.dart';

var mediaDetailsProvider = AsyncNotifierProvider.autoDispose
    .family<SeriesDetailData, SeriesDetails, String>(SeriesDetailData.new);

class SeriesDetailData
    extends AutoDisposeFamilyAsyncNotifier<SeriesDetails, String> {
  String? id;
  @override
  FutureOr<SeriesDetails> build(String arg) async {
    id = arg;
    final dio = await APIs.getDio();
    var resp = await dio.get("${APIs.seriesDetailUrl}$arg");
    var rsp = ServerResponse.fromJson(resp.data);
    if (rsp.code != 0) {
      throw rsp.message;
    }
    return SeriesDetails.fromJson(rsp.data);
  }

  Future<void> delete() async {
    final dio = await APIs.getDio();
    var resp = await dio.delete("${APIs.seriesDetailUrl}$id");
    var rsp = ServerResponse.fromJson(resp.data);
    if (rsp.code != 0) {
      throw rsp.message;
    }
  }

  Future<String> searchAndDownload(
      String seriesId, int seasonNum, int episodeNum) async {
    final dio = await APIs.getDio();
    var resp = await dio.post(APIs.searchAndDownloadUrl, data: {
      "id": int.parse(seriesId),
      "season": seasonNum,
      "episode": episodeNum,
    });
    var sp = ServerResponse.fromJson(resp.data);
    if (sp.code != 0) {
      throw sp.message;
    }
    ref.invalidateSelf();
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
  String? posterPath;
  String? createdAt;
  List<Episodes>? episodes;
  String? resolution;
  int? storageId;
  String? airDate;
  String? mediaType;
  Storage? storage;
  String? targetDir;

  SeriesDetails(
      {this.id,
      this.tmdbId,
      this.name,
      this.originalName,
      this.overview,
      this.posterPath,
      this.createdAt,
      this.resolution,
      this.storageId,
      this.airDate,
      this.episodes,
      this.mediaType,
      this.targetDir,
      this.storage});

  SeriesDetails.fromJson(Map<String, dynamic> json) {
    id = json['id'];
    tmdbId = json['tmdb_id'];
    name = json['name_cn'];
    originalName = json['original_name'];
    overview = json['overview'];
    posterPath = json['poster_path'];
    createdAt = json['created_at'];
    resolution = json["resolution"];
    storageId = json["storage_id"];
    airDate = json["air_date"];
    mediaType = json["media_type"];
    storage = Storage.fromJson(json["storage"]);
    targetDir = json["target_dir"];
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
  int? mediaId;
  int? episodeNumber;
  String? title;
  String? airDate;
  int? seasonNumber;
  String? overview;
  String? status;

  Episodes(
      {this.id,
      this.mediaId,
      this.episodeNumber,
      this.title,
      this.airDate,
      this.seasonNumber,
      this.status,
      this.overview});

  Episodes.fromJson(Map<String, dynamic> json) {
    id = json['id'];
    mediaId = json['media_id'];
    episodeNumber = json['episode_number'];
    title = json['title'];
    airDate = json['air_date'];
    seasonNumber = json['season_number'];
    status = json['status'];
    overview = json['overview'];
  }
}

var mediaTorrentsDataProvider = AsyncNotifierProvider.autoDispose
    .family<MediaTorrentResource, List<TorrentResource>, TorrentQuery>(
        MediaTorrentResource.new);

// class TorrentQuery {
//   final String mediaId;
//   final int seasonNumber;
//   final int episodeNumber;
//   TorrentQuery(
//       {required this.mediaId, this.seasonNumber = 0, this.episodeNumber = 0});
//   Map<String, dynamic> toJson() {
//     final Map<String, dynamic> data = <String, dynamic>{};
//     data["id"] = int.parse(mediaId);
//     data["season"] = seasonNumber;
//     data["episode"] = episodeNumber;
//     return data;
//   }
// }

typedef TorrentQuery = ({String mediaId, int seasonNumber, int episodeNumber});

class MediaTorrentResource extends AutoDisposeFamilyAsyncNotifier<
    List<TorrentResource>, TorrentQuery> {
  @override
  FutureOr<List<TorrentResource>> build(TorrentQuery arg) async {
    final dio = await APIs.getDio();
    var resp = await dio.post(APIs.availableTorrentsUrl, data: {
      "id": int.parse(arg.mediaId),
      "season": arg.seasonNumber,
      "episode": arg.episodeNumber
    });
    var rsp = ServerResponse.fromJson(resp.data);
    if (rsp.code != 0) {
      throw rsp.message;
    }
    return (rsp.data as List).map((v) => TorrentResource.fromJson(v)).toList();
  }

  Future<void> download(TorrentResource res) async {
    final data = res.toJson();
    data.addAll({
      "id": int.parse(arg.mediaId),
      "season": arg.seasonNumber,
      "episode": arg.episodeNumber
    });
    final dio = await APIs.getDio();
    var resp = await dio.post(APIs.downloadTorrentUrl, data: data);
    var rsp = ServerResponse.fromJson(resp.data);
    if (rsp.code != 0) {
      throw rsp.message;
    }
  }
}

class TorrentResource {
  TorrentResource({this.name, this.size, this.seeders, this.peers, this.link});

  String? name;
  int? size;
  int? seeders;
  int? peers;
  String? link;

  factory TorrentResource.fromJson(Map<String, dynamic> json) {
    return TorrentResource(
        name: json["name"],
        size: json["size"],
        seeders: json["seeders"],
        peers: json["peers"],
        link: json["link"]);
  }
  Map<String, dynamic> toJson() {
    final Map<String, dynamic> data = <String, dynamic>{};
    data['name'] = name;
    data['size'] = size;
    data["link"] = link;
    return data;
  }
}
