import 'package:flutter/material.dart';
import 'package:flutter_form_builder/flutter_form_builder.dart';
import 'package:flutter_riverpod/flutter_riverpod.dart';
import 'package:form_builder_validators/form_builder_validators.dart';
import 'package:ui/providers/settings.dart';
import 'package:ui/settings/dialog.dart';
import 'package:ui/widgets/progress_indicator.dart';
import 'package:ui/widgets/widgets.dart';

class Importlist extends ConsumerStatefulWidget {
  const Importlist({super.key});

  @override
  ConsumerState<ConsumerStatefulWidget> createState() {
    return _ImportlistState();
  }
}

class _ImportlistState extends ConsumerState<Importlist> {
  @override
  Widget build(BuildContext context) {
    var importlists = ref.watch(importlistProvider);

    return importlists.when(
        data: (value) => Wrap(
              children: List.generate(value.length + 1, (i) {
                if (i < value.length) {
                  var indexer = value[i];
                  return SettingsCard(
                      onTap: () => showImportlistDetails(indexer),
                      child: Text(indexer.name ?? ""));
                }
                return SettingsCard(
                    onTap: () => showImportlistDetails(ImportList()),
                    child: const Icon(Icons.add));
              }),
            ),
        error: (err, trace) => PoError(msg: "网络错误", err: err),
        loading: () => const MyProgressIndicator());
  }

  Future<void> showImportlistDetails(ImportList list) {
    final _formKey = GlobalKey<FormBuilderState>();
    String? _selectedType = list.type;

    var body = StatefulBuilder(builder: (context, setState) {
      return FormBuilder(
        key: _formKey,
        initialValue: {
          "name": list.name,
          "url": list.url,
          "qulity": list.qulity,
          "type": list.type,
          "storage_id": list.storageId
        },
        child: Column(
          children: [
            FormBuilderDropdown(
              name: "type",
              decoration: const InputDecoration(
                  labelText: "类型",
                  hintText:
                      "Plex Watchlist: https://support.plex.tv/articles/universal-watchlist/"),
              items: const [
                DropdownMenuItem(value: "plex", child: Text("Plex Watchlist")),
              ],
              onChanged: (value) {
                setState(() {
                  _selectedType = value;
                });
              },
            ),
            _selectedType == "plex"
                ? const Text(
                    "Plex Watchlist: https://support.plex.tv/articles/universal-watchlist/")
                : const Text(""),
            FormBuilderTextField(
              name: "name",
              decoration: Commons.requiredTextFieldStyle(text: "名称"),
              autovalidateMode: AutovalidateMode.onUserInteraction,
              validator: FormBuilderValidators.required(),
            ),
            FormBuilderTextField(
              name: "url",
              decoration: Commons.requiredTextFieldStyle(text: "地址"),
              autovalidateMode: AutovalidateMode.onUserInteraction,
              validator: FormBuilderValidators.required(),
            ),
            FormBuilderDropdown(
              name: "qulity",
              decoration: const InputDecoration(labelText: "清晰度"),
              items: const [
                DropdownMenuItem(value: "720p", child: Text("720p")),
                DropdownMenuItem(value: "1080p", child: Text("1080p")),
                DropdownMenuItem(value: "2160p", child: Text("2160p")),
              ],
            ),
            Consumer(
              builder: (context, ref, child) {
                var storage = ref.watch(storageSettingProvider);
                return storage.when(
                    data: (v) {
                      return FormBuilderDropdown(
                        name: "storage_id",
                        decoration: const InputDecoration(labelText: "存储"),
                        items: List.generate(
                            v.length,
                            (i) => DropdownMenuItem(
                                  value: v[i].id,
                                  child: Text(v[i].name!),
                                )),
                      );
                    },
                    error: (err, trace) => Text("$err"),
                    loading: () => const MyProgressIndicator());
              },
            ),
          ],
        ),
      );
    });

    onDelete() async {
      return ref.read(importlistProvider.notifier).deleteimportlist(list.id!);
    }

    onSubmit() async {
      if (_formKey.currentState!.saveAndValidate()) {
        var values = _formKey.currentState!.value;

        return ref.read(importlistProvider.notifier).addImportlist(ImportList(
              id: list.id,
              name: values["name"],
              url: values["url"],
              type: values["type"],
              qulity: values["qulity"],
              storageId: values["storage_id"],
            ));
      } else {
        throw "validation_error";
      }
    }

    return showSettingDialog(
        context, "监控列表", list.id != null, body, onSubmit, onDelete);
  }
}
