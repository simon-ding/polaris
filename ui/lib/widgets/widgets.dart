import 'package:flutter/material.dart';
import 'package:ui/providers/APIs.dart';

class Commons {
  static InputDecoration requiredTextFieldStyle({
    required String text,
    Widget? icon,
    String? helperText,
  }) {
    return InputDecoration(
      helperText: helperText,
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

showLoadingWithFuture(Future f) {
  final context = APIs.navigatorKey.currentContext;
  if (context == null) {
    return;
  }
  showDialog(
    context: context,
    barrierDismissible: false, //点击遮罩不关闭对话框
    builder: (context) {
      return FutureBuilder(
          future: f,
          builder: (context, snapshot) {
            if (snapshot.connectionState == ConnectionState.done) {
              if (snapshot.hasError) {
                return AlertDialog(
                  content: Text("处理失败：${snapshot.error}"),
                  actions: [
                    TextButton(
                        onPressed: () => Navigator.of(context).pop(),
                        child: const Text("好"))
                  ],
                );
              }
              Navigator.of(context).pop();
              return Container();
            } else {
              return const AlertDialog(
                content: Column(
                  mainAxisSize: MainAxisSize.min,
                  children: <Widget>[
                    CircularProgressIndicator(),
                    Padding(
                      padding: EdgeInsets.only(top: 26.0),
                      child: Text("正在处理，请稍后..."),
                    )
                  ],
                ),
              );
            }
          });
    },
  );
}
