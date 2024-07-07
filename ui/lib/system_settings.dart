import 'package:dio/dio.dart';
import 'package:flutter/material.dart';
import 'package:ui/APIs.dart';
import 'package:ui/server_response.dart';
import 'package:ui/utils.dart';

class SystemSettingsPage extends StatefulWidget {
  static const route = "/systemsettings";

  const SystemSettingsPage({super.key});
  @override
  State<StatefulWidget> createState() {
    return _SystemSettingsPageState();
  }
}

class _SystemSettingsPageState extends State<SystemSettingsPage> {
  final GlobalKey _formKey = GlobalKey<FormState>();
  final TextEditingController _tmdbApiKeyController = TextEditingController();

  @override
  void initState() {
    super.initState();
    _handleRefresh();
  }

  @override
  Widget build(BuildContext context) {
    return Container(
      padding: const EdgeInsets.fromLTRB(40, 10, 40, 0),
      child: RefreshIndicator(
          onRefresh: _handleRefresh,
          child: Form(
            key: _formKey, //设置globalKey，用于后面获取FormState
            autovalidateMode: AutovalidateMode.onUserInteraction,
            child: Column(
              children: [
                TextFormField(
                  autofocus: true,
                  controller: _tmdbApiKeyController,
                  decoration: const InputDecoration(
                    labelText: "TMDB Api Key",
                    icon: Icon(Icons.key),
                  ),
                  //
                  validator: (v) {
                    return v!.trim().isNotEmpty ? null : "ApiKey 不能为空";
                  },
                ),
                Center(
                  child: Padding(
                    padding: const EdgeInsets.only(top: 28.0),
                    child: ElevatedButton(
                      child: const Padding(
                        padding: EdgeInsets.all(16.0),
                        child: Text("保存"),
                      ),
                      onPressed: () {
                        // 通过_formKey.currentState 获取FormState后，
                        // 调用validate()方法校验用户名密码是否合法，校验
                        // 通过后再提交数据。
                        if ((_formKey.currentState as FormState).validate()) {
                          _submitSettings(context, _tmdbApiKeyController.text);
                        }
                      },
                    ),
                  ),
                )
              ],
            ),
          )),
    );
  }

  Future<void> _handleRefresh() async {
    final dio = Dio();
    var resp = await dio
        .get(APIs.settingsUrl, queryParameters: {"key": APIs.tmdbApiKey});
    var rrr = resp.data as Map<String, dynamic>;
    var data = rrr["data"] as Map<String, dynamic>;
    var key = data[APIs.tmdbApiKey] as String;
    _tmdbApiKeyController.text = key;

    // Fetch new data and update the UI
  }

  void _submitSettings(BuildContext context, String v) async {
    var resp = await Dio().post(APIs.settingsUrl, data: {APIs.tmdbApiKey: v});
    var sp = ServerResponse.fromJson(resp.data as Map<String, dynamic>);
    if (sp.code != 0) {
      if (context.mounted) {
        Utils.showAlertDialog(context, sp.message);
      }
    }
  }
}
