import json
import sys
from datetime import datetime
from collections import defaultdict

def parse_test_results(json_file):
    """Parse Go test JSON output and extract test results."""
    tests = defaultdict(lambda: {
        'name': '',
        'package': '',
        'status': 'unknown',
        'elapsed': 0,
        'output': []
    })

    package_results = {}

    with open(json_file, 'r') as f:
        for line in f:
            if not line.strip():
                continue

            try:
                data = json.loads(line)
                action = data.get('Action', '')
                package = data.get('Package', '')
                test_name = data.get('Test', '')
                elapsed = data.get('Elapsed', 0)
                output = data.get('Output', '')

                # Handle package-level results
                if not test_name and package:
                    if action == 'pass':
                        package_results[package] = {
                            'status': 'PASS',
                            'elapsed': elapsed,
                            'coverage': extract_coverage(output)
                        }
                    elif action == 'fail':
                        package_results[package] = {
                            'status': 'FAIL',
                            'elapsed': elapsed,
                            'coverage': None
                        }
                    elif action == 'skip':
                        package_results[package] = {
                            'status': 'SKIP',
                            'elapsed': elapsed,
                            'coverage': None
                        }

                # Handle test-level results
                if test_name and package:
                    test_key = f"{package}::{test_name}"

                    if action == 'run':
                        tests[test_key]['name'] = test_name
                        tests[test_key]['package'] = package
                        tests[test_key]['status'] = 'RUNNING'
                    elif action == 'pass':
                        tests[test_key]['status'] = 'PASS'
                        tests[test_key]['elapsed'] = elapsed
                    elif action == 'fail':
                        tests[test_key]['status'] = 'FAIL'
                        tests[test_key]['elapsed'] = elapsed
                    elif action == 'skip':
                        tests[test_key]['status'] = 'SKIP'
                        tests[test_key]['elapsed'] = elapsed
                    elif action == 'output' and output:
                        tests[test_key]['output'].append(output.strip())

            except json.JSONDecodeError:
                continue

    return dict(tests), package_results

def extract_coverage(output):
    """Extract coverage percentage from output string."""
    if not output:
        return None

    # Look for coverage pattern like "coverage: 76.9% of statements"
    import re
    match = re.search(r'coverage:\s+(\d+\.?\d*)%', output)
    if match:
        return float(match.group(1))
    return None

def generate_html_report(tests, package_results, output_file):
    """Generate HTML test report."""

    # Count test results
    total_tests = len(tests)
    passed_tests = sum(1 for t in tests.values() if t['status'] == 'PASS')
    failed_tests = sum(1 for t in tests.values() if t['status'] == 'FAIL')
    skipped_tests = sum(1 for t in tests.values() if t['status'] == 'SKIP')

    # Calculate total elapsed time
    total_time = sum(t['elapsed'] for t in tests.values())

    html_content = f"""
<!DOCTYPE html>
<html lang="en">
<head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <title>Test Report</title>
    <style>
        body {{ font-family: Arial, sans-serif; margin: 20px; }}
        .header {{ background: #f5f5f5; padding: 20px; border-radius: 5px; margin-bottom: 20px; }}
        .summary {{ display: flex; gap: 20px; margin-bottom: 20px; }}
        .summary-card {{ padding: 15px; border-radius: 5px; flex: 1; text-align: center; }}
        .pass {{ background: #d4edda; color: #155724; }}
        .fail {{ background: #f8d7da; color: #721c24; }}
        .skip {{ background: #fff3cd; color: #856404; }}
        .total {{ background: #e2e3e5; color: #383d41; }}
        table {{ width: 100%; border-collapse: collapse; margin-top: 20px; }}
        th, td {{ padding: 12px; text-align: left; border: 1px solid #ddd; }}
        th {{ background: #f8f9fa; }}
        .status-pass {{ color: #28a745; font-weight: bold; }}
        .status-fail {{ color: #dc3545; font-weight: bold; }}
        .status-skip {{ color: #ffc107; font-weight: bold; }}
        .test-output {{ cursor: pointer; }}
        .collapsible {{ 
            cursor: pointer; 
            padding: 8px 15px;
            background-color: #f1f1f1;
            border: none;
            text-align: left;
            outline: none;
            border-radius: 4px;
            margin-bottom: 5px;
            font-size: 14px;
        }}
        .active, .collapsible:hover {{
            background-color: #ddd;
        }}
        .content {{ 
            padding: 10px;
            display: none;
            overflow: hidden;
            background-color: #f9f9f9;
            border-radius: 4px;
            white-space: pre-wrap;
            font-family: monospace;
            font-size: 13px;
            max-height: 400px;
            overflow-y: auto;
            border: 1px solid #ddd;
        }}
        .output-preview {{
            color: #666;
            font-family: monospace;
            white-space: nowrap;
            overflow: hidden;
            text-overflow: ellipsis;
            max-width: 500px;
            display: inline-block;
            vertical-align: middle;
        }}
        .badge {{
            display: inline-block;
            padding: 3px 7px;
            font-size: 12px;
            font-weight: bold;
            border-radius: 10px;
            margin-left: 8px;
        }}
        .badge-pass {{ background-color: #d4edda; color: #155724; }}
        .badge-fail {{ background-color: #f8d7da; color: #721c24; }}
        .badge-skip {{ background-color: #fff3cd; color: #856404; }}
    </style>
</head>
<body>
    <div class="header">
        <h1>Go Test Report</h1>
        <p>Generated on: {datetime.now().strftime('%Y-%m-%d %H:%M:%S')}</p>
        <p>Total execution time: {total_time:.2f}s</p>
    </div>
    
    <div class="summary">
        <div class="summary-card total">
            <h3>{total_tests}</h3>
            <p>Total Tests</p>
        </div>
        <div class="summary-card pass">
            <h3>{passed_tests}</h3>
            <p>Passed</p>
        </div>
        <div class="summary-card fail">
            <h3>{failed_tests}</h3>
            <p>Failed</p>
        </div>
        <div class="summary-card skip">
            <h3>{skipped_tests}</h3>
            <p>Skipped</p>
        </div>
    </div>
    
    <h2>Test Details</h2>
    <table>
        <thead>
            <tr>
                <th>Test Name</th>
                <th>Package</th>
                <th>Status</th>
                <th>Duration (s)</th>
                <th>Output</th>
            </tr>
        </thead>
        <tbody>
"""

    # Add test rows
    for i, (test_key, test) in enumerate(sorted(tests.items())):
        if not test['name']:  # Skip empty test entries
            continue

        status_class = f"status-{test['status'].lower()}"
        badge_class = f"badge-{test['status'].lower()}"

        # Prepare preview text
        output_preview = test['output'][0][:100] if test['output'] else "No output"
        if len(output_preview) >= 100:
            output_preview += "..."

        # Determine if logs should be auto-expanded
        auto_expand = test['status'] in ['FAIL']
        display_style = 'block' if auto_expand else 'none'
        active_class = 'active' if auto_expand else ''

        # Count lines of output
        output_count = len(test['output'])

        html_content += f"""
            <tr>
                <td>{test['name']}</td>
                <td>{test['package']}</td>
                <td class="{status_class}">{test['status']}</td>
                <td>{test['elapsed']:.3f}</td>
                <td>
                    <button class="collapsible {active_class}" id="btn-{i}">
                        <span class="output-preview">{output_preview}</span>
                        <span class="badge {badge_class}">{output_count} lines</span>
                    </button>
                    <div class="content" id="content-{i}" style="display: {display_style};">
"""

        # Add each line of output
        for line in test['output']:
            html_content += f"{line}\n"

        html_content += """
                    </div>
                </td>
            </tr>
"""

    html_content += """
        </tbody>
    </table>
    
    <script>
        // Add click handlers for collapsible content
        document.querySelectorAll('.collapsible').forEach(button => {
            button.addEventListener('click', function() {
                this.classList.toggle("active");
                const content = this.nextElementSibling;
                content.style.display = content.style.display === 'block' ? 'none' : 'block';
            });
        });
    </script>
</body>
</html>
"""

    with open(output_file, 'w') as f:
        f.write(html_content)

def main():
    if len(sys.argv) != 3:
        print("Usage: python generate_test_report.py <json_file> <output_html>")
        sys.exit(1)

    json_file = sys.argv[1]
    output_file = sys.argv[2]

    try:
        tests, package_results = parse_test_results(json_file)
        generate_html_report(tests, package_results, output_file)
        print(f"HTML report generated: {output_file}")
    except Exception as e:
        print(f"Error generating report: {e}")
        sys.exit(1)

if __name__ == "__main__":
    main()