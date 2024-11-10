import 'package:flutter/material.dart';
import 'package:flutter_form_builder/flutter_form_builder.dart';
import 'package:flutter_riverpod/flutter_riverpod.dart';
import 'package:form_builder_validators/form_builder_validators.dart';
import 'package:ui/settings/prowlarr.dart';
import 'package:ui/widgets/widgets.dart';

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
    final _formKey = GlobalKey<FormBuilderState>();
    var _enableAuth = false;
    String selectImpl = "transmission";

    return ExpansionTile(
      childrenPadding: EdgeInsets.only(left: 100, right: 20),
      initiallyExpanded: true,
      title: Text("第二步：下载客户端", style: TextStyle(fontWeight: FontWeight.bold)),
      children: [
        StatefulBuilder(builder: (BuildContext context, StateSetter setState) {
          return FormBuilder(
              key: _formKey,
              initialValue: {
                "name": "client.name",
                "url": "client.url",
                "user": "client.user",
                "password": "client.password",
                "impl": selectImpl,
              },
              child: Column(
                children: [
                  FormBuilderDropdown<String>(
                    name: "impl",
                    decoration: const InputDecoration(labelText: "类型"),
                    items: const [
                      DropdownMenuItem(
                          value: "transmission", child: Text("Transmission")),
                      DropdownMenuItem(
                          value: "qbittorrent", child: Text("qBittorrent")),
                    ],
                    validator: FormBuilderValidators.required(),
                  ),
                  FormBuilderTextField(
                      name: "name",
                      decoration: const InputDecoration(labelText: "名称"),
                      validator: FormBuilderValidators.required(),
                      autovalidateMode: AutovalidateMode.onUserInteraction),
                  FormBuilderTextField(
                    name: "url",
                    decoration: const InputDecoration(
                        labelText: "地址", hintText: "http://127.0.0.1:9091"),
                    autovalidateMode: AutovalidateMode.onUserInteraction,
                    validator: FormBuilderValidators.required(),
                  ),
                  StatefulBuilder(
                      builder: (BuildContext context, StateSetter setState) {
                    return Column(
                      children: [
                        FormBuilderSwitch(
                            name: "auth",
                            title: const Text("需要认证"),
                            initialValue: _enableAuth,
                            onChanged: (v) {
                              setState(() {
                                _enableAuth = v!;
                              });
                            }),
                        _enableAuth
                            ? Column(
                                children: [
                                  FormBuilderTextField(
                                      name: "user",
                                      decoration:
                                          Commons.requiredTextFieldStyle(
                                              text: "用户"),
                                      validator:
                                          FormBuilderValidators.required(),
                                      autovalidateMode:
                                          AutovalidateMode.onUserInteraction),
                                  FormBuilderTextField(
                                      name: "password",
                                      decoration:
                                          Commons.requiredTextFieldStyle(
                                              text: "密码"),
                                      validator:
                                          FormBuilderValidators.required(),
                                      obscureText: true,
                                      autovalidateMode:
                                          AutovalidateMode.onUserInteraction),
                                ],
                              )
                            : Container()
                      ],
                    );
                  })
                ],
              ));
        }),
      ],
    );
  }
}
