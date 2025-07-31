#!/usr/bin/env python3
import re
import sys
import os
from datetime import datetime

def parse_log_file(log_file_path):
    """Parse the test log file and extract test information."""
    with open(log_file_path, 'r', encoding='utf-8') as f:
        content = f.read()

    tests = {}
    lines = content.split('\n')

    # First pass: find all tests and their results
    for line in lines:
        # Match test start
        run_match = re.match(r'^=== RUN\s+(\S+)', line)
        if run_match:
            test_name = run_match.group(1)
            if test_name not in tests:
                tests[test_name] = {
                    'name': test_name,
                    'status': 'UNKNOWN',
                    'duration': '0s',
                    'logs': []
                }
            continue

        # Match test result
        result_match = re.match(r'^--- (PASS|FAIL): (\S+) \(([\d\.]+s)\)', line)
        if result_match:
            status, test_name, duration = result_match.groups()
            if test_name in tests:
                tests[test_name]['status'] = status
                tests[test_name]['duration'] = duration
            continue

    # Second pass: collect all logs for each test
    current_test = None
    for line in lines:
        # Track current test context
        run_match = re.match(r'^=== RUN\s+(\S+)', line)
        if run_match:
            current_test = run_match.group(1)
            if current_test in tests:
                tests[current_test]['logs'].append(line)
            continue

        cont_match = re.match(r'^=== CONT\s+(\S+)', line)
        if cont_match:
            current_test = cont_match.group(1)
            if current_test in tests:
                tests[current_test]['logs'].append(line)
            continue

        pause_match = re.match(r'^=== PAUSE\s+(\S+)', line)
        if pause_match:
            test_name = pause_match.group(1)
            if test_name in tests:
                tests[test_name]['logs'].append(line)
            continue

        name_match = re.match(r'^=== NAME\s+(\S+)', line)
        if name_match:
            current_test = name_match.group(1)
            if current_test in tests:
                tests[current_test]['logs'].append(line)
            continue

        # Match test result
        result_match = re.match(r'^--- (PASS|FAIL): (\S+) \(([\d\.]+s)\)', line)
        if result_match:
            test_name = result_match.group(2)
            if test_name in tests:
                tests[test_name]['logs'].append(line)
            current_test = None
            continue

        # Collect logs that contain test names or are in current test context
        if current_test and current_test in tests:
            tests[current_test]['logs'].append(line)
        else:
            # Check if this line mentions any test name
            for test_name in tests.keys():
                if test_name in line:
                    tests[test_name]['logs'].append(line)
                    break

    return list(tests.values())

def generate_html_report(tests, output_file):
    """Generate HTML report from test data."""
    total_tests = len(tests)
    failed_tests = len([t for t in tests if t['status'] == 'FAIL'])
    passed_tests = len([t for t in tests if t['status'] == 'PASS'])

    html_template = '''<!DOCTYPE html>
<html lang="en">
<head>
  <meta charset="UTF-8">
  <title>UAT Test Report</title>
  <style>
    body {{ font-family: Arial, sans-serif; margin: 2em; }}
    h1 {{ color: #333; }}
    .summary {{ margin-bottom: 2em; }}
    .summary span {{ font-weight: bold; }}
    .summary .pass {{ color: green; }}
    .summary .fail {{ color: red; }}
    table {{ border-collapse: collapse; width: 100%; }}
    th, td {{ border: 1px solid #ccc; padding: 0.5em; text-align: left; vertical-align: top; }}
    th {{ background: #f0f0f0; }}
    .pass {{ color: green; font-weight: bold; }}
    .fail {{ color: red; font-weight: bold; }}
    .test-name {{ font-weight: bold; margin-bottom: 10px; }}
    .test-logs {{ margin-top: 10px; }}
    .logs-toggle {{ cursor: pointer; background: #f5f5f5; padding: 5px 10px; border: 1px solid #ddd; display: inline-block; border-radius: 3px; font-size: 12px; }}
    .logs-content {{ display: none; background: #f9f9f9; border: 1px solid #ddd; padding: 10px; margin-top: 5px; font-family: monospace; font-size: 11px; white-space: pre-wrap; max-height: 400px; overflow-y: auto; }}
    .logs-content.show {{ display: block; }}
    .failed-test {{ border: 3px solid #ff0000 !important; background-color: #fff5f5; }}
    .failed-test .test-name {{ color: #cc0000; }}
    .passed-test {{ border: 3px solid #00cc00 !important; background-color: #f5fff5; }}
    .passed-test .test-name {{ color: #006600; }}
  </style>
</head>
<body>
  <h1>UAT Test Report</h1>
  <div class="summary">
    <div>Generated: <span>{timestamp}</span></div>
    <div>Total executed tests: <span class="total">{total_tests}</span></div>
    <div>Failures: <span class="fail">{failed_tests}</span></div>
    <div>Passes: <span class="pass">{passed_tests}</span></div>
  </div>
  <table>
    <thead>
      <tr>
        <th>Test Name & Logs</th>
        <th>Status</th>
        <th>Duration</th>
      </tr>
    </thead>
    <tbody>
{test_rows}
    </tbody>
  </table>

  <script>
    function toggleLogs(element) {{
      const content = element.nextElementSibling;
      const isVisible = content.classList.contains('show');

      if (isVisible) {{
        content.classList.remove('show');
        element.textContent = 'Show Logs';
      }} else {{
        content.classList.add('show');
        element.textContent = 'Hide Logs';
      }}
    }}
  </script>
</body>
</html>'''

    # Generate test rows
    test_rows = []
    for test in tests:
        status_class = 'fail' if test['status'] == 'FAIL' else 'pass'
        row_class = 'failed-test' if test['status'] == 'FAIL' else 'passed-test'
        logs_visible = 'show' if test['status'] == 'FAIL' else ''
        logs_toggle_text = 'Hide Logs' if test['status'] == 'FAIL' else 'Show Logs'

        # Escape HTML in logs
        logs_content = '\n'.join(test['logs']).replace('&', '&amp;').replace('<', '&lt;').replace('>', '&gt;')

        row = f'''      <tr class="{row_class}">
        <td>
          <div class="test-name">{test['name']}</div>
          <div class="test-logs">
            <div class="logs-toggle" onclick="toggleLogs(this)">{logs_toggle_text}</div>
            <div class="logs-content {logs_visible}">{logs_content}</div>
          </div>
        </td>
        <td class="{status_class}">{test['status']}</td>
        <td>{test['duration']}</td>
      </tr>'''
        test_rows.append(row)

    # Fill template
    html_content = html_template.format(
        timestamp=datetime.now().strftime('%Y-%m-%d %H:%M:%S'),
        total_tests=total_tests,
        failed_tests=failed_tests,
        passed_tests=passed_tests,
        test_rows='\n'.join(test_rows)
    )

    with open(output_file, 'w', encoding='utf-8') as f:
        f.write(html_content)

def main():
    if len(sys.argv) < 2:
        print("Usage: python generate_test_report.py <log_file> [output_file]")
        print("Example: python generate_test_report.py uat_test_report.log")
        sys.exit(1)

    log_file = sys.argv[1]
    if not os.path.exists(log_file):
        print(f"Error: Log file '{log_file}' not found")
        sys.exit(1)

    # Determine output file
    if len(sys.argv) >= 3:
        output_file = sys.argv[2]
    else:
        base_name = os.path.splitext(os.path.basename(log_file))[0]
        output_dir = os.path.dirname(log_file)
        output_file = os.path.join(output_dir, f"{base_name}.html")

    print(f"Parsing log file: {log_file}")
    tests = parse_log_file(log_file)

    print(f"Found {len(tests)} tests")
    print(f"Generating HTML report: {output_file}")
    generate_html_report(tests, output_file)

    print("Report generated successfully!")
    print(f"Open {output_file} in your browser to view the report")

if __name__ == "__main__":
    main()
