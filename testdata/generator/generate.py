import json
from datetime import datetime, timedelta
import random
import pandas as pd
import numpy as np

class PipesAuditGenerator:
    def __init__(self):
        self.base_time = datetime.now() - timedelta(days=90)
        
        # Expanded locations with geographic distribution
        self.locations = {
            'normal': [
                'pipes.turbot.com:pipes-prod',
                'pipes.turbot.com:pipes-staging',
                'pipes.turbot.com:pipes-dev'
            ],
            'suspicious': [
                'pipes.turbot.com:pipes-testing',
                'cloud.turbot.com:cloud-prod',
                'cloud.turbot.com:cloud-staging',
                'api.turbot.com:api-prod',
                'api.turbot.com:api-staging',
                'steampipe.io:cloud-prod',
                'steampipe.io:cloud-staging'
            ]
        }
        
        # Users with different risk profiles
        self.users = {
            'normal': {
                'vhadianto': {
                    'identity_id': 'o_cgsfhcah6homf3cnq5pg',
                    'actor_id': 'u_c6bc4h9e4mvcg7bf9c9g',
                    'display_name': 'Victor Hadianto',
                    'avatar_url': 'https://cloud.steampipe.io/api/v0/identity/vhadianto/avatar?v=57'
                },
                'jsmith': {
                    'identity_id': 'o_cgsfhcah6homf3cnq5p1',
                    'actor_id': 'u_c6bc4h9e4mvcg7bf9c91',
                    'display_name': 'John Smith',
                    'avatar_url': 'https://cloud.steampipe.io/api/v0/identity/jsmith/avatar?v=12'
                }
            },
            'suspicious': {
                'awhite': {
                    'identity_id': ['o_cgsfhcah6homf3cnq5p5', 'o_cgsfhcah6homf3cnq5p6'],
                    'actor_id': 'u_c6bc4h9e4mvcg7bf9c92',
                    'display_name': 'Alice White',
                    'avatar_url': 'https://cloud.steampipe.io/api/v0/identity/awhite/avatar?v=34'
                },
                'dzhou': {
                    'identity_id': ['o_cgsfhcah6homf3cnq5p7', 'o_cgsfhcah6homf3cnq5p8'],
                    'actor_id': 'u_c6bc4h9e4mvcg7bf9c93',
                    'display_name': 'David Zhou',
                    'avatar_url': 'https://cloud.steampipe.io/api/v0/identity/dzhou/avatar?v=45'
                }
            }
        }

        self.tenant_id = 't_00000000000000000000'
        self.process_id = None  # Empty as per schema
        self.tp_domains = None  # Empty as per schema
        self.tp_emails = None   # Empty as per schema
        self.tp_tags = None     # Empty as per schema

        # Action patterns for generating suspicious sequences
        self.suspicious_patterns = {
            'settings_changes': {
                'actions': ['settings.update'],
                'targets': ['database', 'security', 'access'],
                'min_sequence': 5,
                'max_sequence': 10
            },
            'failed_actions': {
                'actions': ['login.failed', 'security.failed'],
                'targets': ['critical_resource'],
                'min_sequence': 5,
                'max_sequence': 8
            },
            'sensitive_access': {
                'actions': ['workspace.delete', 'database.modify', 'security.update'],
                'targets': ['database', 'workspace', 'sensitive_resource'],
                'min_sequence': 5,
                'max_sequence': 12
            }
        }
        
    def generate_id(self, prefix='a'):
        return f"{prefix}_{''.join(random.choices('abcdefghijklmnopqrstuvwxyz0123456789', k=20))}"

    def generate_suspicious_sequence(self, user_handle, user_data, pattern_type):
        """Generate a sequence of suspicious actions matching a specific pattern"""
        pattern = self.suspicious_patterns[pattern_type]
        sequence_length = random.randint(pattern['min_sequence'], pattern['max_sequence'])
        sequence = []
        
        for _ in range(sequence_length):
            action_type = random.choice(pattern['actions'])
            target = random.choice(pattern['targets'])
            location = random.choice(self.locations['suspicious'])
            
            record = self.generate_record(
                user_handle, 
                user_data, 
                location, 
                action_type,
                'suspicious',
                target
            )
            sequence.append(record)
            
        return sequence

    def generate_record(self, user_handle, user_data, location, action_type, user_type='normal', forced_target=None):
        timestamp = self.base_time + timedelta(
            days=random.randint(0, 89),
            hours=random.randint(0, 23),
            minutes=random.randint(0, 59)
        )

        # For suspicious users, reuse IPs more frequently
        if user_type == 'suspicious':
            ip = f"159.196.{random.choice([100,101,102])}.{random.randint(1,255)}"
        else:
            ip = f"159.196.{random.randint(1,255)}.{random.randint(1,255)}"
        
        # Handle multiple identity_ids for suspicious users
        identity_id = random.choice(user_data['identity_id']) if isinstance(user_data['identity_id'], list) else user_data['identity_id']
        
        # Determine target_handle based on action type or forced target
        if forced_target:
            target_handle = forced_target
        elif 'security' in action_type or 'database' in action_type:
            target_handle = random.choice(['database', 'security', 'access', 'critical_resource'])
        elif action_type == 'workspace.delete':
            target_handle = random.choice(['workspace', 'sensitive_resource'])
        else:
            target_handle = f"workspace-{random.randint(1, 1000)}"

        # Generate base record matching schema
        record = {
            'tp_domains': self.tp_domains,
            'identity_id': identity_id,
            'tenant_id': self.tenant_id,
            'tp_id': self.generate_id('cs'),
            'tp_source_type': 'pipes_audit_log_api',
            'actor_display_name': user_data['display_name'],
            'data': json.dumps({
                'api_version': '1.8.13',
                'cli_version': '0.19.5',
                'created_at': timestamp.strftime('%Y-%m-%dT%H:%M:%SZ'),
                'created_by': {
                    'avatar_url': user_data['avatar_url'],
                    'created_at': '2021-11-18T21:14:45Z',
                    'display_name': user_data['display_name'],
                    'handle': user_handle,
                    'id': user_data['actor_id'],
                    'status': 'accepted',
                    'updated_at': '2023-03-28T22:44:59Z',
                    'version_id': random.randint(1, 100)
                },
                'created_by_id': user_data['actor_id'],
                'database_name': f'c{random.randint(1000,9999)}',
                'desired_state': 'enabled',
                'handle': target_handle,
                'state': 'enabled',
                'version_id': 1
            }),
            'target_handle': target_handle,
            'tp_destination_ip': None,
            'tp_index': 'acme-tank',
            'process_id': self.process_id,
            'target_id': self.generate_id('w'),
            'actor_avatar_url': user_data['avatar_url'],
            'tp_akas': [self.generate_id('w')],
            'tp_date': timestamp.strftime('%Y-%m-%d'),
            'actor_ip': ip,
            'tp_usernames': [user_handle, user_data['actor_id']],
            'tp_ingest_timestamp': int(timestamp.timestamp() * 1000) + random.randint(1000, 9999),
            'tp_emails': self.tp_emails,
            'tp_source_ip': ip,
            'actor_id': user_data['actor_id'],
            'id': self.generate_id('a'),
            'action_type': action_type,
            'created_at': timestamp.strftime('%Y-%m-%d %H:%M:%S'),
            'tp_source_name': 'pipes_audit_log_api',
            'tp_timestamp': int(timestamp.timestamp() * 1000),
            'tp_partition': location.split(':')[1],
            'tp_source_location': location,
            'tp_ips': [ip],
            'tp_tags': self.tp_tags,
            'actor_handle': user_handle,
            'identity_handle': 'acme-tank'
        }
        
        return record

    def generate_logs(self, total_records=500000):
        records = []
        
        # Calculate distribution
        suspicious_records = int(total_records * 0.15)  # 15% suspicious
        normal_records = total_records - suspicious_records
        
        # Generate normal activity
        for user_handle, user_data in self.users['normal'].items():
            user_records = int(normal_records / len(self.users['normal']))
            for _ in range(user_records):
                location = random.choice(self.locations['normal'])
                action_type = random.choice([
                    'workspace.create', 'workspace.update', 'workspace.delete',
                    'connection.create', 'user.login', 'settings.update'
                ])
                records.append(self.generate_record(user_handle, user_data, location, action_type, 'normal'))
        
        # Generate suspicious activity
        for user_handle, user_data in self.users['suspicious'].items():
            # Generate suspicious sequences for each pattern type
            for pattern_type in self.suspicious_patterns.keys():
                records.extend(self.generate_suspicious_sequence(user_handle, user_data, pattern_type))
            
            # Generate additional random suspicious activity
            remaining_records = int(suspicious_records / len(self.users['suspicious'])) - \
                              sum(p['max_sequence'] for p in self.suspicious_patterns.values())
            
            for _ in range(remaining_records):
                location = random.choice(self.locations['suspicious'])
                action_type = random.choice([
                    'workspace.delete', 'login.failed', 'login.suspicious',
                    'security.failed', 'settings.delete'
                ])
                records.append(self.generate_record(user_handle, user_data, location, action_type, 'suspicious'))
        
        # Shuffle records to mix timestamps
        random.shuffle(records)
        
        # Convert to DataFrame and save
        df = pd.DataFrame(records)
        df.to_parquet("pipes_audit_log.parquet", index=False)
        print(f"Generated {len(records)} records")

if __name__ == "__main__":
    generator = PipesAuditGenerator()
    generator.generate_logs(500000)