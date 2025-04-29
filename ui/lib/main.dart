import 'package:flutter/foundation.dart';
import 'package:flutter/material.dart';
import 'package:flutter_adaptive_scaffold/flutter_adaptive_scaffold.dart';
import 'package:flutter_riverpod/flutter_riverpod.dart';
import 'package:go_router/go_router.dart';
import 'package:intl/date_symbol_data_local.dart';
import 'package:ui/activity.dart';
import 'package:ui/ffi/lib_polaris_boot.dart';
import 'package:ui/init_wizard.dart';
import 'package:ui/login_page.dart';
import 'package:ui/movie_watchlist.dart';
import 'package:ui/providers/APIs.dart';
import 'package:ui/search_page/search.dart';
import 'package:ui/settings/settings.dart';
import 'package:ui/system_page.dart';
import 'package:ui/tv_details.dart';
import 'package:ui/welcome_page.dart';
import 'package:ui/widgets/utils.dart';

void main() async {
  if (!kIsWeb) {
    var port = await LibPolarisBoot.instance.start("");
    APIs.port = port;
  }

  initializeDateFormatting()
      .then((_) => runApp(const ProviderScope(child: MyApp())));
}

class MyApp extends ConsumerStatefulWidget {
  const MyApp({super.key});

  @override
  ConsumerState<ConsumerStatefulWidget> createState() {
    return _MyAppState();
  }
}

CustomTransitionPage buildPageWithDefaultTransition<T>({
  required BuildContext context,
  required GoRouterState state,
  required Widget child,
}) {
  return CustomTransitionPage<T>(
    transitionDuration: Duration.zero,
    reverseTransitionDuration: Duration.zero,
    key: state.pageKey,
    child: child,
    transitionsBuilder: (context, animation, secondaryAnimation, child) =>
        child,
  );
}

class _MyAppState extends ConsumerState<MyApp> {
  // This widget is the root of your application.

  @override
  Widget build(BuildContext context) {
    // GoRouter configuration
    final shellRoute = StatefulShellRoute.indexedStack(
      builder: (BuildContext context, GoRouterState state,
          StatefulNavigationShell navigationShell) {
        return MainSkeleton(
          body: navigationShell,
        );
      },
      branches: [
        StatefulShellBranch(initialLocation: WelcomePage.routeTv, routes: [
          GoRoute(
            path: WelcomePage.routeTv,
            pageBuilder: (context, state) => buildPageWithDefaultTransition(
                context: context, state: state, child: const WelcomePage()),
          ),
          GoRoute(
            path: TvDetailsPage.route,
            pageBuilder: (context, state) => buildPageWithDefaultTransition(
                context: context,
                state: state,
                child: TvDetailsPage(seriesId: state.pathParameters['id']!)),
          ),
          GoRoute(
            path: SearchPage.route,
            pageBuilder: (context, state) => buildPageWithDefaultTransition(
                context: context,
                state: state,
                child: SearchPage(query: state.uri.queryParameters["query"])),
          ),
        ]),
        StatefulShellBranch(initialLocation: WelcomePage.routeMoivie, routes: [
          GoRoute(
            path: WelcomePage.routeMoivie,
            pageBuilder: (context, state) => buildPageWithDefaultTransition(
                context: context, state: state, child: const WelcomePage()),
          ),
          GoRoute(
            path: MovieDetailsPage.route,
            pageBuilder: (context, state) => buildPageWithDefaultTransition(
                context: context,
                state: state,
                child: MovieDetailsPage(id: state.pathParameters['id']!)),
          ),
        ]),
        StatefulShellBranch(routes: [
          GoRoute(
            path: ActivityPage.route,
            pageBuilder: (context, state) => buildPageWithDefaultTransition(
                context: context, state: state, child: const ActivityPage()),
          ),
        ]),
        StatefulShellBranch(routes: [
          GoRoute(
            path: SystemSettingsPage.route,
            pageBuilder: (context, state) => buildPageWithDefaultTransition(
                context: context,
                state: state,
                child: const SystemSettingsPage()),
          ),
        ]),
        StatefulShellBranch(routes: [
          GoRoute(
            path: SystemPage.route,
            pageBuilder: (context, state) => buildPageWithDefaultTransition(
                context: context, state: state, child: const SystemPage()),
          )
        ]),
      ],
    );

    final router = GoRouter(
      navigatorKey: APIs.navigatorKey,
      initialLocation: WelcomePage.routeTv,
      routes: [
        shellRoute,
        GoRoute(
          path: "/",
          redirect: (context, state) => WelcomePage.routeTv,
        ),
        GoRoute(
          path: LoginScreen.route,
          builder: (context, state) => const LoginScreen(),
        ),
        GoRoute(
          path: InitWizard.route,
          builder: (context, state) => const InitWizard(),
        )
      ],
    );

    return MaterialApp.router(
      title: 'Polaris 影视追踪下载',
      theme: ThemeData(
          fontFamily: "NotoSansSC",
          colorScheme: ColorScheme.fromSeed(
              seedColor: Colors.blueAccent,
              brightness: Brightness.dark,
              surface: Colors.black87),
          useMaterial3: true,
          //scaffoldBackgroundColor: Color.fromARGB(255, 26, 24, 24)
          tooltipTheme: TooltipThemeData(
            textStyle: const TextStyle(
              color: Colors.grey,
            ),
            decoration: BoxDecoration(
              color: Colors.black54,
              borderRadius: BorderRadius.circular(20),
            ),
          )),
      routerConfig: router,
    );
  }
}

class MainSkeleton extends StatefulWidget {
  final StatefulNavigationShell body;
  const MainSkeleton({super.key, required this.body});

  @override
  State<StatefulWidget> createState() {
    return _MainSkeletonState();
  }
}

class _MainSkeletonState extends State<MainSkeleton> {
  @override
  Widget build(BuildContext context) {
    var padding = isSmallScreen(context) ? 5.0 : 20.0;
    return AdaptiveScaffold(
      appBarBreakpoint: Breakpoints.standard,
      appBar: appBar(),
      useDrawer: false,
      selectedIndex: widget.body.currentIndex,
      onSelectedIndexChange: (p0) => widget.body
          .goBranch(p0, initialLocation: p0 == widget.body.currentIndex),

      destinations: const <NavigationDestination>[
        NavigationDestination(
          icon: Icon(Icons.live_tv_outlined),
          selectedIcon: Icon(Icons.live_tv),
          label: '剧集',
        ),
        NavigationDestination(
          icon: Icon(Icons.movie_outlined),
          selectedIcon: Icon(Icons.movie),
          label: '电影',
        ),
        NavigationDestination(
          icon: Icon(Icons.download_outlined),
          selectedIcon: Icon(Icons.download),
          label: '活动',
        ),
        NavigationDestination(
          icon: Icon(Icons.settings_outlined),
          selectedIcon: Icon(Icons.settings),
          label: '设置',
        ),
        NavigationDestination(
          icon: Icon(Icons.computer_outlined),
          selectedIcon: Icon(Icons.computer),
          label: '系统',
        ),
      ],
      body: (context) => SafeArea(
          child: Padding(
              padding: EdgeInsets.only(
                  left: padding, right: padding, top: 5, bottom: 5),
              child: widget.body)),
      // Define a default secondaryBody.
      // Override the default secondaryBody during the smallBreakpoint to be
      // empty. Must use AdaptiveScaffold.emptyBuilder to ensure it is properly
      // overridden.
    );
  }

  showDonate(BuildContext context) {
    return showDialog<void>(
      context: context,
      barrierDismissible: true,
      builder: (BuildContext context) {
        return AlertDialog(
          title: Text("项目开发不易，给开发者加个鸡腿："),
          content: SizedBox(
              width: 350,
              height: 400,
              child: Ink.image(
                  fit: BoxFit.fitWidth,
                  image: AssetImage("assets/wechat.jpg"))),
        );
      },
    );
  }

  AppBar appBar() {
    return AppBar(
      // TRY THIS: Try changing the color here to a specific color (to
      // Colors.amber, perhaps?) and trigger a hot reload to see the AppBar
      // change color while the other colors stay the same.
      backgroundColor: Theme.of(context).colorScheme.inversePrimary,
      // Here we take the value from the MyHomePage object that was created by
      // the App.build method, and use it to set our appbar title.
      leading: Container(
        alignment: Alignment.centerLeft,
        child: TextButton(
          onPressed: () => context.go(WelcomePage.routeTv),
          child: const Text(
            "Polaris",
            overflow: TextOverflow.clip,
            style: TextStyle(fontSize: 28),
          ),
        ),
      ),
      leadingWidth: isSmallScreen(context) ? 0 : 190,
      title: Container(
        alignment: Alignment.bottomLeft,
        child: SearchAnchor(
            builder: (BuildContext context, SearchController controller) {
          return Container(
            constraints: const BoxConstraints(maxWidth: 250, maxHeight: 40),
            child: Opacity(
              opacity: 0.8,
              child: SearchBar(
                hintText: "在此搜索...",
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
      ),

      actions: [
        // IconButton(
        //     onPressed: () => showCalendar(context),
        //     icon: Icon(Icons.calendar_month)),
        IconButton(
            onPressed: () => showDonate(context),
            icon: Icon(
              Icons.favorite_rounded,
              color: Colors.red,
            )),

        MenuAnchor(
          menuChildren: [
            MenuItemButton(
              leadingIcon: const Icon(Icons.exit_to_app),
              child: const Text("登出"),
              onPressed: () async {
                await APIs.logout();
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
        ),
      ],
    );
  }
}
