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
                    child: Column(
                      mainAxisAlignment: MainAxisAlignment.center,
                      children: [
                        Text(
                          client.name!,
                          style: TextStyle(fontSize: 20, height: 3),
                        ),
                        Opacity(
                          opacity: 0.5,
                          child: Text(client.service!),
                        )
                      ],
                    ),
                    onTap: () => showNotifierAccordingToService(client),
                  );
                }
                return SettingsCard(
                    onTap: () => showSelections(),
                    child: const Icon(Icons.add));
              }),
            ),
        error: (err, trace) => PoError(msg: "网络错误", err: err),
        loading: () => const MyProgressIndicator());
  }

  Future<void> showSelections() {
    return showDialog<void>(
        context: context,
        barrierDismissible: true,
        builder: (BuildContext context) {
          return AlertDialog(
            content: SizedBox(
              height: 500,
              width: 500,
              child: Wrap(
                children: [
                  SettingsCard(
                    child: InkWell(
                      child: const Center(
                        child: Text("Pushover"),
                      ),
                      onTap: () {
                        Navigator.of(context).pop();
                        showPushoverNotifierDetails(NotifierData());
                      },
                    ),
                  ),
                  SettingsCard(
                    child: InkWell(
                      child: const Center(
                        child: Text("Bark"),
                      ),
                      onTap: () {
                        Navigator.of(context).pop();
                        showBarkNotifierDetails(NotifierData());
                      },
                    ),
                  )
                ],
              ),
            ),
          );
        });
  }

  Future<void> showNotifierAccordingToService(NotifierData notifier) {
    switch (notifier.service) {
      case "bark":
        return showBarkNotifierDetails(notifier);
      case "pushover":
        return showPushoverNotifierDetails(notifier);
    }
    return Future<void>.value();
  }

  Future<void> showBarkNotifierDetails(NotifierData notifier) {
    final _formKey = GlobalKey<FormBuilderState>();

    var body = FormBuilder(
      key: _formKey,
      initialValue: {
        "name": notifier.name,
        "enabled": notifier.enabled ?? true,
        "device_key":
            notifier.settings != null ? notifier.settings!["device_key"] : "",
        "url": notifier.settings != null ? notifier.settings!["url"] : "",
      },
      child: Column(
        children: [
          const Text("https://bark.day.app/#/"),
          FormBuilderTextField(
            name: "name",
            decoration: Commons.requiredTextFieldStyle(text: "名称"),
            autovalidateMode: AutovalidateMode.onUserInteraction,
            validator: FormBuilderValidators.required(),
          ),
          FormBuilderTextField(
            name: "url",
            decoration: const InputDecoration(
                labelText: "服务器地址", helperText: "留空使用默认地址"),
          ),
          FormBuilderTextField(
            name: "device_key",
            decoration: Commons.requiredTextFieldStyle(text: "Device Key"),
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
                service: "bark",
                enabled: values["enabled"],
                settings: {
                  "device_key": values["device_key"],
                  "url": values["url"]
                }));
      } else {
        throw "validation_error";
      }
    }

    return showSettingDialog(
        context, "Bark", notifier.id != null, body, onSubmit, onDelete);
  }

  Future<void> showPushoverNotifierDetails(NotifierData notifier) {
    final _formKey = GlobalKey<FormBuilderState>();

    var body = FormBuilder(
      key: _formKey,
      initialValue: {
        "name": notifier.name,
        "enabled": notifier.enabled ?? true,
        "app_token":
            notifier.settings != null ? notifier.settings!["app_token"] : "",
        "user_key":
            notifier.settings != null ? notifier.settings!["user_key"] : "",
      },
      child: Column(
        children: [
          const Text("https://pushover.net/"),
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
                service: "pushover",
                enabled: values["enabled"],
                settings: {
                  "app_token": values["app_token"],
                  "user_key": values["user_key"]
                }));
      } else {
        throw "validation_error";
      }
    }

    return showSettingDialog(
        context, "Pushover", notifier.id != null, body, onSubmit, onDelete);
  }
}
