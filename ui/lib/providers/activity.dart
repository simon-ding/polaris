import 'dart:async';

import 'package:flutter_riverpod/flutter_riverpod.dart';
import 'package:ui/providers/APIs.dart';
import 'package:ui/providers/server_response.dart';

var activitiesDataProvider =
    AsyncNotifierProvider.autoDispose<ActivityData, List<Activity>>(
        ActivityData.new);

class ActivityData extends AutoDisposeAsyncNotifier<List<Activity>> {
  @override
  FutureOr<List<Activity>> build() async {
    final dio = await APIs.getDio();
    var resp = await dio.get(APIs.activityUrl);
    final sp = ServerResponse.fromJson(resp.data);
    if (sp.code != 0) {
      throw sp.message;
    }
    List<Activity> activities = List.empty(growable: true);
    for (final a in sp.data as List) {
      activities.add(Activity.fromJson(a));
    }
    return activities;
  }

  Future<void> deleteActivity(int id) async {
    final dio = await APIs.getDio();
    var resp = await dio.delete("${APIs.activityUrl}$id");
    final sp = ServerResponse.fromJson(resp.data);
    if (sp.code != 0) {
      throw sp.message;
    }
    ref.invalidateSelf();
  }
}

class Activity {
  Activity({
    required this.id,
    required this.seriesId,
    required this.episodeId,
    required this.sourceTitle,
    required this.date,
    required this.targetDir,
    required this.completed,
    required this.saved,
    required this.inBackgroud,
    required this.progress
  });

  final int? id;
  final int? seriesId;
  final int? episodeId;
  final String? sourceTitle;
  final DateTime? date;
  final String? targetDir;
  final bool? completed;
  final String? saved;
  final bool? inBackgroud;
  final int? progress;

  factory Activity.fromJson(Map<String, dynamic> json) {
    return Activity(
      id: json["id"],
      seriesId: json["series_id"],
      episodeId: json["episode_id"],
      sourceTitle: json["source_title"],
      date: DateTime.tryParse(json["date"] ?? ""),
      targetDir: json["target_dir"],
      completed: json["completed"],
      saved: json["saved"],
      inBackgroud: json["in_backgroud"],
      progress: json["progress"]
    );
  }
}
