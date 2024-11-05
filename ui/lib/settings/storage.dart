import 'package:flutter/material.dart';
import 'package:flutter_form_builder/flutter_form_builder.dart';
import 'package:flutter_riverpod/flutter_riverpod.dart';
import 'package:form_builder_validators/form_builder_validators.dart';
import 'package:ui/providers/settings.dart';
import 'package:ui/settings/dialog.dart';
import 'package:ui/widgets/progress_indicator.dart';
import 'package:ui/widgets/widgets.dart';

class StorageSettings extends ConsumerStatefulWidget {
  static const route = "/settings";

  const StorageSettings({super.key});
  @override
  ConsumerState<ConsumerStatefulWidget> createState() {
    return _StorageState();
  }
}

class _StorageState extends ConsumerState<StorageSettings> {
  @override
  Widget build(BuildContext context) {
    var storageSettingData = ref.watch(storageSettingProvider);
    return storageSettingData.when(
        data: (value) => Wrap(
              children: List.generate(value.length + 1, (i) {
                if (i < value.length) {
                  var storage = value[i];
                  return SettingsCard(
                      onTap: () => showStorageDetails(storage),
                      child: Text(storage.name ?? ""));
                }
                return SettingsCard(
                    onTap: () => showStorageDetails(Storage()),
                    child: const Icon(Icons.add));
              }),
            ),
        error: (err, trace) => PoNetworkError(err: err),
        loading: () => const MyProgressIndicator());
  }

  Future<void> showStorageDetails(Storage s) {
    final _formKey = GlobalKey<FormBuilderState>();

    String selectImpl = s.implementation == null ? "local" : s.implementation!;
    final widgets =
        StatefulBuilder(builder: (BuildContext context, StateSetter setState) {
      return FormBuilder(
          key: _formKey,
          autovalidateMode: AutovalidateMode.disabled,
          initialValue: {
            "name": s.name,
            "impl": s.implementation == null ? "local" : s.implementation!,
            "user": s.settings != null ? s.settings!["user"] ?? "" : "",
            "password": s.settings != null ? s.settings!["password"] ?? "" : "",
            "tv_path": s.tvPath,
            "url": s.settings != null ? s.settings!["url"] ?? "" : "",
            "movie_path": s.moviePath,
            "change_file_hash": s.settings != null
                ? s.settings!["change_file_hash"] == "true"
                    ? true
                    : false
                : false,
          },
          child: Column(
            crossAxisAlignment: CrossAxisAlignment.start,
            children: <Widget>[
              FormBuilderDropdown<String>(
                name: "impl",
                autovalidateMode: AutovalidateMode.onUserInteraction,
                decoration: const InputDecoration(labelText: "类型"),
                onChanged: (value) {
                  setState(() {
                    selectImpl = value!;
                  });
                },
                items: const [
                  DropdownMenuItem(
                    value: "local",
                    child: Text("本地存储"),
                  ),
                  DropdownMenuItem(
                    value: "webdav",
                    child: Text("webdav"),
                  )
                ],
                validator: FormBuilderValidators.required(),
              ),
              FormBuilderTextField(
                name: "name",
                autovalidateMode: AutovalidateMode.onUserInteraction,
                initialValue: s.name,
                decoration: const InputDecoration(labelText: "名称"),
                validator: FormBuilderValidators.required(),
              ),
              selectImpl != "local"
                  ? Column(
                      crossAxisAlignment: CrossAxisAlignment.start,
                      children: [
                        FormBuilderTextField(
                          name: "url",
                          autovalidateMode: AutovalidateMode.onUserInteraction,
                          decoration:
                              const InputDecoration(labelText: "Webdav地址"),
                          validator: FormBuilderValidators.required(),
                        ),
                        FormBuilderTextField(
                          name: "user",
                          autovalidateMode: AutovalidateMode.onUserInteraction,
                          decoration: const InputDecoration(labelText: "用户"),
                        ),
                        FormBuilderTextField(
                          name: "password",
                          autovalidateMode: AutovalidateMode.onUserInteraction,
                          decoration: const InputDecoration(labelText: "密码"),
                          obscureText: true,
                        ),
                        FormBuilderCheckbox(
                          name: "change_file_hash",
                          title: const Text(
                            "上传时更改文件哈希",
                            style: TextStyle(fontSize: 14),
                          ),
                        ),
                      ],
                    )
                  : Container(),
              FormBuilderTextField(
                name: "tv_path",
                autovalidateMode: AutovalidateMode.onUserInteraction,
                decoration: const InputDecoration(labelText: "电视剧路径"),
                validator: FormBuilderValidators.required(),
              ),
              FormBuilderTextField(
                name: "movie_path",
                autovalidateMode: AutovalidateMode.onUserInteraction,
                decoration: const InputDecoration(labelText: "电影路径"),
                validator: FormBuilderValidators.required(),
              )
            ],
          ));
    });
    onSubmit() async {
      if (_formKey.currentState!.saveAndValidate()) {
        final values = _formKey.currentState!.value;
        return ref.read(storageSettingProvider.notifier).addStorage(Storage(
              name: values["name"],
              implementation: selectImpl,
              tvPath: values["tv_path"],
              moviePath: values["movie_path"],
              settings: {
                "url": values["url"],
                "user": values["user"],
                "password": values["password"],
                "change_file_hash":
                    (values["change_file_hash"] ?? false) as bool
                        ? "true"
                        : "false"
              },
            ));
      } else {
        throw "validation_error";
      }
    }

    onDelete() async {
      return ref.read(storageSettingProvider.notifier).deleteStorage(s.id!);
    }

    return showSettingDialog(
        context, '存储', s.id != null, widgets, onSubmit, onDelete);
  }
}
