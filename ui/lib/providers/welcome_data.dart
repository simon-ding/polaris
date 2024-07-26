import 'dart:async';

import 'package:flutter_riverpod/flutter_riverpod.dart';
import 'package:quiver/strings.dart';
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

final suggestNameDataProvider = FutureProvider.autoDispose.family(
  (ref, int arg) async {
    final dio = await APIs.getDio();
    var resp = await dio.get(APIs.suggestedTvName + arg.toString());
    var sp = ServerResponse.fromJson(resp.data);
    if (sp.code != 0) {
      throw sp.message;
    }
    return sp.data["name"] as String;
  },
);

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

var searchPageDataProvider = AsyncNotifierProvider.autoDispose
    .family<SearchPageData, List<SearchResult>, String>(SearchPageData.new);

var movieTorrentsDataProvider = AsyncNotifierProvider.autoDispose
    .family<MovieTorrentResource, List<TorrentResource>, String>(
        MovieTorrentResource.new);

class SearchPageData
    extends AutoDisposeFamilyAsyncNotifier<List<SearchResult>, String> {
  List<SearchResult> list = List.empty(growable: true);
  String? q;
  int page = 1;
  @override
  FutureOr<List<SearchResult>> build(String arg) async {
    q = arg;
    if (isBlank(arg)) {
      return List.empty();
    }
    return query(arg, 1);
  }

  FutureOr<List<SearchResult>> query(String q, int page) async {
    final dio = await APIs.getDio();
    var resp = await dio
        .get(APIs.searchUrl, queryParameters: {"query": q, "page": page});

    var rsp = ServerResponse.fromJson(resp.data as Map<String, dynamic>);
    if (rsp.code != 0) {
      throw rsp.message;
    }
    var sp = SearchResponse.fromJson(rsp.data);
    return sp.results ?? List.empty();
  }

  FutureOr<void> queryNextPage() async {
    //state = const AsyncLoading();
    final newState = await AsyncValue.guard(
      () async {
        page++;
        final awaiteddata = await query(q!, page);
        return [...?state.value, ...awaiteddata];
      },
    );
    state = newState;
  }

  Future<void> submit2Watchlist(int tmdbId, int storageId, String resolution,
      String mediaType, String folder) async {
    final dio = await APIs.getDio();
    if (mediaType == "tv") {
      var resp = await dio.post(APIs.watchlistTvUrl, data: {
        "tmdb_id": tmdbId,
        "storage_id": storageId,
        "resolution": resolution,
        "folder": folder
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
        "resolution": resolution,
        "folder": folder
      });
      var sp = ServerResponse.fromJson(resp.data);
      if (sp.code != 0) {
        throw sp.message;
      }
    }
  }
}

class SearchResponse {
  int? page;
  int? totalResults;
  int? totalPage;
  List<SearchResult>? results;

  SearchResponse({this.page, this.totalResults, this.totalPage, this.results});

  factory SearchResponse.fromJson(Map<String, dynamic> json) {
    return SearchResponse(
        page: json["page"],
        totalPage: json["total_page"],
        totalResults: json["total_results"],
        results: json["results"] == null
            ? []
            : (json["results"] as List)
                .map((v) => SearchResult.fromJson(v))
                .toList());
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
  String? status;

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
    this.status,
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
    status = json["status"];
  }
}

class SearchResult {
  SearchResult(
      {required this.backdropPath,
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
      this.inWatchlist});

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
  final bool? inWatchlist;

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
      inWatchlist: json["in_watchlist"],
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

  Future<void> download(TorrentResource res) async {
    var m = res.toJson();
    m["media_id"] = int.parse(mediaId!);

    final dio = await APIs.getDio();
    var resp = await dio.post(APIs.availableMoviesUrl, data: m);
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
