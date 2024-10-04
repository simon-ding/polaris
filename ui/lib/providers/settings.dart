import 'dart:async';
import 'dart:convert';

import 'package:flutter_riverpod/flutter_riverpod.dart';
import 'package:quiver/strings.dart';
import 'package:ui/providers/APIs.dart';
import 'package:ui/providers/server_response.dart';

var settingProvider =
    AsyncNotifierProvider.autoDispose<EditSettingData, GeneralSetting>(
        EditSettingData.new);

var indexersProvider =
    AsyncNotifierProvider.autoDispose<IndexerSetting, List<Indexer>>(
        IndexerSetting.new);

var dwonloadClientsProvider = AsyncNotifierProvider.autoDispose<
    DownloadClientSetting, List<DownloadClient>>(DownloadClientSetting.new);

var storageSettingProvider =
    AsyncNotifierProvider.autoDispose<StorageSettingData, List<Storage>>(
        StorageSettingData.new);

var importlistProvider =
    AsyncNotifierProvider.autoDispose<ImportListData, List<ImportList>>(
        ImportListData.new);

class EditSettingData extends AutoDisposeAsyncNotifier<GeneralSetting> {
  @override
  FutureOr<GeneralSetting> build() async {
    final dio = await APIs.getDio();

    var resp = await dio.get(APIs.settingsGeneralUrl);
    var rrr = ServerResponse.fromJson(resp.data);
    if (rrr.code != 0) {
      throw rrr.message;
    }
    final ss = GeneralSetting.fromJson(rrr.data);
    return ss;
  }

  Future<void> updateSettings(GeneralSetting gs) async {
    final dio = await APIs.getDio();
    var resp = await dio.post(APIs.settingsGeneralUrl, data: gs.toJson());
    var sp = ServerResponse.fromJson(resp.data);
    if (sp.code != 0) {
      throw sp.message;
    }
    ref.invalidateSelf();
  }
}

class GeneralSetting {
  String? tmdbApiKey;
  String? downloadDIr;
  String? logLevel;
  String? proxy;
  bool? enablePlexmatch;
  bool? allowQiangban;
  bool? enableNfo;
  bool? enableAdult;
  String? tvNamingFormat;
  String? movieNamingFormat;

  GeneralSetting(
      {this.tmdbApiKey,
      this.downloadDIr,
      this.logLevel,
      this.proxy,
      this.enablePlexmatch,
      this.enableNfo,
      this.allowQiangban,
      this.tvNamingFormat,
      this.movieNamingFormat,
      this.enableAdult});

  factory GeneralSetting.fromJson(Map<String, dynamic> json) {
    return GeneralSetting(
        tmdbApiKey: json["tmdb_api_key"],
        downloadDIr: json["download_dir"],
        logLevel: json["log_level"],
        proxy: json["proxy"],
        enableAdult: json["enable_adult_content"] ?? false,
        allowQiangban: json["allow_qiangban"] ?? false,
        enableNfo: json["enable_nfo"] ?? false,
        tvNamingFormat: json["tv_naming_format"],
        movieNamingFormat: json["movie_naming_format"],
        enablePlexmatch: json["enable_plexmatch"] ?? false);
  }

  Map<String, dynamic> toJson() {
    final Map<String, dynamic> data = <String, dynamic>{};
    data['tmdb_api_key'] = tmdbApiKey;
    data['download_dir'] = downloadDIr;
    data["log_level"] = logLevel;
    data["proxy"] = proxy;
    data["enable_plexmatch"] = enablePlexmatch;
    data["allow_qiangban"] = allowQiangban;
    data["enable_nfo"] = enableNfo;
    data["enable_adult_content"] = enableAdult;
    data["tv_naming_format"] = tvNamingFormat;
    data["movie_naming_format"] = movieNamingFormat;
    return data;
  }
}

class IndexerSetting extends AutoDisposeAsyncNotifier<List<Indexer>> {
  @override
  FutureOr<List<Indexer>> build() async {
    final dio = await APIs.getDio();
    var resp = await dio.get(APIs.allIndexersUrl);
    var sp = ServerResponse.fromJson(resp.data);
    if (sp.code != 0) {
      throw sp.message;
    }
    List<Indexer> indexers = List.empty(growable: true);
    if (sp.data == null) {
      return indexers;
    }
    for (final item in sp.data as List) {
      indexers.add(Indexer.fromJson(item));
    }
    return indexers;
  }

  Future<void> addIndexer(Indexer indexer) async {
    if (isBlank(indexer.name) ||
        isBlank(indexer.url) ||
        isBlank(indexer.apiKey)) {
      return;
    }
    final dio = await APIs.getDio();
    var resp = await dio.post(APIs.addIndexerUrl, data: indexer.toJson());
    var sp = ServerResponse.fromJson(resp.data);
    if (sp.code != 0) {
      throw sp.message;
    }
    ref.invalidateSelf();
  }

  Future<void> deleteIndexer(int id) async {
    final dio = await APIs.getDio();
    var resp = await dio.delete("${APIs.delIndexerUrl}$id");
    var sp = ServerResponse.fromJson(resp.data);
    if (sp.code != 0) {
      throw sp.message;
    }
    ref.invalidateSelf();
  }
}

class Indexer {
  String? name;
  String? url;
  String? apiKey;
  int? id;
  int? priority;
  double? seedRatio;
  bool? disabled;

  Indexer(
      {this.name,
      this.url,
      this.apiKey,
      this.id,
      this.priority = 50,
      this.seedRatio = 0,
      this.disabled});

  Indexer.fromJson(Map<String, dynamic> json) {
    name = json['name'];
    url = json['url'];
    apiKey = json['api_key'];
    id = json["id"];
    priority = json["priority"];
    seedRatio = json["seed_ratio"] ?? 0;
    disabled = json["disabled"] ?? false;
  }
  Map<String, dynamic> toJson() {
    final Map<String, dynamic> data = <String, dynamic>{};
    data['name'] = name;
    data['url'] = url;
    data['api_key'] = apiKey;
    data["id"] = id;
    data["priority"] = priority;
    data["seed_ratio"] = seedRatio;
    data["disabled"] = disabled;
    return data;
  }
}

class DownloadClientSetting
    extends AutoDisposeAsyncNotifier<List<DownloadClient>> {
  @override
  FutureOr<List<DownloadClient>> build() async {
    final dio = await APIs.getDio();
    var resp = await dio.get(APIs.allDownloadClientsUrl);
    var sp = ServerResponse.fromJson(resp.data);
    if (sp.code != 0) {
      throw sp.message;
    }
    List<DownloadClient> indexers = List.empty(growable: true);
    if (sp.data == null) {
      return indexers;
    }
    for (final item in sp.data as List) {
      indexers.add(DownloadClient.fromJson(item));
    }
    return indexers;
  }

  Future<void> addDownloadClients(DownloadClient client) async {
    if (isBlank(client.name) || isBlank(client.url)) {
      return;
    }
    final dio = await APIs.getDio();
    var resp = await dio.post(APIs.addDownloadClientUrl, data: client.toJson());
    var sp = ServerResponse.fromJson(resp.data);
    if (sp.code != 0) {
      throw sp.message;
    }
    ref.invalidateSelf();
  }

  Future<void> deleteDownloadClients(int id) async {
    final dio = await APIs.getDio();
    var resp = await dio.delete("${APIs.delDownloadClientUrl}$id");
    var sp = ServerResponse.fromJson(resp.data);
    if (sp.code != 0) {
      throw sp.message;
    }
    ref.invalidateSelf();
  }
}

class DownloadClient {
  int? id;
  bool? enable;
  String? name;
  String? implementation;
  String? url;
  String? user;
  String? password;
  bool? removeCompletedDownloads;
  bool? removeFailedDownloads;
  int? priority;
  DownloadClient(
      {this.id,
      this.enable,
      this.name,
      this.implementation,
      this.url,
      this.user,
      this.password,
      this.removeCompletedDownloads = true,
      this.priority = 1,
      this.removeFailedDownloads = true});

  DownloadClient.fromJson(Map<String, dynamic> json) {
    id = json['id'];
    enable = json['enable'];
    name = json['name'];
    implementation = json['implementation'];
    url = json['url'];
    user = json['user'];
    password = json['password'];
    priority = json["priority1"];
    removeCompletedDownloads = json["remove_completed_downloads"] ?? false;
    removeFailedDownloads = json["remove_failed_downloads"] ?? false;
  }

  Map<String, dynamic> toJson() {
    final Map<String, dynamic> data = <String, dynamic>{};
    data['id'] = id;
    data['enable'] = enable;
    data['name'] = name;
    data['implementation'] = implementation;
    data['url'] = url;
    data['user'] = user;
    data['password'] = password;
    data["priority"] = priority;
    data["remove_completed_downloads"] = removeCompletedDownloads;
    data["remove_failed_downloads"] = removeFailedDownloads;
    return data;
  }
}

class StorageSettingData extends AutoDisposeAsyncNotifier<List<Storage>> {
  @override
  FutureOr<List<Storage>> build() async {
    final dio = await APIs.getDio();
    var resp = await dio.get(APIs.storageUrl);
    var sp = ServerResponse.fromJson(resp.data);
    if (sp.code != 0) {
      throw sp.message;
    }
    var data = sp.data as List<dynamic>;
    List<Storage> list = List.empty(growable: true);
    for (final d in data) {
      list.add(Storage.fromJson(d));
    }
    return list;
  }

  Future<void> deleteStorage(int id) async {
    final dio = await APIs.getDio();
    var resp = await dio.delete("${APIs.storageUrl}$id");
    var sp = ServerResponse.fromJson(resp.data);
    if (sp.code != 0) {
      throw sp.message;
    }
    ref.invalidateSelf();
  }

  Future<void> addStorage(Storage s) async {
    final dio = await APIs.getDio();
    var resp = await dio.post(APIs.storageUrl, data: s.toJson());
    var sp = ServerResponse.fromJson(resp.data);
    if (sp.code != 0) {
      throw sp.message;
    }
    ref.invalidateSelf();
  }
}

class Storage {
  Storage({
    this.id,
    this.name,
    this.implementation,
    this.tvPath,
    this.moviePath,
    this.settings,
    this.isDefault,
  });

  final int? id;
  final String? name;
  final String? implementation;
  final String? tvPath;
  final String? moviePath;
  final Map<String, dynamic>? settings;
  final bool? isDefault;

  factory Storage.fromJson(Map<String, dynamic> json1) {
    return Storage(
        id: json1["id"],
        name: json1["name"],
        implementation: json1["implementation"],
        tvPath: json1["tv_path"],
        moviePath: json1["movie_path"],
        settings: json.decode(json1["settings"]),
        isDefault: json1["default"]);
  }

  Map<String, dynamic> toJson() => {
        "id": id,
        "name": name,
        "implementation": implementation,
        "tv_path": tvPath,
        "movie_path": moviePath,
        "settings": settings,
        "default": isDefault,
      };
}

final logFileDataProvider = FutureProvider.autoDispose((ref) async {
  final dio = await APIs.getDio();
  var resp = await dio.get(APIs.logFilesUrl);
  var sp = ServerResponse.fromJson(resp.data);
  if (sp.code != 0) {
    throw sp.message;
  }
  List<LogFile> favList = List.empty(growable: true);
  for (var item in sp.data as List) {
    var tv = LogFile.fromJson(item);
    favList.add(tv);
  }
  return favList;
});

final aboutDataProvider = FutureProvider.autoDispose((ref) async {
  final dio = await APIs.getDio();
  var resp = await dio.get(APIs.aboutUrl);
  var sp = ServerResponse.fromJson(resp.data);
  if (sp.code != 0) {
    throw sp.message;
  }
  return About.fromJson(sp.data);
});

class LogFile {
  String? name;
  int? size;

  LogFile({this.name, this.size});

  factory LogFile.fromJson(Map<String, dynamic> json1) {
    return LogFile(name: json1["name"], size: json1["size"]);
  }
}

class About {
  About({
    required this.chatGroup,
    required this.goVersion,
    required this.homepage,
    required this.intro,
    required this.uptime,
    required this.version,
  });

  final String? chatGroup;
  final String? goVersion;
  final String? homepage;
  final String? intro;
  final Duration? uptime;
  final String? version;

  factory About.fromJson(Map<String, dynamic> json) {
    return About(
      chatGroup: json["chat_group"],
      goVersion: json["go_version"],
      homepage: json["homepage"],
      intro: json["intro"],
      version: json["version"],
      uptime:
          Duration(microseconds: (json["uptime"] / 1000.0 as double).round()),
    );
  }
}

class ImportList {
  final int? id;
  final String? name;
  final String? url;
  final String? qulity;
  final int? storageId;
  final String? type;
  ImportList({
    this.id,
    this.name,
    this.url,
    this.qulity,
    this.storageId,
    this.type,
  });

  factory ImportList.fromJson(Map<String, dynamic> json) {
    return ImportList(
        id: json["id"],
        name: json["name"],
        url: json["url"],
        qulity: json["qulity"],
        type: json["type"],
        storageId: json["storage_id"]);
  }

  Map<String, dynamic> tojson() => {
        "name": name,
        "url": url,
        "qulity": qulity,
        "type": type,
        "storage_id": storageId
      };
}

class ImportListData extends AutoDisposeAsyncNotifier<List<ImportList>> {
  @override
  FutureOr<List<ImportList>> build() async {
    final dio = APIs.getDio();
    var resp = await dio.get(APIs.getAllImportlists);
    var sp = ServerResponse.fromJson(resp.data);
    if (sp.code != 0) {
      throw sp.message;
    }
    List<ImportList> list = List.empty(growable: true);

    for (var item in sp.data as List) {
      var il = ImportList.fromJson(item);
      list.add(il);
    }
    return list;
  }

  addImportlist(ImportList il) async {
    final dio = APIs.getDio();
    var resp = await dio.post(APIs.addImportlistUrl, data: il.tojson());
    var sp = ServerResponse.fromJson(resp.data);
    if (sp.code != 0) {
      throw sp.message;
    }
    ref.invalidateSelf();
  }

  deleteimportlist(int id) async {
    final dio = APIs.getDio();
    var resp = await dio.post(APIs.deleteImportlistUrl, data: {"id": id});
    var sp = ServerResponse.fromJson(resp.data);
    if (sp.code != 0) {
      throw sp.message;
    }
    ref.invalidateSelf();
  }
}
