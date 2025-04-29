import 'dart:ffi';
import 'dart:io';

import 'package:ffi/ffi.dart';
import 'package:flutter/foundation.dart';
import 'package:quiver/strings.dart';
import 'package:ui/ffi/lib_polaris_boot.dart';

LibPolarisBoot create() => LibpolarisBootNative();

class LibpolarisBootNative implements LibPolarisBoot {
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

  @override
  Future<int> start(String cfg) async {
    var s = lib
        .lookupFunction<StartFunc, StartFunc>('Start')
        ;
    var r = s(cfg.toNativeUtf8());
    if (isNotBlank(r.r1.toDartString())) {
      throw Exception(r.r1.toDartString()); 
    }
    return r.r0;
  }

  @override
  Future<void> stop() async {
        var s = lib
        .lookup<NativeFunction<Void Function()>>('Stop')
        .asFunction<void Function()>();
    return s();
  }
}

typedef StartFunc = StartReturn Function(Pointer<Utf8> cfg);

final class StartReturn extends Struct {
  @Int32()
  external int r0;
  external Pointer<Utf8> r1;
}
