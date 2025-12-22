import 'package:client/app_shell.dart';
import 'package:client/theme/theme.dart';
import 'package:flutter/material.dart';

import 'package:client/state/global_state.dart';
import 'package:client/state/log_state.dart';
import 'package:provider/provider.dart';

void main() {
  runApp(
    MultiProvider(
      providers: [
        ChangeNotifierProvider(create: (_) => GlobalState()),
        ChangeNotifierProvider(create: (_) => LogState()),
      ],
      child: const MyApp(),
    ),
  );
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
