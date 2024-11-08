# record_generator.py
import random
from datetime import timedelta
from network import NetworkConfig

class RecordGenerator:
    def __init__(self, base_generator, action_groups, action_patterns, ip_patterns):
        self.base = base_generator
        self.action_groups = action_groups
        self.action_patterns = action_patterns
        self.ip_patterns = ip_patterns

    def get_action_for_user_type(self, user_type, time_of_day):
        """Get weighted random action based on user type and time of day."""
        if user_type == 'suspicious':
            if time_of_day not in range(9, 18):
                weights = {**self.action_groups['destructive_actions'], **self.action_groups['org_management']}
            else:
                weights = {**self.action_groups['common_actions'], **self.action_groups['destructive_actions']}
        elif user_type == 'admin':
            if time_of_day in range(9, 18):
                weights = {**self.action_groups['org_management'], **self.action_groups['maintenance_actions']}
            else:
                weights = {**self.action_groups['common_actions'], **self.action_groups['maintenance_actions']}
        elif user_type == 'power':
            weights = {**self.action_groups['common_actions'], **self.action_groups['maintenance_actions']}
        else:
            weights = self.action_groups['common_actions']
            
        actions, weights = zip(*weights.items())
        return random.choices(actions, weights=weights)[0]

    def generate_base_record(self, timestamp, user, action_type, ip):
        """Generate a base audit log record."""
        return {
            "id": self.base.generate_id('a'),
            "created_at": timestamp.strftime("%Y-%m-%d %H:%M:%S"),
            "action_type": action_type,
            "actor_id": user["id"],
            "actor_handle": user["handle"],
            "actor_display_name": user["display_name"],
            "actor_type": user["type"],
            "actor_ip": ip,
            "identity_id": self.base.generate_id('o'),
            "target_id": self.base.generate_id('w'),
            "target_handle": f"resource_{random.randint(1000, 9999)}"
        }

    def generate_related_actions(self, base_time, user, num_actions=3):
        """Generate a sequence of related actions."""
        records = []
        current_time = base_time
        pattern = random.choice(list(self.action_patterns.values()))
        
        for action in pattern:
            records.append(self.generate_base_record(
                timestamp=current_time,
                user=user,
                action_type=action,
                ip=NetworkConfig.get_ip_for_user(user, self.ip_patterns)
            ))
            current_time += timedelta(minutes=random.randint(1, 30))
            
        return records