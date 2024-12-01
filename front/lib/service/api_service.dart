import 'package:dio/dio.dart';
import '../models/file_analysis_result_dto.dart';
import '../models/project_dto.dart';
import '../models/project_overview_response.dart';

class ApiService {
  static const String baseUrl = String.fromEnvironment('API_BASE_URL',
      defaultValue: 'http://localhost:8080');
  final Dio _dio = Dio(BaseOptions(baseUrl: baseUrl));

  Future<List<ProjectDTO>> getAllProjects() async {
    try {
      Response response = await _dio.get('/api/projects/all');
      return (response.data['projects'] as List)
          .map((json) => ProjectDTO.fromJson(json))
          .toList();
    } catch (e) {
      rethrow;
    }
  }

  Future<ProjectOverviewResponse> getProjectOverview(int projectId) async {
    try {
      Response response = await _dio.get('/api/projects/$projectId/overview');
      return ProjectOverviewResponse.fromJson(response.data);
    } catch (e) {
      rethrow;
    }
  }

  Future<List<FileAnalysisResultDTO>> getFileAnalysisResults(int fileId) async {
    try {
      Response response = await _dio.get('/api/files/23/analysis_results');
      print(response.data);
      return (response.data['analysis_results'] as List)
          .map((e) => FileAnalysisResultDTO.fromJson(e))
          .toList();
    } catch (e) {
      rethrow;
    }
  }

  Future<void> uploadProject({
    required String name,
    required int userId,
    required int programmingLanguageId,
    required MultipartFile file,
  }) async {
    try {
      FormData formData = FormData.fromMap({
        'name': name,
        'user_id': userId.toString(),
        'programming_language_id': programmingLanguageId.toString(),
        'file': file,
      });

      await _dio.post('/api/projects/upload_project', data: formData);
    } catch (e) {
      rethrow;
    }
  }

  // Implement other API methods as needed
}
