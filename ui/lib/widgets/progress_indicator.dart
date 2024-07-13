import 'package:flutter/material.dart';

class MyProgressIndicator extends StatelessWidget {
  final double size;
  final Color? color;
  final double? value;

  const MyProgressIndicator({super.key, this.size = 30, this.color, this.value});
  @override
  Widget build(BuildContext context) {
    return Center(
        child: SizedBox(
            width: size,
            height: size,
            child: CircularProgressIndicator(
              color: color,
              value: value,
            )));
  }
}
