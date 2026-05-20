import 'package:flutter/material.dart';
import 'package:flutter_riverpod/flutter_riverpod.dart';
import 'package:go_router/go_router.dart';
import 'package:intl/intl.dart';
import 'package:url_launcher/url_launcher.dart';
import 'package:fir_management/core/providers/fir_provider.dart';
import 'package:fir_management/data/services/blockchain_service.dart';
class FIRDetailPage extends ConsumerWidget {
  const FIRDetailPage({super.key, required this.id});
  final String id;

  @override
  Widget build(BuildContext context, WidgetRef ref) {
    final fir = ref.watch(firDetailProvider(id));
    final chain = BlockchainService();

    return Scaffold(
      appBar: AppBar(
        title: const Text('FIR Details'),
        actions: [
          IconButton(
            icon: const Icon(Icons.upload_file),
            onPressed: () => context.push('/fir/$id/evidence'),
          ),
        ],
      ),
      body: fir.when(
        loading: () => const Center(child: CircularProgressIndicator()),
        error: (e, _) => Center(child: Text('$e')),
        data: (f) => ListView(
          padding: const EdgeInsets.all(24),
          children: [
            Text(f.title, style: Theme.of(context).textTheme.headlineSmall),
            const SizedBox(height: 8),
            Chip(label: Text(f.status.toUpperCase())),
            const SizedBox(height: 16),
            _InfoRow('FIR ID', f.firId),
            _InfoRow('Crime Type', f.crimeType),
            _InfoRow('Location', f.location),
            if (f.createdAt != null) _InfoRow('Created', DateFormat.yMMMd().add_jm().format(f.createdAt!)),
            const Divider(height: 32),
            Text('Description', style: Theme.of(context).textTheme.titleMedium),
            const SizedBox(height: 8),
            Text(f.description),
            const Divider(height: 32),
            Text('Blockchain', style: Theme.of(context).textTheme.titleMedium),
            if (f.transactionHash != null)
              ListTile(
                title: const Text('Transaction Hash'),
                subtitle: Text(f.transactionHash!, style: const TextStyle(fontFamily: 'monospace', fontSize: 12)),
                trailing: IconButton(
                  icon: const Icon(Icons.open_in_new),
                  onPressed: () => launchUrl(Uri.parse(chain.explorerTxUrl(f.transactionHash!))),
                ),
              ),
            if (f.blockNumber != null) _InfoRow('Block', '${f.blockNumber}'),
            if (f.ipfsHash != null) _InfoRow('IPFS CID', f.ipfsHash!),
            const SizedBox(height: 24),
            Text('Timeline', style: Theme.of(context).textTheme.titleMedium),
            const SizedBox(height: 8),
            const _TimelineItem('Submitted', 'FIR registered and hash stored on-chain'),
            const _TimelineItem('Under Review', 'Awaiting police verification'),
          ],
        ),
      ),
    );
  }
}

class _InfoRow extends StatelessWidget {
  const _InfoRow(this.label, this.value);
  final String label;
  final String value;

  @override
  Widget build(BuildContext context) {
    return Padding(
      padding: const EdgeInsets.symmetric(vertical: 4),
      child: Row(
        crossAxisAlignment: CrossAxisAlignment.start,
        children: [
          SizedBox(width: 120, child: Text(label, style: TextStyle(color: Theme.of(context).colorScheme.onSurfaceVariant))),
          Expanded(child: Text(value)),
        ],
      ),
    );
  }
}

class _TimelineItem extends StatelessWidget {
  const _TimelineItem(this.title, this.subtitle);
  final String title;
  final String subtitle;

  @override
  Widget build(BuildContext context) {
    return ListTile(
      leading: const Icon(Icons.circle, size: 12),
      title: Text(title),
      subtitle: Text(subtitle),
    );
  }
}
