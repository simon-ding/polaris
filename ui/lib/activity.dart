import 'package:flutter/material.dart';
import 'package:flutter_riverpod/flutter_riverpod.dart';
import 'package:ui/providers/activity.dart';
import 'package:ui/widgets/progress_indicator.dart';

class ActivityPage extends ConsumerWidget {
  static const route = "/activities";

  const ActivityPage({super.key});
  @override
  Widget build(BuildContext context, WidgetRef ref) {
    var activitiesWatcher = ref.watch(activitiesDataProvider);

    return activitiesWatcher.when(
        data: (activities) {
          return SingleChildScrollView(
            child: DataTable(
              columns: const [
                DataColumn(label: Text("id"), numeric: true),
                DataColumn(label: Text("名称")),
                DataColumn(label: Text("开始时间")),
                DataColumn(label: Text("状态")),
                DataColumn(label: Text("操作"))
              ],
              rows: List<DataRow>.generate(activities.length, (i) {
                var activity = activities[activities.length - i - 1];

                return DataRow(cells: [
                  DataCell(Text("${activity.id}")),
                  DataCell(Text("${activity.sourceTitle}")),
                  DataCell(Text("${activity.date!.toLocal()}")),
                  DataCell(() {
                    if (activity.inBackgroud == true) {
                      return const MyProgressIndicator(
                        size: 20,
                      );
                    }

                    if (activity.completed != true && activity.progress == 0) {
                      return const Icon(
                        Icons.close,
                        color: Colors.red,
                      );
                    }

                    return MyProgressIndicator(
                      value: activity.progress!.toDouble() / 100,
                      size: 20,
                    );
                  }()),
                  DataCell(IconButton(
                      onPressed: () {
                        ref
                            .read(activitiesDataProvider.notifier)
                            .deleteActivity(activity.id!);
                      },
                      icon: const Icon(Icons.delete)))
                ]);
              }),
            ),
          );
        },
        error: (err, trace) => Text("$err"),
        loading: () => MyProgressIndicator());
  }
}
