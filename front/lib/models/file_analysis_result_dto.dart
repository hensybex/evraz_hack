class FileAnalysisResultDTO {
  final int id;
  final String promptName;
  final String compliance;
  final String issues;
  final String recommendations;

  FileAnalysisResultDTO({
    required this.id,
    required this.promptName,
    required this.compliance,
    required this.issues,
    required this.recommendations,
  });

  factory FileAnalysisResultDTO.fromJson(Map<String, dynamic> json) {
    return FileAnalysisResultDTO(
      id: json['id'],
      promptName: json['prompt_name'],
      compliance: json['compliance'],
      issues: json['issues'],
      recommendations: json['recommendations'],
    );
  }
}
