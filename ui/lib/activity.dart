import 'package:flutter/material.dart';
import 'package:flutter_riverpod/flutter_riverpod.dart';
import 'package:ui/providers/activity.dart';
import 'package:ui/widgets/progress_indicator.dart';
import 'package:ui/widgets/utils.dart';
import 'package:ui/widgets/widgets.dart';
import 'package:timeago/timeago.dart' as timeago;

class ActivityPage extends ConsumerStatefulWidget {
  const ActivityPage({super.key});
  static const route = "/activities";

  @override
  ConsumerState<ConsumerStatefulWidget> createState() => _ActivityPageState();
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
          AsyncValue<List<Activity>>? activitiesWatcher;

          if (selectedTab == 1) {
            activitiesWatcher = ref.watch(activitiesDataProvider("archive"));
          } else if (selectedTab == 0) {
            activitiesWatcher = ref.watch(activitiesDataProvider("active"));
          }

          return activitiesWatcher!.when(
              data: (activities) {
                return Flexible(
                    child: ListView.builder(
                  itemCount: activities.length,
                  itemBuilder: (context, index) {
                    final ac = activities[index];
                    return Column(
                      children: [
                        ListTile(
                          dense: true,
                          leading: () {
                            if (ac.status == "uploading") {
                              return const SizedBox(
                                  width: 20,
                                  height: 20,
                                  child: Tooltip(
                                    message: "正在上传到指定存储",
                                    child: CircularProgressIndicator(),
                                  ));
                            } else if (ac.status == "fail") {
                              return const Tooltip(
                                  message: "下载失败",
                                  child: Icon(
                                    Icons.close,
                                    color: Colors.red,
                                  ));
                            } else if (ac.status == "seeding") {
                              //seeding
                              return Tooltip(
                                message: "做种中",
                                child: Icon(
                                  Icons.keyboard_double_arrow_up,
                                  color: Theme.of(context)
                                      .colorScheme
                                      .inversePrimary,
                                ),
                              );
                            } else if (ac.status == "success") {
                              return const Tooltip(
                                message: "下载成功",
                                child: Icon(
                                  Icons.check,
                                  color: Colors.green,
                                ),
                              );
                            }

                            double p = ac.progress == null
                                ? 0
                                : ac.progress!.toDouble() / 100;
                            return Tooltip(
                              message: "${ac.progress}%",
                              child: CircularProgressIndicator(
                                backgroundColor: Colors.black26,
                                value: p,
                              ),
                            );
                          }(),
                          title: Text((ac.sourceTitle ?? "")),
                          subtitle: Opacity(
                            opacity: 0.7,
                            child: Wrap(
                              spacing: 10,
                              children: [
                                Text("开始时间：${timeago.format(ac.date!)}"),
                                Text("大小：${(ac.size ?? 0).readableFileSize()}"),
                                ac.seedRatio > 0
                                    ? Text("分享率：${ac.seedRatio}")
                                    : SizedBox()
                              ],
                            ),
                          ),
                          trailing: selectedTab == 0
                              ? IconButton(
                                  tooltip: "删除任务",
                                  onPressed: () => onDelete()(ac.id!),
                                  icon: const Icon(Icons.delete))
                              : const Text("-"),
                        ),
                        Divider(),
                      ],
                    );
                  },
                ));
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
