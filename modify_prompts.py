import re

filepath = r'g:\\VS-Project\\huobao-drama\\application\\services\\prompt_i18n.go'
with open(filepath, 'r', encoding='utf-8') as f:
    content = f.read()

# 1. Replace the fallback return in IsEnglish blocks.
def replacer(match):
    en_ret = match.group(1)
    return f"\n\tif p.IsEnglish() {{\n\t\treturn {en_ret}\n\t}}\n\n\treturn {en_ret}\n}}"

content = re.sub(r"\n\tif p\.IsEnglish\(\) \{\n\t\treturn (.*?)\n\t\}\n\n\treturn (.*?)\n\}", replacer, content, flags=re.DOTALL)

# 2. FormatUserPrompt replacements
templates_match = re.search(r'templates := map\[string\]map\[string\]string\{\s*"en": \{(.*?)\},\s*"zh": \{(.*?)\},\s*}', content, re.DOTALL)
if templates_match:
    en_dict = templates_match.group(1)
    zh_dict = templates_match.group(2)
    to_replace = f'"zh": {{{zh_dict}}}'
    new_str = f'"zh": {{{en_dict}}}'
    content = content.replace(to_replace, new_str)

# 3. GetStylePrompt replacements
style_match = re.search(r'stylePrompts := map\[string\]map\[string\]string\{\s*"zh": \{(.*?)\},\s*"en": \{(.*?)\},\s*\}', content, re.DOTALL)
if style_match:
    zh_dict = style_match.group(1)
    en_dict = style_match.group(2)
    to_replace = f'"zh": {{{zh_dict}}}'
    new_str = f'"zh": {{{en_dict}}}'
    content = content.replace(to_replace, new_str)

# 4. GetVideoConstraintPrompt replacements
# actionSequencePrompts
action_match = re.search(r'actionSequencePrompts := map\[string\]string\{\s*"zh": (`.*?`),\s*"en": (`.*?`),\s*\}', content, re.DOTALL)
if action_match:
    zh_val = action_match.group(1)
    en_val = action_match.group(2)
    to_replace = f'"zh": {zh_val}'
    new_str = f'"zh": {en_val}'
    content = content.replace(to_replace, new_str)

# generalPrompts
general_match = re.search(r'generalPrompts := map\[string\]string\{\s*"zh": (`.*?`),\s*"en": (`.*?`),\s*\}', content, re.DOTALL)
if general_match:
    zh_val = general_match.group(1)
    en_val = general_match.group(2)
    to_replace = f'"zh": {zh_val}'
    new_str = f'"zh": {en_val}'
    content = content.replace(to_replace, new_str)

with open(filepath, 'w', encoding='utf-8') as f:
    f.write(content)

print("Replacement done successfully.")
