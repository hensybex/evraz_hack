import 'package:flutter/material.dart';
import '../models/file_analysis_result_dto.dart';
import '../service/api_service.dart';

class FileAnalysisProvider with ChangeNotifier {
  final ApiService _apiService = ApiService();

  List<FileAnalysisResultDTO> _analysisResults = [];
  bool _isLoading = false;

  List<FileAnalysisResultDTO> get analysisResults => _analysisResults;
  bool get isLoading => _isLoading;

  Future<void> fetchFileAnalysis(int fileId) async {
    _isLoading = true;
    notifyListeners();

    try {
      _analysisResults = await _apiService.getFileAnalysisResults(fileId);
      print(_analysisResults.length);
    } finally {
      _isLoading = false;
      notifyListeners();
    }
  }
}
