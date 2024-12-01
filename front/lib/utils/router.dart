import 'package:flutter/material.dart';
import 'package:go_router/go_router.dart';
import '../screens/file_analysis_screen.dart';
import '../screens/projects_screen.dart';
import '../screens/project_overview_screen.dart';

final GoRouter router = GoRouter(
  routes: [
    GoRoute(
      path: '/',
      name: 'projects',
      builder: (BuildContext context, GoRouterState state) => ProjectsScreen(),
    ),
    GoRoute(
      path: '/project/:id',
      name: 'projectOverview',
      builder: (BuildContext context, GoRouterState state) {
        final id = int.parse(state.pathParameters['id']!);
        return ProjectOverviewScreen(projectId: id);
      },
    ),
    GoRoute(
      path: '/project/:projectId/file/:fileId',
      name: 'fileAnalysis',
      builder: (BuildContext context, GoRouterState state) {
        final projectId = int.parse(state.pathParameters['projectId']!);
        final fileId = int.parse(state.pathParameters['fileId']!);
        return FileAnalysisScreen(projectId: projectId, fileId: fileId);
      },
    ),
  ],
);
