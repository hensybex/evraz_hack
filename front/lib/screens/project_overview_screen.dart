import 'package:flutter/material.dart';
import 'package:provider/provider.dart';
import 'package:url_launcher/url_launcher.dart';
import '../providers/project_overview_provider.dart';
import '../models/project_dto.dart';
import '../models/project_file_dto.dart';
import '../models/project_analysis_result_dto.dart';
import '../service/api_service.dart';
import '../widgets/project_file_card.dart';

class ProjectOverviewScreen extends StatefulWidget {
  final int projectId;

  const ProjectOverviewScreen({super.key, required this.projectId});

  @override
  ProjectOverviewScreenState createState() => ProjectOverviewScreenState();
}

class ProjectOverviewScreenState extends State<ProjectOverviewScreen>
    with SingleTickerProviderStateMixin {
  late TabController _tabController;
  final List<String> _tabs = [
    'Общее',
    'Файлы',
    'Анализ проекта',
  ];
  bool _isAnalyzeButtonVisible = true;

  @override
  void initState() {
    super.initState();
    _tabController = TabController(length: _tabs.length, vsync: this);
    WidgetsBinding.instance.addPostFrameCallback((_) {
      final provider =
          Provider.of<ProjectOverviewProvider>(context, listen: false);
      provider.fetchProjectOverview(widget.projectId);
      setState(() {
        _isAnalyzeButtonVisible = !provider.overview!.project.wasAnalyzed;
      });
    });
  }

  Widget _buildOverviewTab(ProjectDTO project) {
    return Padding(
      padding: const EdgeInsets.all(16.0),
      child: ListView(
        children: [
          Text('ID: ${project.id}', style: const TextStyle(fontSize: 16)),
          const SizedBox(height: 8),
          Text(
            'Язык программирования: ${getProgrammingLanguageName(project.programmingLanguageId)}',
            style: const TextStyle(fontSize: 16),
          ),
          const SizedBox(height: 8),
          Text('Имя: ${project.name}', style: const TextStyle(fontSize: 16)),
          const SizedBox(height: 8),
          Text('Путь: ${project.path}', style: const TextStyle(fontSize: 16)),
          const SizedBox(height: 8),
          Text('Структура проекта: ${project.tree}',
              style: const TextStyle(fontSize: 16)),
        ],
      ),
    );
  }

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

  Widget _buildFilesTab(List<ProjectFileDTO> files) {
    return ListView.builder(
      itemCount: files.length,
      itemBuilder: (context, index) {
        final file = files[index];
        return ProjectFileCard(file: file, projectId: widget.projectId);
      },
    );
  }

  Widget _buildAnalysisResultsTab(List<ProjectAnalysisResultDTO> results) {
    // Define the prompt names
    final promptNames = [
      'ProjectStructure',
      'KeyFiles',
      'ApplicationArchitecture',
      'DependencyManagement',
      'ProjectSettings',
      'TestingStrategy',
      'AdditionalTechnical',
      'DateTimeHandling',
    ];

// Map for Russian display names
    final Map<String, String> promptDisplayNames = {
      'ProjectStructure': 'Структура проекта',
      'KeyFiles': 'Ключевые файлы',
      'ApplicationArchitecture': 'Архитектура приложения',
      'DependencyManagement': 'Управление зависимостями',
      'ProjectSettings': 'Настройки проекта',
      'TestingStrategy': 'Стратегия тестирования',
      'AdditionalTechnical': 'Дополнительные технические аспекты',
      'DateTimeHandling': 'Работа с датами и временем',
    };

// Group the results by prompt name
    final Map<String, ProjectAnalysisResultDTO> groupedResults = {};
    for (var promptName in promptNames) {
      final result = results.firstWhere(
        (r) => r.promptName == promptName,
        orElse: () => ProjectAnalysisResultDTO(
          id: 0,
          promptName: promptName,
          compliance: 'N/A',
          issues: '',
          recommendations: '',
        ),
      );
      groupedResults[promptName] = result;
    }

    return ListView(
      children: promptNames.map((promptName) {
        final result = groupedResults[promptName]!;
        final correctPrompt = promptDisplayNames[promptName] ??
            promptName; // Use Russian display name or fallback
        return Card(
          margin: const EdgeInsets.all(8.0),
          child: ExpansionTile(
            title: Text(correctPrompt),
            children: [
              ListTile(
                title: Text('Соответствует требованиям: ${result.compliance}'),
              ),
              Padding(
                padding: const EdgeInsets.symmetric(horizontal: 16.0),
                child: Column(
                  crossAxisAlignment: CrossAxisAlignment.start,
                  children: [
                    const Text(
                      'Проблемы:',
                      style: TextStyle(fontWeight: FontWeight.bold),
                    ),
                    ..._parseStringToList(result.issues).map((issue) {
                      return Text('- $issue');
                    }),
                    const SizedBox(height: 8),
                    const Text(
                      'Рекомендации:',
                      style: TextStyle(fontWeight: FontWeight.bold),
                    ),
                    ..._parseStringToList(result.recommendations).map((rec) {
                      return Text('- $rec');
                    }),
                  ],
                ),
              ),
            ],
          ),
        );
      }).toList(),
    );
  }

  List<String> _parseStringToList(String data) {
    return data.split(', ').map((item) => item.trim()).toList();
  }

  @override
  Widget build(BuildContext context) {
    return Consumer<ProjectOverviewProvider>(
      builder: (context, provider, child) {
        if (provider.isLoading || provider.overview == null) {
          return Scaffold(
            appBar: AppBar(
              title: const Text('Обзор проекта'),
              bottom: TabBar(
                controller: _tabController,
                tabs: _tabs.map((tab) => Tab(text: tab)).toList(),
              ),
            ),
            body: const Center(child: CircularProgressIndicator()),
          );
        }

        final project = provider.overview!.project;
        final files = provider.overview!.files;
        final analysisResults = provider.overview!.analysisResults;

        return Scaffold(
          appBar: AppBar(
            title: const Text('Обзор проекта'),
            bottom: TabBar(
              controller: _tabController,
              tabs: _tabs.map((tab) => Tab(text: tab)).toList(),
            ),
            actions: [
              if (_isAnalyzeButtonVisible)
                IconButton(
                  icon: const Icon(
                    Icons.rocket_launch_sharp,
                    color: Colors.white,
                  ),
                  onPressed: () async {
                    try {
                      ScaffoldMessenger.of(context).showSnackBar(
                        const SnackBar(
                            content: Text('Проект отправлен на анализ')),
                      );

                      setState(() {
                        _isAnalyzeButtonVisible = false;
                      });
                      final ApiService apiService = ApiService();
                      await apiService.analyzeProject(
                          projectId: widget.projectId);
                    } catch (e) {
                      ScaffoldMessenger.of(context).showSnackBar(
                        SnackBar(content: Text('Error: ${e.toString()}')),
                      );
                    }
                  },
                ),
              IconButton(
                icon: const Icon(
                  Icons.download,
                  color: Colors.white,
                ),
                onPressed: () async {
                  final url =
                      'http://localhost:8080/api/projects/${widget.projectId}/generate_pdf';
                  if (await canLaunch(url)) {
                    await launch(url);
                  } else {
                    // Handle error
                    ScaffoldMessenger.of(context).showSnackBar(
                      const SnackBar(
                          content: Text('Не получилось загрузить пдф')),
                    );
                  }
                },
              ),
            ],
          ),
          body: TabBarView(
            controller: _tabController,
            children: [
              _buildOverviewTab(project),
              _buildFilesTab(files),
              _buildAnalysisResultsTab(analysisResults),
            ],
          ),
        );
      },
    );
  }
}
