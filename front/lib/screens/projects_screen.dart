import 'package:flutter/material.dart';
import 'package:provider/provider.dart';
import '../providers/project_provider.dart';
import '../widgets/project_card.dart';
import 'package:file_picker/file_picker.dart';
import 'package:dio/dio.dart';

class ProjectsScreen extends StatefulWidget {
  const ProjectsScreen({super.key});

  @override
  _ProjectsScreenState createState() => _ProjectsScreenState();
}

class _ProjectsScreenState extends State<ProjectsScreen> {
  @override
  void initState() {
    super.initState();
    WidgetsBinding.instance.addPostFrameCallback((_) {
      Provider.of<ProjectProvider>(context, listen: false).fetchProjects();
    });
  }

  void _showCreateProjectDialog() {
    final nameController = TextEditingController();
    int? selectedLanguageId;
    PlatformFile? selectedFile;

    showDialog(
      context: context,
      builder: (context) {
        return AlertDialog(
          title: const Text("Загрузить новый проект"),
          content: SingleChildScrollView(
            child: Column(
              children: [
                TextField(
                  controller: nameController,
                  decoration: const InputDecoration(labelText: "Имя проекта"),
                ),
                DropdownButtonFormField<int>(
                  decoration:
                      const InputDecoration(labelText: "Язык программирования"),
                  items: const [
                    DropdownMenuItem(value: 1, child: Text('Python')),
                    DropdownMenuItem(value: 2, child: Text('C#')),
                    DropdownMenuItem(value: 3, child: Text('TypeScript')),
                  ],
                  onChanged: (value) {
                    selectedLanguageId = value;
                  },
                ),
                const SizedBox(height: 20),
                ElevatedButton(
                  onPressed: () async {
                    FilePickerResult? result =
                        await FilePicker.platform.pickFiles();

                    if (result != null) {
                      selectedFile = result.files.first;
                    }
                  },
                  child: const Text("Загрузить архив"),
                ),
              ],
            ),
          ),
          actions: [
            TextButton(
              child: const Text("Загрузить проект"),
              onPressed: () async {
                if (nameController.text.isNotEmpty &&
                    selectedLanguageId != null &&
                    selectedFile != null) {
                  final bytes = selectedFile!.bytes!;
                  final filename = selectedFile!.name;

                  MultipartFile file = MultipartFile.fromBytes(
                    bytes,
                    filename: filename,
                  );

                  await Provider.of<ProjectProvider>(context, listen: false)
                      .uploadProject(
                    name: nameController.text,
                    programmingLanguageId: selectedLanguageId!,
                    file: file,
                  );

                  Navigator.of(context).pop();
                }
              },
            ),
          ],
        );
      },
    );
  }

  @override
  Widget build(BuildContext context) {
    final projectProvider = Provider.of<ProjectProvider>(context);

    return Scaffold(
      appBar: AppBar(title: const Text('Проекты')),
      body: projectProvider.isLoading
          ? const Center(child: CircularProgressIndicator())
          : ListView.builder(
              itemCount: projectProvider.projects.length,
              itemBuilder: (context, index) {
                final project = projectProvider.projects[index];
                return ProjectCard(project: project);
              },
            ),
      floatingActionButton: FloatingActionButton(
        onPressed: _showCreateProjectDialog,
        child: const Icon(Icons.add),
      ),
    );
  }
}
