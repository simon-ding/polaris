import 'package:flutter/material.dart';
import 'package:flutter_riverpod/flutter_riverpod.dart';
import 'package:go_router/go_router.dart';
import 'package:ui/movie_watchlist.dart';
import 'package:ui/providers/APIs.dart';
import 'package:ui/providers/welcome_data.dart';
import 'package:ui/tv_details.dart';
import 'package:ui/widgets/progress_indicator.dart';

class WelcomePage extends ConsumerStatefulWidget {
  const WelcomePage({super.key});
  static const routeTv = "/series";
  static const routeMoivie = "/movies";

  @override
  ConsumerState<ConsumerStatefulWidget> createState() {
    return WelcomePageState();
  }
}

class WelcomePageState extends ConsumerState<WelcomePage> {
  //WelcomePageState({super.key});

  bool onlyShowUnfinished = false;

  @override
  Widget build(BuildContext context) {
    var uri = GoRouterState.of(context).uri.toString();

    AsyncValue<List<MediaDetail>> data;
    if (uri == WelcomePage.routeMoivie) {
      data = ref.watch(movieWatchlistDataProvider);
    } else {
      data = ref.watch(tvWatchlistDataProvider);
    }

    return Stack(
      //alignment: Alignment.bottomRight,
      children: [
        () {
          return switch (data) {
            AsyncData(:final value) => SingleChildScrollView(
                child: Wrap(
                  alignment: WrapAlignment.start,
                  spacing: isSmallScreen(context) ? 0 : 10,
                  runSpacing: isSmallScreen(context) ? 10 : 20,
                  children: getMediaAll(value),
                ),
              ),
            _ => const MyProgressIndicator(),
          };
        }(),
        Row(
          children: [
            Expanded(child: Container()),
            Column(
              children: [
                Expanded(child: Container()),
                Padding(
                  padding: EdgeInsets.all(20),
                  child: MenuAnchor(
                    style: MenuStyle(
                      //minimumSize: WidgetStatePropertyAll(Size(400, 300)),

                      backgroundColor: WidgetStatePropertyAll(Theme.of(context)
                          .colorScheme
                          .inversePrimary
                          .withOpacity(0.9)),
                    ),
                    menuChildren: [
                      MenuItemButton(
                        onPressed: null,
                        child: CheckboxListTile(
                          value: onlyShowUnfinished,
                          onChanged: (b) {
                            setState(() {
                              onlyShowUnfinished = b!;
                            });
                          },
                          title: const Text(
                            "未完成",
                            style: TextStyle(fontSize: 16),
                            softWrap: false,
                          ),
                          controlAffinity: ListTileControlAffinity.leading,
                        ),
                      ),
                    ],
                    builder: (context, controller, child) {
                      return Opacity(
                          opacity: 0.7,
                          child: FloatingActionButton(
                            onPressed: () {
                              if (controller.isOpen) {
                                controller.close();
                              } else {
                                controller.open();
                              }
                            },
                            child: const Icon(Icons.more_horiz),
                          ));
                    },
                  ),
                ),
              ],
            )
          ],
        ),
      ],
    );
  }

  bool isSmallScreen(BuildContext context) {
    final screenWidth = MediaQuery.of(context).size.width;
    return screenWidth < 600;
  }

  List<Widget> getMediaAll(List<MediaDetail> list) {
    if (list.isEmpty) {
      return [
        Container(
            height: MediaQuery.of(context).size.height * 0.6,
            alignment: Alignment.center,
            child: const Text(
              "啥都没有...",
              style: TextStyle(fontSize: 16),
            ))
      ];
    }
    if (onlyShowUnfinished) {
      list = list.where((v) => v.downloadedNum != v.monitoredNum).toList();
    }
    return List.generate(list.length, (i) {
      final item = list[i];
      return MediaCard(item: item);
    });
  }
}

class MediaCard extends StatelessWidget {
  final MediaDetail item;
  static const double smallWidth = 110;
  static const double largeWidth = 140;

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
                width: cardWidth(context),
                height: cardWidth(context) / 2 * 3,
                child: Ink.image(
                    fit: BoxFit.cover,
                    image: NetworkImage(
                      "${APIs.imagesUrl}/${item.id}/poster.jpg",
                    )),
              ),
              SizedBox(
                  width: cardWidth(context),
                  child: Column(
                    children: [
                      LinearProgressIndicator(
                        value: 1,
                        color: item.downloadedNum! >= item.monitoredNum!
                            ? Colors.teal
                            : Colors.lightGreen,
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

  double cardWidth(BuildContext context) {
    final screenWidth = MediaQuery.of(context).size.width;
    if (screenWidth < 600) {
      return smallWidth;
    }
    return largeWidth;
  }
}
