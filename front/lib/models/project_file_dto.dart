class ProjectFileDTO {
  final int id;
  final String name;
  final bool wasAnalyzed;

  ProjectFileDTO({
    required this.id,
    required this.name,
    required this.wasAnalyzed,
  });

  factory ProjectFileDTO.fromJson(Map<String, dynamic> json) {
    return ProjectFileDTO(
      id: json['id'],
      name: json['name'],
      wasAnalyzed: json['was_analyzed'],
    );
  }
}
