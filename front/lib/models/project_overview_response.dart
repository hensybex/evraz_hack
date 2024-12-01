import 'project_dto.dart';
import 'project_file_dto.dart';
import 'project_analysis_result_dto.dart';

class ProjectOverviewResponse {
  final ProjectDTO project;
  final List<ProjectFileDTO> files;
  final List<ProjectAnalysisResultDTO> analysisResults;

  ProjectOverviewResponse({
    required this.project,
    required this.files,
    required this.analysisResults,
  });

  factory ProjectOverviewResponse.fromJson(Map<String, dynamic> json) {
    return ProjectOverviewResponse(
      project: ProjectDTO.fromJson(json['project']),
      files: (json['files'] as List)
          .map((e) => ProjectFileDTO.fromJson(e))
          .toList(),
      analysisResults: (json['analysis_results'] as List)
          .map((e) => ProjectAnalysisResultDTO.fromJson(e))
          .toList(),
    );
  }
}
