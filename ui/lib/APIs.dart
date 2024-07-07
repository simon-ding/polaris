import 'package:flutter/foundation.dart';

class APIs {
  static final _baseUrl = baseUrl();
  static final searchUrl = "$_baseUrl/api/v1/tv/search";
  static final settingsUrl = "$_baseUrl/api/v1/setting/do";
  static final watchlistUrl = "$_baseUrl/api/v1/tv/watchlist";
  static final seriesDetailUrl = "$_baseUrl/api/v1/tv/series/";

  static const tmdbImgBaseUrl = "https://image.tmdb.org/t/p/w500/";

  static const tmdbApiKey = "tmdb_api_key";

  static String baseUrl() {
    if (kReleaseMode) {
      return "";
    }
    return "http://127.0.0.1:8080";
  }
}
