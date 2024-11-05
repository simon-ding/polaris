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
        error: (err, trace) => PoNetworkError(err: err),
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
        "impl": "torznab",
        "priority": indexer.priority.toString(),
        "seed_ratio": indexer.seedRatio.toString(),
        "disabled": indexer.disabled
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
          FormBuilderTextField(
            name: "priority",
            decoration: const InputDecoration(
              labelText: "索引优先级",
              helperText: "取值范围1-128， 数值越大，优先级越高",
            ),
            autovalidateMode: AutovalidateMode.onUserInteraction,
            validator: FormBuilderValidators.positiveNumber(),
          ),
          FormBuilderTextField(
            name: "seed_ratio",
            decoration: const InputDecoration(
              labelText: "做种率",
              helperText: "种子的做种率，达到此做种率后，种子才会被删除, 0表示不做种",
              hintText: "1.0",
            ),
            autovalidateMode: AutovalidateMode.onUserInteraction,
            validator: FormBuilderValidators.numeric(),
          ),
          FormBuilderSwitch(name: "disabled", title: const Text("禁用此索引器"))
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
            id: indexer.id,
            name: values["name"],
            url: values["url"],
            apiKey: values["api_key"],
            priority: int.parse(values["priority"]),
            seedRatio: double.parse(values["seed_ratio"]),
            disabled: values["disabled"]));
      } else {
        throw "validation_error";
      }
    }

    return showSettingDialog(
        context, "索引器", indexer.id != null, body, onSubmit, onDelete);
  }
}
