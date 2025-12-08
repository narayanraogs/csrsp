import 'package:flutter/foundation.dart';
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

  // Permissions State
  List<String> _permissions = [];

  // Permission Constants
  static const String permAcquireData = 'ACQUIRE_DATA';
  static const String permOfflineProcessing = 'OFFLINE_PROCESSING';
  static const String permFileProcessing = 'FILE_PROCESSING';
  static const String permResultProfiles = 'RESULT_PROFILES';
  static const String permTrendAnalysis = 'TREND_ANALYSIS';
  static const String permBerLogging = 'BER_LOGGING';
  static const String permDataTransfer = 'DATA_TRANSFER';
  static const String permDatabaseOptions = 'DATABASE_OPTIONS';
  static const String permDeveloperOptions = 'DEVELOPER_OPTIONS';

  static const List<String> _allPermissions = [
    permAcquireData,
    permOfflineProcessing,
    permFileProcessing,
    permResultProfiles,
    permTrendAnalysis,
    permBerLogging,
    permDataTransfer,
    permDatabaseOptions,
    permDeveloperOptions,
  ];

  // Getters to access state
  String get clientId => _clientId;
  String get serverUrl => _serverUrl;
  List<String> get permissions => _permissions;

  // Permission Management
  void setPermissions(List<String> perms) {
    _permissions = perms;
    notifyListeners();
  }

  void grantAllPermissions() {
    _permissions = List.from(_allPermissions);
    notifyListeners();
  }

  void clearPermissions() {
    _permissions = [];
    notifyListeners();
  }

  bool hasPermission(String permission) {
    return _permissions.contains(permission);
  }
}
