import 'package:flutter/material.dart';
import 'package:flutter_riverpod/flutter_riverpod.dart';
import 'package:go_router/go_router.dart';
import 'package:ui/providers/APIs.dart';
import 'package:ui/providers/welcome_data.dart';
import 'package:ui/tv_details.dart';
import 'package:ui/widgets/progress_indicator.dart';

class TvWatchlistPage extends ConsumerWidget {
  static const route = "/series";

  const TvWatchlistPage({super.key});

  @override
  Widget build(BuildContext context, WidgetRef ref) {
    final data = ref.watch(tvWatchlistDataProvider);

    return switch (data) {
      AsyncData(:final value) => SingleChildScrollView(
          child: Wrap(
            spacing: 20,
            children: List.generate(value.length, (i) {
              var item = value[i];
              return Card(
                  margin: const EdgeInsets.all(4),
                  clipBehavior: Clip.hardEdge,
                  child: InkWell(
                    //splashColor: Colors.blue.withAlpha(30),
                    onTap: () {
                      context.go(TvDetailsPage.toRoute(item.id!));
                      //showDialog(context: context, builder: builder)
                    },
                    child: Column(
                      children: <Widget>[
                        SizedBox(
                            width: 160,
                            height: 240,
                            child:Image.network(
                                "${APIs.imagesUrl}/${item.id}/poster.jpg",
                                fit: BoxFit.fill,
                                headers: APIs.authHeaders,
                              ),
                            ),
                        Text(
                          item.name!,
                          style: const TextStyle(
                              fontSize: 14, fontWeight: FontWeight.bold, height: 2.5),
                        ),
                      ],
                    ),
                  ));
            }),
          ),
        ),
      _ => MyProgressIndicator(),
    };
  }
}
