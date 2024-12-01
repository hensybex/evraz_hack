import os
import subprocess

def save_tree_structure(base_dir, output_file):
    """Saves the tree structure of the base directory to a file."""
    with open(output_file, 'w') as f:
        subprocess.run(["tree", base_dir], stdout=f, text=True)
    print(f"Tree structure saved to {output_file}")

def collect_files_by_directory(root_dir, extension, excluded_suffix=None):
    """Collects files grouped by their directories."""
    files_by_directory = {}
    for subdir, _, files in os.walk(root_dir):
        for file in files:
            if file.endswith(extension) and (not excluded_suffix or not file.endswith(excluded_suffix)):
                file_path = os.path.join(subdir, file)
                if subdir not in files_by_directory:
                    files_by_directory[subdir] = []
                files_by_directory[subdir].append(file_path)
    return files_by_directory

def output_files_content(files_by_directory, output_file):
    """Writes the full content of all files to the output file, grouped by directory."""
    with open(output_file, 'w') as f:
        for directory, files in files_by_directory.items():
            f.write(f"Directory: {directory}\n")
            f.write("=" * (len("Directory: ") + len(directory)) + "\n\n")
            for file in files:
                f.write(f"File: {file}\n")
                f.write("-" * (len("File: ") + len(file)) + "\n")
                with open(file, 'r') as file_content:
                    f.write(file_content.read())
                f.write("\n" + "-" * 40 + "\n\n")
    print(f"Files content written to {output_file}")

def process_and_save(targets, tree_output_file, content_output_file):
    """Main function to process targets and save outputs."""
    for target in targets:
        print(f"Processing {target['name']} files...")
        # Save tree structure for the target directory
        save_tree_structure(target['root_directory'], tree_output_file)
        # Collect files grouped by directory
        files_by_directory = collect_files_by_directory(
            target['root_directory'],
            target['file_extension'],
            target.get('excluded_suffix', None)
        )
        # Output files content grouped by directory
        output_files_content(files_by_directory, content_output_file)

if __name__ == "__main__":
    # List of directories and their respective file extensions to process
    targets = [
        {"name": "API", "root_directory": "internal", "file_extension": ".go", "excluded_suffix": ".g.dart"},
        #{"name": "APP", "root_directory": "app/lib", "file_extension": ".dart", "excluded_suffix": ".g.dart"}
    ]

    # Specify output files
    tree_output_file = "tree_structure.txt"
    content_output_file = "files_content.txt"

    # Process and save outputs
    process_and_save(targets, tree_output_file, content_output_file)
