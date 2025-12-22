import 'package:client/state/log_state.dart';
import 'package:flutter/material.dart';
import 'package:provider/provider.dart';
import 'package:intl/intl.dart';

class LogsSheet extends StatelessWidget {
  final VoidCallback onClose;

  const LogsSheet({super.key, required this.onClose});

  @override
  Widget build(BuildContext context) {
    return Container(
      height: 300,
      decoration: BoxDecoration(
        color: Theme.of(context).cardColor,
        border: Border(top: BorderSide(color: Theme.of(context).dividerColor)),
        boxShadow: [
          BoxShadow(
            color: Colors.black.withValues(alpha: 0.1),
            blurRadius: 4,
            offset: const Offset(0, -2),
          ),
        ],
      ),
      child: Column(
        children: [
          Padding(
            padding: const EdgeInsets.symmetric(
              horizontal: 16.0,
              vertical: 8.0,
            ),
            child: Row(
              mainAxisAlignment: MainAxisAlignment.spaceBetween,
              children: [
                const Text(
                  'System Logs',
                  style: TextStyle(fontWeight: FontWeight.bold, fontSize: 16),
                ),
                Row(
                  children: [
                    IconButton(
                      icon: const Icon(Icons.delete_outline),
                      onPressed: () {
                        context.read<LogState>().clearLogs();
                      },
                      tooltip: 'Clear Logs',
                    ),
                    IconButton(
                      icon: const Icon(Icons.close),
                      onPressed: onClose,
                      tooltip: 'Close Logs',
                    ),
                  ],
                ),
              ],
            ),
          ),
          const Divider(height: 1),
          Expanded(
            child: Consumer<LogState>(
              builder: (context, logState, child) {
                if (logState.logs.isEmpty) {
                  return const Center(
                    child: Text(
                      'No logs yet.',
                      style: TextStyle(color: Colors.grey),
                    ),
                  );
                }
                return ListView.builder(
                  itemCount: logState.logs.length,
                  itemBuilder: (context, index) {
                    final log = logState.logs[index];
                    return ListTile(
                      dense: true,
                      leading: _getIconForType(log.type),
                      title: Text(
                        log.message,
                        style: const TextStyle(fontSize: 13),
                      ),
                      subtitle: Text(
                        DateFormat('yyyy-MM-dd HH:mm:ss').format(log.timestamp),
                        style: const TextStyle(fontSize: 11),
                      ),
                    );
                  },
                );
              },
            ),
          ),
        ],
      ),
    );
  }

  Icon _getIconForType(LogType type) {
    switch (type) {
      case LogType.success:
        return const Icon(Icons.check_circle, color: Colors.green, size: 16);
      case LogType.warning:
        return const Icon(Icons.warning, color: Colors.orange, size: 16);
      case LogType.error:
        return const Icon(Icons.error, color: Colors.red, size: 16);
      case LogType.info:
        return const Icon(Icons.info, color: Colors.blue, size: 16);
    }
  }
}
