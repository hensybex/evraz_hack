import 'package:flutter/material.dart';
import 'package:front/providers/file_analysis_provider.dart';
import 'package:provider/provider.dart';
import 'providers/project_overview_provider.dart';
import 'providers/project_provider.dart';
import 'utils/router.dart';

void main() {
  runApp(
    MultiProvider(
      providers: [
        ChangeNotifierProvider(create: (_) => ProjectProvider()),
        ChangeNotifierProvider(create: (_) => ProjectOverviewProvider()),
        ChangeNotifierProvider(create: (_) => FileAnalysisProvider()),
      ],
      child: MyApp(),
    ),
  );
}

final ThemeData lightTheme = ThemeData(
  brightness: Brightness.light,
  primarySwatch: Colors.blue,
  scaffoldBackgroundColor: Colors.white,
  appBarTheme: AppBarTheme(
    backgroundColor: Colors.blue,
    foregroundColor: Colors.white,
  ),
);

final ThemeData darkTheme = ThemeData(
  brightness: Brightness.dark,
  primarySwatch: Colors.deepPurple,
  scaffoldBackgroundColor: Colors.black,
  appBarTheme: AppBarTheme(
    backgroundColor: Colors.deepPurple,
    foregroundColor: Colors.white,
  ),
);

class MyApp extends StatelessWidget {
  const MyApp({super.key});

  @override
  Widget build(BuildContext context) {
    return MaterialApp.router(
      routerConfig: router,
      title: 'Programming Projects',
      theme: lightTheme, // Light theme configuration
      darkTheme: darkTheme, // Dark theme configuration
    );
  }
}
