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
    _nestedTabController = new TabController(length: 4, vsync: this);
  }

  @override
  void dispose() {
    super.dispose();
    _nestedTabController.dispose();
  }

  int selectedTab = 0;

  @override
  Widget build(BuildContext context) {
    return SelectionArea(
        child: Column(
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
              text: "做种中",
            ),
            Tab(
              text: "历史记录",
            ),
            Tab(
              text: "黑名单",
            )
          ],
        ),
        Builder(builder: (context) {
          AsyncValue<List<Activity>>? activitiesWatcher;

          if (selectedTab == 2) {
            activitiesWatcher =
                ref.watch(activitiesDataProvider(ActivityStatus.archive));
          } else if (selectedTab == 1) {
            activitiesWatcher =
                ref.watch(activitiesDataProvider(ActivityStatus.seeding));
          } else if (selectedTab == 0) {
            activitiesWatcher =
                ref.watch(activitiesDataProvider(ActivityStatus.active));
          } else if (selectedTab == 3) {
            return showBlacklistTab();
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
                            } else if (ac.status == "removed") {
                              return const Tooltip(
                                message: "已删除",
                                child: Icon(
                                  Icons.remove,
                                  //color: Colors.orange,
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
                                (ac.seedRatio ?? 0) > 0
                                    ? Text("分享率：${ac.seedRatio}")
                                    : SizedBox()
                              ],
                            ),
                          ),
                          trailing: selectedTab == 0
                              ? IconButton(
                                  tooltip: "删除任务",
                                  onPressed: () =>
                                      showConfirmDialog(context, ac.id!),
                                  icon: const Icon(Icons.delete))
                              : (selectedTab == 1
                                  ? IconButton(
                                      tooltip: "完成做种",
                                      onPressed: () => {
                                            ref
                                                .read(activitiesDataProvider(
                                                        ActivityStatus.active)
                                                    .notifier)
                                                .deleteActivity(ac.id!, false)
                                          },
                                      icon: const Icon(Icons.check))
                                  : const Text("-")),
                        ),
                        Divider(),
                      ],
                    );
                  },
                ));
              },
              error: (err, trace) => PoNetworkError(err: err),
              loading: () => const MyProgressIndicator());
        })
      ],
    ));
  }

  Future<void> showConfirmDialog(BuildContext oriContext, int id) {
    var add2Blacklist = false;
    return showDialog<void>(
      context: context,
      barrierDismissible: true,
      builder: (BuildContext context) {
        return AlertDialog(
          title: const Text("确认删除"),
          content: StatefulBuilder(builder: (context, setState) {
            return CheckboxListTile(
                value: add2Blacklist,
                title: Text("加入黑名单"),
                onChanged: (v) {
                  setState(
                    () {
                      add2Blacklist = v!;
                    },
                  );
                });
          }),
          actions: [
            TextButton(
                onPressed: () => Navigator.of(context).pop(),
                child: const Text("取消")),
            TextButton(
                child: const Text("确认"),
                onPressed: () {
                  final f = ref
                      .read(activitiesDataProvider(ActivityStatus.active)
                          .notifier)
                      .deleteActivity(id, add2Blacklist)
                      .then((value) {
                    Navigator.of(context).pop();
                  });
                  showLoadingWithFuture(f);
                }),
          ],
        );
      },
    );
  }

  Widget showBlacklistTab() {
    var blacklistDataWacher = ref.watch(blacklistDataProvider);

    return blacklistDataWacher.when(
      data: (blacklists) {
        return Flexible(
            child: SelectionArea(
                child: ListView.builder(
                    itemCount: blacklists.length,
                    itemBuilder: (context, index) {
                      final item = blacklists[index];
                      return ListTile(
                        dense: true,
                        title: Text(item.torrentName ?? ""),
                        subtitle: Opacity(
                            opacity: 0.7,
                            child: Text("hash: ${item.torrentHash ?? ""}")),
                        trailing: IconButton(
                            onPressed: () {
                              final f = ref
                                  .read(blacklistDataProvider.notifier)
                                  .deleteBlacklist(item.id!)
                                  .then((value) {
                                //Navigator.of(context).pop();
                              });
                              showLoadingWithFuture(f);
                            },
                            icon: const Icon(Icons.delete)),
                      );
                    })));
      },
      error: (err, trace) => PoNetworkError(err: err),
      loading: () => const MyProgressIndicator(),
    );
  }
}
