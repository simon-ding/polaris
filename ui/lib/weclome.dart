import 'package:flutter/material.dart';
import 'package:flutter_riverpod/flutter_riverpod.dart';
import 'package:go_router/go_router.dart';
import 'package:ui/APIs.dart';
import 'package:ui/providers/welcome_data.dart';
import 'package:ui/tv_details.dart';

class WelcomePage extends ConsumerWidget {
  static const route = "/welcome";

  const WelcomePage({super.key});

  @override
  Widget build(BuildContext context, WidgetRef ref) {
    final data = ref.watch(welcomePageDataProvider);

    return switch (data) {
      AsyncData(:final value) => GridView.builder(
        itemCount: value.length,
        gridDelegate:
            const SliverGridDelegateWithFixedCrossAxisCount(crossAxisCount: 6),
        itemBuilder: (context, i) {
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
                    Flexible(
                      child: SizedBox(
                        width: 300,
                        height: 600,
                        child: Image.network(
                          APIs.tmdbImgBaseUrl + item.posterPath!,
                          fit: BoxFit.contain,
                        ),
                      ),
                    ),
                    Flexible(
                      child: Text(
                        item.name!,
                        style: const TextStyle(
                            fontSize: 14, fontWeight: FontWeight.bold),
                      ),
                    )
                  ],
                ),
              ));
        }),
      _ => const CircularProgressIndicator(),
      
    };
  }
}
