import 'package:client/pages/acquire_data_page.dart';
import 'package:client/pages/ber_logging_page.dart';
import 'package:client/pages/data_transfer_page.dart';
import 'package:client/pages/database_options_page.dart';
import 'package:client/pages/developer_options_page.dart';
import 'package:client/pages/file_based_processing_page.dart';
import 'package:client/pages/home_page.dart';
import 'package:client/pages/offline_processing_page.dart';
import 'package:client/pages/result_profiles_page.dart';
import 'package:client/pages/trend_analysis_page.dart';
import 'package:client/state/global_state.dart';
import 'package:flutter/material.dart';
import 'package:provider/provider.dart';

class _DrawerItemConfig {
  final String title;
  final IconData icon;
  final Widget page;
  final String? requiredPermission;

  const _DrawerItemConfig({
    required this.title,
    required this.icon,
    required this.page,
    this.requiredPermission,
  });
}

class AppDrawer extends StatelessWidget {
  final ValueChanged<Widget> onItemSelected;

  const AppDrawer({super.key, required this.onItemSelected});

  static const List<_DrawerItemConfig> _menuItems = [
    _DrawerItemConfig(title: 'Home', icon: Icons.home, page: HomeView()),
    _DrawerItemConfig(
      title: 'Acquire Data',
      icon: Icons.input,
      page: AcquireDataView(),
      requiredPermission: GlobalState.permAcquireData,
    ),
    _DrawerItemConfig(
      title: 'Offline Processing',
      icon: Icons.sync,
      page: OfflineProcessingView(),
      requiredPermission: GlobalState.permOfflineProcessing,
    ),
    _DrawerItemConfig(
      title: 'File Based Processing',
      icon: Icons.insert_drive_file,
      page: FileBasedProcessingView(),
      requiredPermission: GlobalState.permFileProcessing,
    ),
    _DrawerItemConfig(
      title: 'Result Profiles',
      icon: Icons.bar_chart,
      page: ResultProfilesView(),
      requiredPermission: GlobalState.permResultProfiles,
    ),
    _DrawerItemConfig(
      title: 'Trend Analysis',
      icon: Icons.trending_up,
      page: TrendAnalysisView(),
      requiredPermission: GlobalState.permTrendAnalysis,
    ),
    _DrawerItemConfig(
      title: 'BER Logging',
      icon: Icons.analytics,
      page: BerLoggingView(),
      requiredPermission: GlobalState.permBerLogging,
    ),
    _DrawerItemConfig(
      title: 'Data Transfer',
      icon: Icons.swap_horiz,
      page: DataTransferView(),
      requiredPermission: GlobalState.permDataTransfer,
    ),
    _DrawerItemConfig(
      title: 'Database Options',
      icon: Icons.storage,
      page: DatabaseOptionsView(),
      requiredPermission: GlobalState.permDatabaseOptions,
    ),
    _DrawerItemConfig(
      title: 'Developer Options',
      icon: Icons.developer_mode,
      page: DeveloperOptionsView(),
      requiredPermission: GlobalState.permDeveloperOptions,
    ),
  ];

  @override
  Widget build(BuildContext context) {
    final globalState = context.watch<GlobalState>();

    // Filter items based on permissions
    final visibleItems = _menuItems.where((item) {
      if (item.requiredPermission == null) return true;
      return globalState.hasPermission(item.requiredPermission!);
    }).toList();

    // Group items: Primary nav vs Secondary/Admin nav (separated by divider)
    // For simplicity efficiently, we'll just check if we hit the "Database Options"
    // or similar admin-like tools to insert a divider if desired, but for now
    // let's just render the list cleanly.

    // To maintain the divider before Database Options as in the original design:
    // We can check indices or just do a simple map.
    // Let's stick to the previous grouping logic dynamically.
    // If the item is "Database Options", we insert a Divider before it.

    List<Widget> drawerChildren = [
      DrawerHeader(
        decoration: BoxDecoration(color: Theme.of(context).colorScheme.primary),
        child: const Text(
          'CSRSP Menu',
          style: TextStyle(color: Colors.white, fontSize: 24),
        ),
      ),
    ];

    for (var item in visibleItems) {
      // Add divider before Database Options if it's visible
      if (item.title == 'Database Options') {
        drawerChildren.add(const Divider());
      }

      drawerChildren.add(
        _createDrawerItem(
          icon: item.icon,
          text: item.title,
          onTap: () => onItemSelected(item.page),
        ),
      );
    }

    return Drawer(
      child: ListView(padding: EdgeInsets.zero, children: drawerChildren),
    );
  }

  Widget _createDrawerItem({
    required IconData icon,
    required String text,
    required GestureTapCallback onTap,
  }) {
    return ListTile(
      title: Row(
        children: [
          Icon(icon),
          Padding(padding: const EdgeInsets.only(left: 8.0), child: Text(text)),
        ],
      ),
      onTap: onTap,
    );
  }
}
