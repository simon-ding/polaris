import 'package:flutter/material.dart';
import 'package:ui/widgets/widgets.dart';

Future<void> showSettingDialog(
    BuildContext context,
    String title,
    bool showDelete,
    Widget body,
    Future Function() onSubmit,
    Future Function() onDelete) {
  return showDialog<void>(
      context: context,
      barrierDismissible: true,
      builder: (BuildContext context) {
        return AlertDialog(
          title: Text(title),
          content: SingleChildScrollView(
            child: SizedBox(
              width: 400,
              child: SelectionArea(child: body),
            ),
          ),
          actions: <Widget>[
            showDelete
                ? LoadingTextButton(
                    onPressed: () async {
                      await onDelete().then((v) => Navigator.of(context).pop());
                    },
                    label: const Text(
                      '删除',
                      style: TextStyle(color: Colors.red),
                    ))
                : const Text(""),
            TextButton(
                onPressed: () => Navigator.of(context).pop(),
                child: const Text('取消')),
            LoadingTextButton(
              label: const Text('确定'),
              onPressed: () async {
                await onSubmit().then((v) => Navigator.of(context).pop());
              },
            ),
          ],
        );
      });
}
