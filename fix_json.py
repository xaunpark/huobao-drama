import sqlite3

db_path = 'data/drama_generator.db'
con = sqlite3.connect(db_path)
cursor = con.cursor()

try:
    cursor.execute("SELECT id, prompts FROM prompt_templates")
    rows = cursor.fetchall()
    
    for row in rows:
        id, prompts_raw = row
        if not prompts_raw or prompts_raw == '' or prompts_raw == b'' or str(prompts_raw).strip() == '':
            print(f"Fixing empty JSON for template ID {id}")
            cursor.execute("UPDATE prompt_templates SET prompts = '{}' WHERE id = ?", (id,))
        elif isinstance(prompts_raw, str) and not (prompts_raw.startswith('{') or prompts_raw.startswith('[')):
            print(f"Fixing invalid JSON string for template ID {id}: {prompts_raw[:20]}")
            cursor.execute("UPDATE prompt_templates SET prompts = '{}' WHERE id = ?", (id,))
            
    con.commit()
    print("Database JSON patch completed.")
except Exception as e:
    print(f"Failed to update: {e}")

con.close()
