import 'dart:async';

import 'package:flutter_riverpod/flutter_riverpod.dart';
import 'package:ui/providers/APIs.dart';
import 'package:ui/providers/server_response.dart';

var mediaSizeLimiterDataProvider =
    AsyncNotifierProvider.autoDispose<MediaSizeLimiterData, MediaSizeLimiter>(
        MediaSizeLimiterData.new);

class MediaSizeLimiterData extends AutoDisposeAsyncNotifier<MediaSizeLimiter> {
  @override
  FutureOr<MediaSizeLimiter> build() async {
    final dio = APIs.getDio();
    var resp = await dio.get(APIs.mediaSizeLimiterUrl);
    var sp = ServerResponse.fromJson(resp.data);
    if (sp.code != 0) {
      throw sp.message;
    }
    return MediaSizeLimiter.fromJson(sp.data);
  }

  Future<void> submit(MediaSizeLimiter limiter) async {
    final dio = APIs.getDio();
    var resp = await dio.post(APIs.mediaSizeLimiterUrl, data: limiter.toJson());
    var sp = ServerResponse.fromJson(resp.data);
    if (sp.code != 0) {
      throw sp.message;
    }
    ref.invalidateSelf();
  }
}

class MediaSizeLimiter {
  SizeLimiter? tvLimiter;
  SizeLimiter? movieLimiter;

  MediaSizeLimiter({this.tvLimiter, this.movieLimiter});

  MediaSizeLimiter.fromJson(Map<String, dynamic> json) {
    tvLimiter = json['tv_limiter'] != null
        ? SizeLimiter.fromJson(json['tv_limiter'])
        : null;
    movieLimiter = json['movie_limiter'] != null
        ? SizeLimiter.fromJson(json['movie_limiter'])
        : null;
  }

  Map<String, dynamic> toJson() {
    final Map<String, dynamic> data = <String, dynamic>{};
    if (tvLimiter != null) {
      data['tv_limiter'] = tvLimiter!.toJson();
    }
    if (movieLimiter != null) {
      data['movie_limiter'] = movieLimiter!.toJson();
    }
    return data;
  }
}

class SizeLimiter {
  ResLimiter? p720p;
  ResLimiter? p1080p;
  ResLimiter? p2160p;

  SizeLimiter({this.p720p, this.p1080p, this.p2160p});

  SizeLimiter.fromJson(Map<String, dynamic> json) {
    p720p = json['720p'] != null ? ResLimiter.fromJson(json['720p']) : null;
    p1080p = json['1080p'] != null ? ResLimiter.fromJson(json['1080p']) : null;
    p2160p = json['2160p'] != null ? ResLimiter.fromJson(json['2160p']) : null;
  }

  Map<String, dynamic> toJson() {
    final Map<String, dynamic> data = <String, dynamic>{};
    if (p720p != null) {
      data['720p'] = p720p!.toJson();
    }
    if (p1080p != null) {
      data['1080p'] = p1080p!.toJson();
    }
    if (p2160p != null) {
      data['2160p'] = p2160p!.toJson();
    }
    return data;
  }
}

class ResLimiter {
  int? maxSize;
  int? minSize;
  int? preferSize;

  ResLimiter({this.maxSize, this.minSize, this.preferSize});

  ResLimiter.fromJson(Map<String, dynamic> json) {
    maxSize = json['max_size'];
    minSize = json['min_size'];
    preferSize = json['prefer_size'];
  }

  Map<String, dynamic> toJson() {
    final Map<String, dynamic> data = <String, dynamic>{};
    data['max_size'] = maxSize;
    data['min_size'] = minSize;
    data['prefer_size'] = preferSize;
    return data;
  }
}
