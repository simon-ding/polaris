import 'dart:async';
import 'dart:convert';

import 'package:flutter_riverpod/flutter_riverpod.dart';
import 'package:ui/providers/APIs.dart';
import 'package:ui/providers/server_response.dart';

final notifiersDataProvider =
    AsyncNotifierProvider.autoDispose<NotifiersData, List<NotifierData>>(
        NotifiersData.new);

class NotifierData {
  int? id;
  String? name;
  String? service;
  Map<String, dynamic>? settings;
  bool? enabled;

  NotifierData({this.id, this.name, this.service, this.enabled, this.settings});

  factory NotifierData.fromJson(Map<String, dynamic> json) {
    return NotifierData(
        id: json["id"],
        name: json["name"],
        service: json["service"],
        enabled: json["enabled"] ?? true,
        settings: json["settings"]);
  }

  Map<String, dynamic> toJson() {
    final data = <String, dynamic>{};
    data["name"] = name;
    data["service"] = service;
    data["enabled"] = enabled;

    data["settings"] = jsonEncode(settings);
    return data;
  }
}

class NotifiersData extends AutoDisposeAsyncNotifier<List<NotifierData>> {
  @override
  FutureOr<List<NotifierData>> build() async {
    final dio = APIs.getDio();
    final resp = await dio.get(APIs.notifierAllUrl);
    final sp = ServerResponse.fromJson(resp.data);
    if (sp.code != 0) {
      throw sp.message;
    }

    return sp.data == null
        ? List.empty()
        : (sp.data as List).map((e) => NotifierData.fromJson(e)).toList();
  }

  Future<void> delete(int id) async {
    final dio = APIs.getDio();
    final resp = await dio.delete(APIs.notifierDeleteUrl + id.toString());
    final sp = ServerResponse.fromJson(resp.data);
    if (sp.code != 0) {
      throw sp.message;
    }
    ref.invalidateSelf();
  }

  Future<void> add(NotifierData n) async {
    final dio = APIs.getDio();
    final resp = await dio.post(APIs.notifierAddUrl, data: n.toJson());
    final sp = ServerResponse.fromJson(resp.data);
    if (sp.code != 0) {
      throw sp.message;
    }
    ref.invalidateSelf();
  }
}
