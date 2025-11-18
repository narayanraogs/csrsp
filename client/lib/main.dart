import 'package:client/app_shell.dart';
import 'package:client/theme/theme.dart';
import 'package:flutter/material.dart';

void main() {
  runApp(const MyApp());
}

class MyApp extends StatelessWidget {
  const MyApp({super.key});

  @override
  Widget build(BuildContext context) {
    return MaterialApp(
      title: 'CSRSP',
      theme: appTheme,
      debugShowCheckedModeBanner: false,
      home: const AppShell(),
    );
  }
}
