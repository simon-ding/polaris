import 'package:flutter/material.dart';
import 'package:flutter_riverpod/flutter_riverpod.dart';
import 'package:ui/settings/auth.dart';
import 'package:ui/settings/downloader.dart';
import 'package:ui/settings/general.dart';
import 'package:ui/settings/importlist.dart';
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
    return SelectionArea(
        child: ListView(
      children: [
        getExpansionTile("常规", const GeneralSettings()),
        getExpansionTile("索引器", const IndexerSettings()),
        getExpansionTile("下载器", const DownloaderSettings()),
        getExpansionTile("存储", const StorageSettings()),
        getExpansionTile("通知客户端", const NotifierSettings()),
        getExpansionTile("监控列表", const Importlist()),
        getExpansionTile("认证", const AuthSettings())
      ],
    ));
  }

  Widget getExpansionTile(String name, Widget body) {
    return ExpansionTile(
      childrenPadding: const EdgeInsets.fromLTRB(20, 0, 20, 0),
      expandedAlignment: Alignment.topLeft,
      initiallyExpanded: false,
      title: Text(name),
      children: [body],
    );
  }
}
