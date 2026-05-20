import 'package:flutter/material.dart';
import 'package:flutter_riverpod/flutter_riverpod.dart';
import 'package:go_router/go_router.dart';
import 'package:fir_management/core/providers/fir_provider.dart';
import 'package:fir_management/features/shared/app_shell.dart';

class CitizenDashboard extends ConsumerWidget {
  const CitizenDashboard({super.key});

  @override
  Widget build(BuildContext context, WidgetRef ref) {
    final firs = ref.watch(firListProvider);
    return AppShell(
      title: 'Citizen Dashboard',
      child: firs.when(
        loading: () => const Center(child: CircularProgressIndicator()),
        error: (e, _) => Center(child: Text('Error: $e')),
        data: (items) => ListView.builder(
          padding: const EdgeInsets.all(16),
          itemCount: items.length,
          itemBuilder: (_, i) {
            final fir = items[i];
            return Card(
              child: ListTile(
                title: Text(fir.title),
                subtitle: Text('${fir.firId} • ${fir.status}'),
                trailing: const Icon(Icons.chevron_right),
                onTap: () => context.push('/fir/${fir.firId}'),
              ),
            );
          },
        ),
      ),
    );
  }
}
