import 'package:flutter/material.dart';
import 'package:flutter_form_builder/flutter_form_builder.dart';
import 'package:flutter_riverpod/flutter_riverpod.dart';
import 'package:form_builder_validators/form_builder_validators.dart';
import 'package:ui/providers/settings.dart';
import 'package:ui/settings/dialog.dart';
import 'package:ui/widgets/progress_indicator.dart';
import 'package:ui/widgets/widgets.dart';

class IndexerSettings extends ConsumerStatefulWidget {

  const IndexerSettings({super.key});
  @override
  ConsumerState<ConsumerStatefulWidget> createState() {
    return _IndexerState();
  }
}

class _IndexerState extends ConsumerState<IndexerSettings> {
  @override
  Widget build(BuildContext context) {
    var indexers = ref.watch(indexersProvider);
    return indexers.when(
        data: (value) => Wrap(
              children: List.generate(value.length + 1, (i) {
                if (i < value.length) {
                  var indexer = value[i];
                  return SettingsCard(
                      onTap: () => showIndexerDetails(indexer),
                      child: Text(indexer.name ?? ""));
                }
                return SettingsCard(
                    onTap: () => showIndexerDetails(Indexer()),
                    child: const Icon(Icons.add));
              }),
            ),
        error: (err, trace) => Text("$err"),
        loading: () => const MyProgressIndicator());
  }

  Future<void> showIndexerDetails(Indexer indexer) {
    final _formKey = GlobalKey<FormBuilderState>();

    var body = FormBuilder(
      key: _formKey,
      initialValue: {
        "name": indexer.name,
        "url": indexer.url,
        "api_key": indexer.apiKey,
        "impl": "torznab"
      },
      child: Column(
        children: [
          FormBuilderDropdown(
            name: "impl",
            decoration: const InputDecoration(labelText: "类型"),
            items: const [
              DropdownMenuItem(value: "torznab", child: Text("Torznab")),
            ],
          ),
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
          FormBuilderTextField(
            name: "api_key",
            decoration: Commons.requiredTextFieldStyle(text: "API Key"),
            autovalidateMode: AutovalidateMode.onUserInteraction,
            validator: FormBuilderValidators.required(),
          ),
        ],
      ),
    );
    onDelete() async {
      return ref.read(indexersProvider.notifier).deleteIndexer(indexer.id!);
    }

    onSubmit() async {
      if (_formKey.currentState!.saveAndValidate()) {
        var values = _formKey.currentState!.value;
        return ref.read(indexersProvider.notifier).addIndexer(Indexer(
            name: values["name"],
            url: values["url"],
            apiKey: values["api_key"]));
      } else {
        throw "validation_error";
      }
    }

    return showSettingDialog(
        context, "索引器", indexer.id != null, body, onSubmit, onDelete);
  }
}
