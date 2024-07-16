import 'dart:async';

import 'package:flutter_riverpod/flutter_riverpod.dart';
import 'package:ui/providers/APIs.dart';
import 'package:ui/providers/server_response.dart';

final tvWatchlistDataProvider = FutureProvider.autoDispose((ref) async {
  final dio = await APIs.getDio();
  var resp = await dio.get(APIs.watchlistTvUrl);
  var sp = ServerResponse.fromJson(resp.data);
  List<MediaDetail> favList = List.empty(growable: true);
  for (var item in sp.data as List) {
    var tv = MediaDetail.fromJson(item);
    favList.add(tv);
  }
  return favList;
});

final movieWatchlistDataProvider = FutureProvider.autoDispose((ref) async {
  final dio = await APIs.getDio();
  var resp = await dio.get(APIs.watchlistMovieUrl);
  var sp = ServerResponse.fromJson(resp.data);
  List<MediaDetail> favList = List.empty(growable: true);
  for (var item in sp.data as List) {
    var tv = MediaDetail.fromJson(item);
    favList.add(tv);
  }
  return favList;
});

var searchPageDataProvider =
    AsyncNotifierProvider.autoDispose<SearchPageData, List<SearchResult>>(
        SearchPageData.new);

var movieTorrentsDataProvider = AsyncNotifierProvider.autoDispose
    .family<MovieTorrentResource, List<TorrentResource>, String>(
        MovieTorrentResource.new);

class SearchPageData extends AutoDisposeAsyncNotifier<List<SearchResult>> {
  List<SearchResult> list = List.empty(growable: true);

  @override
  FutureOr<List<SearchResult>> build() async {
    return list;
  }

  Future<void> submit2Watchlist(
      int tmdbId, int storageId, String resolution, String mediaType) async {
    final dio = await APIs.getDio();
    if (mediaType == "tv") {
      var resp = await dio.post(APIs.watchlistTvUrl, data: {
        "tmdb_id": tmdbId,
        "storage_id": storageId,
        "resolution": resolution
      });
      var sp = ServerResponse.fromJson(resp.data);
      if (sp.code != 0) {
        throw sp.message;
      }
      ref.invalidate(tvWatchlistDataProvider);
    } else {
      var resp = await dio.post(APIs.watchlistMovieUrl, data: {
        "tmdb_id": tmdbId,
        "storage_id": storageId,
        "resolution": resolution
      });
      var sp = ServerResponse.fromJson(resp.data);
      if (sp.code != 0) {
        throw sp.message;
      }
      ref.invalidate(movieWatchlistDataProvider);
    }
  }

  Future<void> queryResults(String q) async {
    list = List.empty(growable: true);
    final dio = await APIs.getDio();
    var resp = await dio.get(APIs.searchUrl, queryParameters: {"query": q});
    //var dy = jsonDecode(resp.data.toString());

    print("search page results: ${resp.data}");
    var rsp = ServerResponse.fromJson(resp.data as Map<String, dynamic>);
    if (rsp.code != 0) {
      throw rsp.message;
    }

    var data = rsp.data as Map<String, dynamic>;
    var results = data["results"] as List<dynamic>;
    for (final r in results) {
      var res = SearchResult.fromJson(r);
      list.add(res);
    }
    ref.invalidateSelf();
  }
}

class MediaDetail {
  int? id;
  int? tmdbId;
  String? mediaType;
  String? name;
  String? originalName;
  String? overview;
  String? posterPath;
  String? createdAt;
  String? resolution;
  int? storageId;
  String? airDate;

  MediaDetail({
    this.id,
    this.tmdbId,
    this.mediaType,
    this.name,
    this.originalName,
    this.overview,
    this.posterPath,
    this.createdAt,
    this.resolution,
    this.storageId,
    this.airDate,
  });

  MediaDetail.fromJson(Map<String, dynamic> json) {
    id = json['id'];
    tmdbId = json['tmdb_id'];
    mediaType = json["media_type"];
    name = json['name_cn'];
    originalName = json['original_name'];
    overview = json['overview'];
    posterPath = json['poster_path'];
    createdAt = json['created_at'];
    resolution = json["resolution"];
    storageId = json["storage_id"];
    airDate = json["air_date"];
  }
}

class SearchResult {
  SearchResult({
    required this.backdropPath,
    required this.id,
    required this.name,
    required this.originalName,
    required this.overview,
    required this.posterPath,
    required this.mediaType,
    required this.adult,
    required this.originalLanguage,
    required this.genreIds,
    required this.popularity,
    required this.firstAirDate,
    required this.voteAverage,
    required this.voteCount,
    required this.originCountry,
  });

  final String? backdropPath;
  final int? id;
  final String? name;
  final String? originalName;
  final String? overview;
  final String? posterPath;
  final String? mediaType;
  final bool? adult;
  final String? originalLanguage;
  final List<int> genreIds;
  final double? popularity;
  final DateTime? firstAirDate;
  final double? voteAverage;
  final int? voteCount;
  final List<String> originCountry;

  factory SearchResult.fromJson(Map<String, dynamic> json) {
    return SearchResult(
      backdropPath: json["backdrop_path"],
      id: json["id"],
      name: json["name"],
      originalName: json["original_name"],
      overview: json["overview"],
      posterPath: json["poster_path"],
      mediaType: json["media_type"],
      adult: json["adult"],
      originalLanguage: json["original_language"],
      genreIds: json["genre_ids"] == null
          ? []
          : List<int>.from(json["genre_ids"]!.map((x) => x)),
      popularity: json["popularity"],
      firstAirDate: DateTime.tryParse(json["first_air_date"] ?? ""),
      voteAverage: json["vote_average"],
      voteCount: json["vote_count"],
      originCountry: json["origin_country"] == null
          ? []
          : List<String>.from(json["origin_country"]!.map((x) => x)),
    );
  }
}

class MovieTorrentResource
    extends AutoDisposeFamilyAsyncNotifier<List<TorrentResource>, String> {
  String? mediaId;
  @override
  FutureOr<List<TorrentResource>> build(String id) async {
    mediaId = id;
    final dio = await APIs.getDio();
    var resp = await dio.get(APIs.availableMoviesUrl + id);
    var rsp = ServerResponse.fromJson(resp.data);
    if (rsp.code != 0) {
      throw rsp.message;
    }
    return (rsp.data as List).map((v) => TorrentResource.fromJson(v)).toList();
  }

  Future<void> download(String link) async {
    final dio = await APIs.getDio();
    var resp = await dio.post(APIs.availableMoviesUrl,
        data: {"media_id": int.parse(mediaId!), "link": link});
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
}
