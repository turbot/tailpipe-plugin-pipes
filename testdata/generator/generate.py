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

class PipesAuditGenerator(BaseAuditGenerator):
    def __init__(self):
        super().__init__()
        self.action_types = ['workspace.delete', 'workspace.create', 'workspace.update', 
                           'connection.create', 'connection.delete', 'user.login']
        
        # Define users with their patterns
        self.users = [
            {
                'handle': 'vhadianto',
                'display_name': 'Victor Hadianto',
                'type': 'normal'
            },
            {
                'handle': 'jsmith',
                'display_name': 'John Smith',
                'type': 'normal'
            },
            {
                'handle': 'awhite',
                'display_name': 'Alice White',
                'type': 'suspicious'
            },
            {
                'handle': 'rjohnson',
                'display_name': 'Robert Johnson',
                'type': 'normal'
            }
        ]
        
        # Generate a pool of IPs
        self.ips = [
            f"159.196.{random.randint(1, 255)}.{random.randint(1, 255)}" 
            for _ in range(20)
        ]

    def generate_data_object(self, timestamp, user, workspace_id):
        """Generate the nested data object."""
        return {
            "api_version": "1.8.13",
            "cli_version": "0.19.5",
            "created_at": timestamp.strftime("%Y-%m-%dT%H:%M:%SZ"),
            "created_by": {
                "avatar_url": f"https://cloud.steampipe.io/api/v0/identity/{user['handle']}/avatar?v=57",
                "display_name": user['display_name'],
                "handle": user['handle'],
                "id": self.generate_id('u'),
                "status": "accepted",
                "updated_at": "2023-03-28T22:44:59Z",
                "version_id": random.randint(1, 100)
            },
            "database_name": "c6bcy0",
            "handle": f"workspace-{random.randint(1, 1000)}",
            "id": workspace_id,
            "state": random.choice(["enabled", "disabled"])
        }

    def generate_record(self):
        """Generate a single audit log record."""
        user = random.choice(self.users)
        timestamp = self.generate_timestamp(user['type'])
        
        # Generate IDs
        workspace_id = self.generate_id('w')
        user_id = self.generate_id('u')
        ip = random.choice(self.ips)
        
        data_obj = self.generate_data_object(timestamp, user, workspace_id)
        
        return {
            'identity_id': self.generate_id('o'),
            'tenant_id': 't_00000000000000000000',
            'tp_id': self.generate_id('cs'),
            'tp_source_type': 'pipes_audit_log_api',
            'actor_display_name': user['display_name'],
            'data': json.dumps(data_obj),
            'target_handle': data_obj['handle'],
            'tp_destination_ip': '',
            'tp_index': 'acme-tank',
            'process_id': '',
            'target_id': workspace_id,
            'actor_avatar_url': f"https://cloud.steampipe.io/api/v0/identity/{user_id}/avatar?v=57",
            'tp_akas': [workspace_id],
            'tp_date': timestamp.strftime('%Y-%m-%d'),
            'actor_ip': ip,
            'tp_usernames': [user['handle'], user_id],
            'tp_ingest_timestamp': int(timestamp.timestamp() * 1000),
            'tp_emails': '',
            'tp_source_ip': ip,
            'actor_id': user_id,
            'id': self.generate_id('a'),
            'action_type': random.choice(self.action_types),
            'created_at': timestamp.strftime('%Y-%m-%d %H:%M:%S'),
            'tp_source_name': 'pipes_audit_log_api',
            'tp_timestamp': int(timestamp.timestamp() * 1000),
            'tp_partition': 'pipes_testing',
            'tp_source_location': 'pipes.turbot.com:pipes-testing',
            'tp_ips': [ip],
            'tp_tags': '',
            'actor_handle': user['handle'],
            'identity_handle': 'acme-tank'
        }

    def generate_logs(self, num_records=10000):
        """Generate multiple audit log records."""
        records = [self.generate_record() for _ in range(num_records)]
        self.save_to_parquet(records, "pipes_audit_log.parquet")
        return records

if __name__ == "__main__":
    generator = PipesAuditGenerator()
    generator.generate_logs(10000)