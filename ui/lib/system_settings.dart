import 'package:dio/dio.dart';
import 'package:flutter/material.dart';
import 'package:flutter_riverpod/flutter_riverpod.dart';
import 'package:ui/APIs.dart';
import 'package:ui/providers/welcome_data.dart';
import 'package:ui/server_response.dart';
import 'package:ui/utils.dart';

class SystemSettingsPage extends ConsumerStatefulWidget {
  static const route = "/systemsettings";

  const SystemSettingsPage({super.key});
  @override
  ConsumerState<ConsumerStatefulWidget> createState() {
    return _SystemSettingsPageState();
  }
}

class _SystemSettingsPageState extends ConsumerState<SystemSettingsPage> {
  final GlobalKey _formKey = GlobalKey<FormState>();

  List<dynamic> indexers = List.empty();

  @override
  void initState() {
    super.initState();
  }

  @override
  Widget build(BuildContext context) {
    var key = ref.watch(tmdbApiSettingProvider);

      return key.when(
        data: (data ) => Container(
          padding: const EdgeInsets.fromLTRB(40, 10, 40, 0),
          child: Form(
            key: _formKey, //设置globalKey，用于后面获取FormState
            autovalidateMode: AutovalidateMode.onUserInteraction,
            child: Column(
              children: [
                TextFormField(
                  autofocus: true,
                  initialValue: data,
                  decoration: const InputDecoration(
                    labelText: "TMDB Api Key",
                    icon: Icon(Icons.key),
                  ),
                  //
                  validator: (v) {
                    return v!.trim().isNotEmpty ? null : "ApiKey 不能为空";
                  },
                  onSaved: (newValue) {
                    _submitSettings(context, newValue!);
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
                          (_formKey.currentState as FormState).save();
                        }
                      },
                    ),
                  ),
                )
              ],
            ),
          ),
        ),
        error: (err, trace) => Text("$err"),
        loading: () => const CircularProgressIndicator());
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
