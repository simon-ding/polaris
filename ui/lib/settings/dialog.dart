import 'package:flutter/material.dart';
import 'package:ui/widgets/utils.dart';

Future<void> showSettingDialog(BuildContext context,String title, bool showDelete, Widget body,
    Future Function() onSubmit, Future Function() onDelete) {
  return showDialog<void>(
      context: context,
      barrierDismissible: true,
      builder: (BuildContext context) {
        return AlertDialog(
          title: Text(title),
          content: SingleChildScrollView(
            child: SizedBox(
              width: 300,
              child: body,
            ),
          ),
          actions: <Widget>[
            showDelete
                ? TextButton(
                    onPressed: () {
                      final f = onDelete();
                      f.then((v) {
                        Utils.showSnakeBar("删除成功");
                        Navigator.of(context).pop();
                      }).onError((e, s) {
                        Utils.showSnakeBar("删除失败：$e");
                      });
                    },
                    child: const Text(
                      '删除',
                      style: TextStyle(color: Colors.red),
                    ))
                : const Text(""),
            TextButton(
                onPressed: () => Navigator.of(context).pop(),
                child: const Text('取消')),
            TextButton(
              child: const Text('确定'),
              onPressed: () {
                final f = onSubmit();
                f.then((v) {
                  Utils.showSnakeBar("操作成功");
                  Navigator.of(context).pop();
                }).onError((e, s) {
                  if (e.toString() != "validation_error") {
                    Utils.showSnakeBar("操作失败：$e");
                  }
                });
              },
            ),
          ],
        );
      });
}
