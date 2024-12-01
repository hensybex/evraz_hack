class ProjectDTO {
  final int id;
  final int programmingLanguageId;
  final String name;
  final String description;
  final String path;
  final String tree;

  ProjectDTO({
    required this.id,
    required this.programmingLanguageId,
    required this.name,
    required this.description,
    required this.path,
    required this.tree,
  });

  factory ProjectDTO.fromJson(Map<String, dynamic> json) {
    return ProjectDTO(
      id: json['id'],
      programmingLanguageId: json['programming_language_id'],
      name: json['name'],
      description: json['description'],
      path: json['path'],
      tree: json['tree'],
    );
  }
}
