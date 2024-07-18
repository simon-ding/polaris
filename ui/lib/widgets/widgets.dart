import 'package:flutter/material.dart';

class Commons {
  static InputDecoration requiredTextFieldStyle({
    required String text,
    Widget? icon,
  }) {
    return InputDecoration(
      label: Row(
        children: [
          Text(text),
          const Text(
            "*",
            style: TextStyle(color: Colors.red),
          )
        ],
      ),
      icon: icon,
    );
  }
}

class SettingsCard extends StatelessWidget {
  final Widget? child;
  final GestureTapCallback? onTap;

  const SettingsCard({super.key, required this.child, this.onTap});

  @override
  Widget build(BuildContext context) {
    return Card(
      margin: const EdgeInsets.all(4),
      clipBehavior: Clip.hardEdge,
      child: InkWell(
          //splashColor: Colors.blue.withAlpha(30),
          onTap: onTap,
          child:
              SizedBox(width: 150, height: 150, child: Center(child: child))),
    );
  }
}
