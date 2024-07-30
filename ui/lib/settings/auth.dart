import 'package:flutter/material.dart';
import 'package:flutter_form_builder/flutter_form_builder.dart';
import 'package:flutter_riverpod/flutter_riverpod.dart';
import 'package:form_builder_validators/form_builder_validators.dart';
import 'package:ui/providers/login.dart';
import 'package:ui/widgets/progress_indicator.dart';
import 'package:ui/widgets/utils.dart';
import 'package:ui/widgets/widgets.dart';

class AuthSettings extends ConsumerStatefulWidget {
  static const route = "/settings";

  const AuthSettings({super.key});
  @override
  ConsumerState<ConsumerStatefulWidget> createState() {
    return _AuthState();
  }
}

class _AuthState extends ConsumerState<AuthSettings> {
  final _formKey2 = GlobalKey<FormBuilderState>();
  bool? _enableAuth;

  @override
  Widget build(BuildContext context) {
    var authData = ref.watch(authSettingProvider);
    return authData.when(
        data: (data) {
          if (_enableAuth == null) {
            setState(() {
              _enableAuth = data.enable;
            });
          }
          return FormBuilder(
              key: _formKey2,
              initialValue: {
                "user": data.user,
                "password": data.password,
                "enable": data.enable
              },
              child: Column(
                children: [
                  FormBuilderSwitch(
                      name: "enable",
                      title: const Text("开启认证"),
                      onChanged: (v) {
                        setState(() {
                          _enableAuth = v;
                        });
                      }),
                  _enableAuth!
                      ? Column(
                          children: [
                            FormBuilderTextField(
                                name: "user",
                                autovalidateMode:
                                    AutovalidateMode.onUserInteraction,
                                validator: FormBuilderValidators.required(),
                                decoration: Commons.requiredTextFieldStyle(
                                  text: "用户名",
                                  icon: const Icon(Icons.account_box),
                                )),
                            FormBuilderTextField(
                                name: "password",
                                obscureText: true,
                                enableSuggestions: false,
                                autocorrect: false,
                                autovalidateMode:
                                    AutovalidateMode.onUserInteraction,
                                validator: FormBuilderValidators.required(),
                                decoration: Commons.requiredTextFieldStyle(
                                  text: "密码",
                                  icon: const Icon(Icons.password),
                                ))
                          ],
                        )
                      : const Column(),
                  Center(
                      child: ElevatedButton(
                          child: const Text("保存"),
                          onPressed: () {
                            if (_formKey2.currentState!.saveAndValidate()) {
                              var values = _formKey2.currentState!.value;
                              var f = ref
                                  .read(authSettingProvider.notifier)
                                  .updateAuthSetting(_enableAuth!,
                                      values["user"], values["password"])
                                  .then((v) {
                                showSnakeBar("更新成功");
                              });
                              showLoadingWithFuture(f);
                            }
                          }))
                ],
              ));
        },
        error: (err, trace) => Text("$err"),
        loading: () => const MyProgressIndicator());
  }
}
