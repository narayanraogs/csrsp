import 'package:flutter/foundation.dart';
import 'package:flutter/material.dart';
import 'package:uuid/uuid.dart';

class GlobalState extends ChangeNotifier {
  // Unique Client ID for this session
  final String _clientId;
  late final String _serverUrl;

  GlobalState() : _clientId = const Uuid().v4() {
    _serverUrl = _determineServerUrl();
  }

  String _determineServerUrl() {
    // 1. Debug Mode (Development)
    if (kDebugMode) {
      return 'http://localhost:8080';
    }

    // 2. Production Web (Deployed)
    if (kIsWeb) {
      // Uri.base gives the current URL in the browser
      // We extract just the scheme (http), host (IP), and port.
      // Note: Uri.base.port might be 80 or 443 which are often omitted in strings,
      // but explicitly including them is safer for gRPC.
      // However, if port is 0 or null (unlikely in browser), we might need care.
      // Uri.base always returns a valid port (80/443 if implicit).
      return '${Uri.base.scheme}://${Uri.base.host}:${Uri.base.port}';
    }

    // 3. Production Linux/Other (Standalone App)
    // Default to localhost for now.
    return 'http://localhost:8080';
  }

  // Example state: A simple counter and a status message
  int _counter = 0;
  String _status = "Initial Status";

  // Getters to access state
  String get clientId => _clientId;
  String get serverUrl => _serverUrl;
  int get counter => _counter;
  String get status => _status;

  // Methods to modify state
  void incrementCounter() {
    _counter++;
    // notifyListeners() tells Flutter to rebuild any widgets listening to this state
    notifyListeners();
  }

  void updateStatus(String newStatus) {
    _status = newStatus;
    notifyListeners();
  }
}
