import 'package:flutter/material.dart';
import 'package:flutter_riverpod/flutter_riverpod.dart';
import 'package:go_router/go_router.dart';
import 'package:shared_preferences/shared_preferences.dart';
import 'package:ui/activity.dart';
import 'package:ui/login_page.dart';
import 'package:ui/movie_watchlist.dart';
import 'package:ui/navdrawer.dart';
import 'package:ui/providers/APIs.dart';
import 'package:ui/search.dart';
import 'package:ui/system_settings.dart';
import 'package:ui/tv_details.dart';
import 'package:ui/welcome_page.dart';

void main() {
  runApp(const MyApp());
}

class MyApp extends StatelessWidget {
  const MyApp({super.key});

  // This widget is the root of your application.
  @override
  Widget build(BuildContext context) {
    // GoRouter configuration
    final shellRoute = ShellRoute(
      builder: (BuildContext context, GoRouterState state, Widget child) {
        return SelectionArea(
          child: Scaffold(
              appBar: AppBar(
                // TRY THIS: Try changing the color here to a specific color (to
                // Colors.amber, perhaps?) and trigger a hot reload to see the AppBar
                // change color while the other colors stay the same.
                backgroundColor: Theme.of(context).colorScheme.inversePrimary,
                // Here we take the value from the MyHomePage object that was created by
                // the App.build method, and use it to set our appbar title.
                title: Row(
                  children: [
                    const Text("Polaris"),
                    const SizedBox(
                      width: 100,
                    ),
                    IconButton(
                        tooltip: "搜索剧集",
                        onPressed: () => context.go(SearchPage.route),
                        icon: const Icon(Icons.search)),
                  ],
                ),
                actions: [
                  FutureBuilder(
                      future: APIs.isLoggedIn(),
                      builder: (context, snapshot) {
                        if (snapshot.hasData && snapshot.data == true) {
                          return MenuAnchor(
                            menuChildren: [
                              MenuItemButton(
                                leadingIcon: const Icon(Icons.exit_to_app),
                                child: const Text("登出"),
                                onPressed: () async {
                                  final SharedPreferences prefs =
                                      await SharedPreferences.getInstance();
                                  await prefs.remove('token');
                                  if (context.mounted) {
                                    context.go(LoginScreen.route);
                                  }
                                },
                              ),
                            ],
                            builder: (context, controller, child) {
                              return TextButton(
                                onPressed: () {
                                  if (controller.isOpen) {
                                    controller.close();
                                  } else {
                                    controller.open();
                                  }
                                },
                                child: const Icon(Icons.account_circle),
                              );
                            },
                          );
                        }
                        return Container();
                      })
                ],
              ),
              body: Center(
                  // Center is a layout widget. It takes a single child and positions it
                  // in the middle of the parent.
                  child: Flex(direction: Axis.horizontal, children: <Widget>[
                const Flexible(
                  flex: 1,
                  child: NavDrawer(),
                ),
                const VerticalDivider(thickness: 1, width: 1),
                Flexible(
                  flex: 7,
                  child:
                      Padding(padding: const EdgeInsets.all(20), child: child),
                )
              ]))),
        );
      },
      routes: [
        GoRoute(
          path: "/",
          redirect: (context, state) => WelcomePage.routeTv,
        ),
        GoRoute(
          path: WelcomePage.routeTv,
          builder: (context, state) => const WelcomePage(),
        ),
        GoRoute(
          path: TvDetailsPage.route,
          builder: (context, state) =>
              TvDetailsPage(seriesId: state.pathParameters['id']!),
        ),
        GoRoute(
          path: WelcomePage.routeMoivie,
          builder: (context, state) => const WelcomePage(),
        ),
        GoRoute(
          path: MovieDetailsPage.route,
          builder: (context, state) =>
              MovieDetailsPage(id: state.pathParameters['id']!),
        ),
        GoRoute(
          path: SearchPage.route,
          builder: (context, state) => const SearchPage(),
        ),
        GoRoute(
          path: SystemSettingsPage.route,
          builder: (context, state) => const SystemSettingsPage(),
        ),
        GoRoute(
          path: ActivityPage.route,
          builder: (context, state) => const ActivityPage(),
        )
      ],
    );

    final router = GoRouter(
      navigatorKey: APIs.navigatorKey,
      routes: [
        shellRoute,
        GoRoute(
          path: LoginScreen.route,
          builder: (context, state) => const LoginScreen(),
        )
      ],
    );

    return ProviderScope(
      child: MaterialApp.router(
        title: 'Polaris 电影电视剧追踪',
        theme: ThemeData(
          fontFamily: "NotoSansSC",
          colorScheme: ColorScheme.fromSeed(
              seedColor: Colors.blue, brightness: Brightness.dark),
          useMaterial3: true,
        ),
        routerConfig: router,
      ),
    );
  }
}

CustomTransitionPage buildPageWithDefaultTransition<T>({
  required BuildContext context,
  required GoRouterState state,
  required Widget child,
}) {
  return CustomTransitionPage<T>(
    key: state.pageKey,
    child: child,
    transitionsBuilder: (context, animation, secondaryAnimation, child) =>
        FadeTransition(opacity: animation, child: child),
  );
}
