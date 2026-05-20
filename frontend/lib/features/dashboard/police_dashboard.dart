import 'package:fl_chart/fl_chart.dart';
import 'package:flutter/material.dart';
import 'package:flutter_riverpod/flutter_riverpod.dart';
import 'package:go_router/go_router.dart';
import 'package:fir_management/core/providers/fir_provider.dart';
import 'package:fir_management/features/shared/app_shell.dart';

class PoliceDashboard extends ConsumerWidget {
  const PoliceDashboard({super.key});

  @override
  Widget build(BuildContext context, WidgetRef ref) {
    final analytics = ref.watch(firAnalyticsProvider);
    final firs = ref.watch(firListProvider);

    return AppShell(
      title: 'Police Dashboard',
      child: analytics.when(
        loading: () => const Center(child: CircularProgressIndicator()),
        error: (e, _) => Center(child: Text('$e')),
        data: (stats) {
          final byStatus = (stats['byStatus'] as Map?)?.cast<String, dynamic>() ?? {};
          final sections = byStatus.entries.map((e) {
            return PieChartSectionData(
              value: (e.value as num).toDouble().clamp(1, double.infinity),
              title: e.key,
              radius: 50,
            );
          }).toList();

          return CustomScrollView(
            slivers: [
              SliverToBoxAdapter(
                child: Padding(
                  padding: const EdgeInsets.all(16),
                  child: Wrap(
                    spacing: 12,
                    runSpacing: 12,
                    children: [
                      _StatCard('Total', '${stats['totalFirs'] ?? 0}'),
                      _StatCard('Pending', '${stats['pendingFirs'] ?? 0}'),
                      _StatCard('Approved', '${stats['approvedFirs'] ?? 0}'),
                    ],
                  ),
                ),
              ),
              if (sections.isNotEmpty)
                SliverToBoxAdapter(
                  child: SizedBox(height: 220, child: PieChart(PieChartData(sections: sections))),
                ),
              firs.when(
                loading: () => const SliverFillRemaining(child: Center(child: CircularProgressIndicator())),
                error: (e, _) => SliverToBoxAdapter(child: Text('$e')),
                data: (items) => SliverList(
                  delegate: SliverChildBuilderDelegate(
                    (context, i) {
                      final fir = items[i];
                      return ListTile(
                        title: Text(fir.title),
                        subtitle: Text('${fir.firId} • ${fir.status}'),
                        onTap: () => context.push('/fir/${fir.firId}'),
                      );
                    },
                    childCount: items.length,
                  ),
                ),
              ),
            ],
          );
        },
      ),
    );
  }
}

class _StatCard extends StatelessWidget {
  const _StatCard(this.label, this.value);
  final String label;
  final String value;

  @override
  Widget build(BuildContext context) {
    return Card(
      child: Padding(
        padding: const EdgeInsets.symmetric(horizontal: 24, vertical: 16),
        child: Column(
          children: [Text(label), const SizedBox(height: 4), Text(value, style: Theme.of(context).textTheme.headlineSmall)],
        ),
      ),
    );
  }
}
