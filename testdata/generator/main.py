# main.py
from base_generator import BaseAuditGenerator
from action_types import ActionTypes
from network import NetworkConfig
from user import UserProfileGenerator
from record_generator import RecordGenerator

def main():
    # Initialize base components
    base_gen = BaseAuditGenerator()
    ip_patterns = NetworkConfig.get_ip_patterns()
    action_groups = ActionTypes.get_action_groups()
    action_patterns = ActionTypes.get_action_patterns()
    
    # Generate user profiles
    user_gen = UserProfileGenerator(ip_patterns)
    users = user_gen.generate_profiles(50)
    
    # Initialize record generator
    record_gen = RecordGenerator(base_gen, action_groups, action_patterns, ip_patterns)
    
    # Generate records
    records = []
    for user_type, user_list in users.items():
        for user in user_list:
            # Regular activity
            for _ in range(100 if user_type == 'normal' else 300):
                timestamp = base_gen.generate_timestamp(user_type)
                action = record_gen.get_action_for_user_type(user_type, timestamp.hour)
                records.append(record_gen.generate_base_record(
                    timestamp=timestamp,
                    user=user,
                    action_type=action,
                    ip=NetworkConfig.get_ip_for_user(user, ip_patterns)
                ))
            
            # Related action sequences
            for _ in range(10 if user_type == 'normal' else 30):
                records.extend(record_gen.generate_related_actions(
                    base_gen.generate_timestamp(user_type),
                    user
                ))
    
    # Save to parquet
    base_gen.save_to_parquet(records)

if __name__ == "__main__":
    main()