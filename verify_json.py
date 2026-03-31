import sqlite3
import json

db_path = 'data/drama_generator.db'
con = sqlite3.connect(db_path)
cursor = con.cursor()

try:
    cursor.execute("SELECT id, name, prompts FROM prompt_templates")
    rows = cursor.fetchall()
    
    for row in rows:
        id, name, prompts = row
        print(f"Template ID: {id}, Name: {name}")
        try:
            parsed = json.loads(prompts)
            # print("  Valid JSON!")
        except Exception as e:
            print(f"  INVALID JSON: {e}")
            print(f"  Raw: {prompts[:100]}")
except Exception as e:
    print(f"Failed to query: {e}")

con.close()
