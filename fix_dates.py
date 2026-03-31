import sqlite3
from datetime import datetime
import re

db_path = 'data/drama_generator.db'
con = sqlite3.connect(db_path)
cursor = con.cursor()

def convert_to_sqlite_datetime(timestamp):
    if not timestamp: return None
    # Assuming the string might be some ISO8601 that Go's parser dislikes like "2026-03-31T06:48:58Z"
    # or just missing values
    ts_str = str(timestamp).strip()
    if not ts_str: return None
    
    # Simple regex replacing T and Z
    ts_str = re.sub(r'T', ' ', ts_str)
    ts_str = re.sub(r'Z.*', '', ts_str)
    ts_str = re.sub(r'\+.*', '', ts_str)
    
    # Truncate microseconds if they are weird
    parts = ts_str.split('.')
    if len(parts) > 1:
        ts_str = parts[0] + '.' + parts[1][:6]
        
    return ts_str

try:
    cursor.execute("SELECT id, created_at, updated_at, deleted_at FROM prompt_templates")
    rows = cursor.fetchall()
    
    for row in rows:
        id, created_at, updated_at, deleted_at = row
        print(f"ID {id} | created: {created_at} | updated: {updated_at} | deleted: {deleted_at}")
        
        new_c = convert_to_sqlite_datetime(created_at)
        new_u = convert_to_sqlite_datetime(updated_at)
        new_d = convert_to_sqlite_datetime(deleted_at)
        
        if new_c != created_at or new_u != updated_at or new_d != deleted_at:
            cursor.execute("UPDATE prompt_templates SET created_at=?, updated_at=?, deleted_at=? WHERE id=?", 
                           (new_c, new_u, new_d, id))

    con.commit()
    print("Dates checked/patched.")
except Exception as e:
    print(f"Error: {e}")

con.close()
