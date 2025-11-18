import 'package:client/pages/acquire_data_page.dart';
import 'package:client/pages/ber_logging_page.dart';
import 'package:client/pages/data_transfer_page.dart';
import 'package:client/pages/database_options_page.dart';
import 'package:client/pages/developer_options_page.dart';
import 'package:client/pages/file_based_processing_page.dart';
import 'package:client/pages/home_page.dart';
import 'package:client/pages/offline_processing_page.dart';
import 'package:client/pages/result_profiles_page.dart';
import 'package:flutter/material.dart';

class AppDrawer extends StatelessWidget {
  final ValueChanged<Widget> onItemSelected;

  const AppDrawer({super.key, required this.onItemSelected});

  @override
  Widget build(BuildContext context) {
    return Drawer(
      child: ListView(
        padding: EdgeInsets.zero,
        children: [
          DrawerHeader(
            decoration: BoxDecoration(
              color: Theme.of(context).colorScheme.primary,
            ),
            child: const Text(
              'CSRSP Menu',
              style: TextStyle(color: Colors.white, fontSize: 24),
            ),
          ),
          _createDrawerItem(
            icon: Icons.home,
            text: 'Home',
            onTap: () => onItemSelected(const HomeView()),
          ),
          _createDrawerItem(
            icon: Icons.input,
            text: 'Acquire Data',
            onTap: () => onItemSelected(const AcquireDataView()),
          ),
          _createDrawerItem(
            icon: Icons.sync,
            text: 'Offline Processing',
            onTap: () => onItemSelected(const OfflineProcessingView()),
          ),
          _createDrawerItem(
            icon: Icons.insert_drive_file,
            text: 'File Based Processing',
            onTap: () => onItemSelected(const FileBasedProcessingView()),
          ),
          _createDrawerItem(
            icon: Icons.bar_chart,
            text: 'Result Profiles',
            onTap: () => onItemSelected(const ResultProfilesView()),
          ),
          _createDrawerItem(
            icon: Icons.analytics,
            text: 'BER Logging',
            onTap: () => onItemSelected(const BerLoggingView()),
          ),
          _createDrawerItem(
            icon: Icons.swap_horiz,
            text: 'Data Transfer',
            onTap: () => onItemSelected(const DataTransferView()),
          ),
          const Divider(),
          _createDrawerItem(
            icon: Icons.storage,
            text: 'Database Options',
            onTap: () => onItemSelected(const DatabaseOptionsView()),
          ),
          _createDrawerItem(
            icon: Icons.developer_mode,
            text: 'Developer Options',
            onTap: () => onItemSelected(const DeveloperOptionsView()),
          ),
        ],
      ),
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
