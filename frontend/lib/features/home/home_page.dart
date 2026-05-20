import 'package:flutter/material.dart';
import 'package:go_router/go_router.dart';

class HomePage extends StatelessWidget {
  const HomePage({super.key});

  @override
  Widget build(BuildContext context) {
    final theme = Theme.of(context);
    return Scaffold(
      body: SafeArea(
        child: Padding(
          padding: const EdgeInsets.all(24),
          child: Column(
            crossAxisAlignment: CrossAxisAlignment.start,
            children: [
              Text('Blockchain FIR', style: theme.textTheme.headlineLarge?.copyWith(fontWeight: FontWeight.bold)),
              const SizedBox(height: 8),
              Text(
                'Tamper-proof First Information Reports with IPFS evidence and on-chain verification.',
                style: theme.textTheme.bodyLarge?.copyWith(color: theme.colorScheme.onSurfaceVariant),
              ),
              const Spacer(),
              _FeatureCard(
                icon: Icons.security,
                title: 'Immutable Records',
                subtitle: 'FIR metadata hashes stored on Polygon Mumbai.',
              ),
              const SizedBox(height: 12),
              _FeatureCard(
                icon: Icons.folder_shared,
                title: 'IPFS Evidence',
                subtitle: 'Documents pinned securely with CID tracking.',
              ),
              const SizedBox(height: 12),
              _FeatureCard(
                icon: Icons.badge,
                title: 'Role-Based Access',
                subtitle: 'Citizens, police officers, and admins.',
              ),
              const Spacer(),
              FilledButton(onPressed: () => context.go('/login'), child: const Text('Sign In')),
              const SizedBox(height: 12),
              OutlinedButton(onPressed: () => context.go('/register'), child: const Text('Register as Citizen')),
            ],
          ),
        ),
      ),
    );
  }
}

class _FeatureCard extends StatelessWidget {
  const _FeatureCard({required this.icon, required this.title, required this.subtitle});
  final IconData icon;
  final String title;
  final String subtitle;

  @override
  Widget build(BuildContext context) {
    return Card(
      child: ListTile(
        leading: Icon(icon),
        title: Text(title, style: const TextStyle(fontWeight: FontWeight.w600)),
        subtitle: Text(subtitle),
      ),
    );
  }
}
