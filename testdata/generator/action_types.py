class ActionTypes:
    @staticmethod
    def get_action_groups():
        return {
            'common_actions': {
                'workspace.create': 0.08,
                'workspace.update': 0.1,
                'connection.create': 0.08,
                'connection.update': 0.1,
                'pipeline.create': 0.06,
                'pipeline.update': 0.08,
                'pipeline.command.run': 0.15,
                'workspace.schema.create': 0.05,
                'integration.refresh': 0.1,
            },
            'maintenance_actions': {
                'workspace.reboot': 0.02,
                'workspace.upgrade': 0.02,
                'workspace.usage.enable': 0.01,
                'workspace.usage.disable': 0.01,
                'workspace.inactive.warn': 0.03,
                'workspace.inactive.disable': 0.01,
            },
            'org_management': {
                'org.member.add': 0.02,
                'org.member.update': 0.02,
                'org.member.delete': 0.01,
                'org.workspace.member.add': 0.02,
                'org.workspace.member.update': 0.02,
                'org.update': 0.01,
                'org.billing.payment_method.create': 0.005,
            },
            'destructive_actions': {
                'workspace.delete': 0.01,
                'connection.delete': 0.01,
                'pipeline.delete': 0.01,
                'workspace.schema.delete': 0.01,
                'integration.delete': 0.005,
                'datatank.delete': 0.005,
                'workspace.mod.uninstall': 0.01,
            }
        }

    @staticmethod
    def get_action_patterns():
        return {
            'workspace_lifecycle': [
                'workspace.create',
                'workspace.schema.create',
                'connection.create',
                'workspace.update',
                'workspace.delete'
            ],
            'pipeline_lifecycle': [
                'pipeline.create',
                'pipeline.update',
                'pipeline.command.run',
                'pipeline.delete'
            ],
            'maintenance_cycle': [
                'workspace.inactive.warn',
                'workspace.usage.disable',
                'workspace.inactive.disable'
            ]
        }
