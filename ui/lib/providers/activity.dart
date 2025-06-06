import 'dart:async';

import 'package:flutter_riverpod/flutter_riverpod.dart';
import 'package:ui/providers/APIs.dart';
import 'package:ui/providers/server_response.dart';

var activitiesDataProvider = AsyncNotifierProvider.autoDispose
    .family<ActivityData, List<Activity>, ActivityStatus>(ActivityData.new);

var blacklistDataProvider =
    AsyncNotifierProvider.autoDispose<BlacklistData, List<Blacklist>>(
        BlacklistData.new);

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

  Future<void> deleteActivity(int id, bool add2Blacklist) async {
    final dio = APIs.getDio();
    var resp = await dio.post(APIs.activityDeleteUrl, data: {
      "id": id,
      "add_2_blacklist": add2Blacklist,
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
  final num? seedRatio;
  final num? uploadProgress;

  factory Activity.fromJson(Map<String, dynamic> json) {
    return Activity(
        id: json["id"],
        mediaId: json["media_id"],
        episodeId: json["episode_id"],
        sourceTitle: json["source_title"],
        date: DateTime.tryParse(json["date"] ?? DateTime.now().toString()),
        targetDir: json["target_dir"],
        status: json["status"],
        saved: json["saved"],
        progress: json["progress"] ?? 0,
        seedRatio: json["seed_ratio"] ?? 0,
        size: json["size"] ?? 0,
        uploadProgress: json["upload_progress"] ?? 0);
  }
}

class BlacklistData extends AutoDisposeAsyncNotifier<List<Blacklist>> {
  @override
  FutureOr<List<Blacklist>> build() async {
    final dio = APIs.getDio();
    var resp = await dio.get(APIs.blacklistUrl);
    final sp = ServerResponse.fromJson(resp.data);
    if (sp.code != 0) {
      throw sp.message;
    }
    List<Blacklist> blaclklists = List.empty(growable: true);
    for (final a in sp.data as List) {
      blaclklists.add(Blacklist.fromJson(a));
    }
    return blaclklists;
  }

  Future<void> deleteBlacklist(int id) async {
    final dio = APIs.getDio();
    var resp = await dio.delete("${APIs.blacklistUrl}/$id");
    final sp = ServerResponse.fromJson(resp.data);
    if (sp.code != 0) {
      throw sp.message;
    }
  }
}

class Blacklist {
  int? id;
  String? type;
  String? torrentHash;
  String? torrentName;
  int? mediaId;
  String? createTime;

  Blacklist(
      {this.id,
      this.type,
      this.torrentHash,
      this.torrentName,
      this.mediaId,
      this.createTime});

  Blacklist.fromJson(Map<String, dynamic> json) {
    id = json['id'];
    type = json['type'];
    torrentHash = json['torrent_hash'];
    torrentName = json['torrent_name'];
    mediaId = json['media_id'];
    createTime = json['create_time'];
  }

  Map<String, dynamic> toJson() {
    final Map<String, dynamic> data = new Map<String, dynamic>();
    data['id'] = this.id;
    data['type'] = this.type;
    data['torrent_hash'] = this.torrentHash;
    data['torrent_name'] = this.torrentName;
    data['media_id'] = this.mediaId;
    data['create_time'] = this.createTime;
    return data;
  }
}
