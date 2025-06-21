cat > consolidate_report.sh << 'EOF'
#!/bin/bash
python3 << 'PY'
import xml.etree.ElementTree as ET
tree = ET.parse('uat_test_report_raw.xml')
new_root = ET.Element('testsuites')
new_suite = ET.SubElement(new_root, 'testsuite', name='Combined Tests')
tests = failures = time = 0
for case in tree.findall('.//testcase'):
    new_suite.append(case)
    tests += 1
    if case.find('failure') is not None: failures += 1
    time += float(case.get('time', 0))
new_suite.set('tests', str(tests))
new_suite.set('failures', str(failures))
new_suite.set('time', str(round(time, 3)))
ET.ElementTree(new_root).write('uat_test_report.xml')
PY
EOF
chmod +x consolidate_report.sh