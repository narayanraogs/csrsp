import 'dart:async';
import 'dart:math';

import 'package:client/pages/home_page.dart';
import 'package:client/pages/login_view.dart';
import 'package:client/widgets/app_drawer.dart';
import 'package:client/widgets/status_bar.dart';
import 'package:flutter/material.dart';

import 'package:grpc/grpc_web.dart';
import 'package:client/communication/communication.pbgrpc.dart';
import 'package:client/state/global_state.dart';
import 'package:provider/provider.dart';

class AppShell extends StatefulWidget {
  const AppShell({super.key});

  @override
  State<AppShell> createState() => _AppShellState();
}

class _AppShellState extends State<AppShell> {
  // App states
  bool _isLoading = true;
  bool _isLoggedIn = false;

  // State for the main content
  Widget _selectedContent = const HomeView();

  // State for App Bar
  String _satelliteName = "Sat-01";
  String _testPhase = "Integration Testing";

  // State for the status bar
  bool _isConnected = false;
  double _memoryUsage = 0.0;
  double _cpuUsage = 0.0;

  // gRPC
  GrpcWebClientChannel? _channel;
  CommunicationClient? _client;

  @override
  void initState() {
    super.initState();
    _checkIpWhitelist();
    _fetchServerDetails();
    _startStatusUpdates();
  }

  Future<void> _fetchServerDetails() async {
    // Initialize gRPC channel and client if not already done
    if (_channel == null) {
      final serverUrl = context.read<GlobalState>().serverUrl;
      _channel = GrpcWebClientChannel.xhr(Uri.parse(serverUrl));
      _client = CommunicationClient(_channel!);
    }

    try {
      final request = ClientID()..id = context.read<GlobalState>().clientId;
      final response = await _client!.getServerDetails(request);

      if (mounted) {
        setState(() {
          _satelliteName = response.satelliteName;
          _testPhase = response.testPhase;
        });
      }
    } catch (e) {
      debugPrint("Server Details Fetch Error: $e");
    }
  }

  @override
  void dispose() {
    _channel?.shutdown();
    super.dispose();
  }

  Future<void> _checkIpWhitelist() async {
    // Initialize gRPC channel and client if not already done
    if (_channel == null) {
      final serverUrl = context.read<GlobalState>().serverUrl;
      _channel = GrpcWebClientChannel.xhr(Uri.parse(serverUrl));
      _client = CommunicationClient(_channel!);
    }

    try {
      final request = ClientID()..id = context.read<GlobalState>().clientId;
      final response = await _client!.isWhitelisted(request);

      if (mounted) {
        setState(() {
          _isLoggedIn = response.whitelisted;
          _isLoading = false;
        });
      }
    } catch (e) {
      debugPrint("Whitelist Check Error: $e");
      if (mounted) {
        setState(() {
          _isLoading = false;
          // Optionally show an error message or keep _isLoggedIn as false
        });
      }
    }
  }

  void _startStatusUpdates() {
    // Initialize gRPC channel and client if not already done
    if (_channel == null) {
      final serverUrl = context.read<GlobalState>().serverUrl;
      _channel = GrpcWebClientChannel.xhr(Uri.parse(serverUrl));
      _client = CommunicationClient(_channel!);
    }

    final request = ClientID()..id = context.read<GlobalState>().clientId;

    try {
      final stream = _client!.getServerStatus(request);
      stream.listen(
        (status) {
          setState(() {
            _memoryUsage = status.memory;
            _cpuUsage = status.cpu;
            _isConnected = true;
            debugPrint("RPC Status: $status");
          });
        },
        onError: (error) {
          debugPrint("RPC Error: $error");
          setState(() {
            _isConnected = false;
          });
          // Retry logic could go here
        },
        onDone: () {
          debugPrint("RPC Stream Closed");
          setState(() {
            _isConnected = false;
          });
        },
      );
    } catch (e) {
      debugPrint("Caught error: $e");
      setState(() {
        _isConnected = false;
      });
    }
  }

  void _onItemSelected(Widget content) {
    setState(() {
      _selectedContent = content;
    });
    Navigator.of(context).pop(); // Close the drawer
  }

  void _onLogin(String username, String password) {
    // Simulate a login attempt
    ScaffoldMessenger.of(context).showSnackBar(
      const SnackBar(
        content: Text('Validating credentials...'),
        duration: Duration(seconds: 1),
      ),
    );
    Future.delayed(const Duration(seconds: 1), () {
      // Simple validation for demonstration
      if (username == 'admin' && password == 'password') {
        setState(() {
          _isLoggedIn = true;
        });
      } else {
        ScaffoldMessenger.of(context).showSnackBar(
          const SnackBar(
            content: Text('Invalid username or password.'),
            backgroundColor: Colors.red,
          ),
        );
      }
    });
  }

  void _reconnect() {
    // Simulate a reconnection attempt
    setState(() {
      _isConnected = false;
    });
    ScaffoldMessenger.of(context).showSnackBar(
      const SnackBar(
        content: Text('Attempting to reconnect...'),
        duration: Duration(seconds: 2),
      ),
    );
    Future.delayed(const Duration(seconds: 2), () {
      setState(() {
        _isConnected = Random().nextDouble() > 0.2;
      });
      ScaffoldMessenger.of(context).showSnackBar(
        SnackBar(
          content: Text(
            _isConnected ? 'Reconnected successfully!' : 'Reconnect failed.',
          ),
          backgroundColor: _isConnected ? Colors.green : Colors.red,
        ),
      );
    });
  }

  @override
  Widget build(BuildContext context) {
    return Scaffold(
      appBar: AppBar(
        title: Column(
          crossAxisAlignment: CrossAxisAlignment.start,
          children: [
            Text(
              'CSRSP: $_satelliteName',
              style: const TextStyle(fontSize: 18, fontWeight: FontWeight.bold),
            ),
            Text(
              _testPhase,
              style: const TextStyle(fontSize: 12, fontStyle: FontStyle.italic),
            ),
          ],
        ),
      ),
      drawer: _isLoggedIn ? AppDrawer(onItemSelected: _onItemSelected) : null,
      body: _buildBody(),
      bottomNavigationBar: StatusBar(
        isConnected: _isConnected,
        memoryUsage: _memoryUsage,
        cpuUsage: _cpuUsage,
        onReconnect: _reconnect,
      ),
    );
  }

  Widget _buildBody() {
    if (_isLoading) {
      return const Center(child: CircularProgressIndicator());
    }
    if (!_isLoggedIn) {
      return LoginView(onLogin: _onLogin);
    }
    return _selectedContent;
  }
}
