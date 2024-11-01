import 'package:flutter/material.dart';
import 'package:flutter_form_builder/flutter_form_builder.dart';
import 'package:flutter_riverpod/flutter_riverpod.dart';
import 'package:go_router/go_router.dart';
import 'package:ui/movie_watchlist.dart';
import 'package:ui/providers/APIs.dart';
import 'package:ui/providers/welcome_data.dart';
import 'package:ui/tv_details.dart';
import 'package:ui/widgets/progress_indicator.dart';
import 'package:ui/widgets/utils.dart';
import 'package:ui/widgets/widgets.dart';

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
  final _formKey = GlobalKey<FormBuilderState>();
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
        getMoreButtonAndActions(uri)
      ],
    );
  }

  Widget getMoreButtonAndActions(String uri) {
    return Row(
      children: [
        Expanded(child: Container()),
        Column(
          children: [
            Expanded(child: Container()),
            Padding(
              padding: EdgeInsets.all(20),
              child: MenuAnchor(
                style: MenuStyle(
                  alignment: Alignment.topLeft,
                  backgroundColor: WidgetStatePropertyAll(Theme.of(context)
                      .colorScheme
                      .inversePrimary
                      .withOpacity(0.7)),
                ),
                menuChildren: [parseName(), onlyUnfinished(), refreshAll(uri)],
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
    );
  }

  Widget onlyUnfinished() {
    return CheckboxListTile(
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
    );
  }

  Widget refreshAll(String uri) {
    return LoadingListTile(
      icon: Icons.refresh,
      text: "全部更新",
      onPressed: () async {
        if (uri == WelcomePage.routeMoivie) {
          await APIs.downloadAllMovies().then((v) {
            showSnakeBar("开始下载电影：$v");
          });
        } else {
          await APIs.downloadAllTv().then((v) {
            showSnakeBar("开始下载剧集：$v");
          });
        }
      },
    );
  }

  Widget parseName() {
    return ListTile(
      leading: Icon(Icons.calculate),
      title: Text("测试解析"),
      onTap: () => _showNameParsingDialog(),
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

  Future<void> _showNameParsingDialog() async {
    final resultController = TextEditingController();
    return showDialog<void>(
      context: context,
      barrierDismissible: true,
      builder: (BuildContext context) {
        return AlertDialog(
          title: const Text('测试名称解析'),
          content: SizedBox(
            width: 500,
            height: 400,
            child: FormBuilder(
              key: _formKey,
              initialValue: {"name": "", "type": "tv"},
              child: Column(
                children: [
                  FormBuilderTextField(
                    name: "name",
                    decoration: InputDecoration(labelText: "要解析的名字"),
                  ),
                  FormBuilderDropdown(
                    name: "type",
                    items: [
                      DropdownMenuItem(
                        value: "tv",
                        child: const Text("电视剧"),
                      ),
                      DropdownMenuItem(value: "movie", child: const Text("电影"))
                    ],
                  ),
                  Center(
                    child: Padding(
                      padding: EdgeInsets.all(10),
                      child: LoadingTextButton(
                          onPressed: () async {
                            if (_formKey.currentState!.saveAndValidate()) {
                              final values = _formKey.currentState!.value;
                              //print(values);
                              if (values["type"] == "tv") {
                                var s = await APIs.parseTvName(values["name"]);
                                resultController.text = s;
                              } else {
                                var s =
                                    await APIs.parseMovieName(values["name"]);
                                resultController.text = s;
                              }
                            }
                            return;
                          },
                          label: Text("解析")),
                    ),
                  ),
                  TextField(
                    maxLines: 8,
                    controller: resultController,
                  )
                ],
              ),
            ),
          ),
        );
      },
    );
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
