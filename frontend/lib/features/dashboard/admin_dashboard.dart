import 'package:flutter/material.dart';
import 'package:flutter_riverpod/flutter_riverpod.dart';
import 'package:fir_management/core/providers/fir_provider.dart';
import 'package:fir_management/features/shared/app_shell.dart';

class AdminDashboard extends ConsumerWidget {
  const AdminDashboard({super.key});

  @override
  Widget build(BuildContext context, WidgetRef ref) {
    final analytics = ref.watch(firAnalyticsProvider);
    return AppShell(
      title: 'Admin Dashboard',
      child: analytics.when(
        loading: () => const Center(child: CircularProgressIndicator()),
        error: (e, _) => Center(child: Text('$e')),
        data: (stats) => ListView(
          padding: const EdgeInsets.all(16),
          children: [
            Text('System Overview', style: Theme.of(context).textTheme.headlineSmall),
            const SizedBox(height: 16),
            ListTile(
              leading: const Icon(Icons.folder),
              title: Text('Total FIRs: ${stats['totalFirs']}'),
            ),
            ListTile(
              leading: const Icon(Icons.pending_actions),
              title: Text('Pending: ${stats['pendingFirs']}'),
            ),
            ListTile(
              leading: const Icon(Icons.verified),
              title: Text('Approved: ${stats['approvedFirs']}'),
            ),
            const Divider(),
            Text('Crime Categories', style: Theme.of(context).textTheme.titleMedium),
            ...((stats['byCrimeType'] as Map?) ?? {}).entries.map(
              (e) => ListTile(title: Text('${e.key}: ${e.value}')),
            ),
          ],
        ),
      ),
    );
  }
}
