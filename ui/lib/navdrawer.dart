import 'package:flutter/material.dart';
import 'package:go_router/go_router.dart';
import 'package:ui/activity.dart';
import 'package:ui/search.dart';
import 'package:ui/system_settings.dart';
import 'package:ui/welcome_page.dart';

class NavDrawer extends StatefulWidget {
  const NavDrawer({super.key});

  @override
  State<StatefulWidget> createState() {
    return _NavDrawerState();
  }
}

class _NavDrawerState extends State<NavDrawer> {
  var _counter = 0;

  @override
  Widget build(BuildContext context) {
    return NavigationRail(
      selectedIndex: _counter,
      onDestinationSelected: (value) {
        setState(() {
          _counter = value;
        });
        if (value == 0) {
          context.go(WelcomePage.routeMoivie);
        } else if (value == 1) {
          context.go(WelcomePage.routeTv);
        } else if (value == 2) {
          context.go(SearchPage.route);
        } else if (value == 3) {
          context.go(ActivityPage.route);
        } else if (value == 4) {
          context.go(SystemSettingsPage.route);
        }
      },
      extended: MediaQuery.of(context).size.width >= 850,
      unselectedIconTheme: const IconThemeData(color: Colors.grey),
      destinations: const <NavigationRailDestination>[
        NavigationRailDestination(
          icon: Icon(Icons.movie),
          label: Text('电影'),
        ),
        NavigationRailDestination(
          icon: Icon(Icons.live_tv),
          label: Text('电视剧'),
        ),
        NavigationRailDestination(
          icon: Icon(Icons.search),
          label: Text('搜索'),
        ),
        NavigationRailDestination(
          icon: Icon(Icons.download),
          label: Text('活动'),
        ),
        NavigationRailDestination(
          icon: Icon(Icons.settings),
          label: Text('设置'),
        ),
      ],
    );
  }
}
