import 'package:flutter/material.dart';

class MyProgressIndicator extends StatelessWidget {
  double size;
  Color? color;
  double? value;

  MyProgressIndicator({super.key, this.size = 30, this.color, this.value});
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
