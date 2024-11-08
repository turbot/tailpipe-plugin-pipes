import json
from datetime import datetime, timedelta
import random
import pandas as pd

class BaseAuditGenerator:
    def __init__(self):
        self.base_time = datetime.now() - timedelta(days=90)
        
    def generate_id(self, prefix='a'):
        """Generate a random ID with given prefix."""
        return f"{prefix}_{''.join(random.choices('abcdefghijklmnopqrstuvwxyz0123456789', k=20))}"
        
    def save_to_parquet(self, records, filename="audit_logs.parquet"):
        """Save logs to Parquet format for DuckDB."""
        records.sort(key=lambda x: x["created_at"])
        df = pd.DataFrame(records)
        df.to_parquet(filename, index=False)
        print(f"Generated {len(records)} records")

    def generate_timestamp(self, user_type):
        """Generate timestamp based on user type and patterns."""
        if user_type == 'suspicious':
            # Higher chance of off-hours activity
            hour = random.choices(
                range(24),
                weights=[4 if h in range(1,5) else 1 for h in range(24)]
            )[0]
        else:
            # Normal business hours with some after-hours activity
            hour = random.choices(
                range(24),
                weights=[3 if h in range(9,18) else 1 for h in range(24)]
            )[0]

        days_ago = random.randint(0, 89)  # Up to 90 days of history
        return self.base_time - timedelta(days=days_ago, hours=-hour, minutes=random.randint(0,59))
