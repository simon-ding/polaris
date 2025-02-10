import 'dart:async';

import 'package:flutter_riverpod/flutter_riverpod.dart';
import 'package:ui/providers/APIs.dart';
import 'package:ui/providers/server_response.dart';

var activitiesDataProvider = AsyncNotifierProvider.autoDispose
    .family<ActivityData, List<Activity>, ActivityStatus>(ActivityData.new);

var mediaHistoryDataProvider = FutureProvider.autoDispose.family(
  (ref, arg) async {
    final dio = await APIs.getDio();
    var resp = await dio.get("${APIs.activityMediaUrl}$arg");
    final sp = ServerResponse.fromJson(resp.data);
    if (sp.code != 0) {
      throw sp.message;
    }
    List<Activity> activities = List.empty(growable: true);
    for (final a in sp.data as List) {
      activities.add(Activity.fromJson(a));
    }
    return activities;
  },
);

enum ActivityStatus {
  active,
  seeding,
  archive,
}

class ActivityData
    extends AutoDisposeFamilyAsyncNotifier<List<Activity>, ActivityStatus> {
  @override
  FutureOr<List<Activity>> build(ActivityStatus arg) async {
    final dio = APIs.getDio();

    var status = arg == ActivityStatus.archive
        ? "archive"
        : "active"; //archive or active
    var resp =
        await dio.get(APIs.activityUrl, queryParameters: {"status": status});
    final sp = ServerResponse.fromJson(resp.data);
    if (sp.code != 0) {
      throw sp.message;
    }
    List<Activity> activities = List.empty(growable: true);
    for (final a in sp.data as List) {
      var activity = Activity.fromJson(a);
      if (arg == ActivityStatus.archive) {
        activities.add(activity);
      } else {
        if (arg == ActivityStatus.active && activity.status != "seeding") {
          activities.add(activity);
        } else if (arg == ActivityStatus.seeding &&
            activity.status == "seeding") {
          activities.add(activity);
        }
      }
    }

    if (status == "active") {
      //refresh active downloads
      final _timer = Timer(const Duration(seconds: 5),
          ref.invalidateSelf); //Periodically Refresh
      ref.onDispose(_timer.cancel);
    }
    return activities;
  }

  Future<void> deleteActivity(int id) async {
    final dio = APIs.getDio();
    var resp = await dio.post(APIs.activityDeleteUrl, data: {
      "id": id,
      "add_2_blacklist": false,
    });
    final sp = ServerResponse.fromJson(resp.data);
    if (sp.code != 0) {
      throw sp.message;
    }
    ref.invalidateSelf();
  }
}

class Activity {
  Activity(
      {required this.id,
      required this.mediaId,
      required this.episodeId,
      required this.sourceTitle,
      required this.date,
      required this.targetDir,
      required this.status,
      required this.saved,
      required this.progress,
      required this.size,
      required this.seedRatio,
      required this.uploadProgress});

  final int? id;
  final int? mediaId;
  final int? episodeId;
  final String? sourceTitle;
  final DateTime? date;
  final String? targetDir;
  final String? status;
  final String? saved;
  final int? progress;
  final int? size;
  final double seedRatio;
  final double uploadProgress;

  factory Activity.fromJson(Map<String, dynamic> json) {
    return Activity(
        id: json["id"],
        mediaId: json["media_id"],
        episodeId: json["episode_id"],
        sourceTitle: json["source_title"],
        date: DateTime.tryParse(json["date"] ?? ""),
        targetDir: json["target_dir"],
        status: json["status"],
        saved: json["saved"],
        progress: json["progress"],
        seedRatio: json["seed_ratio"],
        size: json["size"],
        uploadProgress: json["upload_progress"]);
  }
}
