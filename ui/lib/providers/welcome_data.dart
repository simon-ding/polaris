import 'dart:async';

import 'package:flutter_riverpod/flutter_riverpod.dart';
import 'package:ui/providers/APIs.dart';
import 'package:ui/providers/server_response.dart';

final welcomePageDataProvider = FutureProvider((ref) async {
  final dio = await APIs.getDio();
  var resp = await dio.get(APIs.watchlistUrl);
  var sp = ServerResponse.fromJson(resp.data);
  List<TvSeries> favList = List.empty(growable: true);
  for (var item in sp.data as List) {
    var tv = TvSeries.fromJson(item);
    favList.add(tv);
  }
  return favList;
});

var searchPageDataProvider = AsyncNotifierProvider.autoDispose
    <SearchPageData, List<SearchResult>>(SearchPageData.new);

class SearchPageData extends AutoDisposeAsyncNotifier<List<SearchResult>> {
  
  List<SearchResult> list = List.empty(growable: true);

  @override
  FutureOr<List<SearchResult>> build() async {
    return list;
  }

  Future<void> submit2Watchlist(int id) async {
    final dio = await APIs.getDio();
    var resp = await dio
        .post(APIs.watchlistUrl, data: {"id": id, "folder": "/downloads"});
    var sp = ServerResponse.fromJson(resp.data);
    if (sp.code != 0) {
      throw sp.message;
    }
    ref.invalidate(welcomePageDataProvider);
  }

  void queryResults(String q) async {
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

class SearchResult {
  String? originalName;
  int? id;
  String? name;
  int? voteCount;
  double? voteAverage;
  String? posterPath;
  String? firstAirDate;
  double? popularity;
  List<int>? genreIds;
  String? originalLanguage;
  String? backdropPath;
  String? overview;
  List<String>? originCountry;

  SearchResult(
      {this.originalName,
      this.id,
      this.name,
      this.voteCount,
      this.voteAverage,
      this.posterPath,
      this.firstAirDate,
      this.popularity,
      this.genreIds,
      this.originalLanguage,
      this.backdropPath,
      this.overview,
      this.originCountry});

  SearchResult.fromJson(Map<String, dynamic> json) {
    originalName = json['original_name'];
    id = json['id'];
    name = json['name'];
    voteCount = json['vote_count'];
    voteAverage = json['vote_average'];
    posterPath = json['poster_path'];
    firstAirDate = json['first_air_date'];
    popularity = json['popularity'];
    genreIds = json['genre_ids'].cast<int>();
    originalLanguage = json['original_language'];
    backdropPath = json['backdrop_path'];
    overview = json['overview'];
    originCountry = json['origin_country'].cast<String>();
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
