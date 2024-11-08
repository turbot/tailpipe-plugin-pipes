# user.py
import random

class UserProfileGenerator:
    def __init__(self, ip_patterns):
        self.ip_patterns = ip_patterns

    def generate_profiles(self, num_users=50):
        """Generate different types of user profiles."""
        users = {
            'normal': [],      # Regular users
            'power': [],       # DevOps/Engineers
            'admin': [],       # Organization admins
            'suspicious': []   # Potentially compromised accounts
        }

        # Distribution of user types (50 total)
        user_types = (
            ['normal'] * 35 +  # 70% regular users
            ['power'] * 10 +   # 20% power users
            ['admin'] * 3 +    # 6% admins
            ['suspicious'] * 2  # 4% suspicious
        )

        for i in range(num_users):
            user_type = user_types[i] if i < len(user_types) else 'normal'
            user_id = f"u_{''.join(random.choices('abcdefghijklmnopqrstuvwxyz0123456789', k=20))}"
            handle = f"user_{i:03d}"
            
            user = {
                "id": user_id,
                "handle": handle,
                "display_name": f"User {i:03d}",
                "avatar_url": f"https://example.com/avatars/{handle}",
                "type": user_type,
                "primary_ip": random.choice(
                    self.ip_patterns['internal_office'] if user_type != 'suspicious'
                    else self.ip_patterns['suspicious']
                )
            }
            
            users[user_type].append(user)
            
        return users