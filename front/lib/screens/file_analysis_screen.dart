import 'package:flutter/material.dart';
import 'package:provider/provider.dart';
import '../providers/file_analysis_provider.dart';
import '../models/file_analysis_result_dto.dart';

class FileAnalysisScreen extends StatefulWidget {
  final int projectId;
  final int fileId;

  FileAnalysisScreen({required this.projectId, required this.fileId});

  @override
  _FileAnalysisScreenState createState() => _FileAnalysisScreenState();
}

class _FileAnalysisScreenState extends State<FileAnalysisScreen> {
  @override
  void initState() {
    super.initState();
    WidgetsBinding.instance.addPostFrameCallback((_) {
      Provider.of<FileAnalysisProvider>(context, listen: false)
          .fetchFileAnalysis(widget.fileId);
    });
  }

  @override
  Widget build(BuildContext context) {
    return Consumer<FileAnalysisProvider>(
      builder: (context, provider, child) {
        if (provider.isLoading) {
          return Scaffold(
            appBar: AppBar(title: Text('File Analysis')),
            body: Center(child: CircularProgressIndicator()),
          );
        }

        final analysisResults = provider.analysisResults;

        return Scaffold(
          appBar: AppBar(title: Text('File Analysis')),
          body: ListView(
            children: _buildAnalysisContainers(analysisResults),
          ),
        );
      },
    );
  }

  List<Widget> _buildAnalysisContainers(
      List<FileAnalysisResultDTO> analysisResults) {
    final promptNames = [
      'ApplicationLayerCode',
      'CodingStandards',
      'DateTimeHandlingFile',
      'ErrorHandlingAndLogging',
      'AdditionalTechnicalFile',
    ];

    final Map<String, FileAnalysisResultDTO> groupedResults = {};
    for (var promptName in promptNames) {
      final result = analysisResults.firstWhere(
        (r) => r.promptName == promptName,
        orElse: () => FileAnalysisResultDTO(
          id: 0,
          promptName: promptName,
          compliance: 'N/A',
          issues: '',
          recommendations: '',
        ),
      );
      groupedResults[promptName] = result;
    }

    return promptNames.map((promptName) {
      final result = groupedResults[promptName]!;
      return Card(
        margin: const EdgeInsets.all(8.0),
        child: ExpansionTile(
          title: Text(promptName),
          children: [
            ListTile(
              title: Text('Compliance: ${result.compliance}'),
            ),
            Padding(
              padding: const EdgeInsets.symmetric(horizontal: 16.0),
              child: Column(
                crossAxisAlignment: CrossAxisAlignment.start,
                children: [
                  const Text(
                    'Issues:',
                    style: TextStyle(fontWeight: FontWeight.bold),
                  ),
                  ..._parseStringToList(result.issues).map((issue) {
                    return Text('- $issue');
                  }),
                  const SizedBox(height: 8),
                  const Text(
                    'Recommendations:',
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
    }).toList();
  }

  List<String> _parseStringToList(String data) {
    return data.split(', ').map((item) => item.trim()).toList();
  }
}
