import os

dir_path = r'g:\VS-Project\huobao-drama\application\services'

for filename in os.listdir(dir_path):
    if not filename.endswith('.go'):
        continue
    filepath = os.path.join(dir_path, filename)
    with open(filepath, 'r', encoding='utf-8') as f:
        content = f.read()
    
    # We want to force the English branch.
    # Replace:
    # `if s.promptI18n.IsEnglish() {` with `if true || s.promptI18n.IsEnglish() {`
    # `if p.IsEnglish() {` with `if true || p.IsEnglish() {`
    
    new_content = content.replace('if s.promptI18n.IsEnglish() {', 'if true || s.promptI18n.IsEnglish() {')
    new_content = new_content.replace('if p.IsEnglish() {', 'if true || p.IsEnglish() {')
    
    if new_content != content:
        with open(filepath, 'w', encoding='utf-8') as f:
            f.write(new_content)
        print(f"Updated {filename}")
