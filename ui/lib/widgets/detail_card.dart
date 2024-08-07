import 'package:flutter/material.dart';
import 'package:flutter_form_builder/flutter_form_builder.dart';
import 'package:flutter_riverpod/flutter_riverpod.dart';
import 'package:form_builder_validators/form_builder_validators.dart';
import 'package:go_router/go_router.dart';
import 'package:ui/providers/APIs.dart';
import 'package:ui/providers/series_details.dart';
import 'package:ui/welcome_page.dart';
import 'package:ui/widgets/utils.dart';

import 'widgets.dart';

class DetailCard extends ConsumerStatefulWidget {
  final SeriesDetails details;

  const DetailCard({super.key, required this.details});

  @override
  ConsumerState<ConsumerStatefulWidget> createState() {
    return _DetailCardState();
  }
}

class _DetailCardState extends ConsumerState<DetailCard> {
  @override
  Widget build(BuildContext context) {
    return Card(
      margin: const EdgeInsets.all(4),
      clipBehavior: Clip.hardEdge,
      child: Container(
        constraints:
            BoxConstraints(maxHeight: MediaQuery.of(context).size.height * 0.4),
        decoration: BoxDecoration(
            image: DecorationImage(
                fit: BoxFit.cover,
                opacity: 0.3,
                colorFilter: ColorFilter.mode(
                    Colors.black.withOpacity(0.3), BlendMode.dstATop),
                image: NetworkImage(
                    "${APIs.imagesUrl}/${widget.details.id}/backdrop.jpg"))),
        child: Padding(
          padding: const EdgeInsets.all(10),
          child: Row(
            children: <Widget>[
              Flexible(
                flex: 2,
                child: Padding(
                  padding: const EdgeInsets.all(10),
                  child: Image.network(
                      "${APIs.imagesUrl}/${widget.details.id}/poster.jpg",
                      fit: BoxFit.contain),
                ),
              ),
              Flexible(
                flex: 4,
                child: Row(
                  children: [
                    Expanded(
                        child: Column(
                      crossAxisAlignment: CrossAxisAlignment.start,
                      children: [
                        const Text(""),
                        Wrap(
                          children: [
                            Text("${widget.details.resolution}"),
                            const SizedBox(
                              width: 30,
                            ),
                            Text("${widget.details.storage!.name}:"),
                            Text(
                                "${widget.details.mediaType == "tv" ? widget.details.storage!.tvPath : widget.details.storage!.moviePath}"
                                "${widget.details.targetDir}"),
                            const SizedBox(
                              width: 30,
                            ),
                            widget.details.limiter != null &&
                                    widget.details.limiter!.sizeMax > 0
                                ? Text(
                                    "${(widget.details.limiter!.sizeMin).readableFileSize()} - ${(widget.details.limiter!.sizeMax).readableFileSize()}")
                                : const SizedBox()
                          ],
                        ),
                        const Divider(thickness: 1, height: 1),
                        Text(
                          "${widget.details.name} ${widget.details.name != widget.details.originalName ? widget.details.originalName : ''} (${widget.details.airDate!.split("-")[0]})",
                          style: const TextStyle(
                              fontSize: 20, fontWeight: FontWeight.bold),
                        ),
                        const Text(""),
                        Expanded(
                            child: Text(
                          overflow: TextOverflow.visible,
                          maxLines: 9,
                          widget.details.overview ?? "",
                        )),
                        Row(
                          mainAxisAlignment: MainAxisAlignment.spaceEvenly,
                          children: [
                            downloadButton(),
                            editIcon(),
                            deleteIcon(),
                          ],
                        )
                      ],
                    )),
                  ],
                ),
              ),
            ],
          ),
        ),
      ),
    );
  }

  Widget deleteIcon() {
    return IconButton(
        tooltip: widget.details.mediaType == "tv" ? "删除剧集" : "删除电影",
        onPressed: () => showConfirmDialog(),
        icon: const Icon(Icons.delete));
  }

  Future<void> showConfirmDialog() {
    return showDialog<void>(
      context: context,
      barrierDismissible: true,
      builder: (BuildContext context) {
        return AlertDialog(
          title: const Text("确认删除："),
          content: Text("${widget.details.name}"),
          actions: [
            TextButton(
                onPressed: () => Navigator.of(context).pop(),
                child: const Text("取消")),
            TextButton(
                onPressed: () {
                  ref
                      .read(mediaDetailsProvider(widget.details.id.toString())
                          .notifier)
                      .delete()
                      .then((v) => context.go(widget.details.mediaType == "tv"
                          ? WelcomePage.routeTv
                          : WelcomePage.routeMoivie));
                  Navigator.of(context).pop();
                },
                child: const Text("确认"))
          ],
        );
      },
    );
  }

  Widget editIcon() {
    return IconButton(
        tooltip: "编辑",
        onPressed: () => showEditDialog(),
        icon: const Icon(Icons.edit));
  }

  showEditDialog() {
    final _formKey = GlobalKey<FormBuilderState>();
    return showDialog<void>(
      context: context,
      barrierDismissible: true,
      builder: (BuildContext context) {
        return AlertDialog(
          title: Text("编辑 ${widget.details.name}"),
          content: SelectionArea(
            child: SizedBox(
              width: MediaQuery.of(context).size.width * 0.3,
              height: MediaQuery.of(context).size.height * 0.4,
              child: SingleChildScrollView(
                  child: FormBuilder(
                key: _formKey,
                initialValue: {
                  "resolution": widget.details.resolution,
                  "target_dir": widget.details.targetDir,
                  "limiter": widget.details.limiter != null
                      ? RangeValues(
                          widget.details.limiter!.sizeMin.toDouble() /
                              1000 /
                              1000,
                          widget.details.limiter!.sizeMax.toDouble() /
                              1000 /
                              1000)
                      : const RangeValues(0, 0)
                },
                child: Column(
                  children: [
                    FormBuilderDropdown(
                      name: "resolution",
                      decoration: const InputDecoration(labelText: "清晰度"),
                      items: const [
                        DropdownMenuItem(value: "720p", child: Text("720p")),
                        DropdownMenuItem(value: "1080p", child: Text("1080p")),
                        DropdownMenuItem(value: "2160p", child: Text("2160p")),
                      ],
                    ),
                    FormBuilderTextField(
                      name: "target_dir",
                      decoration: const InputDecoration(labelText: "存储路径"),
                      validator: FormBuilderValidators.required(),
                    ),
                    const MyRangeSlider(name: "limiter"),
                  ],
                ),
              )),
            ),
          ),
          actions: [
            TextButton(
                onPressed: () => Navigator.of(context).pop(),
                child: const Text("取消")),
            LoadingTextButton(
                onPressed: () async {
                  if (_formKey.currentState!.saveAndValidate()) {
                    final values = _formKey.currentState!.value;
                    await ref
                        .read(mediaDetailsProvider(widget.details.id.toString())
                            .notifier)
                        .edit(values["resolution"], values["target_dir"],
                            values["limiter"])
                        .then((v) => Navigator.of(context).pop());
                  }
                },
                label: const Text("确认"))
          ],
        );
      },
    );
  }

  Widget downloadButton() {
    return IconButton(
        tooltip: widget.details.mediaType == "tv" ? "查找并下载所有监控剧集" : "查找并下载此电影",
        onPressed: () {},
        icon: const Icon(Icons.download_rounded));
  }
}
