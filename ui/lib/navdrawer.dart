import 'package:flutter/material.dart';

class NavDrawer extends StatelessWidget {
  var _counter = 0;

  @override
  Widget build(BuildContext context) {
    return ConstrainedBox(
        constraints: const BoxConstraints(maxWidth: 300),
        child: Row(
            mainAxisAlignment: MainAxisAlignment.start,
            crossAxisAlignment: CrossAxisAlignment.start,
            children: [
              Flexible(
                  child: NavigationRail(
                selectedIndex: _counter,
                onDestinationSelected: (value) {},
                extended: MediaQuery.of(context).size.width >= 850,
                unselectedIconTheme: const IconThemeData(color: Colors.grey),
                destinations: const <NavigationRailDestination>[
                  NavigationRailDestination(
                    icon: Icon(Icons.search),
                    label: Text('Buscar '),
                  ),
                  NavigationRailDestination(
                    icon: Icon(Icons.engineering),
                    label: Text('Ingeniería '),
                  ),
                  NavigationRailDestination(
                    icon: Icon(Icons.business),
                    label: Text('Sociales '),
                  ),
                  NavigationRailDestination(
                    icon: Icon(Icons.local_hospital),
                    label: Text('Salud'),
                  ),
                  NavigationRailDestination(
                    icon: Icon(Icons.school),
                    label: Text('Iniciales'),
                  ),
                  NavigationRailDestination(
                    icon: Icon(Icons.design_services),
                    label: Text('Talleres y Extracurriculares'),
                  ),
                ],
                selectedLabelTextStyle:
                    TextStyle(color: Colors.lightBlue, fontSize: 20),
                unselectedLabelTextStyle:
                    TextStyle(color: Colors.grey, fontSize: 18),
                groupAlignment: -1,
                minWidth: 56,
              ))
            ]));
  }
}
