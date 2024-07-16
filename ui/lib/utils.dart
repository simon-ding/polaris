import 'dart:math';

import 'package:flutter/material.dart';
import 'package:intl/intl.dart';

class Utils {
  static Future<void> showAlertDialog(BuildContext context, String msg) async {
    return showDialog<void>(
      context: context,
      barrierDismissible: true, // user must tap button!
      builder: (BuildContext context) {
        return AlertDialog(
          title: const Text('警告 ⚠️'),
          content: SingleChildScrollView(
            child: ListBody(
              children: <Widget>[
                Text(msg),
              ],
            ),
          ),
          actions: <Widget>[
            TextButton(
              child: const Text('确定'),
              onPressed: () {
                Navigator.of(context).pop();
              },
            ),
          ],
        );
      },
    );
  }

  static showSnakeBar(BuildContext context, String msg) {
    ScaffoldMessenger.of(context).showSnackBar(SnackBar(content: Text(msg), showCloseIcon: true,));
  }

  static bool showError(BuildContext context, AsyncSnapshot snapshot) {
    final isErrored = snapshot.hasError &&
        snapshot.connectionState != ConnectionState.waiting;
    if (isErrored) {
      Utils.showSnakeBar(context, "当前操作出错: ${snapshot.error}");
      return true;
    }
    return false;
  }
}


extension FileFormatter on num {
  String readableFileSize({bool base1024 = true}) {
    final base = base1024 ? 1024 : 1000;
    if (this <= 0) return "0";
    final units = ["B", "kB", "MB", "GB", "TB"];
    int digitGroups = (log(this) / log(base)).round();
    return "${NumberFormat("#,##0.#").format(this / pow(base, digitGroups))} ${units[digitGroups]}";
  }
}
