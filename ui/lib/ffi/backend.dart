import 'dart:ffi';
import 'dart:io';
import 'dart:isolate';

import 'package:flutter/foundation.dart';


class FFIBackend {
  final lib = DynamicLibrary.open(libname());

  static String libname() {
    if (!kIsWeb) {
      if (Platform.isWindows) {
        return 'libpolaris.dll';
      } else if (Platform.isLinux) {
        return 'libpolaris.so';
      } else if (Platform.isMacOS) {
        return 'libpolaris.dylib';
      } else {
        throw UnsupportedError(
            'Unsupported platform: ${Platform.operatingSystem}');
      }
    } else {
      return "";
    }
  }

  Future<void> start() async {
    var s = lib
        .lookup<NativeFunction<Void Function()>>('Start')
        .asFunction<void Function()>();
        
    return Isolate.run(s);
  }
  Future<void> stop() async {
    var s = lib
        .lookup<NativeFunction<Void Function()>>('Stop')
        .asFunction<void Function()>();
    return s();
  }
}
