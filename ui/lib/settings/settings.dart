import 'package:flutter/material.dart';
import 'package:flutter_riverpod/flutter_riverpod.dart';
import 'package:ui/settings/auth.dart';
import 'package:ui/settings/downloader.dart';
import 'package:ui/settings/general.dart';
import 'package:ui/settings/indexer.dart';
import 'package:ui/settings/notifier.dart';
import 'package:ui/settings/storage.dart';

class SystemSettingsPage extends ConsumerStatefulWidget {
  static const route = "/settings";

  const SystemSettingsPage({super.key});
  @override
  ConsumerState<ConsumerStatefulWidget> createState() {
    return _SystemSettingsPageState();
  }
}

class _SystemSettingsPageState extends ConsumerState<SystemSettingsPage> {
  @override
  Widget build(BuildContext context) {
    return ListView(
      children: const [
        ExpansionTile(
          expandedAlignment: Alignment.centerLeft,
          childrenPadding: EdgeInsets.fromLTRB(20, 0, 20, 0),
          initiallyExpanded: true,
          title: Text("常规"),
          children: [GeneralSettings()],
        ),
        ExpansionTile(
          expandedAlignment: Alignment.centerLeft,
          childrenPadding: EdgeInsets.fromLTRB(20, 0, 20, 0),
          initiallyExpanded: false,
          title: Text("索引器"),
          children: [IndexerSettings()],
        ),
        ExpansionTile(
          expandedAlignment: Alignment.centerLeft,
          childrenPadding: EdgeInsets.fromLTRB(20, 0, 20, 0),
          initiallyExpanded: false,
          title: Text("下载器"),
          children: [DownloaderSettings()],
        ),
        ExpansionTile(
          expandedAlignment: Alignment.centerLeft,
          childrenPadding: EdgeInsets.fromLTRB(20, 0, 20, 0),
          initiallyExpanded: false,
          title: Text("存储"),
          children: [StorageSettings()],
        ),
        ExpansionTile(
          expandedAlignment: Alignment.centerLeft,
          childrenPadding: EdgeInsets.fromLTRB(20, 0, 20, 0),
          initiallyExpanded: false,
          title: Text("通知客户端"),
          children: [NotifierSettings()],
        ),
        ExpansionTile(
          childrenPadding: EdgeInsets.fromLTRB(20, 0, 20, 0),
          initiallyExpanded: false,
          title: Text("认证"),
          children: [AuthSettings()],
        ),
      ],
    );
  }
}
