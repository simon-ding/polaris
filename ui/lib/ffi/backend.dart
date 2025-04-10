import 'dart:ffi';
import 'dart:io';
import 'dart:isolate';

import 'package:ui/widgets/utils.dart' as Utils;


class FFIBackend {
  final lib = DynamicLibrary.open(libname());

  static String libname() {
    if (Utils.isDesktop()) {
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
