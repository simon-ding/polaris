import 'package:flutter/material.dart';
import 'package:flutter_form_builder/flutter_form_builder.dart';
import 'package:flutter_riverpod/flutter_riverpod.dart';
import 'package:form_builder_validators/form_builder_validators.dart';
import 'package:ui/providers/settings.dart';
import 'package:ui/widgets/progress_indicator.dart';
import 'package:ui/widgets/utils.dart';
import 'package:ui/widgets/widgets.dart';

class ProwlarrSettingPage extends ConsumerStatefulWidget {
  const ProwlarrSettingPage({super.key});
  @override
  ConsumerState<ConsumerStatefulWidget> createState() {
    return ProwlarrSettingState();
  }
}

class ProwlarrSettingState extends ConsumerState<ProwlarrSettingPage> {
  final _formKey = GlobalKey<FormBuilderState>();

  @override
  Widget build(BuildContext context) {
    var ps = ref.watch(prowlarrSettingDataProvider);
    return ps.when(
        data: (v) => FormBuilder(
              key: _formKey, //设置globalKey，用于后面获取FormState
              autovalidateMode: AutovalidateMode.onUserInteraction,
              initialValue: {
                "api_key": v.apiKey,
                "url": v.url,
                "disabled": v.disabled
              },
              child: Column(
                children: [
                  FormBuilderTextField(
                    name: "url",
                    decoration: const InputDecoration(
                        labelText: "Prowlarr URL",
                        icon: Icon(Icons.web),
                        hintText: "http://10.0.0.8:9696"),
                    validator: FormBuilderValidators.required(),
                  ),
                  FormBuilderTextField(
                    name: "api_key",
                    decoration: const InputDecoration(
                        labelText: "API Key",
                        icon: Icon(Icons.key),
                        helperText: "Prowlarr 设置 -> 通用 -> API 密钥"),
                    validator: FormBuilderValidators.required(),
                  ),
                  FormBuilderSwitch(
                    name: "disabled",
                    title: const Text("禁用 Prowlarr"),
                    decoration:
                        InputDecoration(icon: Icon(Icons.do_not_disturb)),
                  ),
                  Center(
                    child: Padding(
                      padding: const EdgeInsets.all(10),
                      child: ElevatedButton(
                          onPressed: () {
                            if (_formKey.currentState!.saveAndValidate()) {
                              var values = _formKey.currentState!.value;
                              var f = ref
                                  .read(prowlarrSettingDataProvider.notifier)
                                  .save(ProwlarrSetting(
                                      apiKey: values["api_key"],
                                      url: values["url"],
                                      disabled: values["disabled"]))
                                  .then((v) => showSnakeBar("更新成功"));
                              showLoadingWithFuture(f);
                            }
                          },
                          child: const Padding(
                            padding: EdgeInsets.all(10),
                            child: Text("保存"),
                          )),
                    ),
                  )
                ],
              ),
            ),
        error: (err, trace) => PoNetworkError(err: err),
        loading: () => const MyProgressIndicator());
  }
}
