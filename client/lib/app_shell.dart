import 'dart:async';
import 'dart:math';

import 'package:client/pages/home_page.dart';
import 'package:client/pages/login_view.dart';
import 'package:client/widgets/app_drawer.dart';
import 'package:client/widgets/status_bar.dart';
import 'package:flutter/material.dart';

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
  bool _isConnected = true;
  double _memoryUsage = 0.0;
  double _cpuUsage = 0.0;
  Timer? _timer;

  @override
  void initState() {
    super.initState();
    _simulateIpWhitelistCheck();
    _startStatusUpdates();
  }

  @override
  void dispose() {
    _timer?.cancel();
    super.dispose();
  }

  void _simulateIpWhitelistCheck() {
    // Simulate a network call to check if the IP is whitelisted
    Future.delayed(const Duration(seconds: 2), () {
      // 20% chance of being whitelisted for demonstration
      if (Random().nextDouble() < 0.2) {
        setState(() {
          _isLoggedIn = true;
          _isLoading = false;
        });
      } else {
        setState(() {
          _isLoading = false;
        });
      }
    });
  }

  void _startStatusUpdates() {
    // Start a timer to simulate live data updates for the status bar
    _timer = Timer.periodic(const Duration(seconds: 2), (timer) {
      setState(() {
        _memoryUsage = 128.0 + Random().nextDouble() * 50;
        _cpuUsage = 5.0 + Random().nextDouble() * 15;
        if (_isLoggedIn && Random().nextDouble() < 0.1) {
          _isConnected = false;
        }
      });
    });
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
              _isConnected ? 'Reconnected successfully!' : 'Reconnect failed.'),
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
