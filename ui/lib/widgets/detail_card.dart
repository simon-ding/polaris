import 'package:flutter/material.dart';
import 'package:flutter_form_builder/flutter_form_builder.dart';
import 'package:flutter_riverpod/flutter_riverpod.dart';
import 'package:form_builder_validators/form_builder_validators.dart';
import 'package:go_router/go_router.dart';
import 'package:quiver/strings.dart';
import 'package:ui/providers/APIs.dart';
import 'package:ui/providers/series_details.dart';
import 'package:ui/welcome_page.dart';
import 'package:ui/widgets/utils.dart';
import 'package:url_launcher/url_launcher.dart';

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
  final tmdbBase = "https://www.themoviedb.org/";
  final imdbBase = "https://www.imdb.com/title/";
  @override
  Widget build(BuildContext context) {
    final screenWidth = MediaQuery.of(context).size.width;
    final url = Uri.parse(tmdbBase +
        (widget.details.mediaType == "tv" ? "tv/" : "movie/") +
        widget.details.tmdbId.toString());

    final imdbUrl = Uri.parse(imdbBase + (widget.details.imdbid ?? ""));
    return Card(
      margin: const EdgeInsets.all(4),
      clipBehavior: Clip.hardEdge,
      child: Container(
        constraints: const BoxConstraints(maxHeight: 400),
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
              screenWidth < 600
                  ? const SizedBox()
                  : Flexible(
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
                        //const Text(""),
                        Text(
                          "${widget.details.name} ${widget.details.name != widget.details.originalName ? widget.details.originalName : ''} ${widget.details.airDate == null ? "" : (widget.details.airDate!.split("-")[0])}",
                          style: const TextStyle(
                              fontSize: 24,
                              fontWeight: FontWeight.bold,
                              height: 2.5),
                        ),
                        const Divider(thickness: 1, height: 1),
                        const Text(
                          "",
                          style: TextStyle(height: 0.2),
                        ),
                        Wrap(
                          spacing: 10,
                          children: [
                            Chip(
                              clipBehavior: Clip.hardEdge,
                              shape: ContinuousRectangleBorder(
                                  borderRadius: BorderRadius.circular(100.0)),
                              label: Text(
                                  "${widget.details.storage!.name}: ${widget.details.mediaType == "tv" ? widget.details.storage!.tvPath : widget.details.storage!.moviePath}"
                                  "${widget.details.targetDir}"),
                              padding: const EdgeInsets.all(0),
                            ),
                            Chip(
                              clipBehavior: Clip.hardEdge,
                              shape: ContinuousRectangleBorder(
                                  borderRadius: BorderRadius.circular(100.0)),
                              label: Text("${widget.details.resolution}"),
                              padding: const EdgeInsets.all(0),
                            ),
                            widget.details.limiter != null &&
                                    widget.details.limiter!.sizeMax > 0
                                ? Chip(
                                    clipBehavior: Clip.hardEdge,
                                    shape: ContinuousRectangleBorder(
                                        borderRadius:
                                            BorderRadius.circular(100.0)),
                                    padding: const EdgeInsets.all(0),
                                    label: Text(
                                        "${(widget.details.limiter!.sizeMin).readableFileSize()} - ${(widget.details.limiter!.sizeMax).readableFileSize()}"))
                                : const SizedBox(),
                            MenuAnchor(
                              style: const MenuStyle(
                                  alignment: Alignment.bottomRight),
                              menuChildren: [
                                ActionChip.elevated(
                                    onPressed: () => launchUrl(url),
                                    clipBehavior: Clip.hardEdge,
                                    backgroundColor: Colors.indigo[700],
                                    shape: ContinuousRectangleBorder(
                                        borderRadius:
                                            BorderRadius.circular(100.0)),
                                    padding: const EdgeInsets.all(0),
                                    label: const Text("TMDB")),
                                isBlank(widget.details.imdbid)
                                    ? const SizedBox()
                                    : ActionChip.elevated(
                                        onPressed: () => launchUrl(imdbUrl),
                                        backgroundColor: Colors.indigo[700],
                                        clipBehavior: Clip.hardEdge,
                                        shape: ContinuousRectangleBorder(
                                            borderRadius:
                                                BorderRadius.circular(100.0)),
                                        padding: const EdgeInsets.all(0),
                                        label: const Text("IMDB"),
                                      )
                              ],
                              builder: (context, controller, child) {
                                return ActionChip.elevated(
                                    onPressed: () {
                                      if (controller.isOpen) {
                                        controller.close();
                                      } else {
                                        controller.open();
                                      }
                                    },
                                    clipBehavior: Clip.hardEdge,
                                    shape: ContinuousRectangleBorder(
                                        borderRadius:
                                            BorderRadius.circular(100.0)),
                                    padding: const EdgeInsets.all(0),
                                    label: const Text("外部链接"));
                              },
                            ),
                          ],
                        ),
                        const Text("", style: TextStyle(height: 1)),
                        Expanded(
                            child: Text(
                          overflow: TextOverflow.ellipsis,
                          maxLines: 7,
                          widget.details.overview ?? "",
                        )),
                        Row(
                          mainAxisAlignment: MainAxisAlignment.spaceEvenly,
                          children: [
                            downloadButton(),
                            editIcon(),
                            deleteIcon(context),
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

  Widget deleteIcon(BuildContext oriContext) {
    return IconButton(
        tooltip: widget.details.mediaType == "tv" ? "删除剧集" : "删除电影",
        onPressed: () => showConfirmDialog(oriContext),
        icon: const Icon(Icons.delete));
  }

  Future<void> showConfirmDialog(BuildContext oriContext) {
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
                      .then((v) {
                    if (oriContext.mounted) {
                      oriContext.go(widget.details.mediaType == "tv"
                          ? WelcomePage.routeTv
                          : WelcomePage.routeMoivie);
                    }
                  });
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
                        DropdownMenuItem(value: "any", child: Text("不限")),
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
                        .then((v) {
                      if (context.mounted) {
                        Navigator.of(context).pop();
                      }
                    });
                  }
                },
                label: const Text("确认"))
          ],
        );
      },
    );
  }

  Widget downloadButton() {
    return LoadingIconButton(
        tooltip: widget.details.mediaType == "tv" ? "查找并下载所有监控剧集" : "查找并下载此电影",
        onPressed: () async {
          await ref
              .read(mediaDetailsProvider(widget.details.id.toString()).notifier)
              .downloadall()
              .then((list) => {
                    if (list != null) {showSnakeBar("开始下载：$list")}
                  });
        },
        icon: Icons.download_rounded);
  }
}
