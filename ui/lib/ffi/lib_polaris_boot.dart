import 'package:ui/ffi/lib_polaris_boot_stub.dart'
    if (dart.library.html) 'package:ui/ffi/entry/libpolaris_boot_browser.dart'
    if (dart.library.io) 'package:ui/ffi/entry/libpolaris_boot_native.dart';

abstract class LibPolarisBoot {
  static LibPolarisBoot? _instance;

  static LibPolarisBoot get instance {
    _instance ??= LibPolarisBoot();
    return _instance!;
  }

  factory LibPolarisBoot() => create();

  Future<int> start(String cfg);

  Future<void> stop();
}
