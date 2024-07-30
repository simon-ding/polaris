import 'package:flutter/material.dart';
import 'package:flutter_riverpod/flutter_riverpod.dart';
import 'package:percent_indicator/circular_percent_indicator.dart';
import 'package:ui/providers/activity.dart';
import 'package:ui/widgets/progress_indicator.dart';
import 'package:ui/widgets/widgets.dart';

class ActivityPage extends ConsumerStatefulWidget {
  const ActivityPage({super.key});
  static const route = "/activities";

  @override
  _ActivityPageState createState() => _ActivityPageState();
}

class _ActivityPageState extends ConsumerState<ActivityPage>
    with TickerProviderStateMixin {
  late TabController _nestedTabController;
  @override
  void initState() {
    super.initState();
    _nestedTabController = new TabController(length: 2, vsync: this);
  }

  @override
  void dispose() {
    super.dispose();
    _nestedTabController.dispose();
  }

  int selectedTab = 0;

  @override
  Widget build(BuildContext context) {
    return Column(
      crossAxisAlignment: CrossAxisAlignment.stretch,
      children: [
        TabBar(
          controller: _nestedTabController,
          isScrollable: true,
          onTap: (value) {
            setState(() {
              selectedTab = value;
            });
          },
          tabs: const <Widget>[
            Tab(
              text: "下载中",
            ),
            Tab(
              text: "历史记录",
            ),
          ],
        ),
        Builder(builder: (context) {
          var activitiesWatcher = ref.watch(activitiesDataProvider("active"));
          if (selectedTab == 1) {
            activitiesWatcher = ref.watch(activitiesDataProvider("archive"));
          }

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
                    source: ActivityDataSource(
                        activities: activities,
                        onDelete: selectedTab == 0 ? onDelete() : null),
                  ),
                );
              },
              error: (err, trace) => Text("$err"),
              loading: () => const MyProgressIndicator());
        })
      ],
    );
  }

  Function(int) onDelete() {
    return (id) {
      final f = ref
          .read(activitiesDataProvider("active").notifier)
          .deleteActivity(id);
      showLoadingWithFuture(f);
    };
  }
}

class ActivityDataSource extends DataTableSource {
  List<Activity> activities;
  Function(int)? onDelete;
  ActivityDataSource({required this.activities, this.onDelete});

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
              width: 20,
              height: 20,
              child: Tooltip(
                message: "正在上传到指定存储",
                child: CircularProgressIndicator(),
              ));
        } else if (activity.status == "fail") {
          return const Tooltip(
              message: "下载失败",
              child: Icon(
                Icons.close,
                color: Colors.red,
              ));
        } else if (activity.status == "success") {
          return const Tooltip(
            message: "下载成功",
            child: Icon(
              Icons.check,
              color: Colors.green,
            ),
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
      onDelete != null
          ? DataCell(Tooltip(
              message: "删除任务",
              child: IconButton(
                  onPressed: () => onDelete!(activity.id!),
                  icon: const Icon(Icons.delete))))
          : const DataCell(Text("-"))
    ]);
  }

  @override
  bool get isRowCountApproximate => false;

  @override
  int get selectedRowCount => 0;
}
