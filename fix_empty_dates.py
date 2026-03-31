import sqlite3

db_path = 'data/drama_generator.db'
con = sqlite3.connect(db_path)
cursor = con.cursor()

def fix_table_dates(table_name):
    # Get columns
    cursor.execute(f"PRAGMA table_info({table_name})")
    cols = [col[1] for col in cursor.fetchall()]
    
    date_cols = [c for c in cols if c in ('created_at', 'updated_at', 'deleted_at', 'completed_at')]
    
    if not date_cols:
        return
        
    for col in date_cols:
        print(f"Fixing {col} for {table_name}")
        # Fix empty strings to NULL if nullable, or a safe date if not nullable
        # Gorm deleted_at is always nullable. created_at/updated_at are usually not nullable, but checking them anyway.
        
        # 1. Update completely empty or string literal "NULL" or "None" to SQLite NULL
        cursor.execute(f"UPDATE {table_name} SET {col} = NULL WHERE {col} = '' OR {col} = 'NULL' OR {col} = 'None' OR {col} IS NULL;")
        
        # 2. For created_at / updated_at, if they ended up NULL, set them to a default safe timestamp
        if col in ('created_at', 'updated_at'):
            cursor.execute(f"UPDATE {table_name} SET {col} = '2026-01-01 00:00:00' WHERE {col} IS NULL;")

try:
    cursor.execute("SELECT name FROM sqlite_master WHERE type='table'")
    tables = cursor.fetchall()
    
    for table_row in tables:
        table = table_row[0]
        # Ignore system tables
        if table.startswith('sqlite_'): continue
        fix_table_dates(table)
        
    con.commit()
    print("Empty dates successfully patched.")
except Exception as e:
    print(f"Error: {e}")

con.close()
