import 'package:flutter/material.dart';

/// A permanent status bar for the bottom of the screen.
class StatusBar extends StatelessWidget {
  final bool isConnected;
  final double memoryUsage;
  final double cpuUsage;
  final VoidCallback onReconnect;

  const StatusBar({
    super.key,
    required this.isConnected,
    required this.memoryUsage,
    required this.cpuUsage,
    required this.onReconnect,
  });

  @override
  Widget build(BuildContext context) {
    final theme = Theme.of(context);
    final statusBarColor = theme.colorScheme.primary.withOpacity(0.1);
    final textColor = theme.textTheme.bodySmall?.color?.withOpacity(0.8);

    return Container(
      color: statusBarColor,
      padding: const EdgeInsets.symmetric(horizontal: 16.0, vertical: 4.0),
      child: Row(
        mainAxisAlignment: MainAxisAlignment.spaceBetween,
        children: [
          _buildStatusIndicator(textColor),
          _buildUsageIndicator(
              'MEM', memoryUsage, 'MB', textColor),
          _buildUsageIndicator(
              'CPU', cpuUsage, '%', textColor),
        ],
      ),
    );
  }

  Widget _buildStatusIndicator(Color? textColor) {
    return Row(
      children: [
        Text(
          'Server:',
          style: TextStyle(color: textColor, fontSize: 12),
        ),
        const SizedBox(width: 8),
        Icon(
          isConnected ? Icons.check_circle : Icons.error,
          color: isConnected ? Colors.green : Colors.red,
          size: 16,
        ),
        const SizedBox(width: 4),
        Text(
          isConnected ? 'Connected' : 'Offline',
          style: TextStyle(
              color: textColor, fontSize: 12, fontWeight: FontWeight.bold),
        ),
        const SizedBox(width: 8),
        SizedBox(
          height: 24,
          child: TextButton(
            onPressed: onReconnect,
            style: TextButton.styleFrom(
              padding: const EdgeInsets.symmetric(horizontal: 8),
              foregroundColor: textColor,
            ),
            child: const Row(
              children: [
                Icon(Icons.sync, size: 16),
                SizedBox(width: 4),
                Text('Reconnect', style: TextStyle(fontSize: 12)),
              ],
            ),
          ),
        ),
      ],
    );
  }

  Widget _buildUsageIndicator(
      String label, double value, String unit, Color? textColor) {
    return Row(
      mainAxisSize: MainAxisSize.min,
      children: [
        Text(
          '$label:',
          style: TextStyle(color: textColor, fontSize: 12),
        ),
        const SizedBox(width: 4),
        Text(
          '${value.toStringAsFixed(2)} $unit',
          style: TextStyle(
              color: textColor, fontSize: 12, fontWeight: FontWeight.bold),
        ),
      ],
    );
  }
}
