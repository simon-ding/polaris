import 'package:flutter/material.dart';
import 'package:flutter_riverpod/flutter_riverpod.dart';
import 'package:percent_indicator/circular_percent_indicator.dart';
import 'package:ui/providers/activity.dart';
import 'package:ui/utils.dart';
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
            child: PaginatedDataTable(
              rowsPerPage: 10,
              columns: const [
                DataColumn(label: Text("#"), numeric: true),
                DataColumn(label: Text("名称")),
                DataColumn(label: Text("开始时间")),
                DataColumn(label: Text("状态")),
                DataColumn(label: Text("操作"))
              ],
              source: ActivityDataSource(activities: activities, onDelete: onDelete(ref)),
            ),
          );
        },
        error: (err, trace) => Text("$err"),
        loading: () => const MyProgressIndicator());
  }

  Function(int) onDelete(WidgetRef ref) {
    return (id) {
      ref
          .read(activitiesDataProvider.notifier)
          .deleteActivity(id)
          .whenComplete(() => Utils.showSnakeBar("删除成功"))
          .onError((error, trace) => Utils.showSnakeBar("删除失败：$error"));
    };
  }
}

class ActivityDataSource extends DataTableSource {
  List<Activity> activities;
  Function(int) onDelete;
  ActivityDataSource({required this.activities, required this.onDelete});

  @override
  int get rowCount => activities.length;

  @override
  DataRow? getRow(int index) {
    final activity = activities[index];
    return DataRow(cells: [
      DataCell(Text("${activity.id}")),
      DataCell(Text("${activity.sourceTitle}")),
      DataCell(Text("${activity.date!.toLocal()}")),
      DataCell(() {
        if (activity.status == "uploading") {
          return const SizedBox(
              width: 20, height: 20, child: CircularProgressIndicator());
        } else if (activity.status == "fail") {
          return const Icon(
            Icons.close,
            color: Colors.red,
          );
        } else if (activity.status == "success") {
          return const Icon(
            Icons.check,
            color: Colors.green,
          );
        }

        double p =
            activity.progress == null ? 0 : activity.progress!.toDouble() / 100;
        return CircularPercentIndicator(
          radius: 15.0,
          lineWidth: 5.0,
          percent: p,
          center: Text("${p * 100}"),
          progressColor: Colors.green,
        );
      }()),
      DataCell(IconButton(onPressed:() => onDelete(activity.id!), icon: const Icon(Icons.delete)))
    ]);
  }

  @override
  bool get isRowCountApproximate => false;

  @override
  int get selectedRowCount => 0;
}
