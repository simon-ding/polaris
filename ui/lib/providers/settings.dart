import 'dart:async';

import 'package:dio/dio.dart';
import 'package:flutter_riverpod/flutter_riverpod.dart';
import 'package:quiver/strings.dart';
import 'package:ui/providers/APIs.dart';
import 'package:ui/server_response.dart';

var tmdbApiSettingProvider =
    AsyncNotifierProvider<TmdbApiSetting, String>(TmdbApiSetting.new);

var indexersProvider =
    AsyncNotifierProvider<IndexerSetting, List<Indexer>>(IndexerSetting.new);

var dwonloadClientsProvider =
    AsyncNotifierProvider<DownloadClientSetting, List<DownloadClient>>(
        DownloadClientSetting.new);

class TmdbApiSetting extends AsyncNotifier<String> {
  @override
  FutureOr<String> build() async {
    final dio = Dio();
    var resp = await dio
        .get(APIs.settingsUrl, queryParameters: {"key": APIs.tmdbApiKey});
    var rrr = ServerResponse.fromJson(resp.data);
    if (rrr.code != 0) {
      throw rrr.message;
    }
    var data = rrr.data as Map<String, dynamic>;
    var key = data[APIs.tmdbApiKey] as String;

    return key;
  }

  Future<void> submitSettings(String v) async {
    var resp = await Dio().post(APIs.settingsUrl, data: {APIs.tmdbApiKey: v});
    var sp = ServerResponse.fromJson(resp.data as Map<String, dynamic>);
    if (sp.code != 0) {
      throw sp.message;
    }
  }
}

class IndexerSetting extends AsyncNotifier<List<Indexer>> {
  final dio = Dio();

  @override
  FutureOr<List<Indexer>> build() async {
    var resp = await dio.get(APIs.allIndexersUrl);
    var sp = ServerResponse.fromJson(resp.data);
    if (sp.code != 0) {
      throw sp.message;
    }
    List<Indexer> indexers = List.empty(growable: true);
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
    var resp = await dio.post(APIs.addIndexerUrl, data: indexer.toJson());
    var sp = ServerResponse.fromJson(resp.data);
    if (sp.code != 0) {
      throw sp.message;
    }
    ref.invalidateSelf();
  }

  Future<void> deleteIndexer(int id) async {
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
    final Map<String, dynamic> data = new Map<String, dynamic>();
    data['name'] = this.name;
    data['url'] = this.url;
    data['api_key'] = this.apiKey;
    return data;
  }
}

class DownloadClientSetting extends AsyncNotifier<List<DownloadClient>> {
  final dio = Dio();

  @override
  FutureOr<List<DownloadClient>> build() async {
    var resp = await dio.get(APIs.allDownloadClientsUrl);
    var sp = ServerResponse.fromJson(resp.data);
    if (sp.code != 0) {
      throw sp.message;
    }
    List<DownloadClient> indexers = List.empty(growable: true);
    for (final item in sp.data as List) {
      indexers.add(DownloadClient.fromJson(item));
    }
    return indexers;
  }

  Future<void> addDownloadClients(String name, String url) async {
    if (name.isEmpty || url.isEmpty) {
      return;
    }
    var dio = Dio();
    var resp = await dio.post(APIs.addDownloadClientUrl, data: {
      "name": name,
      "url": url,
    });
    var sp = ServerResponse.fromJson(resp.data);
    if (sp.code != 0) {
      throw sp.message;
    }
    ref.invalidateSelf();
  }

  Future<void> deleteDownloadClients(int id) async {
    var dio = Dio();
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
  bool? removeCompletedDownloads;
  bool? removeFailedDownloads;

  DownloadClient(
      {this.id,
      this.enable,
      this.name,
      this.implementation,
      this.url,
      this.removeCompletedDownloads,
      this.removeFailedDownloads});

  DownloadClient.fromJson(Map<String, dynamic> json) {
    id = json['id'];
    enable = json['enable'];
    name = json['name'];
    implementation = json['implementation'];
    url = json['url'];
    removeCompletedDownloads = json['remove_completed_downloads'];
    removeFailedDownloads = json['remove_failed_downloads'];
  }

  Map<String, dynamic> toJson() {
    final Map<String, dynamic> data = new Map<String, dynamic>();
    data['id'] = this.id;
    data['enable'] = this.enable;
    data['name'] = this.name;
    data['implementation'] = this.implementation;
    data['url'] = this.url;
    data['remove_completed_downloads'] = this.removeCompletedDownloads;
    data['remove_failed_downloads'] = this.removeFailedDownloads;
    return data;
  }
}
