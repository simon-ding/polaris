import 'package:flutter/material.dart';
import 'package:flutter_riverpod/flutter_riverpod.dart';
import 'package:go_router/go_router.dart';
import 'package:ui/providers/APIs.dart';
import 'package:ui/providers/series_details.dart';
import 'package:ui/welcome_page.dart';

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
                        Row(
                          children: [
                            Text("${widget.details.resolution}"),
                            const SizedBox(
                              width: 30,
                            ),
                            Text(
                                "${widget.details.storage!.name} (${widget.details.storage!.implementation})"),
                            const SizedBox(
                              width: 30,
                            ),
                            Expanded(child: Text(
                                "${widget.details.mediaType == "tv" ? widget.details.storage!.tvPath : widget.details.storage!.moviePath}"
                                "${widget.details.targetDir}"),)
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
                          overflow: TextOverflow.ellipsis,
                          maxLines: 9,
                          widget.details.overview ?? "",
                        )),
                        Row(
                          children: [
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
    return Tooltip(
      message: widget.details.mediaType == "tv" ? "删除剧集" : "删除电影",
      child: IconButton(
          onPressed: () {
            var f = ref
                .read(
                    mediaDetailsProvider(widget.details.id.toString()).notifier)
                .delete()
                .then((v) => context.go(widget.details.mediaType == "tv"
                    ? WelcomePage.routeTv
                    : WelcomePage.routeMoivie));
            showLoadingWithFuture(f);
          },
          icon: const Icon(Icons.delete)),
    );
  }
}
