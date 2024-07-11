import 'package:flutter/foundation.dart';

class APIs {
  static final _baseUrl = baseUrl();
  static final searchUrl = "$_baseUrl/api/v1/tv/search";
  static final settingsUrl = "$_baseUrl/api/v1/setting/do";
  static final watchlistUrl = "$_baseUrl/api/v1/tv/watchlist";
  static final seriesDetailUrl = "$_baseUrl/api/v1/tv/series/";
  static final searchAndDownloadUrl = "$_baseUrl/api/v1/indexer/download";
  static final allIndexersUrl = "$_baseUrl/api/v1/indexer/";
  static final addIndexerUrl = "$_baseUrl/api/v1/indexer/add";
  static final delIndexerUrl = "$_baseUrl/api/v1/indexer/del/";
  static final allDownloadClientsUrl = "$_baseUrl/api/v1/downloader";
  static final addDownloadClientUrl = "$_baseUrl/api/v1/downloader/add";
  static final delDownloadClientUrl = "$_baseUrl/api/v1/downloader/del/";
  static final storageUrl = "$_baseUrl/api/v1/storage/";

  static const tmdbImgBaseUrl = "https://image.tmdb.org/t/p/w500/";

  static const tmdbApiKey = "tmdb_api_key";
  static const downloadDirKey = "download_dir";

  static String baseUrl() {
    if (kReleaseMode) {
      return "";
    }
    return "http://127.0.0.1:8080";
  }
}
