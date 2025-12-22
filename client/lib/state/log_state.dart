import 'package:flutter/material.dart';

enum LogType { info, success, warning, error }

class LogEntry {
  final DateTime timestamp;
  final String message;
  final LogType type;

  LogEntry({
    required this.timestamp,
    required this.message,
    required this.type,
  });
}

class LogState extends ChangeNotifier {
  final List<LogEntry> _logs = [];
  int _unreadCount = 0;

  List<LogEntry> get logs => List.unmodifiable(_logs);
  int get unreadCount => _unreadCount;

  void addLog(
    String message, {
    LogType type = LogType.info,
    bool showSnackbar = false,
    BuildContext? context,
  }) {
    // Add to internal list
    _logs.insert(
      0,
      LogEntry(timestamp: DateTime.now(), message: message, type: type),
    );

    // Limit log size to prevent memory issues
    if (_logs.length > 500) {
      _logs.removeLast();
    }

    _unreadCount++;
    notifyListeners();

    // Optionally show Snackbar
    if (showSnackbar && context != null) {
      ScaffoldMessenger.of(context).showSnackBar(
        SnackBar(
          content: Text(message),
          backgroundColor: _getSnackBarColor(type),
          behavior: SnackBarBehavior.floating,
          duration: const Duration(seconds: 3),
        ),
      );
    }
  }

  void clearLogs() {
    _logs.clear();
    _unreadCount = 0;
    notifyListeners();
  }

  void markAllRead() {
    _unreadCount = 0;
    notifyListeners();
  }

  Color _getSnackBarColor(LogType type) {
    switch (type) {
      case LogType.success:
        return Colors.green;
      case LogType.warning:
        return Colors.orange;
      case LogType.error:
        return Colors.red;
      case LogType.info:
        return Colors.blueGrey;
    }
  }
}
