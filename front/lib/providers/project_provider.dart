import 'package:dio/dio.dart';
import 'package:flutter/material.dart';
import '../models/project_dto.dart';
import '../service/api_service.dart';

class ProjectProvider with ChangeNotifier {
  final ApiService _apiService = ApiService();

  List<ProjectDTO> _projects = [];
  bool _isLoading = false;

  List<ProjectDTO> get projects => _projects;
  bool get isLoading => _isLoading;

  Future<void> fetchProjects() async {
    _isLoading = true;
    notifyListeners();

    try {
      _projects = await _apiService.getAllProjects();
    } finally {
      _isLoading = false;
      notifyListeners();
    }
  }

  Future<void> uploadProject({
    required String name,
    required int programmingLanguageId,
    required MultipartFile file,
  }) async {
    await _apiService.uploadProject(
      name: name,
      userId: 1,
      programmingLanguageId: programmingLanguageId,
      file: file,
    );
    await fetchProjects();
  }
}
