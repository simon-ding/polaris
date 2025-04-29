import 'package:ui/ffi/lib_polaris_boot.dart';

LibPolarisBoot create() {
  return LibpolarisBootBrowser(); 
}

class LibpolarisBootBrowser implements LibPolarisBoot {
  @override
  Future<int> start(String cfg) async{
    return 0;
  }

  @override
  Future<void> stop() async{
  }

}