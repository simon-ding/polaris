import 'package:flutter/material.dart';
import 'package:flutter_adaptive_scaffold/flutter_adaptive_scaffold.dart';
import 'package:flutter_riverpod/flutter_riverpod.dart';
import 'package:go_router/go_router.dart';
import 'package:shared_preferences/shared_preferences.dart';
import 'package:ui/activity.dart';
import 'package:ui/login_page.dart';
import 'package:ui/movie_watchlist.dart';
import 'package:ui/providers/APIs.dart';
import 'package:ui/search.dart';
import 'package:ui/settings.dart';
import 'package:ui/system_page.dart';
import 'package:ui/tv_details.dart';
import 'package:ui/welcome_page.dart';

void main() {
  runApp(const MyApp());
}

class MyApp extends ConsumerStatefulWidget {
  const MyApp({super.key});

  @override
  ConsumerState<ConsumerStatefulWidget> createState() {
    return _MyAppState();
  }
}

class _MyAppState extends ConsumerState<MyApp> {
  // This widget is the root of your application.
  @override
  Widget build(BuildContext context) {
    // GoRouter configuration
    final shellRoute = ShellRoute(
      builder: (BuildContext context, GoRouterState state, Widget child) {
        return SelectionArea(
          child: MainSkeleton(
            body: Padding(padding: const EdgeInsets.all(20), child: child),
          ),
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
          builder: (context, state) =>
              SearchPage(query: state.uri.queryParameters["query"]),
        ),
        GoRoute(
          path: SystemSettingsPage.route,
          builder: (context, state) => const SystemSettingsPage(),
        ),
        GoRoute(
          path: ActivityPage.route,
          builder: (context, state) => const ActivityPage(),
        ),
        GoRoute(
          path: SystemPage.route,
          builder: (context, state) => const SystemPage(),
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
        title: 'Polaris 影视追踪',
        theme: ThemeData(
          fontFamily: "NotoSansSC",
          colorScheme: ColorScheme.fromSeed(
              seedColor: Colors.blueAccent,
              brightness: Brightness.dark,
              surface: Colors.black54),
          useMaterial3: true,
          //scaffoldBackgroundColor: Color.fromARGB(255, 26, 24, 24)
        ),
        routerConfig: router,
      ),
    );
  }
}

class MainSkeleton extends StatefulWidget {
  final Widget body;
  const MainSkeleton({super.key, required this.body});

  @override
  State<StatefulWidget> createState() {
    return _MainSkeletonState();
  }
}

class _MainSkeletonState extends State<MainSkeleton> {
  var _selectedTab;

  @override
  Widget build(BuildContext context) {
    var uri = GoRouterState.of(context).uri.toString();
    if (uri.contains(WelcomePage.routeTv)) {
      _selectedTab = 0;
    } else if (uri.contains(WelcomePage.routeMoivie)) {
      _selectedTab = 1;
    } else if (uri.contains(ActivityPage.route)) {
      _selectedTab = 2;
    } else if (uri.contains(SystemSettingsPage.route)) {
      _selectedTab = 3;
    } else if (uri.contains(SystemPage.route)) {
      _selectedTab = 4;
    }

    return AdaptiveScaffold(
      appBarBreakpoint: Breakpoints.standard,
      appBar: AppBar(
        // TRY THIS: Try changing the color here to a specific color (to
        // Colors.amber, perhaps?) and trigger a hot reload to see the AppBar
        // change color while the other colors stay the same.
        backgroundColor: Theme.of(context).colorScheme.inversePrimary,
        // Here we take the value from the MyHomePage object that was created by
        // the App.build method, and use it to set our appbar title.
        title: const Row(
          children: [
            Text("Polaris"),
          ],
        ),
        actions: [
          SearchAnchor(
              builder: (BuildContext context, SearchController controller) {
            return Container(
              constraints: const BoxConstraints(maxWidth: 300, maxHeight: 40),
              child: Opacity(
                opacity: 0.8,
                child: SearchBar(
                  hintText: "搜索...",
                  leading: const Icon(Icons.search),
                  controller: controller,
                  shadowColor: WidgetStateColor.transparent,
                  backgroundColor: WidgetStatePropertyAll(
                      Theme.of(context).colorScheme.primaryContainer),
                  onSubmitted: (value) => context.go(Uri(
                      path: SearchPage.route,
                      queryParameters: {'query': value}).toString()),
                ),
              ),
            );
          }, suggestionsBuilder:
                  (BuildContext context, SearchController controller) {
            return [Text("dadada")];
          }),
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
      useDrawer: false,
      selectedIndex: _selectedTab,
      onSelectedIndexChange: (int index) {
        setState(() {
          _selectedTab = index;
        });
        if (index == 0) {
          context.go(WelcomePage.routeTv);
        } else if (index == 1) {
          context.go(WelcomePage.routeMoivie);
        } else if (index == 2) {
          context.go(ActivityPage.route);
        } else if (index == 3) {
          context.go(SystemSettingsPage.route);
        } else if (index == 4) {
          context.go(SystemPage.route);
        }
      },
      destinations: const <NavigationDestination>[
        NavigationDestination(
          icon: Icon(Icons.live_tv),
          label: '电视剧',
        ),
        NavigationDestination(
          icon: Icon(Icons.movie),
          label: '电影',
        ),
        NavigationDestination(
          icon: Icon(Icons.download),
          label: '活动',
        ),
        NavigationDestination(
          icon: Icon(Icons.settings),
          label: '设置',
        ),
        NavigationDestination(
          icon: Icon(Icons.computer_rounded),
          label: '系统',
        ),
      ],
      body: (context) => widget.body,
      // Define a default secondaryBody.
      // Override the default secondaryBody during the smallBreakpoint to be
      // empty. Must use AdaptiveScaffold.emptyBuilder to ensure it is properly
      // overridden.
    );
  }
}
