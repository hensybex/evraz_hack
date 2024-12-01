import 'package:flutter/material.dart';
import '../models/project_overview_response.dart';
import '../service/api_service.dart';

class ProjectOverviewProvider with ChangeNotifier {
  final ApiService _apiService = ApiService();

  ProjectOverviewResponse? _overview;
  bool _isLoading = false;

  ProjectOverviewResponse? get overview => _overview;
  bool get isLoading => _isLoading;

  Future<void> fetchProjectOverview(int projectId) async {
    _isLoading = true;
    notifyListeners();

    try {
      _overview = await _apiService.getProjectOverview(projectId);
    } finally {
      _isLoading = false;
      notifyListeners();
    }
  }
}
