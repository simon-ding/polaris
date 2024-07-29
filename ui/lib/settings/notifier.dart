import 'package:flutter/material.dart';
import 'package:flutter_form_builder/flutter_form_builder.dart';
import 'package:flutter_riverpod/flutter_riverpod.dart';
import 'package:form_builder_validators/form_builder_validators.dart';
import 'package:ui/providers/notifier.dart';
import 'package:ui/settings/dialog.dart';
import 'package:ui/widgets/progress_indicator.dart';
import 'package:ui/widgets/widgets.dart';

class NotifierSettings extends ConsumerStatefulWidget {
  static const route = "/settings";

  const NotifierSettings({super.key});
  @override
  ConsumerState<ConsumerStatefulWidget> createState() {
    return _NotifierState();
  }
}

class _NotifierState extends ConsumerState<NotifierSettings> {
  @override
  Widget build(BuildContext context) {
    final notifierData = ref.watch(notifiersDataProvider);
    return notifierData.when(
        data: (v) => Wrap(
              children: List.generate(v.length + 1, (i) {
                if (i < v.length) {
                  final client = v[i];
                  return SettingsCard(
                    child: Text("${client.name!} (${client.service})"),
                    onTap: () => showNotifierDetails(client),
                  );
                }
                return SettingsCard(
                    onTap: () => showNotifierDetails(NotifierData()),
                    child: const Icon(Icons.add));
              }),
            ),
        error: (err, trace) => Text("$err"),
        loading: () => const MyProgressIndicator());
  }

    Future<void> showNotifierDetails(NotifierData notifier) {
    final _formKey = GlobalKey<FormBuilderState>();

    var body = FormBuilder(
      key: _formKey,
      initialValue: {
        "name": notifier.name,
        "service": notifier.service,
        "enabled": notifier.enabled ?? true,
        "app_token":
            notifier.settings != null ? notifier.settings!["app_token"] : "",
        "user_key":
            notifier.settings != null ? notifier.settings!["user_key"] : "",
      },
      child: Column(
        children: [
          FormBuilderDropdown(
            name: "service",
            decoration: const InputDecoration(labelText: "类型"),
            items: const [
              DropdownMenuItem(value: "pushover", child: Text("Pushover")),
            ],
          ),
          FormBuilderTextField(
            name: "name",
            decoration: Commons.requiredTextFieldStyle(text: "名称"),
            autovalidateMode: AutovalidateMode.onUserInteraction,
            validator: FormBuilderValidators.required(),
          ),
          FormBuilderTextField(
            name: "app_token",
            decoration: Commons.requiredTextFieldStyle(text: "APP密钥"),
            autovalidateMode: AutovalidateMode.onUserInteraction,
            validator: FormBuilderValidators.required(),
          ),
          FormBuilderTextField(
            name: "user_key",
            decoration: Commons.requiredTextFieldStyle(text: "用户密钥"),
            autovalidateMode: AutovalidateMode.onUserInteraction,
            validator: FormBuilderValidators.required(),
          ),
          FormBuilderSwitch(name: "enabled", title: const Text("启用"))
        ],
      ),
    );
    onDelete() async {
      return ref.read(notifiersDataProvider.notifier).delete(notifier.id!);
    }

    onSubmit() async {
      if (_formKey.currentState!.saveAndValidate()) {
        var values = _formKey.currentState!.value;
        return ref.read(notifiersDataProvider.notifier).add(NotifierData(
                name: values["name"],
                service: values["service"],
                enabled: values["enabled"],
                settings: {
                  "app_token": values["app_token"],
                  "user_key": values["user_key"]
                }));
      } else {
        throw "validation_error";
      }
    }

    return showSettingDialog(context,
        "通知客户端", notifier.id != null, body, onSubmit, onDelete);
  }


}