import 'package:flutter/material.dart';
import 'package:flutter_riverpod/flutter_riverpod.dart';
import 'package:go_router/go_router.dart';
import 'package:ui/movie_watchlist.dart';
import 'package:ui/providers/APIs.dart';
import 'package:ui/providers/welcome_data.dart';
import 'package:ui/tv_details.dart';
import 'package:ui/widgets/progress_indicator.dart';

class WelcomePage extends ConsumerWidget {
  static const routeTv = "/series";
  static const routeMoivie = "/movies";

  const WelcomePage({super.key});

  @override
  Widget build(BuildContext context, WidgetRef ref) {
    var uri = GoRouterState.of(context).uri.toString();

    AsyncValue<List<MediaDetail>> data;
    if (uri == routeMoivie) {
      data = ref.watch(movieWatchlistDataProvider);
    } else {
      data = ref.watch(tvWatchlistDataProvider);
    }

    return switch (data) {
      AsyncData(:final value) => SingleChildScrollView(
          child: Wrap(
            spacing: 10,
            runSpacing: 20,
            children: value.isEmpty
                ? [
                    Container(
                        height: MediaQuery.of(context).size.height * 0.6,
                        alignment: Alignment.center,
                        child: const Text(
                          "啥都没有...",
                          style: TextStyle(fontSize: 16),
                        ))
                  ]
                : List.generate(value.length, (i) {
                    final item = value[i];
                    return MediaCard(item: item);
                  }),
          ),
        ),
      _ => const MyProgressIndicator(),
    };
  }
}

class MediaCard extends StatelessWidget {
  final MediaDetail item;

  const MediaCard({super.key, required this.item});
  @override
  Widget build(BuildContext context) {
    return Card(
        shape: ContinuousRectangleBorder(
          borderRadius: BorderRadius.circular(10.0),
        ),
        //margin: const EdgeInsets.all(4),
        clipBehavior: Clip.hardEdge,
        elevation: 10,
        child: InkWell(
          //splashColor: Colors.blue.withAlpha(30),
          onTap: () {
            if (item.mediaType == "movie") {
              context.go(MovieDetailsPage.toRoute(item.id!));
            } else {
              context.go(TvDetailsPage.toRoute(item.id!));
            }
          },
          child: Column(
            children: <Widget>[
              SizedBox(
                width: 140,
                height: 210,
                child: Ink.image(
                    fit: BoxFit.cover,
                    image: NetworkImage(
                      "${APIs.imagesUrl}/${item.id}/poster.jpg",
                    )),
              ),
              SizedBox(
                  width: 140,
                  child: Column(
                    children: [
                      LinearProgressIndicator(
                        value: 1,
                        color: item.downloadedNum! >= item.monitoredNum!
                            ? Colors.green
                            : Colors.blue,
                      ),
                      Text(
                        item.name!,
                        overflow: TextOverflow.ellipsis,
                        style: const TextStyle(
                            fontSize: 14,
                            fontWeight: FontWeight.bold,
                            height: 2.5),
                      ),
                    ],
                  )),
            ],
          ),
        ));
  }
}
