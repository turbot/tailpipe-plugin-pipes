# network.py
import random

class NetworkConfig:
    @staticmethod
    def get_ip_patterns():
        return {
            'internal_office': [f"10.0.{i}.{j}" for i in range(0,3) for j in range(1,10)],
            'vpn_pool': [f"172.16.{i}.{j}" for i in range(0,2) for j in range(1,10)],
            'suspicious': [
                f"45.123.{i}.{j}" for i in range(0,3) for j in range(1,10)
            ] + [f"103.{i}.45.{j}" for i in range(100,102) for j in range(1,10)]
        }

    @staticmethod
    def get_ip_for_user(user, ip_patterns):
        """Get IP address based on user type."""
        if user['type'] == 'suspicious':
            # Higher chance of IP hopping
            if random.random() < 0.4:
                return random.choice(ip_patterns['suspicious'])
            return user['primary_ip']
        else:
            # Normal users mostly use their primary IP, sometimes VPN
            if random.random() < 0.2:  # 20% chance of VPN usage
                return random.choice(ip_patterns['vpn_pool'])
            return user['primary_ip']