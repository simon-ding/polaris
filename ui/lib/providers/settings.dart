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

  GeneralSetting({this.tmdbApiKey, this.downloadDIr, this.logLevel});

  factory GeneralSetting.fromJson(Map<String, dynamic> json) {
    return GeneralSetting(
        tmdbApiKey: json["tmdb_api_key"],
        downloadDIr: json["download_dir"],
        logLevel: json["log_level"]);
  }

  Map<String, dynamic> toJson() {
    final Map<String, dynamic> data = <String, dynamic>{};
    data['tmdb_api_key'] = tmdbApiKey;
    data['download_dir'] = downloadDIr;
    data["log_level"] = logLevel;
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

  Indexer({this.name, this.url, this.apiKey});

  Indexer.fromJson(Map<String, dynamic> json) {
    name = json['name'];
    url = json['url'];
    apiKey = json['api_key'];
    id = json["id"];
  }
  Map<String, dynamic> toJson() {
    final Map<String, dynamic> data = <String, dynamic>{};
    data['name'] = name;
    data['url'] = url;
    data['api_key'] = apiKey;
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

  DownloadClient(
      {this.id,
      this.enable,
      this.name,
      this.implementation,
      this.url,
      this.user,
      this.password});

  DownloadClient.fromJson(Map<String, dynamic> json) {
    id = json['id'];
    enable = json['enable'];
    name = json['name'];
    implementation = json['implementation'];
    url = json['url'];
    user = json['user'];
    password = json['password'];
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
    this.settings,
    this.isDefault,
  });

  final int? id;
  final String? name;
  final String? implementation;
  final Map<String, dynamic>? settings;
  final bool? isDefault;

  factory Storage.fromJson(Map<String, dynamic> json1) {
    return Storage(
        id: json1["id"],
        name: json1["name"],
        implementation: json1["implementation"],
        settings: json.decode(json1["settings"]),
        isDefault: json1["default"]);
  }

  Map<String, dynamic> toJson() => {
        "id": id,
        "name": name,
        "implementation": implementation,
        "settings": settings,
        "default": isDefault,
      };
}
