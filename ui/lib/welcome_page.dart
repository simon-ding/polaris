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
                    var item = value[i];
                    return Card(
                        margin: const EdgeInsets.all(4),
                        clipBehavior: Clip.hardEdge,
                        child: InkWell(
                          //splashColor: Colors.blue.withAlpha(30),
                          onTap: () {
                            if (uri == routeMoivie) {
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
                                child: Image.network(
                                    "${APIs.imagesUrl}/${item.id}/poster.jpg",
                                    fit: BoxFit.fill),
                              ),
                              SizedBox(
                                  width: 140,
                                  child: LinearProgressIndicator(
                                    value: 1,
                                    color: item.status == "downloaded"
                                        ? Colors.green
                                        : Colors.blue,
                                  )),
                              Text(
                                item.name!,
                                style: const TextStyle(
                                    fontSize: 14,
                                    fontWeight: FontWeight.bold,
                                    height: 2.5),
                              ),
                            ],
                          ),
                        ));
                  }),
          ),
        ),
      _ => const MyProgressIndicator(),
    };
  }
}
