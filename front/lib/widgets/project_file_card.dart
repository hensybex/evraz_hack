import 'package:flutter/material.dart';
import '../models/project_file_dto.dart';
import 'package:go_router/go_router.dart';

class ProjectFileCard extends StatelessWidget {
  final ProjectFileDTO file;
  final int projectId;

  ProjectFileCard({required this.file, required this.projectId});

  @override
  Widget build(BuildContext context) {
    return Card(
      color: file.wasAnalyzed
          ? Colors.green[400]
          : Colors.red[400], // Change color based on `wasAnalyzed`
      child: ListTile(
        title: Text(file.name),
        subtitle:
            Text('Успешно проанализирован: ${file.wasAnalyzed ? 'Да' : 'Нет'}'),
        onTap: () {
          context.go('/project/$projectId/file/${file.id}');
        },
      ),
    );
  }
}
