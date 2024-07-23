import 'package:dio/dio.dart';
import 'package:flutter/foundation.dart';
import 'package:flutter/material.dart';
import 'package:go_router/go_router.dart';
import 'package:quiver/strings.dart';
import 'package:shared_preferences/shared_preferences.dart';
import 'package:ui/providers/server_response.dart';

class APIs {
  static final _baseUrl = baseUrl();
  static final searchUrl = "$_baseUrl/api/v1/media/search";
  static final settingsUrl = "$_baseUrl/api/v1/setting/do";
  static final settingsGeneralUrl = "$_baseUrl/api/v1/setting/general";
  static final watchlistTvUrl = "$_baseUrl/api/v1/media/tv/watchlist";
  static final watchlistMovieUrl = "$_baseUrl/api/v1/media/movie/watchlist";
  static final availableMoviesUrl = "$_baseUrl/api/v1/media/movie/resources/";
  static final seriesDetailUrl = "$_baseUrl/api/v1/media/record/";
  static final suggestedTvName = "$_baseUrl/api/v1/media/suggest/";
  static final searchAndDownloadUrl = "$_baseUrl/api/v1/indexer/download";
  static final allIndexersUrl = "$_baseUrl/api/v1/indexer/";
  static final addIndexerUrl = "$_baseUrl/api/v1/indexer/add";
  static final delIndexerUrl = "$_baseUrl/api/v1/indexer/del/";
  static final allDownloadClientsUrl = "$_baseUrl/api/v1/downloader";
  static final addDownloadClientUrl = "$_baseUrl/api/v1/downloader/add";
  static final delDownloadClientUrl = "$_baseUrl/api/v1/downloader/del/";
  static final storageUrl = "$_baseUrl/api/v1/storage/";
  static final loginUrl = "$_baseUrl/api/login";
  static final loginSettingUrl = "$_baseUrl/api/v1/setting/auth";
  static final activityUrl = "$_baseUrl/api/v1/activity/";
  static final activityMediaUrl = "$_baseUrl/api/v1/activity/media/";
  static final imagesUrl = "$_baseUrl/api/v1/img";

  static final tmdbImgBaseUrl = "$_baseUrl/api/v1/posters";

  static const tmdbApiKey = "tmdb_api_key";
  static const downloadDirKey = "download_dir";

  static final GlobalKey<NavigatorState> navigatorKey =
      GlobalKey<NavigatorState>();

  static String baseUrl() {
    if (kReleaseMode) {
      return "";
    }
    return "http://127.0.0.1:8080";
  }

  static Dio? gDio;
  static Map<String, String> authHeaders = {};

  static Future<bool> isLoggedIn() async {
    return isNotBlank(await getToken());
  }

  static Future<String> getToken() async {
    var token = authHeaders["Authorization"];
    if (isBlank(token)) {
      final SharedPreferences prefs = await SharedPreferences.getInstance();
      var t = prefs.getString("token");
      if (isNotBlank(t)) {
        authHeaders["Authorization"] = t!;
        token = t;
      }
    }
    return token ?? "";
  }

  static Future<Dio> getDio() async {
    if (gDio != null) {
      return gDio!;
    }
    var token = await getToken();

    var dio = Dio();
    dio.interceptors.add(InterceptorsWrapper(
      onRequest: (options, handler) {
        options.headers['Authorization'] = token;
        return handler.next(options);
      },
      onError: (error, handler) {
        if (error.response?.statusCode != null &&
            error.response?.statusCode! == 403) {
          final context = navigatorKey.currentContext;
          if (context != null) {
            context.go('/login');
            gDio = null;
          }
        }
        return handler.next(error);
      },
    ));
    if (isNotBlank(token)) {
      gDio = dio;
    }
    return dio;
  }

  static Future<void> login(String user, String password) async {
    var resp = await Dio()
        .post(APIs.loginUrl, data: {"user": user, "password": password});

    var sp = ServerResponse.fromJson(resp.data);

    if (sp.code != 0) {
      throw sp.message;
    }
    final SharedPreferences prefs = await SharedPreferences.getInstance();
    var t = sp.data["token"];
    authHeaders["Authorization"] = "Bearer $t";
    prefs.setString("token", "Bearer $t");
  }
}
