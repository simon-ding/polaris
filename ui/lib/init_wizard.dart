import 'package:flutter/material.dart';
import 'package:flutter_form_builder/flutter_form_builder.dart';
import 'package:flutter_riverpod/flutter_riverpod.dart';
import 'package:ui/settings/downloader.dart';
import 'package:ui/settings/prowlarr.dart';
import 'package:ui/settings/storage.dart';

class InitWizard extends ConsumerStatefulWidget {
  const InitWizard({super.key});
  static final String route = "/init_wizard";
  @override
  ConsumerState<ConsumerStatefulWidget> createState() {
    return _InitWizardState();
  }
}

class _InitWizardState extends ConsumerState<InitWizard> {
  @override
  Widget build(BuildContext context) {
    return Scaffold(
      body: SelectionArea(
          child: Container(
        padding: EdgeInsets.all(50),
        child: ListView(
          children: [
            Container(
              alignment: Alignment.center,
              child: Text(
                "Polaris 影视追踪下载",
                style: TextStyle(
                    fontSize: 30,
                    color: Theme.of(context).colorScheme.primary,
                    fontWeight: FontWeight.bold),
              ),
            ),
            Container(
              padding: EdgeInsets.only(left: 10, top: 30, bottom: 30),
              child: Text(
                "设置向导",
                style: TextStyle(
                    fontSize: 20,
                    fontWeight: FontWeight.bold,
                    color: Theme.of(context).colorScheme.primary),
              ),
            ),
            tmdbSetting(),
            downloaderSetting(),
            indexerSetting(),
            storageSetting(),
          ],
        ),
      )),
    );
  }

  Widget tmdbSetting() {
    return ExpansionTile(
      title: Text(
        "第一步：TMDB设置",
        style: TextStyle(fontWeight: FontWeight.bold),
      ),
      childrenPadding: EdgeInsets.only(left: 100, right: 20),
      initiallyExpanded: true,
      children: [
        Container(
          alignment: Alignment.topLeft,
          child: Text("TMDB API Key 设置，用来获取各种影视的信息，API Key获取方式参考官网"),
        ),
        FormBuilder(
            child: Column(
          children: [
            FormBuilderTextField(
              name: "tmdb",
              decoration: InputDecoration(labelText: "TMDB API Key"),
            ),
            Center(
                child: Padding(
              padding: EdgeInsets.all(10),
              child: ElevatedButton(onPressed: null, child: Text("保存")),
            ))
          ],
        ))
      ],
    );
  }

  Widget indexerSetting() {
    return ExpansionTile(
      initiallyExpanded: true,
      childrenPadding: EdgeInsets.only(left: 100, right: 20),
      title: Text(
        "第三步：Prowlarr设置",
        style: TextStyle(fontWeight: FontWeight.bold),
      ),
      children: [ProwlarrSettingPage()],
    );
  }

  Widget downloaderSetting() {
    return ExpansionTile(
      childrenPadding: EdgeInsets.only(left: 100, right: 20),
      initiallyExpanded: true,
      title: Text("第二步：下载客户端", style: TextStyle(fontWeight: FontWeight.bold)),
      children: [
        DownloaderSettings(),
      ],
    );
  }

  Widget storageSetting() {
    return ExpansionTile(
      childrenPadding: EdgeInsets.only(left: 100, right: 20),
      title: Text("第四步：存储设置", style: TextStyle(fontWeight: FontWeight.bold)),
      initiallyExpanded: true,
      children: [StorageSettings()],
    );
  }
}
