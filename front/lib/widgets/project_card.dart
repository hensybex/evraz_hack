import 'package:flutter/material.dart';
import '../models/project_dto.dart';
import 'package:go_router/go_router.dart';

class ProjectCard extends StatelessWidget {
  final ProjectDTO project;

  const ProjectCard({super.key, required this.project});

  String getProgrammingLanguageName(int id) {
    switch (id) {
      case 1:
        return 'Python';
      case 2:
        return 'C#';
      case 3:
        return 'TypeScript';
      default:
        return 'Unknown';
    }
  }

  @override
  Widget build(BuildContext context) {
    final languageName =
        getProgrammingLanguageName(project.programmingLanguageId);
    return Card(
      child: ListTile(
        title: Text(project.name),
        subtitle: Text('Язык Программирования: $languageName'),
        onTap: () {
          context.go('/project/${project.id}');
        },
      ),
    );
  }
}
